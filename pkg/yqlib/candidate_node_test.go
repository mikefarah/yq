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

func TestCandidateNodeGetPath(t *testing.T) {
	// Test root node with no parent
	root := CandidateNode{Value: "root"}
	path := root.GetPath()
	test.AssertResult(t, 0, len(path))

	// Test node with key
	key := createStringScalarNode("myKey")
	node := CandidateNode{Key: key, Value: "myValue"}
	path = node.GetPath()
	test.AssertResult(t, 1, len(path))
	test.AssertResult(t, "myKey", path[0])

	// Test nested path
	parent := CandidateNode{}
	parentKey := createStringScalarNode("parent")
	parent.Key = parentKey
	node.Parent = &parent
	path = node.GetPath()
	test.AssertResult(t, 2, len(path))
	test.AssertResult(t, "parent", path[0])
	test.AssertResult(t, "myKey", path[1])
}

func TestCandidateNodeGetNicePath(t *testing.T) {
	// Test simple key
	key := createStringScalarNode("simple")
	node := CandidateNode{Key: key}
	nicePath := node.GetNicePath()
	test.AssertResult(t, "simple", nicePath)

	// Test array index
	arrayKey := createScalarNode(0, "0")
	arrayNode := CandidateNode{Key: arrayKey}
	nicePath = arrayNode.GetNicePath()
	test.AssertResult(t, "[0]", nicePath)

	dotKey := createStringScalarNode("key.with.dots")
	dotNode := CandidateNode{Key: dotKey}
	nicePath = dotNode.GetNicePath()
	test.AssertResult(t, "key.with.dots", nicePath)

	// Test nested path
	parentKey := createStringScalarNode("parent")
	parent := CandidateNode{Key: parentKey}
	childKey := createStringScalarNode("child")
	child := CandidateNode{Key: childKey, Parent: &parent}
	nicePath = child.GetNicePath()
	test.AssertResult(t, "parent.child", nicePath)
}

func TestCandidateNodeFilterMapContentByKey(t *testing.T) {
	// Create a map with multiple key-value pairs
	key1 := createStringScalarNode("key1")
	value1 := createStringScalarNode("value1")
	key2 := createStringScalarNode("key2")
	value2 := createStringScalarNode("value2")
	key3 := createStringScalarNode("key3")
	value3 := createStringScalarNode("value3")

	mapNode := &CandidateNode{
		Kind:    MappingNode,
		Content: []*CandidateNode{key1, value1, key2, value2, key3, value3},
	}

	// Filter by key predicate that matches key1 and key3
	filtered := mapNode.FilterMapContentByKey(func(key *CandidateNode) bool {
		return key.Value == "key1" || key.Value == "key3"
	})

	// Should return key1, value1, key3, value3
	test.AssertResult(t, 4, len(filtered))
	test.AssertResult(t, "key1", filtered[0].Value)
	test.AssertResult(t, "value1", filtered[1].Value)
	test.AssertResult(t, "key3", filtered[2].Value)
	test.AssertResult(t, "value3", filtered[3].Value)
}

func TestCandidateNodeVisitValues(t *testing.T) {
	// Test mapping node
	key1 := createStringScalarNode("key1")
	value1 := createStringScalarNode("value1")
	key2 := createStringScalarNode("key2")
	value2 := createStringScalarNode("value2")

	mapNode := &CandidateNode{
		Kind:    MappingNode,
		Content: []*CandidateNode{key1, value1, key2, value2},
	}

	var visited []string
	err := mapNode.VisitValues(func(node *CandidateNode) error {
		visited = append(visited, node.Value)
		return nil
	})

	test.AssertResult(t, nil, err)
	test.AssertResult(t, 2, len(visited))
	test.AssertResult(t, "value1", visited[0])
	test.AssertResult(t, "value2", visited[1])

	// Test sequence node
	item1 := createStringScalarNode("item1")
	item2 := createStringScalarNode("item2")

	seqNode := &CandidateNode{
		Kind:    SequenceNode,
		Content: []*CandidateNode{item1, item2},
	}

	visited = []string{}
	err = seqNode.VisitValues(func(node *CandidateNode) error {
		visited = append(visited, node.Value)
		return nil
	})

	test.AssertResult(t, nil, err)
	test.AssertResult(t, 2, len(visited))
	test.AssertResult(t, "item1", visited[0])
	test.AssertResult(t, "item2", visited[1])

	// Test scalar node (should not visit anything)
	scalarNode := &CandidateNode{
		Kind:  ScalarNode,
		Value: "scalar",
	}

	visited = []string{}
	err = scalarNode.VisitValues(func(node *CandidateNode) error {
		visited = append(visited, node.Value)
		return nil
	})

	test.AssertResult(t, nil, err)
	test.AssertResult(t, 0, len(visited))
}

func TestCandidateNodeCanVisitValues(t *testing.T) {
	mapNode := &CandidateNode{Kind: MappingNode}
	seqNode := &CandidateNode{Kind: SequenceNode}
	scalarNode := &CandidateNode{Kind: ScalarNode}

	test.AssertResult(t, true, mapNode.CanVisitValues())
	test.AssertResult(t, true, seqNode.CanVisitValues())
	test.AssertResult(t, false, scalarNode.CanVisitValues())
}

func TestCandidateNodeAddChild(t *testing.T) {
	parent := &CandidateNode{Kind: SequenceNode}
	child := createStringScalarNode("child")

	parent.AddChild(child)

	test.AssertResult(t, 1, len(parent.Content))
	test.AssertResult(t, false, parent.Content[0].IsMapKey)
	test.AssertResult(t, "0", parent.Content[0].Key.Value)
	// Check that parent is set correctly
	if parent.Content[0].Parent != parent {
		t.Errorf("Expected parent to be set correctly")
	}
}

func TestCandidateNodeAddChildren(t *testing.T) {
	// Test sequence node
	parent := &CandidateNode{Kind: SequenceNode}
	child1 := createStringScalarNode("child1")
	child2 := createStringScalarNode("child2")

	parent.AddChildren([]*CandidateNode{child1, child2})

	test.AssertResult(t, 2, len(parent.Content))
	test.AssertResult(t, "child1", parent.Content[0].Value)
	test.AssertResult(t, "child2", parent.Content[1].Value)

	// Test mapping node
	mapParent := &CandidateNode{Kind: MappingNode}
	key1 := createStringScalarNode("key1")
	value1 := createStringScalarNode("value1")
	key2 := createStringScalarNode("key2")
	value2 := createStringScalarNode("value2")

	mapParent.AddChildren([]*CandidateNode{key1, value1, key2, value2})

	test.AssertResult(t, 4, len(mapParent.Content))
	test.AssertResult(t, true, mapParent.Content[0].IsMapKey)  // key1
	test.AssertResult(t, false, mapParent.Content[1].IsMapKey) // value1
	test.AssertResult(t, true, mapParent.Content[2].IsMapKey)  // key2
	test.AssertResult(t, false, mapParent.Content[3].IsMapKey) // value2
}
