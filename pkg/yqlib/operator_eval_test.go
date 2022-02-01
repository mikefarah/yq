package yqlib

import (
	"testing"
)

var evalOperatorScenarios = []expressionScenario{

	{
		description: "Dynamically evaluate a path",
		document:    `{pathExp: '.a.b[] | select(.name == "cat")', a: {b: [{name: dog}, {name: cat}]}}`,
		expression:  `eval(.pathExp)`,
		expected: []string{
			"D0, P[a b 1], (!!map)::{name: cat}\n",
		},
	},
	{
		description:         "Dynamically update a path from an environment variable",
		subdescription:      "The env variable can be any valid yq expression.",
		document:            `{a: {b: [{name: dog}, {name: cat}]}}`,
		environmentVariable: ".a.b[0].name",
		expression:          `eval(strenv(myenv)) = "cow"`,
		expected: []string{
			"D0, P[], (doc)::{a: {b: [{name: cow}, {name: cat}]}}\n",
		},
	},
}

func TestEvalOperatorsScenarios(t *testing.T) {
	for _, tt := range evalOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "eval", evalOperatorScenarios)
}
