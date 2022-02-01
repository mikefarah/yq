package yqlib

import (
	"testing"
)

var envOperatorScenarios = []expressionScenario{
	{
		description:          "Read string environment variable",
		environmentVariables: map[string]string{"myenv": "cat meow"},
		expression:           `.a = env(myenv)`,
		expected: []string{
			"D0, P[], ()::a: cat meow\n",
		},
	},
	{
		description:          "Read boolean environment variable",
		environmentVariables: map[string]string{"myenv": "true"},
		expression:           `.a = env(myenv)`,
		expected: []string{
			"D0, P[], ()::a: true\n",
		},
	},
	{
		description:          "Read numeric environment variable",
		environmentVariables: map[string]string{"myenv": "12"},
		expression:           `.a = env(myenv)`,
		expected: []string{
			"D0, P[], ()::a: 12\n",
		},
	},
	{
		description:          "Read yaml environment variable",
		environmentVariables: map[string]string{"myenv": "{b: fish}"},
		expression:           `.a = env(myenv)`,
		expected: []string{
			"D0, P[], ()::a: {b: fish}\n",
		},
	},
	{
		description:          "Read boolean environment variable as a string",
		environmentVariables: map[string]string{"myenv": "true"},
		expression:           `.a = strenv(myenv)`,
		expected: []string{
			"D0, P[], ()::a: \"true\"\n",
		},
	},
	{
		description:          "Read numeric environment variable as a string",
		environmentVariables: map[string]string{"myenv": "12"},
		expression:           `.a = strenv(myenv)`,
		expected: []string{
			"D0, P[], ()::a: \"12\"\n",
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
	{
		description:          "Dynamic key lookup with environment variable",
		environmentVariables: map[string]string{"myenv": "cat"},
		document:             `{cat: meow, dog: woof}`,
		expression:           `.[env(myenv)]`,
		expected: []string{
			"D0, P[cat], (!!str)::meow\n",
		},
	},
	{
		description:          "Replace strings with envsubst",
		environmentVariables: map[string]string{"myenv": "cat"},
		expression:           `"the ${myenv} meows" | envsubst`,
		expected: []string{
			"D0, P[], (!!str)::the cat meows\n",
		},
	},
	{
		description:          "Replace strings with envsubst, missing variables",
		environmentVariables: map[string]string{"myenv": "cat"},
		expression:           `"the ${myenvnonexisting} meows" | envsubst`,
		expected: []string{
			"D0, P[], (!!str)::the  meows\n",
		},
	},
	{
		description:          "Replace strings with envsubst, missing variables with defaults",
		environmentVariables: map[string]string{"myenv": "cat"},
		expression:           `"the ${myenvnonexisting-dog} meows" | envsubst`,
		expected: []string{
			"D0, P[], (!!str)::the dog meows\n",
		},
	},
	{
		description:          "Replace string environment variable in document",
		environmentVariables: map[string]string{"myenv": "cat meow"},
		document:             "{v: \"${myenv}\"}",
		expression:           `.v |= envsubst`,
		expected: []string{
			"D0, P[], (doc)::{v: \"cat meow\"}\n",
		},
	},
}

func TestEnvOperatorScenarios(t *testing.T) {
	for _, tt := range envOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "env-variable-operators", envOperatorScenarios)
}
