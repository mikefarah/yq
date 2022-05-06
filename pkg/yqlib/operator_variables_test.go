package yqlib

import (
	"testing"
)

var variableOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		document:   `{}`,
		expression: `.a.b as $foo`,
		expected: []string{
			"D0, P[], (doc)::{}\n",
		},
	},
	{
		document:      "a: [cat]",
		skipDoc:       true,
		expression:    "(.[] | {.name: .}) as $item",
		expectedError: `cannot index array with 'name' (strconv.ParseInt: parsing "name": invalid syntax)`,
	},
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
	{
		description: "Using variables to swap values",
		document:    "a: a_value\nb: b_value",
		expression:  `.a as $x  | .b as $y | .b = $x | .a = $y`,
		expected: []string{
			"D0, P[], (doc)::a: b_value\nb: a_value\n",
		},
	},
	{
		description:    "Use ref to reference a path repeatedly",
		subdescription: "Note: You may find the `with` operator more useful.",
		document:       `a: {b: thing, c: something}`,
		expression:     `.a.b ref $x | $x = "new" | $x style="double"`,
		expected: []string{
			"D0, P[], (doc)::a: {b: \"new\", c: something}\n",
		},
	},
}

func TestVariableOperatorScenarios(t *testing.T) {
	for _, tt := range variableOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "variable-operators", variableOperatorScenarios)
}
