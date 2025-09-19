package yqlib

import "testing"

var firstOperatorScenarios = []expressionScenario{
	{
		description: "First matching element from array",
		document:    "[{a: banana},{a: cat},{a: apple}]",
		expression:  `first(.a == "cat")`,
		expected: []string{
			"D0, P[1], (!!map)::{a: cat}\n",
		},
	},
	{
		description: "First matching element from array with multiple matches",
		document:    "[{a: banana},{a: cat},{a: apple},{a: cat}]",
		expression:  `first(.a == "cat")`,
		expected: []string{
			"D0, P[1], (!!map)::{a: cat}\n",
		},
	},
	{
		description: "First matching element from array with numeric condition",
		document:    "[{a: 10},{a: 100},{a: 1}]",
		expression:  `first(.a > 50)`,
		expected: []string{
			"D0, P[1], (!!map)::{a: 100}\n",
		},
	},
	{
		description: "First matching element from array with boolean condition",
		document:    "[{a: false},{a: true},{a: false}]",
		expression:  `first(.a == true)`,
		expected: []string{
			"D0, P[1], (!!map)::{a: true}\n",
		},
	},
	{
		description: "First matching element from array with null values",
		document:    "[{a: null},{a: cat},{a: apple}]",
		expression:  `first(.a != null)`,
		expected: []string{
			"D0, P[1], (!!map)::{a: cat}\n",
		},
	},
	{
		description: "First matching element from array with complex condition",
		document:    "[{a: dog, b: 5},{a: cat, b: 3},{a: apple, b: 7}]",
		expression:  `first(.b > 4)`,
		expected: []string{
			"D0, P[0], (!!map)::{a: dog, b: 5}\n",
		},
	},
	{
		description: "First matching element from map",
		document:    "x: {a: banana}\ny: {a: cat}\nz: {a: apple}",
		expression:  `first(.a == "cat")`,
		expected: []string{
			"D0, P[y], (!!map)::{a: cat}\n",
		},
	},
	{
		description: "First matching element from map with numeric condition",
		document:    "x: {a: 10}\ny: {a: 100}\nz: {a: 1}",
		expression:  `first(.a > 50)`,
		expected: []string{
			"D0, P[y], (!!map)::{a: 100}\n",
		},
	},
	{
		description: "First matching element from nested structure",
		document:    "items: [{a: banana},{a: cat},{a: apple}]",
		expression:  `.items | first(.a == "cat")`,
		expected: []string{
			"D0, P[items 1], (!!map)::{a: cat}\n",
		},
	},
	{
		description: "First matching element with no matches",
		document:    "[{a: banana},{a: cat},{a: apple}]",
		expression:  `first(.a == "dog")`,
		expected:    []string{
			// No output expected when no matches
		},
	},
	{
		description: "First matching element from empty array",
		document:    "[]",
		expression:  `first(.a == "cat")`,
		expected:    []string{
			// No output expected when array is empty
		},
	},
	{
		description: "First matching element from scalar node",
		document:    "hello",
		expression:  `first(. == "hello")`,
		expected:    []string{
			// No output expected when node is scalar (no content to splat)
		},
	},
	{
		description: "First matching element from null node",
		document:    "null",
		expression:  `first(. == "hello")`,
		expected:    []string{
			// No output expected when node is null (no content to splat)
		},
	},
	{
		description: "First matching element with string condition",
		document:    "[{a: banana},{a: cat},{a: apple}]",
		expression:  `first(.a | test("^c"))`,
		expected: []string{
			"D0, P[1], (!!map)::{a: cat}\n",
		},
	},
	{
		description: "First matching element with length condition",
		document:    "[{a: hi},{a: hello},{a: world}]",
		expression:  `first(.a | length > 4)`,
		expected: []string{
			"D0, P[1], (!!map)::{a: hello}\n",
		},
	},
	{
		description: "First matching element from array of strings",
		document:    "[banana, cat, apple]",
		expression:  `first(. == "cat")`,
		expected: []string{
			"D0, P[1], (!!str)::cat\n",
		},
	},
	{
		description: "First matching element from array of numbers",
		document:    "[10, 100, 1]",
		expression:  `first(. > 50)`,
		expected: []string{
			"D0, P[1], (!!int)::100\n",
		},
	},
	{
		description: "First element with no filter from array",
		document:    "[10, 100, 1]",
		expression:  `first`,
		expected: []string{
			"D0, P[0], (!!int)::10\n",
		},
	},
	{
		description: "First element with no filter from array of maps",
		document:    "[{a: 10},{a: 100}]",
		expression:  `first`,
		expected: []string{
			"D0, P[0], (!!map)::{a: 10}\n",
		},
	},
	{
		description: "No filter on empty array returns nothing",
		skipDoc:     true,
		document:    "[]",
		expression:  `first`,
		expected:    []string{},
	},
	{
		description: "No filter on scalar returns nothing",
		skipDoc:     true,
		document:    "hello",
		expression:  `first`,
		expected:    []string{},
	},
	{
		description: "No filter on null returns nothing",
		skipDoc:     true,
		document:    "null",
		expression:  `first`,
		expected:    []string{},
	},
}

func TestFirstOperatorScenarios(t *testing.T) {
	for _, tt := range firstOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "first", firstOperatorScenarios)
}
