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
		skipDoc:    true,
		document:   `[[1]]`,
		expression: `.[0][0]`,
		expected: []string{
			"D0, P[0 0], (!!int)::1\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[[[1]]]`,
		expression: `.[0][0][0]`,
		expected: []string{
			"D0, P[0 0 0], (!!int)::1\n",
		},
	},
	{
		description: "Simple map navigation",
		document:    `{a: {b: apple}}`,
		expression:  `.a`,
		expected: []string{
			"D0, P[a], (!!map)::{b: apple}\n",
		},
	},
	{
		description:    "Splat",
		subdescription: "Often used to pipe children into other operators",
		document:       `[{b: apple}, {c: banana}]`,
		expression:     `.[]`,
		expected: []string{
			"D0, P[0], (!!map)::{b: apple}\n",
			"D0, P[1], (!!map)::{c: banana}\n",
		},
	},
	{
		description:    "Optional Splat",
		subdescription: "Just like splat, but won't error if you run it against scalars",
		document:       `"cat"`,
		expression:     `.[]`,
		expected:       []string{},
	},
	{
		description:    "Special characters",
		subdescription: "Use quotes with brackets around path elements with special characters",
		document:       `{"{}": frog}`,
		expression:     `.["{}"]`,
		expected: []string{
			"D0, P[{}], (!!str)::frog\n",
		},
	},
	{
		description: "Multiple special characters",
		document:    `a: {"key.withdots": {"another.key": apple}}`,
		expression:  `.a["key.withdots"]["another.key"]`,
		expected: []string{
			"D0, P[a key.withdots another.key], (!!str)::apple\n",
		},
	},
	{
		description:    "Keys with spaces",
		subdescription: "Use quotes with brackets around path elements with special characters",
		document:       `{"red rabbit": frog}`,
		expression:     `.["red rabbit"]`,
		expected: []string{
			"D0, P[red rabbit], (!!str)::frog\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{"flying fox": frog}`,
		expression: `.["flying fox"]`,
		expected: []string{
			"D0, P[flying fox], (!!str)::frog\n",
		},
	},
	{
		skipDoc:    true,
		document:   `c: dog`,
		expression: `.[.a.b] as $x`,
		expected: []string{
			"D0, P[], (doc)::c: dog\n",
		},
	},
	{
		description:    "Dynamic keys",
		subdescription: `Expressions within [] can be used to dynamically lookup / calculate keys`,
		document:       `{b: apple, apple: crispy yum, banana: soft yum}`,
		expression:     `.[.b]`,
		expected: []string{
			"D0, P[apple], (!!str)::crispy yum\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{b: apple, fruit: {apple: yum, banana: smooth}}`,
		expression: `.fruit[.b]`,
		expected: []string{
			"D0, P[fruit apple], (!!str)::yum\n",
		},
	},
	{
		description:    "Children don't exist",
		subdescription: "Nodes are added dynamically while traversing",
		document:       `{c: banana}`,
		expression:     `.a.b`,
		expected: []string{
			"D0, P[a b], (!!null)::null\n",
		},
	},
	{
		description:    "Optional identifier",
		subdescription: "Like jq, does not output an error when the yaml is not an array or object as expected",
		document:       `[1,2,3]`,
		expression:     `.a?`,
		expected:       []string{},
	},
	{
		skipDoc:    true,
		document:   `[[1,2,3], {a: frog}]`,
		expression: `.[] | .["a"]?`,
		expected:   []string{"D0, P[1 a], (!!str)::frog\n"},
	},
	{
		skipDoc:    true,
		document:   ``,
		expression: `.[1].a`,
		expected: []string{
			"D0, P[1 a], (!!null)::null\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{}`,
		expression: `.a[1]`,
		expected: []string{
			"D0, P[a 1], (!!null)::null\n",
		},
	},
	{
		description: "Wildcard matching",
		document:    `{a: {cat: apple, mad: things}}`,
		expression:  `.a."*a*"`,
		expected: []string{
			"D0, P[a cat], (!!str)::apple\n",
			"D0, P[a mad], (!!str)::things\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {cat: {b: 3}, mad: {b: 4}, fad: {c: t}}}`,
		expression: `.a."*a*".b`,
		expected: []string{
			"D0, P[a cat b], (!!int)::3\n",
			"D0, P[a mad b], (!!int)::4\n",
			"D0, P[a fad b], (!!null)::null\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {cat: apple, mad: things}}`,
		expression: `.a | (.cat, .mad)`,
		expected: []string{
			"D0, P[a cat], (!!str)::apple\n",
			"D0, P[a mad], (!!str)::things\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {cat: apple, mad: things}}`,
		expression: `.a | (.cat, .mad, .fad)`,
		expected: []string{
			"D0, P[a cat], (!!str)::apple\n",
			"D0, P[a mad], (!!str)::things\n",
			"D0, P[a fad], (!!null)::null\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {cat: apple, mad: things}}`,
		expression: `.a | (.cat, .mad, .fad) | select( (. == null) | not)`,
		expected: []string{
			"D0, P[a cat], (!!str)::apple\n",
			"D0, P[a mad], (!!str)::things\n",
		},
	},
	{
		description: "Aliases",
		document:    `{a: &cat {c: frog}, b: *cat}`,
		expression:  `.b`,
		expected: []string{
			"D0, P[b], (alias)::*cat\n",
		},
	},
	{
		description: "Traversing aliases with splat",
		document:    `{a: &cat {c: frog}, b: *cat}`,
		expression:  `.b[]`,
		expected: []string{
			"D0, P[b c], (!!str)::frog\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: &cat {c: frog}, b: *cat}`,
		expression: `.b[]`,
		expected: []string{
			"D0, P[b c], (!!str)::frog\n",
		},
	},
	{
		description: "Traversing aliases explicitly",
		document:    `{a: &cat {c: frog}, b: *cat}`,
		expression:  `.b.c`,
		expected: []string{
			"D0, P[b c], (!!str)::frog\n",
		},
	},
	{
		description: "Traversing arrays by index",
		document:    `[1,2,3]`,
		expression:  `.[0]`,
		expected: []string{
			"D0, P[0], (!!int)::1\n",
		},
	},
	{
		description: "Traversing nested arrays by index",
		document:    `[[], [cat]]`,
		expression:  `.[1][0]`,
		expected: []string{
			"D0, P[1 0], (!!str)::cat\n",
		},
	},
	{
		description: "Maps with numeric keys",
		document:    `{2: cat}`,
		expression:  `.[2]`,
		expected: []string{
			"D0, P[2], (!!str)::cat\n",
		},
	},
	{
		description: "Maps with non existing numeric keys",
		document:    `{a: b}`,
		expression:  `.[0]`,
		expected: []string{
			"D0, P[0], (!!null)::null\n",
		},
	},
	{
		skipDoc:    true,
		document:   mergeDocSample,
		expression: `.foobar`,
		expected: []string{
			"D0, P[foobar], (!!map)::c: foobar_c\n!!merge <<: *foo\nthing: foobar_thing\n",
		},
	},
	{
		description: "Traversing merge anchors",
		document:    mergeDocSample,
		expression:  `.foobar.a`,
		expected: []string{
			"D0, P[foobar a], (!!str)::foo_a\n",
		},
	},
	{
		description: "Traversing merge anchors with override",
		document:    mergeDocSample,
		expression:  `.foobar.c`,
		expected: []string{
			"D0, P[foobar c], (!!str)::foo_c\n",
		},
	},
	{
		description: "Traversing merge anchors with local override",
		document:    mergeDocSample,
		expression:  `.foobar.thing`,
		expected: []string{
			"D0, P[foobar thing], (!!str)::foobar_thing\n",
		},
	},
	{
		description: "Splatting merge anchors",
		document:    mergeDocSample,
		expression:  `.foobar[]`,
		expected: []string{
			"D0, P[foobar c], (!!str)::foo_c\n",
			"D0, P[foobar a], (!!str)::foo_a\n",
			"D0, P[foobar thing], (!!str)::foobar_thing\n",
		},
	},
	{
		skipDoc:    true,
		document:   mergeDocSample,
		expression: `.foobar[]`,
		expected: []string{
			"D0, P[foobar c], (!!str)::foo_c\n",
			"D0, P[foobar a], (!!str)::foo_a\n",
			"D0, P[foobar thing], (!!str)::foobar_thing\n",
		},
	},
	{
		skipDoc:    true,
		document:   mergeDocSample,
		expression: `.foobarList`,
		expected: []string{
			"D0, P[foobarList], (!!map)::b: foobarList_b\n!!merge <<: [*foo, *bar]\nc: foobarList_c\n",
		},
	},
	{
		skipDoc:    true,
		document:   mergeDocSample,
		expression: `.foobarList.a`,
		expected: []string{
			"D0, P[foobarList a], (!!str)::foo_a\n",
		},
	},
	{
		description:    "Traversing merge anchor lists",
		subdescription: "Note that the later merge anchors override previous",
		document:       mergeDocSample,
		expression:     `.foobarList.thing`,
		expected: []string{
			"D0, P[foobarList thing], (!!str)::bar_thing\n",
		},
	},
	{
		skipDoc:    true,
		document:   mergeDocSample,
		expression: `.foobarList.c`,
		expected: []string{
			"D0, P[foobarList c], (!!str)::foobarList_c\n",
		},
	},
	{
		skipDoc:    true,
		document:   mergeDocSample,
		expression: `.foobarList.b`,
		expected: []string{
			"D0, P[foobarList b], (!!str)::bar_b\n",
		},
	},
	{
		description: "Splatting merge anchor lists",
		document:    mergeDocSample,
		expression:  `.foobarList[]`,
		expected: []string{
			"D0, P[foobarList b], (!!str)::bar_b\n",
			"D0, P[foobarList a], (!!str)::foo_a\n",
			"D0, P[foobarList thing], (!!str)::bar_thing\n",
			"D0, P[foobarList c], (!!str)::foobarList_c\n",
		},
	},
	{
		skipDoc:    true,
		document:   mergeDocSample,
		expression: `.foobarList[]`,
		expected: []string{
			"D0, P[foobarList b], (!!str)::bar_b\n",
			"D0, P[foobarList a], (!!str)::foo_a\n",
			"D0, P[foobarList thing], (!!str)::bar_thing\n",
			"D0, P[foobarList c], (!!str)::foobarList_c\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[a,b,c]`,
		expression: `.[]`,
		expected: []string{
			"D0, P[0], (!!str)::a\n",
			"D0, P[1], (!!str)::b\n",
			"D0, P[2], (!!str)::c\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[a,b,c]`,
		expression: `[]`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: [a,b,c]}`,
		expression: `.a[0]`,
		expected: []string{
			"D0, P[a 0], (!!str)::a\n",
		},
	},
	{
		description: "Select multiple indices",
		document:    `{a: [a,b,c]}`,
		expression:  `.a[0, 2]`,
		expected: []string{
			"D0, P[a 0], (!!str)::a\n",
			"D0, P[a 2], (!!str)::c\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: [a,b,c]}`,
		expression: `.a[0, 2]`,
		expected: []string{
			"D0, P[a 0], (!!str)::a\n",
			"D0, P[a 2], (!!str)::c\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: [a,b,c]}`,
		expression: `.a[-1]`,
		expected: []string{
			"D0, P[a -1], (!!str)::c\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: [a,b,c]}`,
		expression: `.a[-1]`,
		expected: []string{
			"D0, P[a -1], (!!str)::c\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: [a,b,c]}`,
		expression: `.a[-2]`,
		expected: []string{
			"D0, P[a -2], (!!str)::b\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: [a,b,c]}`,
		expression: `.a[-2]`,
		expected: []string{
			"D0, P[a -2], (!!str)::b\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: [a,b,c]}`,
		expression: `.a[]`,
		expected: []string{
			"D0, P[a 0], (!!str)::a\n",
			"D0, P[a 1], (!!str)::b\n",
			"D0, P[a 2], (!!str)::c\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: [a,b,c]}`,
		expression: `.a[]`,
		expected: []string{
			"D0, P[a 0], (!!str)::a\n",
			"D0, P[a 1], (!!str)::b\n",
			"D0, P[a 2], (!!str)::c\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: [a,b,c]}`,
		expression: `.a | .[]`,
		expected: []string{
			"D0, P[a 0], (!!str)::a\n",
			"D0, P[a 1], (!!str)::b\n",
			"D0, P[a 2], (!!str)::c\n",
		},
	},
}

func TestTraversePathOperatorScenarios(t *testing.T) {
	for _, tt := range traversePathOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Traverse (Read)", traversePathOperatorScenarios)
}
