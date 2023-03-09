package yqlib

import (
	"testing"
)

var divideOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		document:   `[{a: foo_bar, b: _}, {a: 4, b: 2}]`,
		expression: ".[] | .a / .b",
		expected: []string{
			"D0, P[0 a], (!!seq)::- foo\n- bar\n",
			"D0, P[1 a], (!!float)::2\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{}`,
		expression: "(.a / .b) as $x | .",
		expected: []string{
			"D0, P[], (doc)::{}\n",
		},
	},
	{
		description: "String split",
		document:    `{a: cat_meow, b: _}`,
		expression:  `.c = .a / .b`,
		expected: []string{
			"D0, P[], (doc)::{a: cat_meow, b: _, c: [cat, meow]}\n",
		},
	},
	{
		description:    "Number division",
		subdescription: "The result during divison is calculated as a float",
		document:       `{a: 12, b: 2.5}`,
		expression:     `.a = .a / .b`,
		expected: []string{
			"D0, P[], (doc)::{a: 4.8, b: 2.5}\n",
		},
	},
	{
		description:    "Number division by zero",
		subdescription: "Dividing by zero results in +Inf or -Inf",
		document:       `{a: 1, b: -1}`,
		expression:     `.a = .a / 0 | .b = .b / 0`,
		expected: []string{
			"D0, P[], (doc)::{a: !!float +Inf, b: !!float -Inf}\n",
		},
	},
	{
		skipDoc:     true,
		description: "Custom types: that are really strings",
		document:    "a: !horse cat_meow\nb: !goat _",
		expression:  `.a = .a / .b`,
		expected: []string{
			"D0, P[], (doc)::a: !horse\n    - cat\n    - meow\nb: !goat _\n",
		},
	},
	{
		skipDoc:     true,
		description: "Custom types: that are really numbers",
		document:    "a: !horse 1.2\nb: !goat 2.3",
		expression:  `.a = .a / .b`,
		expected: []string{
			"D0, P[], (doc)::a: !horse 0.5217391304347826\nb: !goat 2.3\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 2\nb: !goat 2.3",
		expression: `.a = .a / .b`,
		expected: []string{
			"D0, P[], (doc)::a: 0.8695652173913044\nb: !goat 2.3\n",
		},
	},
	{
		skipDoc:     true,
		description: "Custom types: that are really ints",
		document:    "a: !horse 2\nb: !goat 3",
		expression:  `.a = .a / .b`,
		expected: []string{
			"D0, P[], (doc)::a: !horse 0.6666666666666666\nb: !goat 3\n",
		},
	},
	{
		skipDoc:     true,
		description: "Keep anchors",
		document:    "a: &horse [1]",
		expression:  `.a[1] = .a[0] / 2`,
		expected: []string{
			"D0, P[], (doc)::a: &horse [1, 0.5]\n",
		},
	},
	{
		skipDoc:       true,
		description:   "Divide int by string",
		document:      "a: 123\nb: '2'",
		expression:    `.a / .b`,
		expectedError: "!!int cannot be divided by !!str",
	},
	{
		skipDoc:       true,
		description:   "Divide string by int",
		document:      "a: 2\nb: '123'",
		expression:    `.b / .a`,
		expectedError: "!!str cannot be divided by !!int",
	},
	{
		skipDoc:       true,
		description:   "Divide map by int",
		document:      "a: {\"a\":1}\nb: 2",
		expression:    `.a / .b`,
		expectedError: "!!map (a) cannot be divided by !!int (b)",
	},
	{
		skipDoc:       true,
		description:   "Divide array by str",
		document:      "a: [1,2]\nb: '2'",
		expression:    `.a / .b`,
		expectedError: "!!seq (a) cannot be divided by !!str (b)",
	},
}

func TestDivideOperatorScenarios(t *testing.T) {
	for _, tt := range divideOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "divide", divideOperatorScenarios)
}
