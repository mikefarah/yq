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
}

func TestSliceOperatorScenarios(t *testing.T) {
	for _, tt := range sliceArrayScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "slice-array", sliceArrayScenarios)
}
