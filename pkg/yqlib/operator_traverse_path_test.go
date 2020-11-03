package yqlib

import (
	"testing"
)

var mergeDocSample = `
foo: &foo
  a: foo_a
  thing: foo_thing
  c: foo_c

bar: &bar
  b: bar_b
  thing: bar_thing
  c: bar_c

foobarList:
  b: foobarList_b
  <<: [*foo,*bar]
  c: foobarList_c

foobar:
  c: foobar_c
  <<: *foo
  thing: foobar_thing
`

var traversePathOperatorScenarios = []expressionScenario{
	{
		document:   `{a: {b: apple}}`,
		expression: `.a`,
		expected: []string{
			"D0, P[a], (!!map)::{b: apple}\n",
		},
	},
	{
		document:   `[{b: apple}, {c: banana}]`,
		expression: `.[]`,
		expected: []string{
			"D0, P[0], (!!map)::{b: apple}\n",
			"D0, P[1], (!!map)::{c: banana}\n",
		},
	},
	{
		document:   `{}`,
		expression: `.a.b`,
		expected: []string{
			"D0, P[a b], (!!null)::null\n",
		},
	},
	{
		document:   `{}`,
		expression: `.[1].a`,
		expected: []string{
			"D0, P[1 a], (!!null)::null\n",
		},
	},
	{
		document:   `{}`,
		expression: `.a.[1]`,
		expected: []string{
			"D0, P[a 1], (!!null)::null\n",
		},
	},
	{
		document:   `{a: {cat: apple, mad: things}}`,
		expression: `.a."*a*"`,
		expected: []string{
			"D0, P[a cat], (!!str)::apple\n",
			"D0, P[a mad], (!!str)::things\n",
		},
	},
	{
		document:   `{a: {cat: {b: 3}, mad: {b: 4}, fad: {c: t}}}`,
		expression: `.a."*a*".b`,
		expected: []string{
			"D0, P[a cat b], (!!int)::3\n",
			"D0, P[a mad b], (!!int)::4\n",
			"D0, P[a fad b], (!!null)::null\n",
		},
	},
	{
		document:   `{a: {cat: apple, mad: things}}`,
		expression: `.a | (.cat, .mad)`,
		expected: []string{
			"D0, P[a cat], (!!str)::apple\n",
			"D0, P[a mad], (!!str)::things\n",
		},
	},
	{
		document:   `{a: {cat: apple, mad: things}}`,
		expression: `.a | (.cat, .mad, .fad)`,
		expected: []string{
			"D0, P[a cat], (!!str)::apple\n",
			"D0, P[a mad], (!!str)::things\n",
			"D0, P[a fad], (!!null)::null\n",
		},
	},
	{
		document:   `{a: {cat: apple, mad: things}}`,
		expression: `.a | (.cat, .mad, .fad) | select( (. == null) | not)`,
		expected: []string{
			"D0, P[a cat], (!!str)::apple\n",
			"D0, P[a mad], (!!str)::things\n",
		},
	},
	{
		document:   `{a: &cat {c: frog}, b: *cat}`,
		expression: `.b`,
		expected: []string{
			"D0, P[b], (alias)::*cat\n",
		},
	},
	{
		document:   `{a: &cat {c: frog}, b: *cat}`,
		expression: `.b.[]`,
		expected: []string{
			"D0, P[b c], (!!str)::frog\n",
		},
	},
	{
		document:   `{a: &cat {c: frog}, b: *cat}`,
		expression: `.b.c`,
		expected: []string{
			"D0, P[b c], (!!str)::frog\n",
		},
	},
	{
		document:   `[1,2,3]`,
		expression: `.b`,
		expected:   []string{},
	},
	{
		document:   `[1,2,3]`,
		expression: `[0]`,
		expected: []string{
			"D0, P[0], (!!int)::1\n",
		},
	},
	{
		description: `Maps can have numbers as keys, so this default to a non-exisiting key behaviour.`,
		document:    `{a: b}`,
		expression:  `[0]`,
		expected: []string{
			"D0, P[0], (!!null)::null\n",
		},
	},
	{
		document:   mergeDocSample,
		expression: `.foobar`,
		expected: []string{
			"D0, P[foobar], (!!map)::c: foobar_c\n!!merge <<: *foo\nthing: foobar_thing\n",
		},
	},
	{
		document:   mergeDocSample,
		expression: `.foobar.a`,
		expected: []string{
			"D0, P[foobar a], (!!str)::foo_a\n",
		},
	},
	{
		document:   mergeDocSample,
		expression: `.foobar.c`,
		expected: []string{
			"D0, P[foobar c], (!!str)::foo_c\n",
		},
	},
	{
		document:   mergeDocSample,
		expression: `.foobar.thing`,
		expected: []string{
			"D0, P[foobar thing], (!!str)::foobar_thing\n",
		},
	},
	{
		document:   mergeDocSample,
		expression: `.foobar.[]`,
		expected: []string{
			"D0, P[foobar c], (!!str)::foo_c\n",
			"D0, P[foobar a], (!!str)::foo_a\n",
			"D0, P[foobar thing], (!!str)::foobar_thing\n",
		},
	},
	{
		document:   mergeDocSample,
		expression: `.foobarList`,
		expected: []string{
			"D0, P[foobarList], (!!map)::b: foobarList_b\n!!merge <<: [*foo, *bar]\nc: foobarList_c\n",
		},
	},
	{
		document:   mergeDocSample,
		expression: `.foobarList.a`,
		expected: []string{
			"D0, P[foobarList a], (!!str)::foo_a\n",
		},
	},
	{
		document:   mergeDocSample,
		expression: `.foobarList.thing`,
		expected: []string{
			"D0, P[foobarList thing], (!!str)::bar_thing\n",
		},
	},
	{
		document:   mergeDocSample,
		expression: `.foobarList.c`,
		expected: []string{
			"D0, P[foobarList c], (!!str)::foobarList_c\n",
		},
	},
	{
		document:   mergeDocSample,
		expression: `.foobarList.b`,
		expected: []string{
			"D0, P[foobarList b], (!!str)::bar_b\n",
		},
	},
	{
		document:   mergeDocSample,
		expression: `.foobarList.[]`,
		expected: []string{
			"D0, P[foobarList b], (!!str)::bar_b\n",
			"D0, P[foobarList a], (!!str)::foo_a\n",
			"D0, P[foobarList thing], (!!str)::bar_thing\n",
			"D0, P[foobarList c], (!!str)::foobarList_c\n",
		},
	},
}

func TestTraversePathOperatorScenarios(t *testing.T) {
	for _, tt := range traversePathOperatorScenarios {
		testScenario(t, &tt)
	}
}
