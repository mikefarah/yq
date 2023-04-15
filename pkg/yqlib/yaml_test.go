package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var yamlFormatScenarios = []formatScenario{
	{
		description: "basic - null",
		skipDoc:     true,
		input:       "null",
		expected:    "null\n",
	},
	// {
	// 	description: "basic - ~",
	// 	skipDoc:     true,
	// 	input:       "~",
	// 	expected:    "~\n",
	// },
	{
		description: "basic - [null]",
		skipDoc:     true,
		input:       "[null]",
		expected:    "[null]\n",
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
	{
		document: `a: hello # things`,
		expected: []string{
			"D0, P[], (doc)::a: hello # things\n",
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
			"D0, P[], (doc)::a: [1, 2]\n",
		},
	},
}

func testYamlScenario(t *testing.T, s formatScenario) {
	test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewYamlEncoder(2, false, ConfiguredYamlPreferences)), s.description)
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
