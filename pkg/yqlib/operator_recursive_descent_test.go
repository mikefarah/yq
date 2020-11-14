package yqlib

import (
	"testing"
)

var recursiveDescentOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		document:   `cat`,
		expression: `..`,
		expected: []string{
			"D0, P[], (!!str)::cat\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: frog}`,
		expression: `..`,
		expected: []string{
			"D0, P[], (!!map)::{a: frog}\n",
			"D0, P[a], (!!str)::frog\n",
		},
	},
	{
		description: "Map",
		document:    `{a: {b: apple}}`,
		expression:  `..`,
		expected: []string{
			"D0, P[], (!!map)::{a: {b: apple}}\n",
			"D0, P[a], (!!map)::{b: apple}\n",
			"D0, P[a b], (!!str)::apple\n",
		},
	},
	{
		description: "Array",
		document:    `[1,2,3]`,
		expression:  `..`,
		expected: []string{
			"D0, P[], (!!seq)::[1, 2, 3]\n",
			"D0, P[0], (!!int)::1\n",
			"D0, P[1], (!!int)::2\n",
			"D0, P[2], (!!int)::3\n",
		},
	},
	{
		description: "Array of maps",
		document:    `[{a: cat},2,true]`,
		expression:  `..`,
		expected: []string{
			"D0, P[], (!!seq)::[{a: cat}, 2, true]\n",
			"D0, P[0], (!!map)::{a: cat}\n",
			"D0, P[0 a], (!!str)::cat\n",
			"D0, P[1], (!!int)::2\n",
			"D0, P[2], (!!bool)::true\n",
		},
	},
	{
		description: "Aliases are not traversed",
		document:    `{a: &cat {c: frog}, b: *cat}`,
		expression:  `..`,
		expected: []string{
			"D0, P[], (!!map)::{a: &cat {c: frog}, b: *cat}\n",
			"D0, P[a], (!!map)::&cat {c: frog}\n",
			"D0, P[a c], (!!str)::frog\n",
			"D0, P[b], (alias)::*cat\n",
		},
	},
	{
		description: "Merge docs are not traversed",
		document:    mergeDocSample,
		expression:  `.foobar | ..`,
		expected: []string{
			"D0, P[foobar], (!!map)::c: foobar_c\n!!merge <<: *foo\nthing: foobar_thing\n",
			"D0, P[foobar c], (!!str)::foobar_c\n",
			"D0, P[foobar <<], (alias)::*foo\n",
			"D0, P[foobar thing], (!!str)::foobar_thing\n",
		},
	},
	{
		skipDoc:    true,
		document:   mergeDocSample,
		expression: `.foobarList | ..`,
		expected: []string{
			"D0, P[foobarList], (!!map)::b: foobarList_b\n!!merge <<: [*foo, *bar]\nc: foobarList_c\n",
			"D0, P[foobarList b], (!!str)::foobarList_b\n",
			"D0, P[foobarList <<], (!!seq)::[*foo, *bar]\n",
			"D0, P[foobarList << 0], (alias)::*foo\n",
			"D0, P[foobarList << 1], (alias)::*bar\n",
			"D0, P[foobarList c], (!!str)::foobarList_c\n",
		},
	},
}

func TestRecursiveDescentOperatorScenarios(t *testing.T) {
	for _, tt := range recursiveDescentOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Recursive Descent Operator", recursiveDescentOperatorScenarios)
}
