package yqlib

import (
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
