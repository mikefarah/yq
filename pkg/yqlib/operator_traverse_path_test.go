package yqlib

import (
	"testing"
)

var mergeDocSample = `foo: &foo
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

var fixedTraversePathOperatorScenarios = []expressionScenario{
	{
		description:    "Traversing merge anchor lists",
		subdescription: "Note that the keys earlier in the merge anchors sequence override later ones",
		document:       mergeDocSample,
		expression:     `.foobarList.thing`,
		expected: []string{
			"D0, P[foo thing], (!!str)::foo_thing\n",
		},
	},
	{
		description: "Traversing merge anchors with override",
		document:    mergeDocSample,
		expression:  `.foobar.c`,
		expected: []string{
			"D0, P[foobar c], (!!str)::foobar_c\n",
		},
	},

	// The following tests are the same as below, to verify they still works correctly with the flag:
	{
		skipDoc:        true,
		description:    "Duplicate keys",
		subdescription: "outside merge anchor",
		document:       `{a: 1, a: 2}`,
		expression:     `.a`,
		expected: []string{
			"D0, P[a], (!!int)::2\n",
		},
	},
}

var traversePathOperatorScenarios = []expressionScenario{
	{
		skipDoc:     true,
		description: "strange map with key but no value",
		document:    "!!null\n-",
		expression:  ".x",
		expected: []string{
			"D0, P[x], (!!null)::null\n",
		},
		skipForGoccy: true, // throws an error instead, that's fine
	},
	{
		skipDoc:     true,
		description: "access merge anchors",
		document:    "foo: &foo {x: y}\nbar:\n  <<: *foo\n",
		expression:  `.bar["<<"] | alias`,
		expected: []string{
			"D0, P[bar <<], (!!str)::foo\n",
		},
	},
	{
		skipDoc:     true,
		description: "dynamically set parent and key",
		expression:  `.a.b.c = 3 | .a.b.c`,
		expected: []string{
			"D0, P[a b c], (!!int)::3\n",
		},
	},
	{
		skipDoc:     true,
		description: "dynamically set parent and key in array",
		expression:  `.a.b[0] = 3 | .a.b[0]`,
		expected: []string{
			"D0, P[a b 0], (!!int)::3\n",
		},
	},
	{
		skipDoc:     true,
		description: "dynamically set parent and key",
		expression:  `.a.b = ["x","y"] | .a.b[1]`,
		expected: []string{
			"D0, P[a b 1], (!!str)::y\n",
		},
	},
	{
		skipDoc:     true,
		description: "splat empty map",
		document:    "{}",
		expression:  ".[]",
		expected:    []string{},
	},
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
		expression: `.cat["12"] = "things"`,
		expected: []string{
			"D0, P[], ()::cat:\n    \"12\": things\n",
		},
	},
	{
		skipDoc:    true,
		document:   `blah: {}`,
		expression: `.blah.cat = "cool"`,
		expected: []string{
			"D0, P[], (!!map)::blah:\n    cat: cool\n",
		},
	},
	{
		skipDoc:    true,
		document:   `blah: []`,
		expression: `.blah.0 = "cool"`,
		expected: []string{
			"D0, P[], (!!map)::blah:\n    - cool\n",
		},
	},
	{
		skipDoc:    true,
		document:   `b: cat`,
		expression: ".b\n",
		expected: []string{
			"D0, P[b], (!!str)::cat\n",
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
		skipDoc:    true,
		expression: `.["cat"] = "thing"`,
		expected: []string{
			"D0, P[], ()::cat: thing\n",
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
		subdescription: "Use quotes with square brackets around path elements with special characters",
		document:       `{"{}": frog}`,
		expression:     `.["{}"]`,
		expected: []string{
			"D0, P[{}], (!!str)::frog\n",
		},
	},
	{
		description: "Nested special characters",
		document:    `a: {"key.withdots": {"another.key": apple}}`,
		expression:  `.a["key.withdots"]["another.key"]`,
		expected: []string{
			"D0, P[a key.withdots another.key], (!!str)::apple\n",
		},
	},
	{
		description:    "Keys with spaces",
		subdescription: "Use quotes with square brackets around path elements with special characters",
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
		expression: `.[.a.b] as $x | .`,
		expected: []string{
			"D0, P[], (!!map)::c: dog\n",
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
			"D0, P[a c], (!!str)::frog\n",
		},
	},
	{
		description: "Traversing aliases explicitly",
		document:    `{a: &cat {c: frog}, b: *cat}`,
		expression:  `.b.c`,
		expected: []string{
			"D0, P[a c], (!!str)::frog\n",
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
		description:           "Traversing nested arrays by index",
		dontFormatInputForDoc: true,
		document:              `[[], [cat]]`,
		expression:            `.[1][0]`,
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
		skipDoc:     true,
		description: "Merge anchor with inline map",
		document:    `{<<: {a: 42}}`,
		expression:  `.a`,
		expected: []string{
			"D0, P[<< a], (!!int)::42\n",
		},
	},
	{
		skipDoc:     true,
		description: "Merge anchor with sequence with inline map",
		document:    `{<<: [{a: 42}]}`,
		expression:  `.a`,
		expected: []string{
			"D0, P[<< 0 a], (!!int)::42\n",
		},
	},
	{
		skipDoc:     true,
		description: "Merge anchor with aliased sequence with inline map",
		document:    `{s: &s [{a: 42}], m: {<<: *s}}`,
		expression:  `.m.a`,
		expected: []string{
			"D0, P[s 0 a], (!!int)::42\n",
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
			"D0, P[foo a], (!!str)::foo_a\n",
		},
	},
	{
		description:    "Traversing merge anchors with override",
		subdescription: "This is legacy behaviour, see --yaml-fix-merge-anchor-to-spec",
		document:       mergeDocSample,
		expression:     `.foobar.c`,
		expected: []string{
			"D0, P[foo c], (!!str)::foo_c\n",
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
			"D0, P[foo c], (!!str)::foo_c\n",
			"D0, P[foo a], (!!str)::foo_a\n",
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
			"D0, P[foo a], (!!str)::foo_a\n",
		},
	},
	{
		description: "Traversing merge anchor lists",
		subdescription: "Note that the later merge anchors override previous, " +
			"but this is legacy behaviour, see --yaml-fix-merge-anchor-to-spec",
		document:   mergeDocSample,
		expression: `.foobarList.thing`,
		expected: []string{
			"D0, P[bar thing], (!!str)::bar_thing\n",
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
			"D0, P[bar b], (!!str)::bar_b\n",
		},
	},
	{
		description:    "Splatting merge anchor lists",
		subdescription: "With legacy override behaviour, see --yaml-fix-merge-anchor-to-spec",
		document:       mergeDocSample,
		expression:     `.foobarList[]`,
		expected: []string{
			"D0, P[bar b], (!!str)::bar_b\n",
			"D0, P[foo a], (!!str)::foo_a\n",
			"D0, P[bar thing], (!!str)::bar_thing\n",
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
			"D0, P[a 2], (!!str)::c\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: [a,b,c]}`,
		expression: `.a[-2]`,
		expected: []string{
			"D0, P[a 1], (!!str)::b\n",
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
	{
		skipDoc:        true,
		description:    "Duplicate keys",
		subdescription: "outside merge anchor",
		document:       `{a: 1, a: 2}`,
		expression:     `.a`,
		expected: []string{
			"D0, P[a], (!!int)::2\n",
		},
	},
	{
		skipDoc:        true,
		description:    "Traversing map with invalid merge anchor should not fail",
		subdescription: "Otherwise code cannot do anything with it",
		document:       `{a: 42, <<: 37}`,
		expression:     `.a`,
		expected: []string{
			"D0, P[a], (!!int)::42\n",
		},
	},
	{
		skipDoc:     true,
		description: "Directly accessing invalid merge anchor should not fail",
		document:    `{<<: 37}`,
		expression:  `.<<`,
		expected: []string{
			"D0, P[<<], (!!int)::37\n",
		},
	},
	{
		skipDoc:     true,
		description: "!!str << should not be treated as merge anchor",
		document:    `{!!str <<: {a: 37}}`,
		expression:  `.a`,
		expected: []string{
			"D0, P[a], (!!null)::null\n",
		},
	},
}

func TestTraversePathOperatorScenarios(t *testing.T) {
	for _, tt := range traversePathOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "traverse-read", traversePathOperatorScenarios)
}

func TestTraversePathOperatorAlignedToSpecScenarios(t *testing.T) {
	ConfiguredYamlPreferences.FixMergeAnchorToSpec = true
	for _, tt := range fixedTraversePathOperatorScenarios {
		testScenario(t, &tt)
	}
	ConfiguredYamlPreferences.FixMergeAnchorToSpec = false
}
