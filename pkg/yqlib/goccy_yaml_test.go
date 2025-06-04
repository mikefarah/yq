package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var goccyYamlFormatScenarios = []formatScenario{
	{
		description: "basic scalar - integer",
		skipDoc:     true,
		input:       "3",
		expected:    "3\n",
	},
	{
		description: "basic scalar - float",
		skipDoc:     true,
		input:       "3.1",
		expected:    "3.1\n",
	},
	{
		description: "basic scalar - string",
		skipDoc:     true,
		input:       "hello",
		expected:    "hello\n",
	},
	{
		description: "basic scalar - boolean",
		skipDoc:     true,
		input:       "true",
		expected:    "true\n",
	},
	{
		description: "basic scalar - null",
		skipDoc:     true,
		input:       "null",
		expected:    "",
	},
	{
		description: "basic scalar - tilde null",
		skipDoc:     true,
		input:       "~",
		expected:    "",
	},
	{
		description: "basic mapping",
		skipDoc:     true,
		input:       "key: value",
		expected:    "key: value\n",
	},
	{
		description: "mapping with multiple entries",
		skipDoc:     true,
		input:       "name: John\nage: 30",
		expected:    "name: John\nage: 30\n",
	},
	{
		description: "basic sequence",
		skipDoc:     true,
		input:       "- one\n- two\n- three",
		expected:    "- one\n- two\n- three\n",
	},
	{
		description: "flow style sequence",
		skipDoc:     true,
		input:       "[1, 2, 3]",
		expected:    "[1, 2, 3]\n",
	},
	{
		description: "flow style mapping",
		skipDoc:     true,
		input:       "{name: John, age: 30}",
		expected:    "{name: John, age: 30}\n",
	},
	{
		description: "nested structure",
		skipDoc:     true,
		input:       "person:\n  name: John\n  details:\n    age: 30\n    city: NYC",
		expected:    "person:\n  name: John\n  details:\n    age: 30\n    city: NYC\n",
	},
	{
		description: "quoted strings - single",
		skipDoc:     true,
		input:       "message: 'hello world'",
		expected:    "message: 'hello world'\n",
	},
	{
		description: "quoted strings - double",
		skipDoc:     true,
		input:       "message: \"hello world\"",
		expected:    "message: \"hello world\"\n",
	},
	{
		description: "literal block scalar",
		skipDoc:     true,
		input:       "text: |\n  line one\n  line two",
		expected:    "text: |-\n  line one\n  line two\n",
	},
	{
		description: "folded block scalar",
		skipDoc:     true,
		input:       "text: >\n  line one\n  line two",
		expected:    "text: >-\n  line one line two\n",
	},
	{
		description: "custom tag",
		skipDoc:     true,
		input:       "value: !custom tag_content",
		expected:    "value: !custom tag_content\n",
	},
	{
		description: "anchors and aliases",
		skipDoc:     true,
		input:       "default: &default_value\n  key: value\nother: *default_value",
		expected:    "default: &default_value\n  key: value\nother: *default_value\n",
	},
	{
		description: "merge keys",
		skipDoc:     true,
		input:       "default: &default\n  key1: value1\n  key2: value2\nmerged:\n  <<: *default\n  key3: value3",
		expected:    "default: &default\n  key1: value1\n  key2: value2\nmerged:\n  !!merge <<: *default\n  key3: value3\n",
	},
	{
		description: "array with null",
		skipDoc:     true,
		input:       "[null, \"value\", ~]",
		expected:    "[null, \"value\", ~]\n",
	},
	{
		description: "mapping with null value",
		skipDoc:     true,
		input:       "a: null\nb: ~\nc: value",
		expected:    "a: null\nb: ~\nc: value\n",
	},
	{
		description: "comments - line comment",
		skipDoc:     true,
		input:       "key: value # this is a comment",
		expected:    "key: value # this is a comment\n",
	},
	{
		description: "comments - head comment",
		skipDoc:     true,
		input:       "# this is a head comment\nkey: value",
		expected:    "# this is a head comment\nkey: value\n",
	},
	{
		description: "document separator - first doc only",
		skipDoc:     true,
		input:       "doc1: value1\n---\ndoc2: value2",
		expected:    "doc1: value1\n---\ndoc2: value2\n",
	},
	{
		description: "empty document",
		skipDoc:     true,
		input:       "",
		expected:    "",
	},
	{
		description: "whitespace handling",
		skipDoc:     true,
		input:       "   key   :   value   ",
		expected:    "key: value\n",
	},
	{
		description:  "roundtrip - basic",
		skipDoc:      true,
		input:        "name: John\nage: 30",
		expected:     "age: 30\nname: John\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip - with anchors",
		skipDoc:      true,
		input:        "default: &ref\n  key: value\nother: *ref",
		expected:     "default:\n  key: value\nother:\n  key: value\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "roundtrip - complex structure",
		skipDoc:      true,
		input:        "users:\n  - name: Alice\n    age: 25\n  - name: Bob\n    age: 30",
		expected:     "users:\n- age: 25\n  name: Alice\n- age: 30\n  name: Bob\n",
		scenarioType: "roundtrip",
	},
}

func testGoccyYamlScenario(t *testing.T, s formatScenario) {
	switch s.scenarioType {
	case "roundtrip":
		// Test goccy decoder -> goccy encoder roundtrip
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewGoccyYAMLDecoder(ConfiguredYamlPreferences), NewGoccyYamlEncoder(ConfiguredYamlPreferences)), s.description)
	default:
		// Default: goccy decoder -> yaml encoder
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewGoccyYAMLDecoder(ConfiguredYamlPreferences), NewYamlEncoder(ConfiguredYamlPreferences)), s.description)
	}
}

func TestGoccyYmlFormatScenarios(t *testing.T) {
	for _, tt := range goccyYamlFormatScenarios {
		testGoccyYamlScenario(t, tt)
	}
}
