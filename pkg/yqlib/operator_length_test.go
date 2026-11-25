package yqlib

import (
	"testing"
)

var lengthOperatorScenarios = []expressionScenario{
	{
		description:    "String length",
		subdescription: "returns length of string",
		document:       `{a: cat}`,
		expression:     `.a | length`,
		expected: []string{
			"D0, P[a], (!!int)::3\n",
		},
	},
	{
		description:    "Unicode string length",
		subdescription: "counts characters (Unicode code points), not bytes",
		document:       `{a: héllo}`,
		expression:     `.a | length`,
		expected: []string{
			"D0, P[a], (!!int)::5\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: "日本語"}`,
		expression: `.a | length`,
		expected: []string{
			"D0, P[a], (!!int)::3\n",
		},
	},
	{
		skipDoc:     true,
		description: "length agrees with string slicing",
		document:    `{a: héllo}`,
		expression:  `.a | length as $n | .[$n - 1:]`,
		expected: []string{
			"D0, P[a], (!!str)::o\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: 1234}`,
		expression: `.a | length`,
		expected: []string{
			"D0, P[a], (!!int)::4\n",
		},
	},
	{
		description: "null length",
		document:    `{a: null}`,
		expression:  `.a | length`,
		expected: []string{
			"D0, P[a], (!!int)::0\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: ~}`,
		expression: `.a | length`,
		expected: []string{
			"D0, P[a], (!!int)::0\n",
		},
	},
	{
		skipDoc:    true,
		document:   "# abc\n{a: key no exist}",
		expression: `.b | length`,
		expected: []string{
			"D0, P[b], (!!int)::0\n",
		},
	},
	{
		description:    "Map length",
		subdescription: "returns number of entries",
		document:       `{a: cat, c: dog}`,
		expression:     `length`,
		expected: []string{
			"D0, P[], (!!int)::2\n",
		},
	},
	{
		description:    "Array length",
		subdescription: "returns number of elements",
		document:       `[2,4,6,8]`,
		expression:     `length`,
		expected: []string{
			"D0, P[], (!!int)::4\n",
		},
	},
}

func TestLengthOperatorScenarios(t *testing.T) {
	for _, tt := range lengthOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "length", lengthOperatorScenarios)
}
