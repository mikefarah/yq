package yqlib

import (
	"testing"
)

var envsubstOperatorScenarios = []expressionScenario{
	{
		description:         "Replace strings with envsubst",
		environmentVariable: "cat",
		expression:          `"the ${myenv} meows" | envsubst`,
		expected: []string{
			"D0, P[], (!!str)::the cat meows\n",
		},
	},
	{
		description:         "Replace strings with envsubst, missing variables",
		environmentVariable: "cat",
		expression:          `"the ${myenvnonexisting} meows" | envsubst`,
		expected: []string{
			"D0, P[], (!!str)::the  meows\n",
		},
	},
	{
		description:         "Replace strings with envsubst, missing variables with defaults",
		environmentVariable: "cat",
		expression:          `"the ${myenvnonexisting-dog} meows" | envsubst`,
		expected: []string{
			"D0, P[], (!!str)::the dog meows\n",
		},
	},
	{
		description:         "Replace string environment variable in document",
		environmentVariable: "cat meow",
		document:            "{v: \"${myenv}\"}",
		expression:          `.v |= envsubst`,
		expected: []string{
			"D0, P[], (doc)::{v: \"cat meow\"}\n",
		},
	},
}

func TestEnvSubstOperatorScenarios(t *testing.T) {
	for _, tt := range envsubstOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "envsubst", envsubstOperatorScenarios)
}
