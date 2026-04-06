package yqlib

import "testing"

var sliceArrayScenarios = []expressionScenario{
	{
		description: "Slicing arrays",
		document:    `[cat, dog, frog, cow]`,
		expression:  `.[1:3]`,
		expected: []string{
			"D0, P[], (!!seq)::- dog\n- frog\n",
		},
	},
	{
		description:    "Slicing arrays - without the first number",
		subdescription: "Starts from the start of the array",
		document:       `[cat, dog, frog, cow]`,
		expression:     `.[:2]`,
		expected: []string{
			"D0, P[], (!!seq)::- cat\n- dog\n",
		},
	},
	{
		description:    "Slicing arrays - without the second number",
		subdescription: "Finishes at the end of the array",
		document:       `[cat, dog, frog, cow]`,
		expression:     `.[2:]`,
		expected: []string{
			"D0, P[], (!!seq)::- frog\n- cow\n",
		},
	},
	{
		description: "Slicing arrays - use negative numbers to count backwards from the end",
		document:    `[cat, dog, frog, cow]`,
		expression:  `.[1:-1]`,
		expected: []string{
			"D0, P[], (!!seq)::- dog\n- frog\n",
		},
	},
	{
		description:    "Inserting into the middle of an array",
		subdescription: "using an expression to find the index",
		document:       `[cat, dog, frog, cow]`,
		expression:     `(.[] | select(. == "dog") | key + 1) as $pos | .[0:($pos)] + ["rabbit"] + .[$pos:]`,
		expected: []string{
			"D0, P[], (!!seq)::- cat\n- dog\n- rabbit\n- frog\n- cow\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[[cat, dog, frog, cow], [apple, banana, grape, mango]]`,
		expression: `.[] | .[1:3]`,
		expected: []string{
			"D0, P[0], (!!seq)::- dog\n- frog\n",
			"D0, P[1], (!!seq)::- banana\n- grape\n",
		},
	},
	{
		skipDoc:     true,
		description: "second index beyond array clamps",
		document:    `[cat]`,
		expression:  `.[:3]`,
		expected: []string{
			"D0, P[], (!!seq)::- cat\n",
		},
	},
	{
		skipDoc:     true,
		description: "first index beyond array returns nothing",
		document:    `[cat]`,
		expression:  `.[3:]`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[[cat, dog, frog, cow], [apple, banana, grape, mango]]`,
		expression: `.[] | .[-2:-1]`,
		expected: []string{
			"D0, P[0], (!!seq)::- frog\n",
			"D0, P[1], (!!seq)::- grape\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[cat1, cat2, cat3, cat4, cat5, cat6, cat7, cat8, cat9, cat10, cat11]`,
		expression: `.[10:11]`,
		expected: []string{
			"D0, P[], (!!seq)::- cat11\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[cat1, cat2, cat3, cat4, cat5, cat6, cat7, cat8, cat9, cat10, cat11]`,
		expression: `.[-11:-10]`,
		expected: []string{
			"D0, P[], (!!seq)::- cat1\n",
		},
	},
	{
		// Regression test for https://issues.oss-fuzz.com/issues/438776028
		// Negative second index that underflows after adjustment must
		// clamp to zero, yielding an empty sequence.
		skipDoc:    true,
		document:   `[a, b, c]`,
		expression: `.[0:-99999]`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n",
		},
	},
	{
		// First-index underflow: without clamping, the loop starts at a
		// negative index and panics on Content access.
		skipDoc:    true,
		document:   `[a, b, c]`,
		expression: `.[-99999:3]`,
		expected: []string{
			"D0, P[], (!!seq)::- a\n- b\n- c\n",
		},
	},
	{
		// Both indices underflow: both clamp to zero, yielding an empty
		// sequence.
		skipDoc:    true,
		document:   `[a, b, c]`,
		expression: `.[-99999:-99998]`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n",
		},
	},
	{
		description: "Slicing strings",
		document:    `country: Australia`,
		expression:  `.country[0:5]`,
		expected: []string{
			"D0, P[country], (!!str)::Austr\n",
		},
	},
	{
		description:    "Slicing strings - without the second number",
		subdescription: "Finishes at the end of the string",
		document:       `country: Australia`,
		expression:     `.country[5:]`,
		expected: []string{
			"D0, P[country], (!!str)::alia\n",
		},
	},
	{
		description:    "Slicing strings - without the first number",
		subdescription: "Starts from the start of the string",
		document:       `country: Australia`,
		expression:     `.country[:5]`,
		expected: []string{
			"D0, P[country], (!!str)::Austr\n",
		},
	},
	{
		description:    "Slicing strings - use negative numbers to count backwards from the end",
		subdescription: "Negative indices count from the end of the string",
		document:       `country: Australia`,
		expression:     `.country[-5:]`,
		expected: []string{
			"D0, P[country], (!!str)::ralia\n",
		},
	},
	{
		skipDoc:    true,
		document:   `country: Australia`,
		expression: `.country[1:-1]`,
		expected: []string{
			"D0, P[country], (!!str)::ustrali\n",
		},
	},
	{
		skipDoc:    true,
		document:   `country: Australia`,
		expression: `.country[:]`,
		expected: []string{
			"D0, P[country], (!!str)::Australia\n",
		},
	},
	{
		skipDoc:     true,
		description: "second index beyond string length clamps",
		document:    `country: Australia`,
		expression:  `.country[:100]`,
		expected: []string{
			"D0, P[country], (!!str)::Australia\n",
		},
	},
	{
		skipDoc:     true,
		description: "first index beyond string length returns empty string",
		document:    `country: Australia`,
		expression:  `.country[100:]`,
		expected: []string{
			"D0, P[country], (!!str)::\n",
		},
	},
	{
		description:    "Slicing strings - Unicode",
		subdescription: "Indices are rune-based, so multi-byte characters are handled correctly",
		document:       `greeting: héllo`,
		expression:     `.greeting[1:3]`,
		expected: []string{
			"D0, P[greeting], (!!str)::él\n",
		},
	},
}

func TestSliceOperatorScenarios(t *testing.T) {
	for _, tt := range sliceArrayScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "slice-array", sliceArrayScenarios)
}
