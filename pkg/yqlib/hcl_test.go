package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var hclFormatScenarios = []formatScenario{
	{
		description:  "Simple decode",
		input:        `io_mode = "async"`,
		expected:     "io_mode: async\n",
		scenarioType: "decode",
	},
	{
		description:  "number attribute",
		input:        `port = 8080`,
		expected:     "port: 8080\n",
		scenarioType: "decode",
	},
	{
		description:  "float attribute",
		input:        `pi = 3.14`,
		expected:     "pi: 3.14\n",
		scenarioType: "decode",
	},
	{
		description:  "boolean attribute",
		input:        `enabled = true`,
		expected:     "enabled: true\n",
		scenarioType: "decode",
	},
	{
		description:  "list of strings",
		input:        `tags = ["a", "b"]`,
		expected:     "tags:\n  - a\n  - b\n",
		scenarioType: "decode",
	},
	{
		description:  "object/map attribute",
		input:        `obj = { a = 1, b = "two" }`,
		expected:     "obj:\n  a: 1\n  b: two\n",
		scenarioType: "decode",
	},
	{
		description:  "nested block",
		input:        `server { port = 8080 }`,
		expected:     "server:\n  port: 8080\n",
		scenarioType: "decode",
	},
	{
		description:  "multiple attributes",
		input:        "name = \"app\"\nversion = 1\nenabled = true",
		expected:     "name: app\nversion: 1\nenabled: true\n",
		scenarioType: "decode",
	},
	{
		description:  "binary expression",
		input:        `count = 0 - 42`,
		expected:     "count: -42\n",
		scenarioType: "decode",
	},
	{
		description:  "negative number",
		input:        `count = -42`,
		expected:     "count: -42\n",
		scenarioType: "decode",
	},
	{
		description:  "scientific notation",
		input:        `value = 1e-3`,
		expected:     "value: 0.001\n",
		scenarioType: "decode",
	},
	{
		description:  "nested object",
		input:        `config = { db = { host = "localhost", port = 5432 } }`,
		expected:     "config:\n  db:\n    host: localhost\n    port: 5432\n",
		scenarioType: "decode",
	},
	{
		description:  "mixed list",
		input:        `values = [1, "two", true]`,
		expected:     "values:\n  - 1\n  - two\n  - true\n",
		scenarioType: "decode",
	},
	{
		description:  "block with labels",
		input:        `resource "aws_instance" "example" { ami = "ami-12345" }`,
		expected:     "resource aws_instance example:\n  ami: ami-12345\n",
		scenarioType: "decode",
	},
}

func testHclScenario(t *testing.T, s formatScenario) {
	if s.scenarioType == "decode" {
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewHclDecoder(), NewYamlEncoder(ConfiguredYamlPreferences)), s.description)
	}
}

func TestHclFormatScenarios(t *testing.T) {
	for _, tt := range hclFormatScenarios {
		testHclScenario(t, tt)
	}
}
