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
