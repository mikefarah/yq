package yqlib

import (
	"testing"
)

var kindOperatorScenarios = []expressionScenario{
	{
		description: "Get kind",
		document:    `{a: cat, b: 5, c: 3.2, e: true, f: [], g: {}, h: null}`,
		expression:  `.. | kind`,
		expected: []string{
			"D0, P[], (!!str)::map\n",
			"D0, P[a], (!!str)::scalar\n",
			"D0, P[b], (!!str)::scalar\n",
			"D0, P[c], (!!str)::scalar\n",
			"D0, P[e], (!!str)::scalar\n",
			"D0, P[f], (!!str)::seq\n",
			"D0, P[g], (!!str)::map\n",
			"D0, P[h], (!!str)::scalar\n",
		},
	},
	{
		description:    "Get kind, ignores custom tags",
		subdescription: "Unlike tag, kind is not affected by custom tags.",
		document:       `{a: !!thing cat, b: !!foo {}, c: !!bar []}`,
		expression:     `.. | kind`,
		expected: []string{
			"D0, P[], (!!str)::map\n",
			"D0, P[a], (!!str)::scalar\n",
			"D0, P[b], (!!str)::map\n",
			"D0, P[c], (!!str)::seq\n",
		},
	},
	{
		description:    "Add comments only to scalars",
		subdescription: "An example of how you can use kind",
		document:       "a:\n  b: 5\n  c: 3.2\ne: true\nf: []\ng: {}\nh: null",
		expression:     `(.. | select(kind == "scalar")) line_comment = "this is a scalar"`,
		expected:       []string{"D0, P[], (!!map)::a:\n    b: 5 # this is a scalar\n    c: 3.2 # this is a scalar\ne: true # this is a scalar\nf: []\ng: {}\nh: null # this is a scalar\n"},
	},
}

func TestKindOperatorScenarios(t *testing.T) {
	for _, tt := range kindOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "kind", kindOperatorScenarios)
}
