package yqlib

import (
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

func TestGetLogger(t *testing.T) {
	l := GetLogger()
	if l != log {
		t.Fatal("GetLogger should return the yq logger instance, not a copy")
	}
}

type parseSnippetScenario struct {
	snippet       string
	expected      *CandidateNode
	expectedError string
}

var parseSnippetScenarios = []parseSnippetScenario{
	{
		snippet:       ":",
		expectedError: "yaml: did not find expected key",
	},
	{
		snippet: "",
		expected: &CandidateNode{
			Kind: ScalarNode,
			Tag:  "!!null",
		},
	},
	{
		snippet: "null",
		expected: &CandidateNode{
			Kind:   ScalarNode,
			Tag:    "!!null",
			Value:  "null",
			Line:   0,
			Column: 0,
		},
	},
	{
		snippet: "3",
		expected: &CandidateNode{
			Kind:   ScalarNode,
			Tag:    "!!int",
			Value:  "3",
			Line:   0,
			Column: 0,
		},
	},
	{
		snippet: "cat",
		expected: &CandidateNode{
			Kind:   ScalarNode,
			Tag:    "!!str",
			Value:  "cat",
			Line:   0,
			Column: 0,
		},
	},
	{
		snippet: "# things",
		expected: &CandidateNode{
			Kind:        ScalarNode,
			Tag:         "!!null",
			LineComment: "# things",
			Line:        0,
			Column:      0,
		},
	},
	{
		snippet: "3.1",
		expected: &CandidateNode{
			Kind:   ScalarNode,
			Tag:    "!!float",
			Value:  "3.1",
			Line:   0,
			Column: 0,
		},
	},
	{
		snippet: "true",
		expected: &CandidateNode{
			Kind:   ScalarNode,
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
			continue
		}
		if err != nil {
			t.Error(tt.snippet)
			t.Error(err)
		}
		test.AssertResultComplexWithContext(t, tt.expected, actual, tt.snippet)
	}
}

type parseInt64Scenario struct {
	numberString         string
	expectedParsedNumber int64
	expectedFormatString string
}

var parseInt64Scenarios = []parseInt64Scenario{
	{
		numberString:         "34",
		expectedParsedNumber: 34,
	},
	{
		numberString:         "10_000",
		expectedParsedNumber: 10000,
		expectedFormatString: "10000",
	},
	{
		numberString:         "0x10",
		expectedParsedNumber: 16,
	},
	{
		numberString:         "0x10_000",
		expectedParsedNumber: 65536,
		expectedFormatString: "0x10000",
	},
	{
		numberString:         "0o10",
		expectedParsedNumber: 8,
	},
}

func TestParseInt64(t *testing.T) {
	for _, tt := range parseInt64Scenarios {
		format, actualNumber, err := parseInt64(tt.numberString)

		if err != nil {
			t.Error(tt.numberString)
			t.Error(err)
		}
		test.AssertResultComplexWithContext(t, tt.expectedParsedNumber, actualNumber, tt.numberString)
		if tt.expectedFormatString == "" {
			tt.expectedFormatString = tt.numberString
		}

		test.AssertResultComplexWithContext(t, tt.expectedFormatString, fmt.Sprintf(format, actualNumber), fmt.Sprintf("Formatting of: %v", tt.numberString))
	}
}
