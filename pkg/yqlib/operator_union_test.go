package yqlib

import (
	"testing"
)

var unionOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		document:   "{}",
		expression: `(.a, .b.c) as $x | .`,
		expected: []string{
			"D0, P[], (!!map)::{}\n",
		},
	},
	{
		skipDoc:     true,
		description: "clone test",
		expression:  `"abc" as $a | [$a, "cat"]`,
		expected: []string{
			"D0, P[], (!!seq)::- abc\n- cat\n",
		},
	},
	{
		skipDoc:    true,
		expression: `(.foo = "bar"), (.toe = "jam")`,
		expected: []string{
			"D0, P[], ()::foo: bar\ntoe: jam\n",
		},
	},
	{
		description: "Combine scalars",
		expression:  `1, true, "cat"`,
		expected: []string{
			"D0, P[], (!!int)::1\n",
			"D0, P[], (!!bool)::true\n",
			"D0, P[], (!!str)::cat\n",
		},
	},
	{
		description: "Combine selected paths",
		document:    `{a: fieldA, b: fieldB, c: fieldC}`,
		expression:  `.a, .c`,
		expected: []string{
			"D0, P[a], (!!str)::fieldA\n",
			"D0, P[c], (!!str)::fieldC\n",
		},
	},
}

func TestUnionOperatorScenarios(t *testing.T) {
	for _, tt := range unionOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "union", unionOperatorScenarios)
}
