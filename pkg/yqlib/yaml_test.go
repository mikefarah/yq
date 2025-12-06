package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var yamlFormatScenarios = []formatScenario{
	{
		description: "scalar with doc separator",
		skipDoc:     true,
		input:       "--- cat",
		expected:    "---\ncat\n",
	},
	{
		description: "CRLF doc separator",
		skipDoc:     true,
		input:       "---\r\ncat\r\n",
		expected:    "---\r\ncat\r\n",
	},
	{
		description: "yaml directive preserved (LF)",
		skipDoc:     true,
		input:       "%YAML 1.1\n---\ncat\n",
		expected:    "%YAML 1.1\n---\ncat\n",
	},
	{
		description: "yaml directive preserved (CRLF)",
		skipDoc:     true,
		input:       "%YAML 1.1\r\n---\r\ncat\r\n",
		expected:    "%YAML 1.1\r\n---\r\ncat\r\n",
	},
	{
		description: "comment only no trailing newline",
		skipDoc:     true,
		input:       "# hello",
		expected:    "# hello\n",
	},

	{
		description: "scalar with doc separator",
		skipDoc:     true,
		input:       "---cat",
		expected:    "---cat\n",
	},
	{
		description: "basic - null",
		skipDoc:     true,
		input:       "null",
		expected:    "null\n",
	},
	{
		description: "basic - ~",
		skipDoc:     true,
		input:       "~",
		expected:    "~\n",
	},
	{
		description: "octal",
		skipDoc:     true,
		input:       "0o30",
		expression:  "tag",
		expected:    "!!int\n",
	},
	{
		description: "basic - [null]",
		skipDoc:     true,
		input:       "[null]",
		expected:    "[null]\n",
	},
	{
		description: "multi document anchor map",
		skipDoc:     true,
		input:       "a: &remember mike\n---\nb: *remember",
		expression:  "explode(.)",
		expected:    "a: mike\n---\nb: mike\n",
	},
	{
		description: "basic - [~]",
		skipDoc:     true,
		input:       "[~]",
		expected:    "[~]\n",
	},
	{
		description: "basic - null map value",
		skipDoc:     true,
		input:       "a: null",
		expected:    "a: null\n",
	},
	{
		description: "basic - number",
		skipDoc:     true,
		input:       "3",
		expected:    "3\n",
	},
	{
		description: "basic - float",
		skipDoc:     true,
		input:       "3.1",
		expected:    "3.1\n",
	},
	{
		description: "basic - float",
		skipDoc:     true,
		input:       "[1, 2]",
		expected:    "[1, 2]\n",
	},
}

var yamlParseScenarios = []expressionScenario{
	// {
	// 	description: "with a unquoted question mark in the string",
	// 	document:    "foo: {bar: a?bc}",
	// 	expected: []string{
	// 		"D0, P[], (!!map)::a: hello # things\n",
	// 	},
	// },
	{
		document: `a: hello # things`,
		expected: []string{
			"D0, P[], (!!map)::a: hello # things\n",
		},
	},
	{
		document:   "a: &a apple\nb: *a",
		expression: ".b | explode(.)",
		expected: []string{
			"D0, P[b], (!!str)::apple\n",
		},
	},
	{
		document: `a: [1,2]`,
		expected: []string{
			"D0, P[], (!!map)::a: [1, 2]\n",
		},
	},
	{
		document: `a: !horse [a]`,
		expected: []string{
			"D0, P[], (!!map)::a: !horse [a]\n",
		},
	},
}

func testYamlScenario(t *testing.T, s formatScenario) {
	test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewYamlEncoder(ConfiguredYamlPreferences)), s.description)
}

func TestYamlParseScenarios(t *testing.T) {
	for _, tt := range yamlParseScenarios {
		testScenario(t, &tt)
	}
}

func TestYamlFormatScenarios(t *testing.T) {
	for _, tt := range yamlFormatScenarios {
		testYamlScenario(t, tt)
	}
}
