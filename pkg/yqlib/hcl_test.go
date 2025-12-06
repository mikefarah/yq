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
		expected:     "tags: ' [\"a\", \"b\"]'\n",
		scenarioType: "decode",
	},
	{
		description:  "object/map attribute",
		input:        `obj = { a = 1, b = "two" }`,
		expected:     "obj: ' { a = 1, b = \"two\" }'\n",
		scenarioType: "decode",
	},
	{
		description:  "nested block",
		input:        `server { port = 8080 }`,
		expected:     "server:\n  port: 8080\n",
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
