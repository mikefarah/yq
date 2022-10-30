package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v4/test"
	yaml "gopkg.in/yaml.v3"
)

func TestGetLogger(t *testing.T) {
	l := GetLogger()
	if l != log {
		t.Fatal("GetLogger should return the yq logger instance, not a copy")
	}
}

type parseSnippetScenario struct {
	snippet  string
	expected *yaml.Node
}

var parseSnippetScenarios = []parseSnippetScenario{
	{
		snippet: "",
		expected: &yaml.Node{
			Kind: yaml.ScalarNode,
			Tag:  "!!null",
		},
	},
	{
		snippet: "3",
		expected: &yaml.Node{
			Kind:   yaml.ScalarNode,
			Tag:    "!!int",
			Value:  "3",
			Line:   1,
			Column: 1,
		},
	},
	{
		snippet: "cat",
		expected: &yaml.Node{
			Kind:   yaml.ScalarNode,
			Tag:    "!!str",
			Value:  "cat",
			Line:   1,
			Column: 1,
		},
	},
	{
		snippet: "3.1",
		expected: &yaml.Node{
			Kind:   yaml.ScalarNode,
			Tag:    "!!float",
			Value:  "3.1",
			Line:   1,
			Column: 1,
		},
	},
	{
		snippet: "true",
		expected: &yaml.Node{
			Kind:   yaml.ScalarNode,
			Tag:    "!!bool",
			Value:  "true",
			Line:   1,
			Column: 1,
		},
	},
}

func TestParseSnippet(t *testing.T) {
	for _, tt := range parseSnippetScenarios {
		actual, err := parseSnippet(tt.snippet)
		if err != nil {
			t.Error(tt.snippet)
			t.Error(err)
		}
		test.AssertResultComplexWithContext(t, tt.expected, actual, tt.snippet)
	}
}
