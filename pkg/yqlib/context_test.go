package yqlib

import (
	"container/list"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
	logging "gopkg.in/op/go-logging.v1"
)

func TestChildContext(t *testing.T) {

	expectedOriginal := make(map[string]*list.List)
	expectedOriginal["dog"] = list.New()
	expectedOriginal["dog"].PushBack(&CandidateNode{Value: "woof"})

	originalVariables := make(map[string]*list.List)
	originalVariables["dog"] = list.New()
	originalVariables["dog"].PushBack(&CandidateNode{Value: "woof"})

	original := Context{
		DontAutoCreate: true,
		datetimeLayout: "cat",
		Variables:      originalVariables,
	}

	newResults := list.New()
	newResults.PushBack(&CandidateNode{Value: "bar"})

	clone := original.ChildContext(newResults)
	test.AssertResultComplex(t, originalVariables, clone.Variables)

	clone.Variables["dog"].PushBack("bark")
	// ensure this is a separate copy
	test.AssertResultComplex(t, 1, originalVariables["dog"].Len())

}

func TestChildContextNoVariables(t *testing.T) {

	original := Context{
		DontAutoCreate: true,
		datetimeLayout: "cat",
	}

	newResults := list.New()
	newResults.PushBack(&CandidateNode{Value: "bar"})

	clone := original.ChildContext(newResults)
	test.AssertResultComplex(t, make(map[string]*list.List), clone.Variables)

}

func TestSingleReadonlyChildContext(t *testing.T) {
	original := Context{
		DontAutoCreate: false,
		datetimeLayout: "2006-01-02",
	}

	candidate := &CandidateNode{Value: "test"}
	clone := original.SingleReadonlyChildContext(candidate)

	// Should have DontAutoCreate set to true
	test.AssertResultComplex(t, true, clone.DontAutoCreate)

	// Should have the candidate node in MatchingNodes
	test.AssertResultComplex(t, 1, clone.MatchingNodes.Len())
	test.AssertResultComplex(t, candidate, clone.MatchingNodes.Front().Value)
}

func TestSingleChildContext(t *testing.T) {
	original := Context{
		DontAutoCreate: true,
		datetimeLayout: "2006-01-02",
	}

	candidate := &CandidateNode{Value: "test"}
	clone := original.SingleChildContext(candidate)

	// Should preserve DontAutoCreate
	test.AssertResultComplex(t, true, clone.DontAutoCreate)

	// Should have the candidate node in MatchingNodes
	test.AssertResultComplex(t, 1, clone.MatchingNodes.Len())
	test.AssertResultComplex(t, candidate, clone.MatchingNodes.Front().Value)
}

func TestSetDateTimeLayout(t *testing.T) {
	context := Context{}

	// Test setting datetime layout
	context.SetDateTimeLayout("2006-01-02T15:04:05Z07:00")
	test.AssertResultComplex(t, "2006-01-02T15:04:05Z07:00", context.datetimeLayout)
}

func TestGetDateTimeLayout(t *testing.T) {
	// Test with custom layout
	context := Context{datetimeLayout: "2006-01-02"}
	result := context.GetDateTimeLayout()
	test.AssertResultComplex(t, "2006-01-02", result)

	// Test with empty layout (should return default)
	context = Context{}
	result = context.GetDateTimeLayout()
	test.AssertResultComplex(t, "2006-01-02T15:04:05Z07:00", result)
}

func TestGetVariable(t *testing.T) {
	// Test with nil Variables
	context := Context{}
	result := context.GetVariable("nonexistent")
	test.AssertResultComplex(t, (*list.List)(nil), result)

	// Test with existing variable
	variables := make(map[string]*list.List)
	variables["test"] = list.New()
	variables["test"].PushBack(&CandidateNode{Value: "value"})

	context = Context{Variables: variables}
	result = context.GetVariable("test")
	test.AssertResultComplex(t, variables["test"], result)

	// Test with non-existent variable
	result = context.GetVariable("nonexistent")
	test.AssertResultComplex(t, (*list.List)(nil), result)
}

func TestSetVariable(t *testing.T) {
	// Test setting variable when Variables is nil
	context := Context{}
	value := list.New()
	value.PushBack(&CandidateNode{Value: "test"})

	context.SetVariable("key", value)
	test.AssertResultComplex(t, value, context.Variables["key"])

	// Test setting variable when Variables already exists
	context.SetVariable("key2", value)
	test.AssertResultComplex(t, value, context.Variables["key2"])
}

func TestToString(t *testing.T) {
	context := Context{
		DontAutoCreate: true,
		MatchingNodes:  list.New(),
	}

	// Add a node to test the full string representation
	node := &CandidateNode{Value: "test"}
	context.MatchingNodes.PushBack(node)

	// Test with debug logging disabled (default)
	result := context.ToString()
	test.AssertResultComplex(t, "", result)

	// Test with debug logging enabled
	logging.SetLevel(logging.DEBUG, "")
	defer logging.SetLevel(logging.INFO, "") // Reset to default

	result2 := context.ToString()
	test.AssertResultComplex(t, true, len(result2) > 0)
	test.AssertResultComplex(t, true, strings.Contains(result2, "Context"))
	test.AssertResultComplex(t, true, strings.Contains(result2, "DontAutoCreate: true"))
}

func TestDeepClone(t *testing.T) {
	// Create original context with variables and matching nodes
	originalVariables := make(map[string]*list.List)
	originalVariables["test"] = list.New()
	originalVariables["test"].PushBack(&CandidateNode{Value: "original"})

	original := Context{
		DontAutoCreate: true,
		datetimeLayout: "2006-01-02",
		Variables:      originalVariables,
		MatchingNodes:  list.New(),
	}

	// Add a node to MatchingNodes
	node := &CandidateNode{Value: "test"}
	original.MatchingNodes.PushBack(node)

	clone := original.DeepClone()

	// Should preserve DontAutoCreate and datetimeLayout
	test.AssertResultComplex(t, original.DontAutoCreate, clone.DontAutoCreate)
	test.AssertResultComplex(t, original.datetimeLayout, clone.datetimeLayout)

	// Should have copied variables
	test.AssertResultComplex(t, 1, len(clone.Variables))
	test.AssertResultComplex(t, "original", clone.Variables["test"].Front().Value.(*CandidateNode).Value)

	// Should have deep copied MatchingNodes
	test.AssertResultComplex(t, 1, clone.MatchingNodes.Len())

	// Verify it's a deep copy by modifying the original
	original.MatchingNodes.Front().Value.(*CandidateNode).Value = "modified"
	test.AssertResultComplex(t, "test", clone.MatchingNodes.Front().Value.(*CandidateNode).Value)
}

func TestClone(t *testing.T) {
	// Create original context
	original := Context{
		DontAutoCreate: true,
		datetimeLayout: "2006-01-02",
		MatchingNodes:  list.New(),
	}

	node := &CandidateNode{Value: "test"}
	original.MatchingNodes.PushBack(node)

	clone := original.Clone()

	// Should preserve DontAutoCreate and datetimeLayout
	test.AssertResultComplex(t, original.DontAutoCreate, clone.DontAutoCreate)
	test.AssertResultComplex(t, original.datetimeLayout, clone.datetimeLayout)

	// Should have the same MatchingNodes reference
	test.AssertResultComplex(t, original.MatchingNodes, clone.MatchingNodes)
}

func TestReadOnlyClone(t *testing.T) {
	original := Context{
		DontAutoCreate: false,
		datetimeLayout: "2006-01-02",
		MatchingNodes:  list.New(),
	}

	node := &CandidateNode{Value: "test"}
	original.MatchingNodes.PushBack(node)

	clone := original.ReadOnlyClone()

	// Should set DontAutoCreate to true
	test.AssertResultComplex(t, true, clone.DontAutoCreate)

	// Should preserve other fields
	test.AssertResultComplex(t, original.datetimeLayout, clone.datetimeLayout)
	test.AssertResultComplex(t, original.MatchingNodes, clone.MatchingNodes)
}

func TestWritableClone(t *testing.T) {
	original := Context{
		DontAutoCreate: true,
		datetimeLayout: "2006-01-02",
		MatchingNodes:  list.New(),
	}

	node := &CandidateNode{Value: "test"}
	original.MatchingNodes.PushBack(node)

	clone := original.WritableClone()

	// Should set DontAutoCreate to false
	test.AssertResultComplex(t, false, clone.DontAutoCreate)

	// Should preserve other fields
	test.AssertResultComplex(t, original.datetimeLayout, clone.datetimeLayout)
	test.AssertResultComplex(t, original.MatchingNodes, clone.MatchingNodes)
}
