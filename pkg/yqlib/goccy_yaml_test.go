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
		description: "basic - mike",
		skipDoc:     true,
		input:       "mike: 3",
		expected:    "mike: 3\n",
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
		input:       "{\n mike: 3\n}",
		expected:    "{mike: 3}\n",
	},
	{
		description: "basic - tag with number",
		skipDoc:     true,
		input:       "mike: !!cat 3",
		expected:    "mike: !!cat 3\n",
	},
	{
		description: "basic - array of numbers",
		skipDoc:     true,
		input:       "- 3",
		expected:    "- 3\n",
	},
	{
		description: "basic - single line array",
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
	// {
	// 	description: "basic - head comment",
	// 	skipDoc:     true,
	// 	input:       "# head comment\na: meow\n",
	// 	expected:    "# head comment\na: meow\n", // go-yaml does this
	// },
	// {
	// 	description: "basic - head and line comment",
	// 	skipDoc:     true,
	// 	input:       "# head comment\na: #line comment\n  meow\n",
	// 	expected:    "# head comment\na: meow #line comment\n", // go-yaml does this
	// },
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
	{
		description: "basic - ~",
		skipDoc:     true,
		input:       "null\n",
		expected:    "null\n",
	},
	{
		skipDoc:     true,
		description: "blank value round trip",
		input:       "test:",
		expected:    "test:\n",
	},
	{
		skipDoc:     true,
		description: "trailing comment",
		input:       "test: null\n# this comment will be removed",
		expected:    "test: null\n# this comment will be removed\n",
	},
	// {
	// 	description: "doc separator",
	// 	skipDoc:     true,
	// 	input:       "# hi\n---\na: cat\n---",
	// 	expected:    "---\na: cat\n",
	// },
	// {
	// 	description: "scalar with doc separator",
	// 	skipDoc:     true,
	// 	input:       "--- cat",
	// 	expected:    "---\ncat\n",
	// },
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
		description: "multi document",
		skipDoc:     true,
		input:       "a: mike\n---\nb: remember",
		expected:    "a: mike\n---\nb: remember\n",
	},
	{
		description: "single doc anchor map",
		skipDoc:     true,
		input:       "a: &remember mike\nb: *remember",
		expected:    "a: &remember mike\nb: *remember\n",
	},
	{
		description: "explode doc anchor map",
		skipDoc:     true,
		input:       "a: &remember mike\nb: *remember",
		expression:  "explode(.)",
		expected:    "a: mike\nb: mike\n",
	},
	{
		description: "multi document anchor map",
		skipDoc:     true,
		input:       "a: &remember mike\n---\nb: *remember",
		expression:  "explode(.)",
		expected:    "a: mike\n---\nb: mike\n",
	},
	{
		description: "merge anchor",
		skipDoc:     true,
		input:       "a: &remember\n  c: mike\nb:\n  <<: *remember",
		// fine to have !!merge as that's what the current impl does
		expected: "a: &remember\n  c: mike\nb:\n  !!merge <<: *remember\n",
	},
	{
		description: "custom tag",
		skipDoc:     true,
		input:       "a: !cat mike",
		expected:    "a: !cat mike\n",
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

func testGoccyYamlScenario(t *testing.T, s formatScenario) {
	test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewGoccyYAMLDecoder(), NewYamlEncoder(ConfiguredYamlPreferences)), s.description)
}

func TestGoccyYmlFormatScenarios(t *testing.T) {
	for _, tt := range goccyYamlFormatScenarios {
		testGoccyYamlScenario(t, tt)
	}
}
