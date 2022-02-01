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
		description:          "Dynamically update a path from an environment variable",
		subdescription:       "The env variable can be any valid yq expression.",
		document:             `{a: {b: [{name: dog}, {name: cat}]}}`,
		environmentVariables: map[string]string{"pathEnv": ".a.b[0].name", "valueEnv": "moo"},
		expression:           `eval(strenv(pathEnv)) = strenv(valueEnv)`,
		expected: []string{
			"D0, P[], (doc)::{a: {b: [{name: moo}, {name: cat}]}}\n",
		},
	},
}

func TestEvalOperatorsScenarios(t *testing.T) {
	for _, tt := range evalOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "eval", evalOperatorScenarios)
}
