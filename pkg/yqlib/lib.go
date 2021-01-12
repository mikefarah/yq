package yqlib

import (
	"bytes"
	"container/list"
	"fmt"

	logging "gopkg.in/op/go-logging.v1"
	yaml "gopkg.in/yaml.v3"
)

var log = logging.MustGetLogger("yq-lib")

type operationType struct {
	Type       string
	NumArgs    uint // number of arguments to the op
	Precedence uint
	Handler    operatorHandler
}

// operators TODO:
// - mergeEmpty (sets only if the document is empty, do I do that now?)

var orOpType = &operationType{Type: "OR", NumArgs: 2, Precedence: 20, Handler: orOperator}
var andOpType = &operationType{Type: "AND", NumArgs: 2, Precedence: 20, Handler: andOperator}

var unionOpType = &operationType{Type: "UNION", NumArgs: 2, Precedence: 10, Handler: unionOperator}

var pipeOpType = &operationType{Type: "PIPE", NumArgs: 2, Precedence: 30, Handler: pipeOperator}

var assignOpType = &operationType{Type: "ASSIGN", NumArgs: 2, Precedence: 40, Handler: assignUpdateOperator}
var addAssignOpType = &operationType{Type: "ADD_ASSIGN", NumArgs: 2, Precedence: 40, Handler: addAssignOperator}

var assignAttributesOpType = &operationType{Type: "ASSIGN_ATTRIBUTES", NumArgs: 2, Precedence: 40, Handler: assignAttributesOperator}
var assignStyleOpType = &operationType{Type: "ASSIGN_STYLE", NumArgs: 2, Precedence: 40, Handler: assignStyleOperator}
var assignTagOpType = &operationType{Type: "ASSIGN_TAG", NumArgs: 2, Precedence: 40, Handler: assignTagOperator}
var assignCommentOpType = &operationType{Type: "ASSIGN_COMMENT", NumArgs: 2, Precedence: 40, Handler: assignCommentsOperator}
var assignAnchorOpType = &operationType{Type: "ASSIGN_ANCHOR", NumArgs: 2, Precedence: 40, Handler: assignAnchorOperator}
var assignAliasOpType = &operationType{Type: "ASSIGN_ALIAS", NumArgs: 2, Precedence: 40, Handler: assignAliasOperator}

var multiplyOpType = &operationType{Type: "MULTIPLY", NumArgs: 2, Precedence: 45, Handler: multiplyOperator}
var addOpType = &operationType{Type: "ADD", NumArgs: 2, Precedence: 45, Handler: addOperator}
var alternativeOpType = &operationType{Type: "ALTERNATIVE", NumArgs: 2, Precedence: 45, Handler: alternativeOperator}

var equalsOpType = &operationType{Type: "EQUALS", NumArgs: 2, Precedence: 40, Handler: equalsOperator}
var createMapOpType = &operationType{Type: "CREATE_MAP", NumArgs: 2, Precedence: 40, Handler: createMapOperator}

var shortPipeOpType = &operationType{Type: "SHORT_PIPE", NumArgs: 2, Precedence: 45, Handler: pipeOperator}

var lengthOpType = &operationType{Type: "LENGTH", NumArgs: 0, Precedence: 50, Handler: lengthOperator}
var collectOpType = &operationType{Type: "COLLECT", NumArgs: 0, Precedence: 50, Handler: collectOperator}
var getStyleOpType = &operationType{Type: "GET_STYLE", NumArgs: 0, Precedence: 50, Handler: getStyleOperator}
var getTagOpType = &operationType{Type: "GET_TAG", NumArgs: 0, Precedence: 50, Handler: getTagOperator}
var getCommentOpType = &operationType{Type: "GET_COMMENT", NumArgs: 0, Precedence: 50, Handler: getCommentsOperator}
var getAnchorOpType = &operationType{Type: "GET_ANCHOR", NumArgs: 0, Precedence: 50, Handler: getAnchorOperator}
var getAliasOptype = &operationType{Type: "GET_ALIAS", NumArgs: 0, Precedence: 50, Handler: getAliasOperator}
var getDocumentIndexOpType = &operationType{Type: "GET_DOCUMENT_INDEX", NumArgs: 0, Precedence: 50, Handler: getDocumentIndexOperator}
var getFilenameOpType = &operationType{Type: "GET_FILENAME", NumArgs: 0, Precedence: 50, Handler: getFilenameOperator}
var getFileIndexOpType = &operationType{Type: "GET_FILE_INDEX", NumArgs: 0, Precedence: 50, Handler: getFileIndexOperator}
var getPathOpType = &operationType{Type: "GET_PATH", NumArgs: 0, Precedence: 50, Handler: getPathOperator}

var explodeOpType = &operationType{Type: "EXPLODE", NumArgs: 1, Precedence: 50, Handler: explodeOperator}
var sortKeysOpType = &operationType{Type: "SORT_KEYS", NumArgs: 1, Precedence: 50, Handler: sortKeysOperator}

var collectObjectOpType = &operationType{Type: "COLLECT_OBJECT", NumArgs: 0, Precedence: 50, Handler: collectObjectOperator}
var traversePathOpType = &operationType{Type: "TRAVERSE_PATH", NumArgs: 0, Precedence: 50, Handler: traversePathOperator}
var traverseArrayOpType = &operationType{Type: "TRAVERSE_ARRAY", NumArgs: 1, Precedence: 50, Handler: traverseArrayOperator}

var documentFilterOpType = &operationType{Type: "DOCUMENT_FILTER", NumArgs: 0, Precedence: 50, Handler: traversePathOperator}
var selfReferenceOpType = &operationType{Type: "SELF", NumArgs: 0, Precedence: 50, Handler: selfOperator}
var valueOpType = &operationType{Type: "VALUE", NumArgs: 0, Precedence: 50, Handler: valueOperator}
var envOpType = &operationType{Type: "ENV", NumArgs: 0, Precedence: 50, Handler: envOperator}
var notOpType = &operationType{Type: "NOT", NumArgs: 0, Precedence: 50, Handler: notOperator}
var emptyOpType = &operationType{Type: "EMPTY", NumArgs: 50, Handler: emptyOperator}

var recursiveDescentOpType = &operationType{Type: "RECURSIVE_DESCENT", NumArgs: 0, Precedence: 50, Handler: recursiveDescentOperator}

var selectOpType = &operationType{Type: "SELECT", NumArgs: 1, Precedence: 50, Handler: selectOperator}
var hasOpType = &operationType{Type: "HAS", NumArgs: 1, Precedence: 50, Handler: hasOperator}
var deleteChildOpType = &operationType{Type: "DELETE", NumArgs: 1, Precedence: 40, Handler: deleteChildOperator}
var deleteImmediateChildOpType = &operationType{Type: "DELETE_IMMEDIATE_CHILD", NumArgs: 1, Precedence: 40, Handler: deleteImmediateChildOperator}

type Operation struct {
	OperationType *operationType
	Value         interface{}
	StringValue   string
	CandidateNode *CandidateNode // used for Value Path elements
	Preferences   interface{}
	UpdateAssign  bool // used for assign ops, when true it means we evaluate the rhs given the lhs
}

func createValueOperation(value interface{}, stringValue string) *Operation {
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
		OperationType: valueOpType,
		Value:         value,
		StringValue:   stringValue,
		CandidateNode: &CandidateNode{Node: &node},
	}
}

// debugging purposes only
func (p *Operation) toString() string {
	if p.OperationType == traversePathOpType {
		return fmt.Sprintf("%v", p.Value)
	} else if p.OperationType == documentFilterOpType {
		return fmt.Sprintf("d%v", p.Value)
	} else if p.OperationType == selfReferenceOpType {
		return "SELF"
	} else if p.OperationType == valueOpType {
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

// EvaluateNodes takes an expression and one or more yaml nodes, returning a list of matching candidate nodes
func EvaluateNodes(expression string, nodes ...*yaml.Node) (*list.List, error) {
	inputCandidates := list.New()
	for _, node := range nodes {
		inputCandidates.PushBack(&CandidateNode{Node: node})
	}
	return EvaluateCandidateNodes(expression, inputCandidates)
}

// EvaluateCandidateNodes takes an expression and list of candidate nodes, returning a list of matching candidate nodes
func EvaluateCandidateNodes(expression string, inputCandidates *list.List) (*list.List, error) {
	node, err := treeCreator.ParsePath(expression)
	if err != nil {
		return nil, err
	}
	return treeNavigator.GetMatchingNodes(inputCandidates, node)
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
