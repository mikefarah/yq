package yqlib

import (
	"fmt"
	"strings"
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

func TestGetContentValueByKey(t *testing.T) {
	// Create content with key-value pairs
	key1 := createStringScalarNode("key1")
	value1 := createStringScalarNode("value1")
	key2 := createStringScalarNode("key2")
	value2 := createStringScalarNode("value2")

	content := []*CandidateNode{key1, value1, key2, value2}

	// Test finding existing key
	result := getContentValueByKey(content, "key1")
	test.AssertResult(t, value1, result)

	// Test finding another existing key
	result = getContentValueByKey(content, "key2")
	test.AssertResult(t, value2, result)

	// Test finding non-existing key
	result = getContentValueByKey(content, "nonexistent")
	test.AssertResult(t, (*CandidateNode)(nil), result)

	// Test with empty content
	result = getContentValueByKey([]*CandidateNode{}, "key1")
	test.AssertResult(t, (*CandidateNode)(nil), result)
}

func TestRecurseNodeArrayEqual(t *testing.T) {
	// Create two arrays with same content
	array1 := &CandidateNode{
		Kind: SequenceNode,
		Content: []*CandidateNode{
			createStringScalarNode("item1"),
			createStringScalarNode("item2"),
		},
	}

	array2 := &CandidateNode{
		Kind: SequenceNode,
		Content: []*CandidateNode{
			createStringScalarNode("item1"),
			createStringScalarNode("item2"),
		},
	}

	array3 := &CandidateNode{
		Kind: SequenceNode,
		Content: []*CandidateNode{
			createStringScalarNode("item1"),
			createStringScalarNode("different"),
		},
	}

	array4 := &CandidateNode{
		Kind: SequenceNode,
		Content: []*CandidateNode{
			createStringScalarNode("item1"),
		},
	}

	test.AssertResult(t, true, recurseNodeArrayEqual(array1, array2))
	test.AssertResult(t, false, recurseNodeArrayEqual(array1, array3))
	test.AssertResult(t, false, recurseNodeArrayEqual(array1, array4))
}

func TestFindInArray(t *testing.T) {
	item1 := createStringScalarNode("item1")
	item2 := createStringScalarNode("item2")
	item3 := createStringScalarNode("item3")

	array := &CandidateNode{
		Kind:    SequenceNode,
		Content: []*CandidateNode{item1, item2, item3},
	}

	// Test finding existing items
	test.AssertResult(t, 0, findInArray(array, item1))
	test.AssertResult(t, 1, findInArray(array, item2))
	test.AssertResult(t, 2, findInArray(array, item3))

	// Test finding non-existing item
	nonExistent := createStringScalarNode("nonexistent")
	test.AssertResult(t, -1, findInArray(array, nonExistent))
}

func TestFindKeyInMap(t *testing.T) {
	key1 := createStringScalarNode("key1")
	value1 := createStringScalarNode("value1")
	key2 := createStringScalarNode("key2")
	value2 := createStringScalarNode("value2")

	mapNode := &CandidateNode{
		Kind:    MappingNode,
		Content: []*CandidateNode{key1, value1, key2, value2},
	}

	// Test finding existing keys
	test.AssertResult(t, 0, findKeyInMap(mapNode, key1))
	test.AssertResult(t, 2, findKeyInMap(mapNode, key2))

	// Test finding non-existing key
	nonExistent := createStringScalarNode("nonexistent")
	test.AssertResult(t, -1, findKeyInMap(mapNode, nonExistent))
}

func TestRecurseNodeObjectEqual(t *testing.T) {
	// Create two objects with same content
	key1 := createStringScalarNode("key1")
	value1 := createStringScalarNode("value1")
	key2 := createStringScalarNode("key2")
	value2 := createStringScalarNode("value2")

	obj1 := &CandidateNode{
		Kind:    MappingNode,
		Content: []*CandidateNode{key1, value1, key2, value2},
	}

	obj2 := &CandidateNode{
		Kind:    MappingNode,
		Content: []*CandidateNode{key1, value1, key2, value2},
	}

	// Create object with different values
	value3 := createStringScalarNode("value3")
	obj3 := &CandidateNode{
		Kind:    MappingNode,
		Content: []*CandidateNode{key1, value3, key2, value2},
	}

	// Create object with different keys
	key3 := createStringScalarNode("key3")
	obj4 := &CandidateNode{
		Kind:    MappingNode,
		Content: []*CandidateNode{key1, value1, key3, value2},
	}

	test.AssertResult(t, true, recurseNodeObjectEqual(obj1, obj2))
	test.AssertResult(t, false, recurseNodeObjectEqual(obj1, obj3))
	test.AssertResult(t, false, recurseNodeObjectEqual(obj1, obj4))
}

func TestParseInt(t *testing.T) {
	type parseIntScenario struct {
		numberString         string
		expectedParsedNumber int
		expectedError        string
	}

	scenarios := []parseIntScenario{
		{
			numberString:         "34",
			expectedParsedNumber: 34,
		},
		{
			numberString:         "10_000",
			expectedParsedNumber: 10000,
		},
		{
			numberString:         "0x10",
			expectedParsedNumber: 16,
		},
		{
			numberString:         "0o10",
			expectedParsedNumber: 8,
		},
		{
			numberString:  "invalid",
			expectedError: "strconv.ParseInt",
		},
	}

	for _, tt := range scenarios {
		actualNumber, err := parseInt(tt.numberString)
		if tt.expectedError != "" {
			if err == nil {
				t.Errorf("Expected error for '%s' but got none", tt.numberString)
			} else if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("Expected error containing '%s' for '%s', got '%s'", tt.expectedError, tt.numberString, err.Error())
			}
			continue
		}
		if err != nil {
			t.Errorf("Unexpected error for '%s': %v", tt.numberString, err)
		}
		test.AssertResultComplexWithContext(t, tt.expectedParsedNumber, actualNumber, tt.numberString)
	}
}

func TestHeadAndLineComment(t *testing.T) {
	node := &CandidateNode{
		HeadComment: "# head comment",
		LineComment: "# line comment",
	}

	result := headAndLineComment(node)
	test.AssertResult(t, " head comment line comment", result)
}

func TestHeadComment(t *testing.T) {
	node := &CandidateNode{
		HeadComment: "# head comment",
	}

	result := headComment(node)
	test.AssertResult(t, " head comment", result)

	// Test without #
	node.HeadComment = "no hash comment"
	result = headComment(node)
	test.AssertResult(t, "no hash comment", result)
}

func TestLineComment(t *testing.T) {
	node := &CandidateNode{
		LineComment: "# line comment",
	}

	result := lineComment(node)
	test.AssertResult(t, " line comment", result)

	// Test without #
	node.LineComment = "no hash comment"
	result = lineComment(node)
	test.AssertResult(t, "no hash comment", result)
}

func TestFootComment(t *testing.T) {
	node := &CandidateNode{
		FootComment: "# foot comment",
	}

	result := footComment(node)
	test.AssertResult(t, " foot comment", result)

	// Test without #
	node.FootComment = "no hash comment"
	result = footComment(node)
	test.AssertResult(t, "no hash comment", result)
}

func TestKindString(t *testing.T) {
	test.AssertResult(t, "ScalarNode", KindString(ScalarNode))
	test.AssertResult(t, "SequenceNode", KindString(SequenceNode))
	test.AssertResult(t, "MappingNode", KindString(MappingNode))
	test.AssertResult(t, "AliasNode", KindString(AliasNode))
	test.AssertResult(t, "unknown!", KindString(Kind(999))) // Invalid kind
}

type processEscapeCharactersScenario struct {
	input    string
	expected string
}

var processEscapeCharactersScenarios = []processEscapeCharactersScenario{
	{
		input:    "",
		expected: "",
	},
	{
		input:    "hello",
		expected: "hello",
	},
	{
		input:    "\\\"",
		expected: "\"",
	},
	{
		input:    "hello\\\"world",
		expected: "hello\"world",
	},
	{
		input:    "\\n",
		expected: "\n",
	},
	{
		input:    "line1\\nline2",
		expected: "line1\nline2",
	},
	{
		input:    "\\t",
		expected: "\t",
	},
	{
		input:    "hello\\tworld",
		expected: "hello\tworld",
	},
	{
		input:    "\\r",
		expected: "\r",
	},
	{
		input:    "hello\\rworld",
		expected: "hello\rworld",
	},
	{
		input:    "\\f",
		expected: "\f",
	},
	{
		input:    "hello\\fworld",
		expected: "hello\fworld",
	},
	{
		input:    "\\v",
		expected: "\v",
	},
	{
		input:    "hello\\vworld",
		expected: "hello\vworld",
	},
	{
		input:    "\\b",
		expected: "\b",
	},
	{
		input:    "hello\\bworld",
		expected: "hello\bworld",
	},
	{
		input:    "\\a",
		expected: "\a",
	},
	{
		input:    "hello\\aworld",
		expected: "hello\aworld",
	},
	{
		input:    "\\\"\\n\\t\\r\\f\\v\\b\\a",
		expected: "\"\n\t\r\f\v\b\a",
	},
	{
		input:    "multiple\\nlines\\twith\\ttabs",
		expected: "multiple\nlines\twith\ttabs",
	},
	{
		input:    "quote\\\"here",
		expected: "quote\"here",
	},
	{
		input:    "\\\\",
		expected: "\\", // Backslash is processed: "\\\\" becomes "\\"
	},
	{
		input:    "\\\"test\\\"",
		expected: "\"test\"",
	},
	{
		input:    "a\\\\b",
		expected: "a\\b", // Tests roundtrip: "a\\\\b" should become "a\\b"
	},
	{
		input:    "Hi \\\\(.value)",
		expected: "Hi \\\\(.value)",
	},
	{
		input:    `a\\b`,
		expected: "a\\b",
	},
}

func TestProcessEscapeCharacters(t *testing.T) {
	for _, tt := range processEscapeCharactersScenarios {
		actual := processEscapeCharacters(tt.input)
		test.AssertResultComplexWithContext(t, tt.expected, actual, fmt.Sprintf("Input: %q", tt.input))
	}
}
