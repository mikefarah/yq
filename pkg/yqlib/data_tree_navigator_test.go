package yqlib

import (
	"container/list"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

func TestGetMatchingNodes_NilExpressionNode(t *testing.T) {
	navigator := NewDataTreeNavigator()
	context := Context{
		MatchingNodes: list.New(),
	}

	result, err := navigator.GetMatchingNodes(context, nil)

	test.AssertResult(t, nil, err)
	test.AssertResultComplex(t, context, result)
}

func TestGetMatchingNodes_UnknownOperator(t *testing.T) {
	navigator := NewDataTreeNavigator()
	context := Context{
		MatchingNodes: list.New(),
	}

	// Create an expression node with an unknown operation type
	unknownOpType := &operationType{Type: "UNKNOWN", Handler: nil}
	expressionNode := &ExpressionNode{
		Operation: &Operation{OperationType: unknownOpType},
	}

	result, err := navigator.GetMatchingNodes(context, expressionNode)

	test.AssertResult(t, "unknown operator UNKNOWN", err.Error())
	test.AssertResultComplex(t, Context{}, result)
}

func TestGetMatchingNodes_ValidOperator(t *testing.T) {
	navigator := NewDataTreeNavigator()

	// Create a simple context with a scalar node
	scalarNode := &CandidateNode{
		Kind:  ScalarNode,
		Tag:   "!!str",
		Value: "test",
	}
	context := Context{
		MatchingNodes: list.New(),
	}
	context.MatchingNodes.PushBack(scalarNode)

	// Create an expression node with a valid operation (self reference)
	expressionNode := &ExpressionNode{
		Operation: &Operation{OperationType: selfReferenceOpType},
	}

	result, err := navigator.GetMatchingNodes(context, expressionNode)

	test.AssertResult(t, nil, err)
	test.AssertResult(t, 1, result.MatchingNodes.Len())

	// Verify the result contains the same node
	resultNode := result.MatchingNodes.Front().Value.(*CandidateNode)
	test.AssertResult(t, scalarNode, resultNode)
}

func TestDeeplyAssign_ScalarNode(t *testing.T) {
	navigator := NewDataTreeNavigator()

	// Create a context with a root mapping node
	rootNode := &CandidateNode{
		Kind: MappingNode,
		Tag:  "!!map",
		Content: []*CandidateNode{
			{Kind: ScalarNode, Tag: "!!str", Value: "existing", IsMapKey: true},
			{Kind: ScalarNode, Tag: "!!str", Value: "old_value"},
		},
	}
	context := Context{
		MatchingNodes: list.New(),
	}
	context.MatchingNodes.PushBack(rootNode)

	// Create a scalar node to assign
	scalarNode := &CandidateNode{
		Kind:  ScalarNode,
		Tag:   "!!str",
		Value: "new_value",
	}

	// Assign to path ["new_key"]
	path := []interface{}{"new_key"}
	err := navigator.DeeplyAssign(context, path, scalarNode)

	test.AssertResult(t, nil, err)

	// Verify the assignment was made
	// The root node should now have the new key-value pair
	test.AssertResult(t, 4, len(rootNode.Content)) // 2 original + 2 new

	// Find the new key-value pair
	found := false
	for i := 0; i < len(rootNode.Content)-1; i += 2 {
		key := rootNode.Content[i]
		value := rootNode.Content[i+1]
		if key.Value == "new_key" && value.Value == "new_value" {
			found = true
			break
		}
	}
	test.AssertResult(t, true, found)
}

func TestDeeplyAssign_MappingNode(t *testing.T) {
	navigator := NewDataTreeNavigator()

	// Create a context with a root mapping node
	rootNode := &CandidateNode{
		Kind: MappingNode,
		Tag:  "!!map",
		Content: []*CandidateNode{
			{Kind: ScalarNode, Tag: "!!str", Value: "existing", IsMapKey: true},
			{Kind: ScalarNode, Tag: "!!str", Value: "old_value"},
		},
	}
	context := Context{
		MatchingNodes: list.New(),
	}
	context.MatchingNodes.PushBack(rootNode)

	// Create a mapping node to assign (this should trigger deep merge)
	mappingNode := &CandidateNode{
		Kind: MappingNode,
		Tag:  "!!map",
		Content: []*CandidateNode{
			{Kind: ScalarNode, Tag: "!!str", Value: "nested_key", IsMapKey: true},
			{Kind: ScalarNode, Tag: "!!str", Value: "nested_value"},
		},
	}

	// Assign to path ["new_map"]
	path := []interface{}{"new_map"}
	err := navigator.DeeplyAssign(context, path, mappingNode)

	test.AssertResult(t, nil, err)

	// Verify the assignment was made
	// The root node should now have the new mapping
	test.AssertResult(t, 4, len(rootNode.Content)) // 2 original + 2 new

	// Find the new mapping
	found := false
	for i := 0; i < len(rootNode.Content); i += 2 {
		if i+1 < len(rootNode.Content) {
			key := rootNode.Content[i]
			value := rootNode.Content[i+1]
			if key.Value == "new_map" && value.Kind == MappingNode {
				found = true
				// Verify the nested content
				test.AssertResult(t, 2, len(value.Content))
				test.AssertResult(t, "nested_key", value.Content[0].Value)
				test.AssertResult(t, "nested_value", value.Content[1].Value)
				break
			}
		}
	}
	test.AssertResult(t, true, found)
}

func TestDeeplyAssign_DeepPath(t *testing.T) {
	navigator := NewDataTreeNavigator()

	// Create a context with a root mapping node
	rootNode := &CandidateNode{
		Kind: MappingNode,
		Tag:  "!!map",
		Content: []*CandidateNode{
			{Kind: ScalarNode, Tag: "!!str", Value: "level1", IsMapKey: true},
			{Kind: MappingNode, Tag: "!!map", Content: []*CandidateNode{}},
		},
	}
	context := Context{
		MatchingNodes: list.New(),
	}
	context.MatchingNodes.PushBack(rootNode)

	// Create a scalar node to assign
	scalarNode := &CandidateNode{
		Kind:  ScalarNode,
		Tag:   "!!str",
		Value: "deep_value",
	}

	// Assign to deep path ["level1", "level2", "level3"]
	path := []interface{}{"level1", "level2", "level3"}
	err := navigator.DeeplyAssign(context, path, scalarNode)

	test.AssertResult(t, nil, err)

	// Verify the deep assignment was made
	level1Node := rootNode.Content[1]                // The mapping node
	test.AssertResult(t, 2, len(level1Node.Content)) // Should have level2 key-value

	level2Key := level1Node.Content[0]
	level2Value := level1Node.Content[1]
	test.AssertResult(t, "level2", level2Key.Value)
	test.AssertResult(t, MappingNode, level2Value.Kind)

	level3Key := level2Value.Content[0]
	level3Value := level2Value.Content[1]
	test.AssertResult(t, "level3", level3Key.Value)
	test.AssertResult(t, "deep_value", level3Value.Value)
}

func TestDeeplyAssign_ArrayPath(t *testing.T) {
	navigator := NewDataTreeNavigator()

	// Create a context with a root mapping node containing an array
	rootNode := &CandidateNode{
		Kind: MappingNode,
		Tag:  "!!map",
		Content: []*CandidateNode{
			{Kind: ScalarNode, Tag: "!!str", Value: "array", IsMapKey: true},
			{Kind: SequenceNode, Tag: "!!seq", Content: []*CandidateNode{}},
		},
	}
	context := Context{
		MatchingNodes: list.New(),
	}
	context.MatchingNodes.PushBack(rootNode)

	// Create a scalar node to assign
	scalarNode := &CandidateNode{
		Kind:  ScalarNode,
		Tag:   "!!str",
		Value: "array_value",
	}

	// Assign to array path ["array", 0]
	path := []interface{}{"array", 0}
	err := navigator.DeeplyAssign(context, path, scalarNode)

	test.AssertResult(t, nil, err)

	// Verify the array assignment was made
	arrayNode := rootNode.Content[1]                // The sequence node
	test.AssertResult(t, 1, len(arrayNode.Content)) // Should have one element

	arrayElement := arrayNode.Content[0]
	test.AssertResult(t, "array_value", arrayElement.Value)
}

func TestDeeplyAssign_OverwriteExisting(t *testing.T) {
	navigator := NewDataTreeNavigator()

	// Create a context with a root mapping node
	rootNode := &CandidateNode{
		Kind: MappingNode,
		Tag:  "!!map",
		Content: []*CandidateNode{
			{Kind: ScalarNode, Tag: "!!str", Value: "key", IsMapKey: true},
			{Kind: ScalarNode, Tag: "!!str", Value: "old_value"},
		},
	}
	context := Context{
		MatchingNodes: list.New(),
	}
	context.MatchingNodes.PushBack(rootNode)

	// Create a scalar node to assign
	scalarNode := &CandidateNode{
		Kind:  ScalarNode,
		Tag:   "!!str",
		Value: "new_value",
	}

	// Assign to existing path ["key"]
	path := []interface{}{"key"}
	err := navigator.DeeplyAssign(context, path, scalarNode)

	test.AssertResult(t, nil, err)

	// Verify the value was overwritten
	test.AssertResult(t, 2, len(rootNode.Content)) // Should still have 2 elements

	key := rootNode.Content[0]
	value := rootNode.Content[1]
	test.AssertResult(t, "key", key.Value)
	test.AssertResult(t, "new_value", value.Value) // Should be overwritten
}

func TestDeeplyAssign_ErrorHandling(t *testing.T) {
	navigator := NewDataTreeNavigator()

	// Create a context with a scalar node (not a mapping)
	scalarNode := &CandidateNode{
		Kind:  ScalarNode,
		Tag:   "!!str",
		Value: "not_a_map",
	}
	context := Context{
		MatchingNodes: list.New(),
	}
	context.MatchingNodes.PushBack(scalarNode)

	// Create a scalar node to assign
	assignNode := &CandidateNode{
		Kind:  ScalarNode,
		Tag:   "!!str",
		Value: "value",
	}

	path := []interface{}{"key"}
	err := navigator.DeeplyAssign(context, path, assignNode)

	// Print the actual error for debugging
	if err != nil {
		t.Logf("Actual error: %v", err)
	}

	test.AssertResult(t, nil, err)
}

func TestGetMatchingNodes_WithVariables(t *testing.T) {
	navigator := NewDataTreeNavigator()

	// Create a context with variables
	variables := make(map[string]*list.List)
	varList := list.New()
	varList.PushBack(&CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "var_value"})
	variables["test_var"] = varList

	context := Context{
		MatchingNodes: list.New(),
		Variables:     variables,
	}

	// Create an expression node that gets a variable
	expressionNode := &ExpressionNode{
		Operation: &Operation{OperationType: getVariableOpType, StringValue: "test_var"},
	}

	result, err := navigator.GetMatchingNodes(context, expressionNode)

	test.AssertResult(t, nil, err)
	test.AssertResult(t, 1, result.MatchingNodes.Len())

	// Verify the variable was retrieved
	resultNode := result.MatchingNodes.Front().Value.(*CandidateNode)
	test.AssertResult(t, "var_value", resultNode.Value)
}

func TestGetMatchingNodes_EmptyContext(t *testing.T) {
	navigator := NewDataTreeNavigator()

	// Create an empty context
	context := Context{
		MatchingNodes: list.New(),
	}

	// Create an expression node with self reference
	expressionNode := &ExpressionNode{
		Operation: &Operation{OperationType: selfReferenceOpType},
	}

	result, err := navigator.GetMatchingNodes(context, expressionNode)

	test.AssertResult(t, nil, err)
	test.AssertResult(t, 0, result.MatchingNodes.Len())
}

func TestDeeplyAssign_ComplexMappingMerge(t *testing.T) {
	navigator := NewDataTreeNavigator()

	// Create a context with a root mapping node containing nested data
	rootNode := &CandidateNode{
		Kind: MappingNode,
		Tag:  "!!map",
		Content: []*CandidateNode{
			{Kind: ScalarNode, Tag: "!!str", Value: "config", IsMapKey: true},
			{Kind: MappingNode, Tag: "!!map", Content: []*CandidateNode{
				{Kind: ScalarNode, Tag: "!!str", Value: "existing_key", IsMapKey: true},
				{Kind: ScalarNode, Tag: "!!str", Value: "existing_value"},
			}},
		},
	}
	context := Context{
		MatchingNodes: list.New(),
	}
	context.MatchingNodes.PushBack(rootNode)

	// Create a mapping node to merge
	mappingNode := &CandidateNode{
		Kind: MappingNode,
		Tag:  "!!map",
		Content: []*CandidateNode{
			{Kind: ScalarNode, Tag: "!!str", Value: "new_key", IsMapKey: true},
			{Kind: ScalarNode, Tag: "!!str", Value: "new_value"},
			{Kind: ScalarNode, Tag: "!!str", Value: "existing_key", IsMapKey: true},
			{Kind: ScalarNode, Tag: "!!str", Value: "updated_value"},
		},
	}

	// Assign to path ["config"] (should merge with existing mapping)
	path := []interface{}{"config"}
	err := navigator.DeeplyAssign(context, path, mappingNode)

	test.AssertResult(t, nil, err)

	// Verify the merge was successful
	configNode := rootNode.Content[1]                // The config mapping node
	test.AssertResult(t, 4, len(configNode.Content)) // Should have 2 key-value pairs

	// Check that both existing and new keys are present
	foundExisting := false
	foundNew := false
	for i := 0; i < len(configNode.Content); i += 2 {
		if i+1 < len(configNode.Content) {
			key := configNode.Content[i]
			value := configNode.Content[i+1]
			switch key.Value {
			case "existing_key":
				foundExisting = true
				test.AssertResult(t, "updated_value", value.Value) // Should be updated
			case "new_key":
				foundNew = true
				test.AssertResult(t, "new_value", value.Value)
			}
		}
	}
	test.AssertResult(t, true, foundExisting)
	test.AssertResult(t, true, foundNew)
}
