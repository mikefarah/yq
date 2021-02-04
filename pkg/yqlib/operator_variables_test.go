package yqlib

import (
	"testing"
)

var variableOperatorScenarios = []expressionScenario{
	{
		description: "Single value variable",
		document:    `a: cat`,
		expression:  `.a as $foo | $foo`,
		expected: []string{
			"D0, P[a], (!!str)::cat\n",
		},
	},
	{
		description: "Multi value variable",
		document:    `[cat, dog]`,
		expression:  `.[] as $foo | $foo`,
		expected: []string{
			"D0, P[0], (!!str)::cat\n",
			"D0, P[1], (!!str)::dog\n",
		},
	},
	{
		description:    "Using variables as a lookup",
		subdescription: "Example taken from [jq](https://stedolan.github.io/jq/manual/#Variable/SymbolicBindingOperator:...as$identifier|...)",
		document: `{"posts": [{"title": "Frist psot", "author": "anon"},
		{"title": "A well-written article", "author": "person1"}],
"realnames": {"anon": "Anonymous Coward",
				"person1": "Person McPherson"}}`,
		expression: `.realnames as $names | .posts[] | {"title":.title, "author": $names[.author]}`,
		expected: []string{
			"D0, P[], (!!map)::title: \"Frist psot\"\nauthor: \"Anonymous Coward\"\n",
			"D0, P[], (!!map)::title: \"A well-written article\"\nauthor: \"Person McPherson\"\n",
		},
	},
}

func TestVariableOperatorScenarios(t *testing.T) {
	for _, tt := range variableOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Variable Operators", variableOperatorScenarios)
}
