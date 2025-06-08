package yqlib

import (
	"testing"
)

func TestDateTimePreprocessorBasicFunctionality(t *testing.T) {
	preprocessor := NewDateTimePreprocessor(true)

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "RFC3339 with timezone",
			input:    "timestamp: 2021-01-01T00:00:00Z",
			expected: "timestamp: !!timestamp 2021-01-01T00:00:00Z",
		},
		{
			name:     "RFC3339 with offset timezone",
			input:    "timestamp: 2021-01-01T03:10:00+03:00",
			expected: "timestamp: !!timestamp 2021-01-01T03:10:00+03:00",
		},
		{
			name:     "Date only",
			input:    "date: 2021-01-01",
			expected: "date: !!timestamp 2021-01-01",
		},
		{
			name:     "RFC3339 without timezone",
			input:    "timestamp: 2021-01-01T15:04:05",
			expected: "timestamp: !!timestamp 2021-01-01T15:04:05",
		},
		{
			name:     "RFC3339 with milliseconds",
			input:    "timestamp: 2021-01-01T00:00:00.123Z",
			expected: "timestamp: !!timestamp 2021-01-01T00:00:00.123Z",
		},
		{
			name:     "Array item with timestamp",
			input:    "- 2021-01-01T00:00:00Z",
			expected: "- !!timestamp 2021-01-01T00:00:00Z",
		},
		{
			name:     "Multiple timestamps in document",
			input:    "start: 2021-01-01T00:00:00Z\nend: 2021-12-31T23:59:59Z",
			expected: "start: !!timestamp 2021-01-01T00:00:00Z\nend: !!timestamp 2021-12-31T23:59:59Z",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := preprocessor.PreprocessDocument(tc.input)
			if result != tc.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tc.expected, result)
			}
		})
	}
}

func TestDateTimePreprocessorDisabled(t *testing.T) {
	preprocessor := NewDateTimePreprocessor(false)

	input := "timestamp: 2021-01-01T00:00:00Z"
	expected := "timestamp: 2021-01-01T00:00:00Z"

	result := preprocessor.PreprocessDocument(input)
	if result != expected {
		t.Errorf("Expected preprocessing to be disabled. Got: %s", result)
	}
}

func TestDateTimePreprocessorSkipsCases(t *testing.T) {
	preprocessor := NewDateTimePreprocessor(true)

	testCases := []struct {
		name     string
		input    string
		expected string // Should remain unchanged
	}{
		{
			name:     "Already tagged",
			input:    "timestamp: !!timestamp 2021-01-01T00:00:00Z",
			expected: "timestamp: !!timestamp 2021-01-01T00:00:00Z",
		},
		{
			name:     "Quoted string",
			input:    "timestamp: \"2021-01-01T00:00:00Z\"",
			expected: "timestamp: \"2021-01-01T00:00:00Z\"",
		},
		{
			name:     "Single quoted string",
			input:    "timestamp: '2021-01-01T00:00:00Z'",
			expected: "timestamp: '2021-01-01T00:00:00Z'",
		},
		{
			name:     "Comment line",
			input:    "# This is 2021-01-01T00:00:00Z",
			expected: "# This is 2021-01-01T00:00:00Z",
		},
		{
			name:     "Empty line",
			input:    "",
			expected: "",
		},
		{
			name:     "YAML directive",
			input:    "%YAML 1.1",
			expected: "%YAML 1.1",
		},
		{
			name:     "Document separator",
			input:    "---",
			expected: "---",
		},
		{
			name:     "Document end",
			input:    "...",
			expected: "...",
		},
		{
			name:     "Anchor reference",
			input:    "timestamp: *my_timestamp",
			expected: "timestamp: *my_timestamp",
		},
		{
			name:     "Anchor definition",
			input:    "timestamp: &my_timestamp 2021-01-01T00:00:00Z",
			expected: "timestamp: &my_timestamp 2021-01-01T00:00:00Z",
		},
		{
			name:     "Map value",
			input:    "timestamp: {year: 2021, month: 01}",
			expected: "timestamp: {year: 2021, month: 01}",
		},
		{
			name:     "Array value",
			input:    "timestamp: [2021, 01, 01]",
			expected: "timestamp: [2021, 01, 01]",
		},
		{
			name:     "Multi-line scalar literal",
			input:    "description: |",
			expected: "description: |",
		},
		{
			name:     "Multi-line scalar folded",
			input:    "description: >",
			expected: "description: >",
		},
		{
			name:     "Invalid date format",
			input:    "invalid: 2021/01/01",
			expected: "invalid: 2021/01/01",
		},
		{
			name:     "Non-date numeric string",
			input:    "number: 20210101",
			expected: "number: 20210101",
		},
		{
			name:     "Partial date",
			input:    "partial: 2021-01",
			expected: "partial: 2021-01",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := preprocessor.PreprocessDocument(tc.input)
			if result != tc.expected {
				t.Errorf("Expected no change. Input: %s, Got: %s", tc.input, result)
			}
		})
	}
}

func TestDateTimePreprocessorComplexDocuments(t *testing.T) {
	preprocessor := NewDateTimePreprocessor(true)

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Mixed document with comments and timestamps",
			input: `# Configuration file
version: 1.0
created: 2021-01-01T00:00:00Z
# Last modified date
modified: 2021-12-31T23:59:59Z
author: "John Doe"
tags:
  - production
  - 2021-01-01`,
			expected: `# Configuration file
version: 1.0
created: !!timestamp 2021-01-01T00:00:00Z
# Last modified date
modified: !!timestamp 2021-12-31T23:59:59Z
author: "John Doe"
tags:
  - production
  - !!timestamp 2021-01-01`,
		},
		{
			name: "Nested structures with timestamps",
			input: `events:
  start:
    date: 2021-01-01
    time: 2021-01-01T09:00:00Z
  end:
    date: 2021-01-02
    time: 2021-01-02T17:00:00Z`,
			expected: `events:
  start:
    date: !!timestamp 2021-01-01
    time: !!timestamp 2021-01-01T09:00:00Z
  end:
    date: !!timestamp 2021-01-02
    time: !!timestamp 2021-01-02T17:00:00Z`,
		},
		{
			name: "Array of timestamps",
			input: `dates:
  - 2021-01-01T00:00:00Z
  - 2021-06-15T12:30:45Z
  - 2021-12-31T23:59:59Z`,
			expected: `dates:
  - !!timestamp 2021-01-01T00:00:00Z
  - !!timestamp 2021-06-15T12:30:45Z
  - !!timestamp 2021-12-31T23:59:59Z`,
		},
		{
			name: "Mixed array with timestamps and other values",
			input: `mixed:
  - "string value"
  - 2021-01-01T00:00:00Z
  - 42
  - 2021-12-31`,
			expected: `mixed:
  - "string value"
  - !!timestamp 2021-01-01T00:00:00Z
  - 42
  - !!timestamp 2021-12-31`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := preprocessor.PreprocessDocument(tc.input)
			if result != tc.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tc.expected, result)
			}
		})
	}
}

func TestDateTimePreprocessorEdgeCases(t *testing.T) {
	preprocessor := NewDateTimePreprocessor(true)

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Whitespace around timestamp",
			input:    "timestamp:   2021-01-01T00:00:00Z   ",
			expected: "timestamp: !!timestamp 2021-01-01T00:00:00Z",
		},
		{
			name:     "Indented timestamp",
			input:    "  timestamp: 2021-01-01T00:00:00Z",
			expected: "  timestamp: !!timestamp 2021-01-01T00:00:00Z",
		},
		{
			name:     "Deeply indented array item",
			input:    "    - 2021-01-01T00:00:00Z",
			expected: "    - !!timestamp 2021-01-01T00:00:00Z",
		},
		{
			name:     "Empty value",
			input:    "timestamp:",
			expected: "timestamp:",
		},
		{
			name:     "Only key",
			input:    "timestamp",
			expected: "timestamp",
		},
		{
			name:     "Colon in value (not key-value pair)",
			input:    "description: Time is 15:30:45",
			expected: "description: Time is 15:30:45",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := preprocessor.PreprocessDocument(tc.input)
			if result != tc.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tc.expected, result)
			}
		})
	}
}

func TestIsValidDateTime(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"2021-01-01T00:00:00Z", true},
		{"2021-01-01T00:00:00+03:00", true},
		{"2021-01-01T00:00:00", true},
		{"2021-01-01", true},
		{"2021-01-01T00:00:00.123Z", true},
		{"invalid-date", false},
		{"2021/01/01", false},
		{"2021-13-01", false},
		{"2021-01-32", false},
		{"2021-01-01T25:00:00Z", false},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := isValidDateTime(tc.input)
			if result != tc.expected {
				t.Errorf("isValidDateTime(%s) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestMatchesDateTimePattern(t *testing.T) {
	preprocessor := NewDateTimePreprocessor(true)

	testCases := []struct {
		input    string
		expected bool
	}{
		{"2021-01-01T00:00:00Z", true},
		{"2021-01-01T00:00:00+03:00", true},
		{"2021-01-01T00:00:00", true},
		{"2021-01-01", true},
		{"2021-01-01T00:00:00.123Z", true},
		{"  2021-01-01T00:00:00Z  ", true}, // with whitespace
		{"not-a-date", false},
		{"2021/01/01", false},
		{"15:30:45", false}, // time only
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := preprocessor.matchesDateTimePattern(tc.input)
			if result != tc.expected {
				t.Errorf("matchesDateTimePattern(%s) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}
