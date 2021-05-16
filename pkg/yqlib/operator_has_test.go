package yqlib

import (
	"testing"
)

var hasOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		document:   `a: hello`,
		expression: `has("a")`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: hello`,
		expression: `has(.b) as $c`,
		expected: []string{
			"D0, P[], (doc)::a: hello\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: hello`,
		expression: `has(.b)`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		description: "Has map key",
		document: `- a: "yes"
- a: ~
- a: 
- b: nope
`,
		expression: `.[] | has("a")`,
		expected: []string{
			"D0, P[0], (!!bool)::true\n",
			"D0, P[1], (!!bool)::true\n",
			"D0, P[2], (!!bool)::true\n",
			"D0, P[3], (!!bool)::false\n",
		},
	},
	{
		description:    "Select, checking for existence of deep paths",
		subdescription: "Simply pipe in parent expressions into `has`",
		document:       "- {a: {b: {c: cat}}}\n- {a: {b: {d: dog}}}",
		expression:     `.[] | select(.a.b | has("c"))`,
		expected: []string{
			"D0, P[0], (!!map)::{a: {b: {c: cat}}}\n",
		},
	},
	{
		dontFormatInputForDoc: true,
		description:           "Has array index",
		document: `- []
- [1]
- [1, 2]
- [1, null]
- [1, 2, 3]
`,
		expression: `.[] | has(1)`,
		expected: []string{
			"D0, P[0], (!!bool)::false\n",
			"D0, P[1], (!!bool)::false\n",
			"D0, P[2], (!!bool)::true\n",
			"D0, P[3], (!!bool)::true\n",
			"D0, P[4], (!!bool)::true\n",
		},
	},
}

func TestHasOperatorScenarios(t *testing.T) {
	for _, tt := range hasOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Has", hasOperatorScenarios)
}
