package yqlib

import (
	"testing"
)

var envOperatorScenarios = []expressionScenario{
	{
		description:         "Read boolean environment variable as a string",
		environmentVariable: "true",
		expression:          `strenv(myenv)`,
		expected: []string{
			"D0, P[], (!!str)::\"true\"\n",
		},
	},
	{
		description:         "Read numeric environment variable as a string",
		environmentVariable: "12",
		expression:          `strenv(myenv)`,
		expected: []string{
			"D0, P[], (!!str)::\"12\"\n",
		},
	},
}

func TestEnvOperatorScenarios(t *testing.T) {
	for _, tt := range envOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Env Variable Operators", envOperatorScenarios)
}
