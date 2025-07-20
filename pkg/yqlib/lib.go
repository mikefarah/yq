// Use the top level Evaluator or StreamEvaluator to evaluate expressions and return matches.
package yqlib

import (
	"container/list"
	"fmt"
	"math"
	"strconv"
	"strings"

	logging "gopkg.in/op/go-logging.v1"
)

var ExpressionParser ExpressionParserInterface

func InitExpressionParser() {
	if ExpressionParser == nil {
		ExpressionParser = newExpressionParser()
	}
}

var log = logging.MustGetLogger("yq-lib")

var PrettyPrintExp = `(... | (select(tag != "!!str"), select(tag == "!!str") | select(test("(?i)^(y|yes|n|no|on|off)$") | not))  ) style=""`

// GetLogger returns the yq logger instance.
func GetLogger() *logging.Logger {
	return log
}

func recurseNodeArrayEqual(lhs *CandidateNode, rhs *CandidateNode) bool {
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

func findInArray(array *CandidateNode, item *CandidateNode) int {

	for index := 0; index < len(array.Content); index = index + 1 {
		if recursiveNodeEqual(array.Content[index], item) {
			return index
		}
	}
	return -1
}

func findKeyInMap(dataMap *CandidateNode, item *CandidateNode) int {

	for index := 0; index < len(dataMap.Content); index = index + 2 {
		if recursiveNodeEqual(dataMap.Content[index], item) {
			return index
		}
	}
	return -1
}

func recurseNodeObjectEqual(lhs *CandidateNode, rhs *CandidateNode) bool {
	if len(lhs.Content) != len(rhs.Content) {
		return false
	}

	for index := 0; index < len(lhs.Content); index = index + 2 {
		key := lhs.Content[index]
		value := lhs.Content[index+1]

		indexInRHS := findInArray(rhs, key)

		if indexInRHS == -1 || !recursiveNodeEqual(value, rhs.Content[indexInRHS+1]) {
			return false
		}
	}
	return true
}

func parseSnippet(value string) (*CandidateNode, error) {
	if value == "" {
		return &CandidateNode{
			Kind: ScalarNode,
			Tag:  "!!null",
		}, nil
	}
	decoder := NewYamlDecoder(ConfiguredYamlPreferences)
	err := decoder.Init(strings.NewReader(value))
	if err != nil {
		return nil, err
	}
	result, err := decoder.Decode()
	if err != nil {
		return nil, err
	}

	if result.Kind == ScalarNode {
		result.LineComment = result.LeadingContent
	} else {
		result.HeadComment = result.LeadingContent
	}
	result.LeadingContent = ""

	if result.Tag == "!!str" {
		// use the original string value, as
		// decoding drops new lines
		newNode := createScalarNode(value, value)
		newNode.LineComment = result.LineComment
		return newNode, nil
	}
	result.Line = 0
	result.Column = 0
	return result, err
}

func recursiveNodeEqual(lhs *CandidateNode, rhs *CandidateNode) bool {
	if lhs.Kind != rhs.Kind {
		return false
	}

	if lhs.Kind == ScalarNode {
		//process custom tags of scalar nodes.
		//dont worry about matching tags of maps or arrays.

		lhsTag := lhs.guessTagFromCustomType()
		rhsTag := rhs.guessTagFromCustomType()

		if lhsTag != rhsTag {
			return false
		}
	}

	if lhs.Tag == "!!null" {
		return true

	} else if lhs.Kind == ScalarNode {
		return lhs.Value == rhs.Value
	} else if lhs.Kind == SequenceNode {
		return recurseNodeArrayEqual(lhs, rhs)
	} else if lhs.Kind == MappingNode {
		return recurseNodeObjectEqual(lhs, rhs)
	}
	return false
}

// yaml numbers can have underscores, be hex and octal encoded...
func parseInt64(numberString string) (string, int64, error) {
	if strings.Contains(numberString, "_") {
		numberString = strings.ReplaceAll(numberString, "_", "")
	}

	if strings.HasPrefix(numberString, "0x") ||
		strings.HasPrefix(numberString, "0X") {
		num, err := strconv.ParseInt(numberString[2:], 16, 64)
		return "0x%X", num, err
	} else if strings.HasPrefix(numberString, "0o") {
		num, err := strconv.ParseInt(numberString[2:], 8, 64)
		return "0o%o", num, err
	}
	num, err := strconv.ParseInt(numberString, 10, 64)
	return "%v", num, err
}

func parseInt(numberString string) (int, error) {
	_, parsed, err := parseInt64(numberString)

	if err != nil {
		return 0, err
	} else if parsed > math.MaxInt || parsed < math.MinInt {
		return 0, fmt.Errorf("%v is not within [%v, %v]", parsed, math.MinInt, math.MaxInt)
	}

	return int(parsed), err
}

func headAndLineComment(node *CandidateNode) string {
	return headComment(node) + lineComment(node)
}

func headComment(node *CandidateNode) string {
	return strings.Replace(node.HeadComment, "#", "", 1)
}

func lineComment(node *CandidateNode) string {
	return strings.Replace(node.LineComment, "#", "", 1)
}

func footComment(node *CandidateNode) string {
	return strings.Replace(node.FootComment, "#", "", 1)
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
	if node == nil {
		return "-- nil --"
	}
	tag := node.Tag
	if node.Kind == AliasNode {
		tag = "alias"
	}
	valueToUse := node.Value
	if valueToUse == "" {
		valueToUse = fmt.Sprintf("%v kids", len(node.Content))
	}
	return fmt.Sprintf(`D%v, P%v, %v (%v)::%v`, node.GetDocument(), node.GetNicePath(), KindString(node.Kind), tag, valueToUse)
}

func NodeContentToString(node *CandidateNode, depth int) string {
	if !log.IsEnabledFor(logging.DEBUG) {
		return ""
	}
	var sb strings.Builder
	for _, child := range node.Content {
		for i := 0; i < depth; i++ {
			sb.WriteString(" ")
		}
		sb.WriteString("- ")
		sb.WriteString(NodeToString(child))
		sb.WriteString("\n")
		sb.WriteString(NodeContentToString(child, depth+1))
	}
	return sb.String()
}

func KindString(kind Kind) string {
	switch kind {
	case ScalarNode:
		return "ScalarNode"
	case SequenceNode:
		return "SequenceNode"
	case MappingNode:
		return "MappingNode"
	case AliasNode:
		return "AliasNode"
	default:
		return "unknown!"
	}
}
