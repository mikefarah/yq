package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var goccyYamlFormatScenarios = []formatScenario{
	{
		description: "basic - 3",
		skipDoc:     true,
		input:       "3",
		expected:    "3\n",
	},
	{
		description: "basic - 3.1",
		skipDoc:     true,
		input:       "3.1",
		expected:    "3.1\n",
	},
	{
		description: "basic - 3.1",
		skipDoc:     true,
		input:       "mike: 3",
		expected:    "mike: 3\n",
	},
	{
		description: "basic - 3.1",
		skipDoc:     true,
		input:       "{mike: 3}",
		expected:    "{mike: 3}\n",
	},
	{
		description: "basic - 3.1",
		skipDoc:     true,
		input:       "{\nmike: 3\n}",
		expected:    "{mike: 3}\n",
	},
	{
		description: "basic - 3.1",
		skipDoc:     true,
		input:       "mike: !!cat 3",
		expected:    "mike: !!cat 3\n",
	},
	{
		description: "basic - 3.1",
		skipDoc:     true,
		input:       "- 3",
		expected:    "- 3\n",
	},
	{
		description: "basic - 3.1",
		skipDoc:     true,
		input:       "[3]",
		expected:    "[3]\n",
	},
}

func testGoccyYamlScenario(t *testing.T, s formatScenario) {
	test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewGoccyYAMLDecoder(), NewYamlEncoder(2, false, ConfiguredYamlPreferences)), s.description)
}

func TestGoccyYmlFormatScenarios(t *testing.T) {
	for _, tt := range goccyYamlFormatScenarios {
		testGoccyYamlScenario(t, tt)
	}
}
