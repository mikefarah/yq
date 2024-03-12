package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var goccyYamlFormatScenarios = []FormatScenario{
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
		description: "basic - map multiple entries",
		skipDoc:     true,
		input:       "mike: 3\nfred: 12\n",
		expected:    "mike: 3\nfred: 12\n",
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
	{
		description: "basic - plain string",
		skipDoc:     true,
		input:       `a: meow`,
		expected:    "a: meow\n",
	},
	{
		description: "basic - double quoted string",
		skipDoc:     true,
		input:       `a: "meow"`,
		expected:    "a: \"meow\"\n",
	},
	{
		description: "basic - single quoted string",
		skipDoc:     true,
		input:       `a: 'meow'`,
		expected:    "a: 'meow'\n",
	},
	{
		description: "basic - string block",
		skipDoc:     true,
		input:       "a: |\n  meow\n",
		expected:    "a: |\n  meow\n",
	},
	{
		description: "basic - long string",
		skipDoc:     true,
		input:       "a: the cute cat wrote a long sentence that wasn't wrapped at all.\n",
		expected:    "a: the cute cat wrote a long sentence that wasn't wrapped at all.\n",
	},
	{
		description: "basic - string block",
		skipDoc:     true,
		input:       "a: |-\n  meow\n",
		expected:    "a: |-\n  meow\n",
	},
	{
		description: "basic - line comment",
		skipDoc:     true,
		input:       "a: meow # line comment\n",
		expected:    "a: meow # line comment\n",
	},
	{
		description: "basic - line comment",
		skipDoc:     true,
		input:       "# head comment\na: #line comment\n  meow\n",
		expected:    "# head comment\na: meow #line comment\n", // go-yaml does this
	},
	{
		description: "basic - foot comment",
		skipDoc:     true,
		input:       "a: meow\n# foot comment\n",
		expected:    "a: meow\n# foot comment\n",
	},
	{
		description: "basic - foot comment",
		skipDoc:     true,
		input:       "a: meow\nb: woof\n# foot comment\n",
		expected:    "a: meow\nb: woof\n# foot comment\n",
	},
	{
		description: "basic - boolean",
		skipDoc:     true,
		input:       "true\n",
		expected:    "true\n",
	},
	{
		description: "basic - null",
		skipDoc:     true,
		input:       "a: null\n",
		expected:    "a: null\n",
	},
	{
		description: "basic - ~",
		skipDoc:     true,
		input:       "a: ~\n",
		expected:    "a: ~\n",
	},
	// {
	// 	description: "basic - ~",
	// 	skipDoc:     true,
	// 	input:       "null\n",
	// 	expected:    "null\n",
	// },
	// {
	// 	skipDoc:     true,
	// 	description: "trailing comment",
	// 	input:       "test:",
	// 	expected:    "test:",
	// },
	// {
	// 	skipDoc:     true,
	// 	description: "trailing comment",
	// 	input:       "test:\n# this comment will be removed",
	// 	expected:    "test:\n# this comment will be removed",
	// },
}

func testGoccyYamlScenario(t *testing.T, s FormatScenario) {
	test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewGoccyYAMLDecoder(), NewYamlEncoder(ConfiguredYamlPreferences)), s.description)
}

func TestGoccyYmlFormatScenarios(t *testing.T) {
	for _, tt := range goccyYamlFormatScenarios {
		testGoccyYamlScenario(t, tt)
	}
}
