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
