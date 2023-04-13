package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var yamlScenarios = []formatScenario{
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
}

func testYamlScenario(t *testing.T, s formatScenario) {
	test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewYamlEncoder(2, false, ConfiguredYamlPreferences)), s.description)
}

func TestYamlScenarios(t *testing.T) {
	for _, tt := range yamlScenarios {
		testYamlScenario(t, tt)
	}
}
