package yqlib

import (
	"testing"
)

var envOperatorScenarios = []expressionScenario{
	{
		description:         "Read string environment variable",
		environmentVariable: "cat meow",
		expression:          `.a = env(myenv)`,
		expected: []string{
			"D0, P[], ()::a: cat meow\n",
		},
	},
	{
		description:         "Read boolean environment variable",
		environmentVariable: "true",
		expression:          `.a = env(myenv)`,
		expected: []string{
			"D0, P[], ()::a: true\n",
		},
	},
	{
		description:         "Read numeric environment variable",
		environmentVariable: "12",
		expression:          `.a = env(myenv)`,
		expected: []string{
			"D0, P[], ()::a: 12\n",
		},
	},
	{
		description:         "Read yaml environment variable",
		environmentVariable: "{b: fish}",
		expression:          `.a = env(myenv)`,
		expected: []string{
			"D0, P[], ()::a: {b: fish}\n",
		},
	},
	{
		description:         "Read boolean environment variable as a string",
		environmentVariable: "true",
		expression:          `.a = strenv(myenv)`,
		expected: []string{
			"D0, P[], ()::a: \"true\"\n",
		},
	},
	{
		description:         "Read numeric environment variable as a string",
		environmentVariable: "12",
		expression:          `.a = strenv(myenv)`,
		expected: []string{
			"D0, P[], ()::a: \"12\"\n",
		},
	},
	{
		description:         "Dynamic key lookup with environment variable",
		environmentVariable: "cat",
		document:            `{cat: meow, dog: woof}`,
		expression:          `.[env(myenv)]`,
		expected: []string{
			"D0, P[cat], (!!str)::meow\n",
		},
	},
}

func TestEnvOperatorScenarios(t *testing.T) {
	for _, tt := range envOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "env-variable-operators", envOperatorScenarios)
}
