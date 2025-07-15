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
	test.AssertResultWithContext(t, 2, child.GetFileIndex(), "file index")
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

func TestCandidateNodeAddKeyValueChild(t *testing.T) {
	key := CandidateNode{Value: "cool", IsMapKey: true}
	node := CandidateNode{}

	// if we use a key in a new node as a value, it should no longer be marked as a key

	_, keyIsValueNow := node.AddKeyValueChild(&CandidateNode{Value: "newKey"}, &key)

	test.AssertResult(t, keyIsValueNow.IsMapKey, false)
	test.AssertResult(t, key.IsMapKey, true)

}

func TestConvertToNodeInfo(t *testing.T) {
	child := &CandidateNode{
		Kind:   ScalarNode,
		Style:  DoubleQuotedStyle,
		Tag:    "!!str",
		Value:  "childValue",
		Line:   2,
		Column: 3,
	}
	parent := &CandidateNode{
		Kind:        MappingNode,
		Style:       TaggedStyle,
		Tag:         "!!map",
		Value:       "",
		Line:        1,
		Column:      1,
		Content:     []*CandidateNode{child},
		HeadComment: "head",
		LineComment: "line",
		FootComment: "foot",
		Anchor:      "anchor",
	}
	info := parent.ConvertToNodeInfo()
	test.AssertResult(t, "MappingNode", info.Kind)
	test.AssertResult(t, "TaggedStyle", info.Style)
	test.AssertResult(t, "!!map", info.Tag)
	test.AssertResult(t, "head", info.HeadComment)
	test.AssertResult(t, "line", info.LineComment)
	test.AssertResult(t, "foot", info.FootComment)
	test.AssertResult(t, "anchor", info.Anchor)
	test.AssertResult(t, 1, info.Line)
	test.AssertResult(t, 1, info.Column)
	if len(info.Content) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(info.Content))
	}
	childInfo := info.Content[0]
	test.AssertResult(t, "ScalarNode", childInfo.Kind)
	test.AssertResult(t, "DoubleQuotedStyle", childInfo.Style)
	test.AssertResult(t, "!!str", childInfo.Tag)
	test.AssertResult(t, "childValue", childInfo.Value)
	test.AssertResult(t, 2, childInfo.Line)
	test.AssertResult(t, 3, childInfo.Column)
}
