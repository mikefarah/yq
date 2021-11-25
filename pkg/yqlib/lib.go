// Use the top level Evaluator or StreamEvaluator to evaluate expressions and return matches.
//
package yqlib

import (
	"bytes"
	"container/list"
	"fmt"
	"strconv"
	"strings"

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

var (
	orOpType     = &operationType{Type: "OR", NumArgs: 2, Precedence: 20, Handler: orOperator}
	andOpType    = &operationType{Type: "AND", NumArgs: 2, Precedence: 20, Handler: andOperator}
	reduceOpType = &operationType{Type: "REDUCE", NumArgs: 2, Precedence: 35, Handler: reduceOperator}
)

var blockOpType = &operationType{Type: "BLOCK", Precedence: 10, NumArgs: 2, Handler: emptyOperator}

var unionOpType = &operationType{Type: "UNION", NumArgs: 2, Precedence: 10, Handler: unionOperator}

var pipeOpType = &operationType{Type: "PIPE", NumArgs: 2, Precedence: 30, Handler: pipeOperator}

var (
	assignOpType         = &operationType{Type: "ASSIGN", NumArgs: 2, Precedence: 40, Handler: assignUpdateOperator}
	addAssignOpType      = &operationType{Type: "ADD_ASSIGN", NumArgs: 2, Precedence: 40, Handler: addAssignOperator}
	subtractAssignOpType = &operationType{Type: "SUBTRACT_ASSIGN", NumArgs: 2, Precedence: 40, Handler: subtractAssignOperator}
)

var (
	assignAttributesOpType = &operationType{Type: "ASSIGN_ATTRIBUTES", NumArgs: 2, Precedence: 40, Handler: assignAttributesOperator}
	assignStyleOpType      = &operationType{Type: "ASSIGN_STYLE", NumArgs: 2, Precedence: 40, Handler: assignStyleOperator}
	assignVariableOpType   = &operationType{Type: "ASSIGN_VARIABLE", NumArgs: 2, Precedence: 40, Handler: assignVariableOperator}
	assignTagOpType        = &operationType{Type: "ASSIGN_TAG", NumArgs: 2, Precedence: 40, Handler: assignTagOperator}
	assignCommentOpType    = &operationType{Type: "ASSIGN_COMMENT", NumArgs: 2, Precedence: 40, Handler: assignCommentsOperator}
	assignAnchorOpType     = &operationType{Type: "ASSIGN_ANCHOR", NumArgs: 2, Precedence: 40, Handler: assignAnchorOperator}
	assignAliasOpType      = &operationType{Type: "ASSIGN_ALIAS", NumArgs: 2, Precedence: 40, Handler: assignAliasOperator}
)

var (
	multiplyOpType    = &operationType{Type: "MULTIPLY", NumArgs: 2, Precedence: 42, Handler: multiplyOperator}
	addOpType         = &operationType{Type: "ADD", NumArgs: 2, Precedence: 42, Handler: addOperator}
	subtractOpType    = &operationType{Type: "SUBTRACT", NumArgs: 2, Precedence: 42, Handler: subtractOperator}
	alternativeOpType = &operationType{Type: "ALTERNATIVE", NumArgs: 2, Precedence: 42, Handler: alternativeOperator}
)

var (
	equalsOpType    = &operationType{Type: "EQUALS", NumArgs: 2, Precedence: 40, Handler: equalsOperator}
	notEqualsOpType = &operationType{Type: "EQUALS", NumArgs: 2, Precedence: 40, Handler: notEqualsOperator}
)

// createmap needs to be above union, as we use union to build the components of the objects
var createMapOpType = &operationType{Type: "CREATE_MAP", NumArgs: 2, Precedence: 15, Handler: createMapOperator}

var shortPipeOpType = &operationType{Type: "SHORT_PIPE", NumArgs: 2, Precedence: 45, Handler: pipeOperator}

var (
	lengthOpType  = &operationType{Type: "LENGTH", NumArgs: 0, Precedence: 50, Handler: lengthOperator}
	collectOpType = &operationType{Type: "COLLECT", NumArgs: 0, Precedence: 50, Handler: collectOperator}
	encodeOpType  = &operationType{Type: "ENCODE", NumArgs: 0, Precedence: 50, Handler: encodeOperator}
	decodeOpType  = &operationType{Type: "DECODE", NumArgs: 0, Precedence: 50, Handler: decodeOperator}
)

var (
	anyOpType          = &operationType{Type: "ANY", NumArgs: 0, Precedence: 50, Handler: anyOperator}
	allOpType          = &operationType{Type: "ALL", NumArgs: 0, Precedence: 50, Handler: allOperator}
	containsOpType     = &operationType{Type: "CONTAINS", NumArgs: 1, Precedence: 50, Handler: containsOperator}
	anyConditionOpType = &operationType{Type: "ANY_CONDITION", NumArgs: 1, Precedence: 50, Handler: anyOperator}
	allConditionOpType = &operationType{Type: "ALL_CONDITION", NumArgs: 1, Precedence: 50, Handler: allOperator}
)

var (
	toEntriesOpType   = &operationType{Type: "TO_ENTRIES", NumArgs: 0, Precedence: 50, Handler: toEntriesOperator}
	fromEntriesOpType = &operationType{Type: "FROM_ENTRIES", NumArgs: 0, Precedence: 50, Handler: fromEntriesOperator}
	withEntriesOpType = &operationType{Type: "WITH_ENTRIES", NumArgs: 1, Precedence: 50, Handler: withEntriesOperator}
)

var withOpType = &operationType{Type: "WITH", NumArgs: 1, Precedence: 50, Handler: withOperator}

var (
	splitDocumentOpType = &operationType{Type: "SPLIT_DOC", NumArgs: 0, Precedence: 50, Handler: splitDocumentOperator}
	getVariableOpType   = &operationType{Type: "GET_VARIABLE", NumArgs: 0, Precedence: 55, Handler: getVariableOperator}
	getStyleOpType      = &operationType{Type: "GET_STYLE", NumArgs: 0, Precedence: 50, Handler: getStyleOperator}
	getTagOpType        = &operationType{Type: "GET_TAG", NumArgs: 0, Precedence: 50, Handler: getTagOperator}
)

var (
	getKeyOpType    = &operationType{Type: "GET_KEY", NumArgs: 0, Precedence: 50, Handler: getKeyOperator}
	getParentOpType = &operationType{Type: "GET_PARENT", NumArgs: 0, Precedence: 50, Handler: getParentOperator}
)

var (
	getCommentOpType       = &operationType{Type: "GET_COMMENT", NumArgs: 0, Precedence: 50, Handler: getCommentsOperator}
	getAnchorOpType        = &operationType{Type: "GET_ANCHOR", NumArgs: 0, Precedence: 50, Handler: getAnchorOperator}
	getAliasOptype         = &operationType{Type: "GET_ALIAS", NumArgs: 0, Precedence: 50, Handler: getAliasOperator}
	getDocumentIndexOpType = &operationType{Type: "GET_DOCUMENT_INDEX", NumArgs: 0, Precedence: 50, Handler: getDocumentIndexOperator}
	getFilenameOpType      = &operationType{Type: "GET_FILENAME", NumArgs: 0, Precedence: 50, Handler: getFilenameOperator}
	getFileIndexOpType     = &operationType{Type: "GET_FILE_INDEX", NumArgs: 0, Precedence: 50, Handler: getFileIndexOperator}
	getPathOpType          = &operationType{Type: "GET_PATH", NumArgs: 0, Precedence: 50, Handler: getPathOperator}
)

var (
	explodeOpType     = &operationType{Type: "EXPLODE", NumArgs: 1, Precedence: 50, Handler: explodeOperator}
	sortKeysOpType    = &operationType{Type: "SORT_KEYS", NumArgs: 1, Precedence: 50, Handler: sortKeysOperator}
	joinStringOpType  = &operationType{Type: "JOIN", NumArgs: 1, Precedence: 50, Handler: joinStringOperator}
	subStringOpType   = &operationType{Type: "SUBSTR", NumArgs: 1, Precedence: 50, Handler: substituteStringOperator}
	matchOpType       = &operationType{Type: "MATCH", NumArgs: 1, Precedence: 50, Handler: matchOperator}
	captureOpType     = &operationType{Type: "CAPTURE", NumArgs: 1, Precedence: 50, Handler: captureOperator}
	testOpType        = &operationType{Type: "TEST", NumArgs: 1, Precedence: 50, Handler: testOperator}
	splitStringOpType = &operationType{Type: "SPLIT", NumArgs: 1, Precedence: 50, Handler: splitStringOperator}
)

var loadOpType = &operationType{Type: "LOAD", NumArgs: 1, Precedence: 50, Handler: loadYamlOperator}

var keysOpType = &operationType{Type: "KEYS", NumArgs: 0, Precedence: 50, Handler: keysOperator}

var (
	collectObjectOpType = &operationType{Type: "COLLECT_OBJECT", NumArgs: 0, Precedence: 50, Handler: collectObjectOperator}
	traversePathOpType  = &operationType{Type: "TRAVERSE_PATH", NumArgs: 0, Precedence: 55, Handler: traversePathOperator}
	traverseArrayOpType = &operationType{Type: "TRAVERSE_ARRAY", NumArgs: 2, Precedence: 50, Handler: traverseArrayOperator}
)

var (
	selfReferenceOpType = &operationType{Type: "SELF", NumArgs: 0, Precedence: 55, Handler: selfOperator}
	valueOpType         = &operationType{Type: "VALUE", NumArgs: 0, Precedence: 50, Handler: valueOperator}
	envOpType           = &operationType{Type: "ENV", NumArgs: 0, Precedence: 50, Handler: envOperator}
	notOpType           = &operationType{Type: "NOT", NumArgs: 0, Precedence: 50, Handler: notOperator}
	emptyOpType         = &operationType{Type: "EMPTY", Precedence: 50, Handler: emptyOperator}
)

var recursiveDescentOpType = &operationType{Type: "RECURSIVE_DESCENT", NumArgs: 0, Precedence: 50, Handler: recursiveDescentOperator}

var (
	selectOpType      = &operationType{Type: "SELECT", NumArgs: 1, Precedence: 50, Handler: selectOperator}
	hasOpType         = &operationType{Type: "HAS", NumArgs: 1, Precedence: 50, Handler: hasOperator}
	uniqueOpType      = &operationType{Type: "UNIQUE", NumArgs: 0, Precedence: 50, Handler: unique}
	uniqueByOpType    = &operationType{Type: "UNIQUE_BY", NumArgs: 1, Precedence: 50, Handler: uniqueBy}
	groupByOpType     = &operationType{Type: "GROUP_BY", NumArgs: 1, Precedence: 50, Handler: groupBy}
	flattenOpType     = &operationType{Type: "FLATTEN_BY", NumArgs: 0, Precedence: 50, Handler: flattenOp}
	deleteChildOpType = &operationType{Type: "DELETE", NumArgs: 1, Precedence: 40, Handler: deleteChildOperator}
)

type Operation struct {
	OperationType *operationType
	Value         interface{}
	StringValue   string
	CandidateNode *CandidateNode // used for Value Path elements
	Preferences   interface{}
	UpdateAssign  bool // used for assign ops, when true it means we evaluate the rhs given the lhs
}

func recurseNodeArrayEqual(lhs *yaml.Node, rhs *yaml.Node) bool {
	if len(lhs.Content) != len(rhs.Content) {
		return false
	}

	for index := 0; index < len(lhs.Content); index = index + 1 {
		if !recursiveNodeEqual(lhs.Content[index], rhs.Content[index]) {
			return false
		}
	}
	return true
}

func findInArray(array *yaml.Node, item *yaml.Node) int {
	for index := 0; index < len(array.Content); index = index + 1 {
		if recursiveNodeEqual(array.Content[index], item) {
			return index
		}
	}
	return -1
}

func recurseNodeObjectEqual(lhs *yaml.Node, rhs *yaml.Node) bool {
	if len(lhs.Content) != len(rhs.Content) {
		return false
	}

	for index := 0; index < len(lhs.Content); index = index + 2 {
		key := lhs.Content[index]
		value := lhs.Content[index+1]

		indexInRhs := findInArray(rhs, key)

		if indexInRhs == -1 || !recursiveNodeEqual(value, rhs.Content[indexInRhs+1]) {
			return false
		}
	}
	return true
}

func recursiveNodeEqual(lhs *yaml.Node, rhs *yaml.Node) bool {
	if lhs.Kind != rhs.Kind || lhs.Tag != rhs.Tag {
		return false
	} else if lhs.Tag == "!!null" {
		return true
	} else if lhs.Kind == yaml.ScalarNode {
		return lhs.Value == rhs.Value
	} else if lhs.Kind == yaml.SequenceNode {
		return recurseNodeArrayEqual(lhs, rhs)
	} else if lhs.Kind == yaml.MappingNode {
		return recurseNodeObjectEqual(lhs, rhs)
	}
	return false
}

// yaml numbers can be hex encoded...
func parseInt(numberString string) (string, int64, error) {
	if strings.HasPrefix(numberString, "0x") ||
		strings.HasPrefix(numberString, "0X") {
		num, err := strconv.ParseInt(numberString[2:], 16, 64) // nolint
		return "0x%X", num, err
	}
	num, err := strconv.ParseInt(numberString, 10, 64) // nolint
	return "%v", num, err
}

func createScalarNode(value interface{}, stringValue string) *yaml.Node {
	node := &yaml.Node{Kind: yaml.ScalarNode}
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
	return node
}

func createValueOperation(value interface{}, stringValue string) *Operation {
	var node *yaml.Node = createScalarNode(value, stringValue)

	return &Operation{
		OperationType: valueOpType,
		Value:         value,
		StringValue:   stringValue,
		CandidateNode: &CandidateNode{Node: node},
	}
}

// debugging purposes only
func (p *Operation) toString() string {
	if p == nil {
		return "OP IS NIL"
	}
	if p.OperationType == traversePathOpType {
		return fmt.Sprintf("%v", p.Value)
	} else if p.OperationType == selfReferenceOpType {
		return "SELF"
	} else if p.OperationType == valueOpType {
		return fmt.Sprintf("%v (%T)", p.Value, p.Value)
	} else {
		return fmt.Sprintf("%v", p.OperationType.Type)
	}
}

// use for debugging only
func NodesToString(collection *list.List) string {
	if !log.IsEnabledFor(logging.DEBUG) {
		return ""
	}

	result := fmt.Sprintf("%v results\n", collection.Len())
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
	errorClosingEncoder := encoder.Close()
	if errorClosingEncoder != nil {
		log.Error("Error closing encoder: ", errorClosingEncoder.Error())
	}
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
