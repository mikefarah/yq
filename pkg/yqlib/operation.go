package yqlib

import "fmt"

type Operation struct {
	OperationType *operationType
	Value         interface{}
	StringValue   string
	CandidateNode *CandidateNode // used for Value Path elements
	Preferences   interface{}
	UpdateAssign  bool // used for assign ops, when true it means we evaluate the rhs given the lhs
}

type operationType struct {
	Type                 string
	NumArgs              uint // number of arguments to the op
	Precedence           uint
	Handler              operatorHandler
	CheckForPostTraverse bool
	ToString             func(o *Operation) string
}

var valueToStringFunc = func(p *Operation) string {
	return fmt.Sprintf("%v (%T)", p.Value, p.Value)
}

func createValueOperation(value interface{}, stringValue string) *Operation {
	log.Debug("creating value op for string %v", stringValue)
	var node = createScalarNode(value, stringValue)

	return &Operation{
		OperationType: valueOpType,
		Value:         value,
		StringValue:   stringValue,
		CandidateNode: node,
	}
}

var orOpType = &operationType{Type: "OR", NumArgs: 2, Precedence: 20, Handler: orOperator}
var andOpType = &operationType{Type: "AND", NumArgs: 2, Precedence: 20, Handler: andOperator}
var reduceOpType = &operationType{Type: "REDUCE", NumArgs: 2, Precedence: 35, Handler: reduceOperator}

var blockOpType = &operationType{Type: "BLOCK", Precedence: 10, NumArgs: 2, Handler: emptyOperator}

var unionOpType = &operationType{Type: "UNION", NumArgs: 2, Precedence: 10, Handler: unionOperator}

var pipeOpType = &operationType{Type: "PIPE", NumArgs: 2, Precedence: 30, Handler: pipeOperator}

var assignOpType = &operationType{Type: "ASSIGN", NumArgs: 2, Precedence: 40, Handler: assignUpdateOperator}
var addAssignOpType = &operationType{Type: "ADD_ASSIGN", NumArgs: 2, Precedence: 40, Handler: addAssignOperator}
var subtractAssignOpType = &operationType{Type: "SUBTRACT_ASSIGN", NumArgs: 2, Precedence: 40, Handler: subtractAssignOperator}

var assignAttributesOpType = &operationType{Type: "ASSIGN_ATTRIBUTES", NumArgs: 2, Precedence: 40, Handler: assignAttributesOperator}
var assignStyleOpType = &operationType{Type: "ASSIGN_STYLE", NumArgs: 2, Precedence: 40, Handler: assignStyleOperator}
var assignVariableOpType = &operationType{Type: "ASSIGN_VARIABLE", NumArgs: 2, Precedence: 40, Handler: useWithPipe}
var assignTagOpType = &operationType{Type: "ASSIGN_TAG", NumArgs: 2, Precedence: 40, Handler: assignTagOperator}
var assignCommentOpType = &operationType{Type: "ASSIGN_COMMENT", NumArgs: 2, Precedence: 40, Handler: assignCommentsOperator}
var assignAnchorOpType = &operationType{Type: "ASSIGN_ANCHOR", NumArgs: 2, Precedence: 40, Handler: assignAnchorOperator}
var assignAliasOpType = &operationType{Type: "ASSIGN_ALIAS", NumArgs: 2, Precedence: 40, Handler: assignAliasOperator}

var multiplyOpType = &operationType{Type: "MULTIPLY", NumArgs: 2, Precedence: 42, Handler: multiplyOperator}
var multiplyAssignOpType = &operationType{Type: "MULTIPLY_ASSIGN", NumArgs: 2, Precedence: 42, Handler: multiplyAssignOperator}

var divideOpType = &operationType{Type: "DIVIDE", NumArgs: 2, Precedence: 42, Handler: divideOperator}

var moduloOpType = &operationType{Type: "MODULO", NumArgs: 2, Precedence: 42, Handler: moduloOperator}

var addOpType = &operationType{Type: "ADD", NumArgs: 2, Precedence: 42, Handler: addOperator}
var subtractOpType = &operationType{Type: "SUBTRACT", NumArgs: 2, Precedence: 42, Handler: subtractOperator}
var alternativeOpType = &operationType{Type: "ALTERNATIVE", NumArgs: 2, Precedence: 42, Handler: alternativeOperator}

var equalsOpType = &operationType{Type: "EQUALS", NumArgs: 2, Precedence: 40, Handler: equalsOperator}
var notEqualsOpType = &operationType{Type: "NOT_EQUALS", NumArgs: 2, Precedence: 40, Handler: notEqualsOperator}

var compareOpType = &operationType{Type: "COMPARE", NumArgs: 2, Precedence: 40, Handler: compareOperator}
var minOpType = &operationType{Type: "MIN", NumArgs: 0, Precedence: 40, Handler: minOperator}
var maxOpType = &operationType{Type: "MAX", NumArgs: 0, Precedence: 40, Handler: maxOperator}

// createmap needs to be above union, as we use union to build the components of the objects
var createMapOpType = &operationType{Type: "CREATE_MAP", NumArgs: 2, Precedence: 15, Handler: createMapOperator}

var shortPipeOpType = &operationType{Type: "SHORT_PIPE", NumArgs: 2, Precedence: 45, Handler: pipeOperator}

var lengthOpType = &operationType{Type: "LENGTH", NumArgs: 0, Precedence: 50, Handler: lengthOperator}
var lineOpType = &operationType{Type: "LINE", NumArgs: 0, Precedence: 50, Handler: lineOperator}
var columnOpType = &operationType{Type: "LINE", NumArgs: 0, Precedence: 50, Handler: columnOperator}

// Use this expression to create alias/syntactic sugar expressions (in lexer_participle).
var expressionOpType = &operationType{Type: "EXP", NumArgs: 0, Precedence: 50, Handler: expressionOperator}

var collectOpType = &operationType{Type: "COLLECT", NumArgs: 1, Precedence: 50, Handler: collectOperator}
var mapOpType = &operationType{Type: "MAP", NumArgs: 1, Precedence: 52, Handler: mapOperator, CheckForPostTraverse: true}
var filterOpType = &operationType{Type: "FILTER", NumArgs: 1, Precedence: 52, Handler: filterOperator, CheckForPostTraverse: true}
var errorOpType = &operationType{Type: "ERROR", NumArgs: 1, Precedence: 50, Handler: errorOperator}
var pickOpType = &operationType{Type: "PICK", NumArgs: 1, Precedence: 52, Handler: pickOperator, CheckForPostTraverse: true}
var omitOpType = &operationType{Type: "OMIT", NumArgs: 1, Precedence: 52, Handler: omitOperator, CheckForPostTraverse: true}
var evalOpType = &operationType{Type: "EVAL", NumArgs: 1, Precedence: 52, Handler: evalOperator, CheckForPostTraverse: true}
var mapValuesOpType = &operationType{Type: "MAP_VALUES", NumArgs: 1, Precedence: 52, Handler: mapValuesOperator, CheckForPostTraverse: true}

var formatDateTimeOpType = &operationType{Type: "FORMAT_DATE_TIME", NumArgs: 1, Precedence: 50, Handler: formatDateTime}
var withDtFormatOpType = &operationType{Type: "WITH_DATE_TIME_FORMAT", NumArgs: 1, Precedence: 50, Handler: withDateTimeFormat}
var nowOpType = &operationType{Type: "NOW", NumArgs: 0, Precedence: 50, Handler: nowOp}
var tzOpType = &operationType{Type: "TIMEZONE", NumArgs: 1, Precedence: 50, Handler: tzOp}
var fromUnixOpType = &operationType{Type: "FROM_UNIX", NumArgs: 0, Precedence: 50, Handler: fromUnixOp}
var toUnixOpType = &operationType{Type: "TO_UNIX", NumArgs: 0, Precedence: 50, Handler: toUnixOp}

var encodeOpType = &operationType{Type: "ENCODE", NumArgs: 0, Precedence: 50, Handler: encodeOperator}
var decodeOpType = &operationType{Type: "DECODE", NumArgs: 0, Precedence: 50, Handler: decodeOperator}

var anyOpType = &operationType{Type: "ANY", NumArgs: 0, Precedence: 50, Handler: anyOperator}
var allOpType = &operationType{Type: "ALL", NumArgs: 0, Precedence: 50, Handler: allOperator}
var containsOpType = &operationType{Type: "CONTAINS", NumArgs: 1, Precedence: 50, Handler: containsOperator}
var anyConditionOpType = &operationType{Type: "ANY_CONDITION", NumArgs: 1, Precedence: 50, Handler: anyOperator}
var allConditionOpType = &operationType{Type: "ALL_CONDITION", NumArgs: 1, Precedence: 50, Handler: allOperator}

var toEntriesOpType = &operationType{Type: "TO_ENTRIES", NumArgs: 0, Precedence: 52, Handler: toEntriesOperator, CheckForPostTraverse: true}
var fromEntriesOpType = &operationType{Type: "FROM_ENTRIES", NumArgs: 0, Precedence: 50, Handler: fromEntriesOperator}
var withEntriesOpType = &operationType{Type: "WITH_ENTRIES", NumArgs: 1, Precedence: 50, Handler: withEntriesOperator}

var withOpType = &operationType{Type: "WITH", NumArgs: 1, Precedence: 52, Handler: withOperator, CheckForPostTraverse: true}

var splitDocumentOpType = &operationType{Type: "SPLIT_DOC", NumArgs: 0, Precedence: 52, Handler: splitDocumentOperator, CheckForPostTraverse: true}
var getVariableOpType = &operationType{Type: "GET_VARIABLE", NumArgs: 0, Precedence: 55, Handler: getVariableOperator}
var getStyleOpType = &operationType{Type: "GET_STYLE", NumArgs: 0, Precedence: 50, Handler: getStyleOperator}
var getTagOpType = &operationType{Type: "GET_TAG", NumArgs: 0, Precedence: 50, Handler: getTagOperator}
var getKindOpType = &operationType{Type: "GET_KIND", NumArgs: 0, Precedence: 50, Handler: getKindOperator}

var getKeyOpType = &operationType{Type: "GET_KEY", NumArgs: 0, Precedence: 50, Handler: getKeyOperator}
var isKeyOpType = &operationType{Type: "IS_KEY", NumArgs: 0, Precedence: 50, Handler: isKeyOperator}
var getParentOpType = &operationType{Type: "GET_PARENT", NumArgs: 0, Precedence: 50, Handler: getParentOperator}

var getCommentOpType = &operationType{Type: "GET_COMMENT", NumArgs: 0, Precedence: 50, Handler: getCommentsOperator}
var getAnchorOpType = &operationType{Type: "GET_ANCHOR", NumArgs: 0, Precedence: 50, Handler: getAnchorOperator}
var getAliasOpType = &operationType{Type: "GET_ALIAS", NumArgs: 0, Precedence: 50, Handler: getAliasOperator}
var getDocumentIndexOpType = &operationType{Type: "GET_DOCUMENT_INDEX", NumArgs: 0, Precedence: 50, Handler: getDocumentIndexOperator}
var getFilenameOpType = &operationType{Type: "GET_FILENAME", NumArgs: 0, Precedence: 50, Handler: getFilenameOperator}
var getFileIndexOpType = &operationType{Type: "GET_FILE_INDEX", NumArgs: 0, Precedence: 50, Handler: getFileIndexOperator}

var getPathOpType = &operationType{Type: "GET_PATH", NumArgs: 0, Precedence: 52, Handler: getPathOperator, CheckForPostTraverse: true}
var setPathOpType = &operationType{Type: "SET_PATH", NumArgs: 1, Precedence: 50, Handler: setPathOperator}
var delPathsOpType = &operationType{Type: "DEL_PATHS", NumArgs: 1, Precedence: 52, Handler: delPathsOperator, CheckForPostTraverse: true}

var explodeOpType = &operationType{Type: "EXPLODE", NumArgs: 1, Precedence: 52, Handler: explodeOperator, CheckForPostTraverse: true}
var sortByOpType = &operationType{Type: "SORT_BY", NumArgs: 1, Precedence: 52, Handler: sortByOperator, CheckForPostTraverse: true}
var reverseOpType = &operationType{Type: "REVERSE", NumArgs: 0, Precedence: 52, Handler: reverseOperator, CheckForPostTraverse: true}
var sortOpType = &operationType{Type: "SORT", NumArgs: 0, Precedence: 52, Handler: sortOperator, CheckForPostTraverse: true}
var shuffleOpType = &operationType{Type: "SHUFFLE", NumArgs: 0, Precedence: 52, Handler: shuffleOperator, CheckForPostTraverse: true}

var sortKeysOpType = &operationType{Type: "SORT_KEYS", NumArgs: 1, Precedence: 52, Handler: sortKeysOperator, CheckForPostTraverse: true}

var joinStringOpType = &operationType{Type: "JOIN", NumArgs: 1, Precedence: 50, Handler: joinStringOperator}
var subStringOpType = &operationType{Type: "SUBSTR", NumArgs: 1, Precedence: 50, Handler: substituteStringOperator}
var matchOpType = &operationType{Type: "MATCH", NumArgs: 1, Precedence: 50, Handler: matchOperator}
var captureOpType = &operationType{Type: "CAPTURE", NumArgs: 1, Precedence: 50, Handler: captureOperator}
var testOpType = &operationType{Type: "TEST", NumArgs: 1, Precedence: 50, Handler: testOperator}
var splitStringOpType = &operationType{Type: "SPLIT", NumArgs: 1, Precedence: 52, Handler: splitStringOperator, CheckForPostTraverse: true}
var changeCaseOpType = &operationType{Type: "CHANGE_CASE", NumArgs: 0, Precedence: 50, Handler: changeCaseOperator}
var trimOpType = &operationType{Type: "TRIM", NumArgs: 0, Precedence: 50, Handler: trimSpaceOperator}
var toStringOpType = &operationType{Type: "TO_STRING", NumArgs: 0, Precedence: 50, Handler: toStringOperator}
var stringInterpolationOpType = &operationType{Type: "STRING_INT", NumArgs: 0, Precedence: 50, Handler: stringInterpolationOperator, ToString: valueToStringFunc}

var loadOpType = &operationType{Type: "LOAD", NumArgs: 1, Precedence: 52, Handler: loadOperator, CheckForPostTraverse: true}
var loadStringOpType = &operationType{Type: "LOAD_STRING", NumArgs: 1, Precedence: 52, Handler: loadStringOperator}

var keysOpType = &operationType{Type: "KEYS", NumArgs: 0, Precedence: 52, Handler: keysOperator, CheckForPostTraverse: true}

var collectObjectOpType = &operationType{Type: "COLLECT_OBJECT", NumArgs: 0, Precedence: 50, Handler: collectObjectOperator}

var traversePathOpType = &operationType{Type: "TRAVERSE_PATH", NumArgs: 0, Precedence: 55, Handler: traversePathOperator,
	ToString: func(p *Operation) string {
		return fmt.Sprintf("%v", p.Value)
	}}

var traverseArrayOpType = &operationType{Type: "TRAVERSE_ARRAY", NumArgs: 2, Precedence: 50, Handler: traverseArrayOperator}

var selfReferenceOpType = &operationType{Type: "SELF", NumArgs: 0, Precedence: 55, Handler: selfOperator}
var valueOpType = &operationType{Type: "VALUE", NumArgs: 0, Precedence: 50, Handler: valueOperator, ToString: valueToStringFunc}
var referenceOpType = &operationType{Type: "REF", NumArgs: 0, Precedence: 50, Handler: referenceOperator}
var envOpType = &operationType{Type: "ENV", NumArgs: 0, Precedence: 52, Handler: envOperator, CheckForPostTraverse: true}
var notOpType = &operationType{Type: "NOT", NumArgs: 0, Precedence: 50, Handler: notOperator}
var toNumberOpType = &operationType{Type: "TO_NUMBER", NumArgs: 0, Precedence: 50, Handler: toNumberOperator}
var emptyOpType = &operationType{Type: "EMPTY", Precedence: 50, Handler: emptyOperator}

var envsubstOpType = &operationType{Type: "ENVSUBST", NumArgs: 0, Precedence: 50, Handler: envsubstOperator}

var recursiveDescentOpType = &operationType{Type: "RECURSIVE_DESCENT", NumArgs: 0, Precedence: 50, Handler: recursiveDescentOperator}

var selectOpType = &operationType{Type: "SELECT", NumArgs: 1, Precedence: 52, Handler: selectOperator, CheckForPostTraverse: true}
var hasOpType = &operationType{Type: "HAS", NumArgs: 1, Precedence: 50, Handler: hasOperator}
var uniqueOpType = &operationType{Type: "UNIQUE", NumArgs: 0, Precedence: 52, Handler: unique, CheckForPostTraverse: true}
var uniqueByOpType = &operationType{Type: "UNIQUE_BY", NumArgs: 1, Precedence: 52, Handler: uniqueBy, CheckForPostTraverse: true}
var groupByOpType = &operationType{Type: "GROUP_BY", NumArgs: 1, Precedence: 52, Handler: groupBy, CheckForPostTraverse: true}
var flattenOpType = &operationType{Type: "FLATTEN_BY", NumArgs: 0, Precedence: 52, Handler: flattenOp, CheckForPostTraverse: true}
var deleteChildOpType = &operationType{Type: "DELETE", NumArgs: 1, Precedence: 40, Handler: deleteChildOperator}

var pivotOpType = &operationType{Type: "PIVOT", NumArgs: 0, Precedence: 52, Handler: pivotOperator, CheckForPostTraverse: true}

// debugging purposes only
func (p *Operation) toString() string {
	if p == nil {
		return "OP IS NIL"
	}
	if p.OperationType.ToString != nil {
		return p.OperationType.ToString(p)
	}
	return fmt.Sprintf("%v", p.OperationType.Type)
}
