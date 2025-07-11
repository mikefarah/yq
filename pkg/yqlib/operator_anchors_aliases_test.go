package yqlib

import (
	"testing"
)

var specDocument = `- &CENTER { x: 1, y: 2 }
- &LEFT { x: 0, y: 2 }
- &BIG { r: 10 }
- &SMALL { r: 1 }
`

var expectedSpecResult = "D0, P[4], (!!map)::x: 1\ny: 2\nr: 10\n"

var simpleArrayRef = `item_value: &item_value
  value: true

thingOne:
  name: item_1
  <<: *item_value

thingTwo:
  name: item_2
  <<: *item_value
`

var expectedUpdatedArrayRef = `D0, P[], (!!map)::item_value: &item_value
    value: true
thingOne:
    name: item_1
    value: false
thingTwo:
    name: item_2
    !!merge <<: *item_value
`

var explodeMergeAnchorsExpected = `D0, P[], (!!map)::foo:
    a: foo_a
    thing: foo_thing
    c: foo_c
bar:
    b: bar_b
    thing: bar_thing
    c: bar_c
foobarList:
    b: bar_b
    thing: foo_thing
    c: foobarList_c
    a: foo_a
foobar:
    c: foo_c
    a: foo_a
    thing: foobar_thing
`

var anchorMergeDocSample = `foo: &foo
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

var anchorBadAliasSample = `
_common: &common-docker-file
  - FROM ubuntu:18.04

steps:
  <<: *common-docker-file
`

var anchorOperatorScenarios = []expressionScenario{
	{
		skipDoc:       true,
		description:   "merge anchor not map",
		document:      "a: &a\n  - 0\nc:\n  <<: [*a]\n",
		expectedError: "merge anchor only supports maps, got !!seq instead",
		expression:    "explode(.)",
		skipForGoccy:  true, // goccy yaml parser throws a different error
	},
	{
		skipDoc:       true,
		description:   "merge anchor not map",
		document:      "a: &a\n  - 0\nc:\n  <<: [*a]\n",
		expectedError: "merge anchor only supports maps, got !!seq instead a: &a",
		expression:    "explode(.)",
		skipForYamlV3: true, // yaml.v3 parser throws a different error
	},
	{
		description:    "Merge one map",
		subdescription: "see https://yaml.org/type/merge.html",
		document:       specDocument + "- << : *CENTER\n  r: 10\n",
		expression:     ".[4] | explode(.)",
		expected:       []string{expectedSpecResult},
	},
	{
		description:    "Merge multiple maps",
		subdescription: "see https://yaml.org/type/merge.html",
		document:       specDocument + "- << : [ *CENTER, *BIG ]\n",
		expression:     ".[4] | explode(.)",
		expected:       []string{"D0, P[4], (!!map)::r: 10\nx: 1\ny: 2\n"},
	},
	{
		description:    "Override",
		subdescription: "see https://yaml.org/type/merge.html",
		document:       specDocument + "- << : [ *BIG, *LEFT, *SMALL ]\n  x: 1\n",
		expression:     ".[4] | explode(.)",
		expected:       []string{"D0, P[4], (!!map)::r: 10\nx: 1\ny: 2\n"},
	},
	{
		description: "Get anchor",
		document:    `a: &billyBob cat`,
		expression:  `.a | anchor`,
		expected: []string{
			"D0, P[a], (!!str)::billyBob\n",
		},
	},
	{
		description: "Set anchor",
		document:    `a: cat`,
		expression:  `.a anchor = "foobar"`,
		expected: []string{
			"D0, P[], (!!map)::a: &foobar cat\n",
		},
	},
	{
		description: "Set anchor relatively using assign-update",
		document:    `a: {b: cat}`,
		expression:  `.a anchor |= .b`,
		expected: []string{
			"D0, P[], (!!map)::a: &cat {b: cat}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: {c: cat}`,
		expression: `.a anchor |= .b`,
		expected: []string{
			"D0, P[], (!!map)::a: {c: cat}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: {c: cat}`,
		expression: `.a anchor = .b`,
		expected: []string{
			"D0, P[], (!!map)::a: {c: cat}\n",
		},
	},
	{
		description: "Get alias",
		document:    `{b: &billyBob meow, a: *billyBob}`,
		expression:  `.a | alias`,
		expected: []string{
			"D0, P[a], (!!str)::billyBob\n",
		},
	},
	{
		description: "Set alias",
		document:    `{b: &meow purr, a: cat}`,
		expression:  `.a alias = "meow"`,
		expected: []string{
			"D0, P[], (!!map)::{b: &meow purr, a: *meow}\n",
		},
	},
	{
		description: "Set alias to blank does nothing",
		document:    `{b: &meow purr, a: cat}`,
		expression:  `.a alias = ""`,
		expected: []string{
			"D0, P[], (!!map)::{b: &meow purr, a: cat}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{b: &meow purr, a: cat}`,
		expression: `.a alias = .c`,
		expected: []string{
			"D0, P[], (!!map)::{b: &meow purr, a: cat}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{b: &meow purr, a: cat}`,
		expression: `.a alias |= .c`,
		expected: []string{
			"D0, P[], (!!map)::{b: &meow purr, a: cat}\n",
		},
	},
	{
		description: "Set alias relatively using assign-update",
		document:    `{b: &meow purr, a: {f: meow}}`,
		expression:  `.a alias |= .f`,
		expected: []string{
			"D0, P[], (!!map)::{b: &meow purr, a: *meow}\n",
		},
	},
	{
		description: "Dont explode alias and anchor - check alias parent",
		skipDoc:     true,
		document:    `{a: &a [1], b: *a}`,
		expression:  `.b[]`,
		expected: []string{
			"D0, P[a 0], (!!int)::1\n",
		},
	},
	{
		description: "Explode alias and anchor - check alias parent",
		skipDoc:     true,
		document:    `{a: &a cat, b: *a}`,
		expression:  `explode(.) | .b`,
		expected: []string{
			"D0, P[b], (!!str)::cat\n",
		},
	},
	{
		description: "Explode splat",
		skipDoc:     true,
		document:    `{a: &a cat, b: *a}`,
		expression:  `explode(.)[]`,
		expected: []string{
			"D0, P[a], (!!str)::cat\n",
			"D0, P[b], (!!str)::cat\n",
		},
	},
	{
		description: "Explode alias and anchor - check original parent",
		skipDoc:     true,
		document:    `{a: &a cat, b: *a}`,
		expression:  `explode(.) | .a`,
		expected: []string{
			"D0, P[a], (!!str)::cat\n",
		},
	},
	{
		description: "Explode alias and anchor",
		document:    `{f : {a: &a cat, b: *a}}`,
		expression:  `explode(.f)`,
		expected: []string{
			"D0, P[], (!!map)::{f: {a: cat, b: cat}}\n",
		},
	},
	{
		description: "Explode with no aliases or anchors",
		document:    `a: mike`,
		expression:  `explode(.a)`,
		expected: []string{
			"D0, P[], (!!map)::a: mike\n",
		},
	},
	{
		description:    "Explode with alias keys",
		subdescription: "No space between alias",
		skipDoc:        true,
		document:       `{f : {a: &a cat, *a: b}}`,
		expression:     `explode(.f)`,
		expected: []string{
			"D0, P[], (!!map)::{f: {a: cat, cat: b}}\n",
		},
		skipForGoccy: true, // can't handle no space between alias
	},
	{
		skipDoc:        true,
		description:    "flow map with alias keys - goccy helpful error",
		subdescription: "Goccy provides guidance for unsupported flow map alias key syntax",
		document:       `{f : {a: &a cat, *a: b}}`,
		expression:     `explode(.f)`,
		expectedError:  "flow maps with alias keys are not supported in this parser. Consider using block map syntax instead",
		skipForYamlV3:  true, // yaml.v3 supports this syntax
	},
	{
		description: "Explode with alias keys",
		document:    `{f : {a: &a cat, *a : b}}`,
		expression:  `explode(.f)`,
		expected: []string{
			"D0, P[], (!!map)::{f: {a: cat, cat: b}}\n",
		},
	},
	{
		description: "Explode with merge anchors",
		document:    mergeDocSample,
		expression:  `explode(.)`,
		expected:    []string{explodeMergeAnchorsExpected},
	},
	{
		skipDoc:    true,
		document:   mergeDocSample,
		expression: `.foo* | explode(.) | (. style="flow")`,
		expected: []string{
			"D0, P[foo], (!!map)::{a: foo_a, thing: foo_thing, c: foo_c}\n",
			"D0, P[foobarList], (!!map)::{b: bar_b, thing: foo_thing, c: foobarList_c, a: foo_a}\n",
			"D0, P[foobar], (!!map)::{c: foo_c, a: foo_a, thing: foobar_thing}\n",
		},
	},
	{
		skipDoc:    true,
		document:   mergeDocSample,
		expression: `.foo* | explode(explode(.)) | (. style="flow")`,
		expected: []string{
			"D0, P[foo], (!!map)::{a: foo_a, thing: foo_thing, c: foo_c}\n",
			"D0, P[foobarList], (!!map)::{b: bar_b, thing: foo_thing, c: foobarList_c, a: foo_a}\n",
			"D0, P[foobar], (!!map)::{c: foo_c, a: foo_a, thing: foobar_thing}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{f : {a: &a cat, b: &b {foo: *a}, *a: *b}}`,
		expression: `explode(.f)`,
		expected: []string{
			"D0, P[], (!!map)::{f: {a: cat, b: {foo: cat}, cat: {foo: cat}}}\n",
		},
		skipForGoccy: true, // can't handle no space between alias
	},
	{
		skipDoc:    true,
		document:   `{f : {a: &a cat, b: &b {foo: *a}, *a : *b}}`,
		expression: `explode(.f)`,
		expected: []string{
			"D0, P[], (!!map)::{f: {a: cat, b: {foo: cat}, cat: {foo: cat}}}\n",
		},
	},
	{
		description:    "Dereference and update a field",
		subdescription: "Use explode with multiply to dereference an object",
		document:       simpleArrayRef,
		expression:     `.thingOne |= explode(.) * {"value": false}`,
		expected:       []string{expectedUpdatedArrayRef},
	},
	// Additional merge anchor test scenarios demonstrating parser differences
	{
		description:    "Merge anchor display - legacy-v3 behaviour",
		subdescription: "legacy-v3 preserves merge anchor notation in output",
		document:       anchorMergeDocSample,
		expression:     `.foobar`,
		expected: []string{
			"D0, P[foobar], (!!map)::c: foobar_c\n!!merge <<: *foo\nthing: foobar_thing\n",
		},
		skipForGoccy: true,
	},
	{
		description:    "Merge anchor display - goccy behaviour",
		subdescription: "Goccy expands merge anchors in output display",
		document:       anchorMergeDocSample,
		expression:     `.foobar`,
		expected: []string{
			"D0, P[foobar], (!!map)::<<:\n    a: foo_a\n    c: foo_c\n    thing: foo_thing\nc: foobar_c\nthing: foobar_thing\n",
		},
		skipForYamlV3: true,
	},
	{
		description:    "Merge anchor content access - both parsers",
		subdescription: "Both parsers correctly access merged content despite display differences",
		document:       anchorMergeDocSample,
		expression:     `.foobar.a`,
		expected: []string{
			"D0, P[foo a], (!!str)::foo_a\n",
		},
	},
	{
		description:    "Merge anchor with override - both parsers",
		subdescription: "Both parsers handle merge anchor overrides correctly",
		document:       anchorMergeDocSample,
		expression:     `.foobar.thing`,
		expected: []string{
			"D0, P[foobar thing], (!!str)::foobar_thing\n",
		},
	},
	{
		description:    "Invalid merge anchor with array - legacy-v3 runtime error",
		subdescription: "legacy-v3 parses but gives runtime error on access",
		document:       anchorBadAliasSample,
		expression:     ".steps[]",
		expectedError:  "can only use merge anchors with maps (!!map), but got !!seq",
		skipForGoccy:   true, // Goccy gives parse-time error instead
	},
	{
		description:    "Invalid merge anchor with array - goccy parse error",
		subdescription: "Goccy provides stricter validation with parse-time error",
		document:       anchorBadAliasSample,
		expression:     ".steps[]",
		expectedError:  "bad file '-': [2:5] string was used where mapping is expected",
		skipForYamlV3:  true, // legacy-v3 gives runtime error instead
	},
	{
		description:    "Complex merge anchor list - both parsers",
		subdescription: "Both parsers handle multiple merge anchors with proper precedence",
		document:       anchorMergeDocSample,
		expression:     `.foobarList.thing`,
		expected: []string{
			"D0, P[bar thing], (!!str)::bar_thing\n",
		},
	},
	{
		description:    "Merge anchor list display - legacy-v3 behaviour",
		subdescription: "legacy-v3 shows merge anchor list notation",
		document:       anchorMergeDocSample,
		expression:     `.foobarList`,
		expected: []string{
			"D0, P[foobarList], (!!map)::b: foobarList_b\n!!merge <<: [*foo, *bar]\nc: foobarList_c\n",
		},
		skipForGoccy: true,
	},
	{
		description:    "Merge anchor list display - goccy behaviour",
		subdescription: "Goccy expands merge anchor lists in output",
		document:       anchorMergeDocSample,
		expression:     `.foobarList`,
		expected: []string{
			"D0, P[foobarList], (!!map)::<<:\n    - a: foo_a\n      c: foo_c\n      thing: foo_thing\n    - b: bar_b\n      c: bar_c\n      thing: bar_thing\nb: foobarList_b\nc: foobarList_c\n",
		},
		skipForYamlV3: true,
	},
}

func TestAnchorAliasOperatorScenarios(t *testing.T) {
	for _, tt := range anchorOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "anchor-and-alias-operators", anchorOperatorScenarios)
}
