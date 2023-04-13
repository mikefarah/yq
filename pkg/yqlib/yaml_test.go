package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var yamlScenarios = []formatScenario{
	// {
	// 	description: "basic - null",
	// 	skipDoc:     true,
	// 	input:       "null",
	// 	expected:    "null\n",
	// },
	{
		description: "basic - ~",
		skipDoc:     true,
		input:       "~",
		expected:    "~\n",
	},
	// {
	// 	description: "basic - [null]",
	// 	skipDoc:     true,
	// 	input:       "[null]",
	// 	expected:    "[null]\n",
	// },
	// {
	// 	description: "basic - [~]",
	// 	skipDoc:     true,
	// 	input:       "[~]",
	// 	expected:    "[~]\n",
	// },
	// {
	// 	description: "basic - null map value",
	// 	skipDoc:     true,
	// 	input:       "a: null",
	// 	expected:    "a: null\n",
	// },
}

func testYamlScenario(t *testing.T, s formatScenario) {
	// switch s.scenarioType {
	// case "decode":
	test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewYamlEncoder(2, false, ConfiguredYamlPreferences)), s.description)
	// default:
	// panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	// }
}

func TestYamlScenarios(t *testing.T) {
	for _, tt := range yamlScenarios {
		testYamlScenario(t, tt)
	}
	// genericScenarios := make([]interface{}, len(yamlScenarios))
	// for i, s := range yamlScenarios {
	// 	genericScenarios[i] = s
	// }
	// documentScenarios(t, "usage", "convert", genericScenarios, documentJSONScenario)
}
