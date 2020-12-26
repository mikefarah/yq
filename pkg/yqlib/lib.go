package yqlib

import (
	"bytes"
	"container/list"
	"fmt"

	logging "gopkg.in/op/go-logging.v1"
	yaml "gopkg.in/yaml.v3"
)

var log = logging.MustGetLogger("yq-lib")

type OperationType struct {
	Type       string
	NumArgs    uint // number of arguments to the op
	Precedence uint
	Handler    OperatorHandler
}

// operators TODO:
// - cookbook doc for common things
// - mergeEmpty (sets only if the document is empty, do I do that now?)

var Or = &OperationType{Type: "OR", NumArgs: 2, Precedence: 20, Handler: OrOperator}
var And = &OperationType{Type: "AND", NumArgs: 2, Precedence: 20, Handler: AndOperator}

var Union = &OperationType{Type: "UNION", NumArgs: 2, Precedence: 10, Handler: UnionOperator}

var Pipe = &OperationType{Type: "PIPE", NumArgs: 2, Precedence: 30, Handler: PipeOperator}

var Assign = &OperationType{Type: "ASSIGN", NumArgs: 2, Precedence: 40, Handler: AssignUpdateOperator}
var AddAssign = &OperationType{Type: "ADD_ASSIGN", NumArgs: 2, Precedence: 40, Handler: AddAssignOperator}

var AssignAttributes = &OperationType{Type: "ASSIGN_ATTRIBUTES", NumArgs: 2, Precedence: 40, Handler: AssignAttributesOperator}
var AssignStyle = &OperationType{Type: "ASSIGN_STYLE", NumArgs: 2, Precedence: 40, Handler: AssignStyleOperator}
var AssignTag = &OperationType{Type: "ASSIGN_TAG", NumArgs: 2, Precedence: 40, Handler: AssignTagOperator}
var AssignComment = &OperationType{Type: "ASSIGN_COMMENT", NumArgs: 2, Precedence: 40, Handler: AssignCommentsOperator}
var AssignAnchor = &OperationType{Type: "ASSIGN_ANCHOR", NumArgs: 2, Precedence: 40, Handler: AssignAnchorOperator}
var AssignAlias = &OperationType{Type: "ASSIGN_ALIAS", NumArgs: 2, Precedence: 40, Handler: AssignAliasOperator}

var Multiply = &OperationType{Type: "MULTIPLY", NumArgs: 2, Precedence: 45, Handler: MultiplyOperator}
var Add = &OperationType{Type: "ADD", NumArgs: 2, Precedence: 45, Handler: AddOperator}
var Alternative = &OperationType{Type: "ALTERNATIVE", NumArgs: 2, Precedence: 45, Handler: AlternativeOperator}

var Equals = &OperationType{Type: "EQUALS", NumArgs: 2, Precedence: 40, Handler: EqualsOperator}
var CreateMap = &OperationType{Type: "CREATE_MAP", NumArgs: 2, Precedence: 40, Handler: CreateMapOperator}

var ShortPipe = &OperationType{Type: "SHORT_PIPE", NumArgs: 2, Precedence: 45, Handler: PipeOperator}

var Length = &OperationType{Type: "LENGTH", NumArgs: 0, Precedence: 50, Handler: LengthOperator}
var Collect = &OperationType{Type: "COLLECT", NumArgs: 0, Precedence: 50, Handler: CollectOperator}
var GetStyle = &OperationType{Type: "GET_STYLE", NumArgs: 0, Precedence: 50, Handler: GetStyleOperator}
var GetTag = &OperationType{Type: "GET_TAG", NumArgs: 0, Precedence: 50, Handler: GetTagOperator}
var GetComment = &OperationType{Type: "GET_COMMENT", NumArgs: 0, Precedence: 50, Handler: GetCommentsOperator}
var GetAnchor = &OperationType{Type: "GET_ANCHOR", NumArgs: 0, Precedence: 50, Handler: GetAnchorOperator}
var GetAlias = &OperationType{Type: "GET_ALIAS", NumArgs: 0, Precedence: 50, Handler: GetAliasOperator}
var GetDocumentIndex = &OperationType{Type: "GET_DOCUMENT_INDEX", NumArgs: 0, Precedence: 50, Handler: GetDocumentIndexOperator}
var GetFilename = &OperationType{Type: "GET_FILENAME", NumArgs: 0, Precedence: 50, Handler: GetFilenameOperator}
var GetFileIndex = &OperationType{Type: "GET_FILE_INDEX", NumArgs: 0, Precedence: 50, Handler: GetFileIndexOperator}
var GetPath = &OperationType{Type: "GET_PATH", NumArgs: 0, Precedence: 50, Handler: GetPathOperator}

var Explode = &OperationType{Type: "EXPLODE", NumArgs: 1, Precedence: 50, Handler: ExplodeOperator}
var SortKeys = &OperationType{Type: "SORT_KEYS", NumArgs: 1, Precedence: 50, Handler: SortKeysOperator}

var CollectObject = &OperationType{Type: "COLLECT_OBJECT", NumArgs: 0, Precedence: 50, Handler: CollectObjectOperator}
var TraversePath = &OperationType{Type: "TRAVERSE_PATH", NumArgs: 0, Precedence: 50, Handler: TraversePathOperator}
var TraverseArray = &OperationType{Type: "TRAVERSE_ARRAY", NumArgs: 1, Precedence: 50, Handler: TraverseArrayOperator}

var DocumentFilter = &OperationType{Type: "DOCUMENT_FILTER", NumArgs: 0, Precedence: 50, Handler: TraversePathOperator}
var SelfReference = &OperationType{Type: "SELF", NumArgs: 0, Precedence: 50, Handler: SelfOperator}
var ValueOp = &OperationType{Type: "VALUE", NumArgs: 0, Precedence: 50, Handler: ValueOperator}
var Not = &OperationType{Type: "NOT", NumArgs: 0, Precedence: 50, Handler: NotOperator}
var Empty = &OperationType{Type: "EMPTY", NumArgs: 50, Handler: EmptyOperator}

var RecursiveDescent = &OperationType{Type: "RECURSIVE_DESCENT", NumArgs: 0, Precedence: 50, Handler: RecursiveDescentOperator}

var Select = &OperationType{Type: "SELECT", NumArgs: 1, Precedence: 50, Handler: SelectOperator}
var Has = &OperationType{Type: "HAS", NumArgs: 1, Precedence: 50, Handler: HasOperator}
var DeleteChild = &OperationType{Type: "DELETE", NumArgs: 1, Precedence: 40, Handler: DeleteChildOperator}
var DeleteImmediateChild = &OperationType{Type: "DELETE_IMMEDIATE_CHILD", NumArgs: 1, Precedence: 40, Handler: DeleteImmediateChildOperator}

type Operation struct {
	OperationType *OperationType
	Value         interface{}
	StringValue   string
	CandidateNode *CandidateNode // used for Value Path elements
	Preferences   interface{}
}

func CreateValueOperation(value interface{}, stringValue string) *Operation {
	var node yaml.Node = yaml.Node{Kind: yaml.ScalarNode}
	node.Value = stringValue

	switch value.(type) {
	case float32, float64:
		node.Tag = "!!float"
	case int, int64, int32:
		node.Tag = "!!int"
	case bool:
		node.Tag = "!!bool"
	case string:
		node.Tag = "!!str"
	case nil:
		node.Tag = "!!null"
	}

	return &Operation{
		OperationType: ValueOp,
		Value:         value,
		StringValue:   stringValue,
		CandidateNode: &CandidateNode{Node: &node},
	}
}

// debugging purposes only
func (p *Operation) toString() string {
	if p.OperationType == TraversePath {
		return fmt.Sprintf("%v", p.Value)
	} else if p.OperationType == DocumentFilter {
		return fmt.Sprintf("d%v", p.Value)
	} else if p.OperationType == SelfReference {
		return "SELF"
	} else if p.OperationType == ValueOp {
		return fmt.Sprintf("%v (%T)", p.Value, p.Value)
	} else {
		return fmt.Sprintf("%v", p.OperationType.Type)
	}
}

//use for debugging only
func NodesToString(collection *list.List) string {
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
	tag := value.Tag
	if value.Kind == yaml.DocumentNode {
		tag = "doc"
	} else if value.Kind == yaml.AliasNode {
		tag = "alias"
	}
	return fmt.Sprintf(`D%v, P%v, (%v)::%v`, node.Document, node.Path, tag, buf.String())
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
