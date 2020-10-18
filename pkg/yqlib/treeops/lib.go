package treeops

import (
	"bytes"
	"fmt"

	"github.com/elliotchance/orderedmap"
	"gopkg.in/op/go-logging.v1"
	"gopkg.in/yaml.v3"
)

var log = logging.MustGetLogger("yq-treeops")

type PathElementType uint32

const (
	PathKey PathElementType = 1 << iota
	Operation
	SelfReference
	OpenBracket
	CloseBracket
	OpenCollect
	CloseCollect
	Value // e.g. string, int
)

type OperationType struct {
	Type       string
	NumArgs    uint // number of arguments to the op
	Precedence uint
	Handler    OperatorHandler
}

var None = &OperationType{Type: "NONE", NumArgs: 0, Precedence: 0}

var Or = &OperationType{Type: "OR", NumArgs: 2, Precedence: 20, Handler: OrOperator}
var And = &OperationType{Type: "AND", NumArgs: 2, Precedence: 20, Handler: AndOperator}

var Union = &OperationType{Type: "UNION", NumArgs: 2, Precedence: 10, Handler: UnionOperator}
var Intersection = &OperationType{Type: "INTERSECTION", NumArgs: 2, Precedence: 20, Handler: IntersectionOperator}

var Assign = &OperationType{Type: "ASSIGN", NumArgs: 2, Precedence: 40, Handler: AssignOperator}
var Multiply = &OperationType{Type: "MULTIPLY", NumArgs: 2, Precedence: 40, Handler: MultiplyOperator}

var Equals = &OperationType{Type: "EQUALS", NumArgs: 2, Precedence: 40, Handler: EqualsOperator}
var Pipe = &OperationType{Type: "PIPE", NumArgs: 2, Precedence: 45, Handler: PipeOperator}

var Length = &OperationType{Type: "LENGTH", NumArgs: 0, Precedence: 50, Handler: LengthOperator}
var Collect = &OperationType{Type: "COLLECT", NumArgs: 0, Precedence: 50, Handler: CollectOperator}
var RecursiveDescent = &OperationType{Type: "RECURSIVE_DESCENT", NumArgs: 0, Precedence: 50, Handler: RecursiveDescentOperator}

// not sure yet

var Select = &OperationType{Type: "SELECT", NumArgs: 1, Precedence: 50, Handler: SelectOperator}

var DeleteChild = &OperationType{Type: "DELETE", NumArgs: 2, Precedence: 40, Handler: DeleteChildOperator}

// var Splat = &OperationType{Type: "SPLAT", NumArgs: 0, Precedence: 40, Handler: SplatOperator}

// var Exists = &OperationType{Type: "Length", NumArgs: 2, Precedence: 35}
// filters matches if they have the existing path

type PathElement struct {
	PathElementType PathElementType
	OperationType   *OperationType
	Value           interface{}
	StringValue     string
}

// debugging purposes only
func (p *PathElement) toString() string {
	var result string = ``
	switch p.PathElementType {
	case PathKey:
		result = result + fmt.Sprintf("%v", p.Value)
	case SelfReference:
		result = result + fmt.Sprintf("SELF")
	case Operation:
		result = result + fmt.Sprintf("%v", p.OperationType.Type)
	case Value:
		result = result + fmt.Sprintf("%v (%T)", p.Value, p.Value)
	default:
		result = result + "I HAVENT GOT A STRATEGY"
	}
	return result
}

type YqTreeLib interface {
	Get(document int, documentNode *yaml.Node, path string) ([]*CandidateNode, error)
	// GetForMerge(rootNode *yaml.Node, path string, arrayMergeStrategy ArrayMergeStrategy) ([]*NodeContext, error)
	// Update(rootNode *yaml.Node, updateCommand UpdateCommand, autoCreate bool) error
	// New(path string) yaml.Node

	// PathStackToString(pathStack []interface{}) string
	// MergePathStackToString(pathStack []interface{}, arrayMergeStrategy ArrayMergeStrategy) string
}

func NewYqTreeLib() YqTreeLib {
	return &lib{treeCreator: NewPathTreeCreator()}
}

type lib struct {
	treeCreator PathTreeCreator
}

func (l *lib) Get(document int, documentNode *yaml.Node, path string) ([]*CandidateNode, error) {
	nodes := []*CandidateNode{&CandidateNode{Node: documentNode.Content[0], Document: 0}}
	navigator := NewDataTreeNavigator(NavigationPrefs{})
	pathNode, errPath := l.treeCreator.ParsePath(path)
	if errPath != nil {
		return nil, errPath
	}
	return navigator.GetMatchingNodes(nodes, pathNode)
}

//use for debugging only
func NodesToString(collection *orderedmap.OrderedMap) string {
	if !log.IsEnabledFor(logging.DEBUG) {
		return ""
	}

	result := ""
	for el := collection.Front(); el != nil; el = el.Next() {
		result = result + "\n" + NodeToString(el.Value.(*CandidateNode))
	}
	return result
}

func NodeToString(node *CandidateNode) string {
	if !log.IsEnabledFor(logging.DEBUG) {
		return ""
	}
	value := node.Node
	if value == nil {
		return "-- nil --"
	}
	buf := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buf)
	errorEncoding := encoder.Encode(value)
	if errorEncoding != nil {
		log.Error("Error debugging node, %v", errorEncoding.Error())
	}
	encoder.Close()
	return fmt.Sprintf(`D%v, P%v, (%v)::%v`, node.Document, node.Path, value.Tag, buf.String())
}

func KindString(kind yaml.Kind) string {
	switch kind {
	case yaml.ScalarNode:
		return "ScalarNode"
	case yaml.SequenceNode:
		return "SequenceNode"
	case yaml.MappingNode:
		return "MappingNode"
	case yaml.DocumentNode:
		return "DocumentNode"
	case yaml.AliasNode:
		return "AliasNode"
	default:
		return "unknown!"
	}
}
