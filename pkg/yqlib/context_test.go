package yqlib

import (
	"container/list"
	"testing"

	"github.com/mikefarah/yq/v4/test"
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
