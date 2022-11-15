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
	snippet       string
	expected      *yaml.Node
	expectedError string
}

var parseSnippetScenarios = []parseSnippetScenario{
	{
		snippet:       ":",
		expectedError: "yaml: did not find expected key",
	},
	{
		snippet: "",
		expected: &yaml.Node{
			Kind: yaml.ScalarNode,
			Tag:  "!!null",
		},
	},
	{
		snippet: "null",
		expected: &yaml.Node{
			Kind:   yaml.ScalarNode,
			Tag:    "!!null",
			Value:  "null",
			Line:   0,
			Column: 0,
		},
	},
	{
		snippet: "3",
		expected: &yaml.Node{
			Kind:   yaml.ScalarNode,
			Tag:    "!!int",
			Value:  "3",
			Line:   0,
			Column: 0,
		},
	},
	{
		snippet: "cat",
		expected: &yaml.Node{
			Kind:   yaml.ScalarNode,
			Tag:    "!!str",
			Value:  "cat",
			Line:   0,
			Column: 0,
		},
	},
	{
		snippet: "3.1",
		expected: &yaml.Node{
			Kind:   yaml.ScalarNode,
			Tag:    "!!float",
			Value:  "3.1",
			Line:   0,
			Column: 0,
		},
	},
	{
		snippet: "true",
		expected: &yaml.Node{
			Kind:   yaml.ScalarNode,
			Tag:    "!!bool",
			Value:  "true",
			Line:   0,
			Column: 0,
		},
	},
}

func TestParseSnippet(t *testing.T) {
	for _, tt := range parseSnippetScenarios {
		actual, err := parseSnippet(tt.snippet)
		if tt.expectedError != "" {
			if err == nil {
				t.Errorf("Expected error '%v' but it worked!", tt.expectedError)
			} else {
				test.AssertResultComplexWithContext(t, tt.expectedError, err.Error(), tt.snippet)
			}
			return
		}
		if err != nil {
			t.Error(tt.snippet)
			t.Error(err)
		}
		test.AssertResultComplexWithContext(t, tt.expected, actual, tt.snippet)
	}
}
