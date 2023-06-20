package yqlib

import (
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

type valueRepScenario struct {
	input    string
	tag      string
	expected interface{}
}

var valueRepScenarios = []valueRepScenario{
	{
		input:    `"cat"`,
		expected: `"cat"`,
	},
	{
		input:    `3`,
		expected: int64(3),
	},
	{
		input:    `3.1`,
		expected: float64(3.1),
	},
	{
		input:    `true`,
		expected: true,
	},
	{
		input:    `y`,
		tag:      "!!bool",
		expected: true,
	},
	{
		tag:      "!!null",
		expected: nil,
	},
}

func TestCandidateNodeGetValueRepScenarios(t *testing.T) {
	for _, tt := range valueRepScenarios {
		node := CandidateNode{Value: tt.input, Tag: tt.tag}
		actual, err := node.GetValueRep()
		if err != nil {
			t.Error(err)
			return
		}
		test.AssertResult(t, tt.expected, actual)
	}
}

func TestCandidateNodeChildWhenParentUpdated(t *testing.T) {
	parent := CandidateNode{}
	child := parent.CreateChild()
	parent.SetDocument(1)
	parent.SetFileIndex(2)
	parent.SetFilename("meow")
	test.AssertResultWithContext(t, "meow", child.GetFilename(), "filename")
	test.AssertResultWithContext(t, 2, child.GetFileIndex(), "fileindex")
	test.AssertResultWithContext(t, uint(1), child.GetDocument(), "document index")
}

type createScalarNodeScenario struct {
	value       interface{}
	stringValue string
	expectedTag string
}

var createScalarScenarios = []createScalarNodeScenario{
	{
		value:       "mike",
		stringValue: "mike",
		expectedTag: "!!str",
	},
	{
		value:       3,
		stringValue: "3",
		expectedTag: "!!int",
	},
	{
		value:       3.1,
		stringValue: "3.1",
		expectedTag: "!!float",
	},
	{
		value:       true,
		stringValue: "true",
		expectedTag: "!!bool",
	},
	{
		value:       nil,
		stringValue: "~",
		expectedTag: "!!null",
	},
}

func TestCreateScalarNodeScenarios(t *testing.T) {
	for _, tt := range createScalarScenarios {
		actual := createScalarNode(tt.value, tt.stringValue)
		test.AssertResultWithContext(t, tt.stringValue, actual.Value, fmt.Sprintf("Value for: Value: [%v], String: %v", tt.value, tt.stringValue))
		test.AssertResultWithContext(t, tt.expectedTag, actual.Tag, fmt.Sprintf("Value for: Value: [%v], String: %v", tt.value, tt.stringValue))
	}
}

func TestGetKeyForMapValue(t *testing.T) {
	key := createStringScalarNode("yourKey")
	n := CandidateNode{Key: key, Value: "meow", document: 3}
	test.AssertResult(t, "3 - yourKey", n.GetKey())
}

func TestGetKeyForMapKey(t *testing.T) {
	key := createStringScalarNode("yourKey")
	key.IsMapKey = true
	key.document = 3
	test.AssertResult(t, "key-yourKey-3 - ", key.GetKey())
}

func TestGetKeyForValue(t *testing.T) {
	n := CandidateNode{Value: "meow", document: 3}
	test.AssertResult(t, "3 - ", n.GetKey())
}

func TestGetParsedKeyForMapKey(t *testing.T) {
	key := createStringScalarNode("yourKey")
	key.IsMapKey = true
	key.document = 3
	test.AssertResult(t, "yourKey", key.getParsedKey())
}

func TestGetParsedKeyForLooseValue(t *testing.T) {
	n := CandidateNode{Value: "meow", document: 3}
	test.AssertResult(t, nil, n.getParsedKey())
}

func TestGetParsedKeyForMapValue(t *testing.T) {
	key := createStringScalarNode("yourKey")
	n := CandidateNode{Key: key, Value: "meow", document: 3}
	test.AssertResult(t, "yourKey", n.getParsedKey())
}

func TestGetParsedKeyForArrayValue(t *testing.T) {
	key := createScalarNode(4, "4")
	n := CandidateNode{Key: key, Value: "meow", document: 3}
	test.AssertResult(t, 4, n.getParsedKey())
}
