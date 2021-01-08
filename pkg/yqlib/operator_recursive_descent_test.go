package yqlib

import (
	"testing"
)

var recursiveDescentOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		document:   `{}`,
		expression: `..`,
		expected: []string{
			"D0, P[], (!!map)::{}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{}`,
		expression: `...`,
		expected: []string{
			"D0, P[], (!!map)::{}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[]`,
		expression: `..`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[]`,
		expression: `...`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n",
		},
	},
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
		document:   `cat`,
		expression: `...`,
		expected: []string{
			"D0, P[], (!!str)::cat\n",
		},
	},
	{
		description: "Recurse map (values only)",
		document:    `{a: frog}`,
		expression:  `..`,
		expected: []string{
			"D0, P[], (!!map)::{a: frog}\n",
			"D0, P[a], (!!str)::frog\n",
		},
	},
	{
		description: "Recursively find nodes with keys",
		document:    `{a: {name: frog, b: {name: blog, age: 12}}}`,
		expression:  `.. | select(has("name"))`,
		expected: []string{
			"D0, P[a], (!!map)::{name: frog, b: {name: blog, age: 12}}\n",
			"D0, P[a b], (!!map)::{name: blog, age: 12}\n",
		},
	},
	{
		description: "Recursively find nodes with values",
		document:    `{a: {nameA: frog, b: {nameB: frog, age: 12}}}`,
		expression:  `.. | select(. == "frog")`,
		expected: []string{
			"D0, P[a nameA], (!!str)::frog\n",
			"D0, P[a b nameB], (!!str)::frog\n",
		},
	},
	{
		description:    "Recurse map (values and keys)",
		subdescription: "Note that the map key appears in the results",
		document:       `{a: frog}`,
		expression:     `...`,
		expected: []string{
			"D0, P[], (!!map)::{a: frog}\n",
			"D0, P[a], (!!str)::a\n",
			"D0, P[a], (!!str)::frog\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {b: apple}}`,
		expression: `..`,
		expected: []string{
			"D0, P[], (!!map)::{a: {b: apple}}\n",
			"D0, P[a], (!!map)::{b: apple}\n",
			"D0, P[a b], (!!str)::apple\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {b: apple}}`,
		expression: `...`,
		expected: []string{
			"D0, P[], (!!map)::{a: {b: apple}}\n",
			"D0, P[a], (!!str)::a\n",
			"D0, P[a], (!!map)::{b: apple}\n",
			"D0, P[a b], (!!str)::b\n",
			"D0, P[a b], (!!str)::apple\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[1,2,3]`,
		expression: `..`,
		expected: []string{
			"D0, P[], (!!seq)::[1, 2, 3]\n",
			"D0, P[0], (!!int)::1\n",
			"D0, P[1], (!!int)::2\n",
			"D0, P[2], (!!int)::3\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[1,2,3]`,
		expression: `...`,
		expected: []string{
			"D0, P[], (!!seq)::[1, 2, 3]\n",
			"D0, P[0], (!!int)::1\n",
			"D0, P[1], (!!int)::2\n",
			"D0, P[2], (!!int)::3\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[{a: cat},2,true]`,
		expression: `..`,
		expected: []string{
			"D0, P[], (!!seq)::[{a: cat}, 2, true]\n",
			"D0, P[0], (!!map)::{a: cat}\n",
			"D0, P[0 a], (!!str)::cat\n",
			"D0, P[1], (!!int)::2\n",
			"D0, P[2], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[{a: cat},2,true]`,
		expression: `...`,
		expected: []string{
			"D0, P[], (!!seq)::[{a: cat}, 2, true]\n",
			"D0, P[0], (!!map)::{a: cat}\n",
			"D0, P[0 a], (!!str)::a\n",
			"D0, P[0 a], (!!str)::cat\n",
			"D0, P[1], (!!int)::2\n",
			"D0, P[2], (!!bool)::true\n",
		},
	},
	{
		description: "Aliases are not traversed",
		document:    `{a: &cat {c: frog}, b: *cat}`,
		expression:  `[..]`,
		expected: []string{
			"D0, P[a], (!!seq)::- {a: &cat {c: frog}, b: *cat}\n- &cat {c: frog}\n- frog\n- *cat\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: &cat {c: frog}, b: *cat}`,
		expression: `...`,
		expected: []string{
			"D0, P[], (!!map)::{a: &cat {c: frog}, b: *cat}\n",
			"D0, P[a], (!!str)::a\n",
			"D0, P[a], (!!map)::&cat {c: frog}\n",
			"D0, P[a c], (!!str)::c\n",
			"D0, P[a c], (!!str)::frog\n",
			"D0, P[b], (!!str)::b\n",
			"D0, P[b], (alias)::*cat\n",
		},
	},
	{
		description: "Merge docs are not traversed",
		document:    mergeDocSample,
		expression:  `.foobar | [..]`,
		expected: []string{
			"D0, P[foobar], (!!seq)::- c: foobar_c\n  !!merge <<: *foo\n  thing: foobar_thing\n- foobar_c\n- *foo\n- foobar_thing\n",
		},
	},
	{
		skipDoc:    true,
		document:   mergeDocSample,
		expression: `.foobar | [...]`,
		expected: []string{
			"D0, P[foobar], (!!seq)::- c: foobar_c\n  !!merge <<: *foo\n  thing: foobar_thing\n- c\n- foobar_c\n- !!merge <<\n- *foo\n- thing\n- foobar_thing\n",
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
	{
		skipDoc:    true,
		document:   mergeDocSample,
		expression: `.foobarList | ...`,
		expected: []string{
			"D0, P[foobarList], (!!map)::b: foobarList_b\n!!merge <<: [*foo, *bar]\nc: foobarList_c\n",
			"D0, P[foobarList b], (!!str)::b\n",
			"D0, P[foobarList b], (!!str)::foobarList_b\n",
			"D0, P[foobarList <<], (!!merge)::!!merge <<\n",
			"D0, P[foobarList <<], (!!seq)::[*foo, *bar]\n",
			"D0, P[foobarList << 0], (alias)::*foo\n",
			"D0, P[foobarList << 1], (alias)::*bar\n",
			"D0, P[foobarList c], (!!str)::c\n",
			"D0, P[foobarList c], (!!str)::foobarList_c\n",
		},
	},
}

func TestRecursiveDescentOperatorScenarios(t *testing.T) {
	for _, tt := range recursiveDescentOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Recursive Descent", recursiveDescentOperatorScenarios)
}
