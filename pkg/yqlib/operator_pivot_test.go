package yqlib

import "testing"

var pivotOperatorScenarios = []expressionScenario{
	{
		description: "Pivot a sequence of sequences",
		document:    "[[foo, bar, baz], [sis, boom, bah]]\n",
		expression:  `pivot`,
		expected: []string{
			"D0, P[], ()::- - foo\n  - sis\n- - bar\n  - boom\n- - baz\n  - bah\n",
		},
	},
	{
		description:    "Pivot sequence of heterogeneous sequences",
		subdescription: `Missing values are "padded" to null.`,
		document:       "[[foo, bar, baz], [sis, boom, bah, blah]]\n",
		expression:     `pivot`,
		expected: []string{
			"D0, P[], ()::- - foo\n  - sis\n- - bar\n  - boom\n- - baz\n  - bah\n- -\n  - blah\n",
		},
	},
	{
		description: "Pivot sequence of maps",
		document:    "[{foo: a, bar: b, baz: c}, {foo: x, bar: y, baz: z}]\n",
		expression:  `pivot`,
		expected: []string{
			"D0, P[], ()::foo:\n    - a\n    - x\nbar:\n    - b\n    - y\nbaz:\n    - c\n    - z\n",
		},
	},
	{
		description:    "Pivot sequence of heterogeneous maps",
		subdescription: `Missing values are "padded" to null.`,
		document:       "[{foo: a, bar: b, baz: c}, {foo: x, bar: y, baz: z, what: ever}]\n",
		expression:     `pivot`,
		expected: []string{
			"D0, P[], ()::foo:\n    - a\n    - x\nbar:\n    - b\n    - y\nbaz:\n    - c\n    - z\nwhat:\n    -\n    - ever\n",
		},
	},
}

func TestPivotOperatorScenarios(t *testing.T) {
	for _, tt := range pivotOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "pivot", pivotOperatorScenarios)
}
