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
    a: foo_a
    c: foobarList_c
foobar:
    c: foo_c
    a: foo_a
    thing: foobar_thing
`

var explodeWhenKeysExistDocument = `objects:
  - &circle
    name: circle
    shape: round
  - name: ellipse
    !!merge <<: *circle
  - !!merge <<: *circle
    name: egg
`

var explodeWhenKeysExistLegacy = `D0, P[], (!!map)::objects:
    - name: circle
      shape: round
    - name: circle
      shape: round
    - shape: round
      name: egg
`

var explodeWhenKeysExistExpected = `D0, P[], (!!map)::objects:
    - name: circle
      shape: round
    - name: ellipse
      shape: round
    - shape: round
      name: egg
`

var fixedAnchorOperatorScenarios = []expressionScenario{
	{
		skipDoc:        true,
		description:    "merge anchor after existing keys",
		subdescription: "Does not override existing keys",
		document:       explodeWhenKeysExistDocument,
		expression:     "explode(.)",
		expected:       []string{explodeWhenKeysExistExpected},
	},

	// The following tests are the same as below, to verify they still works correctly with the flag:
	{
		description:    "Override",
		subdescription: "see https://yaml.org/type/merge.html",
		document:       specDocument + "- << : [ *BIG, *LEFT, *SMALL ]\n  x: 1\n",
		expression:     ".[4] | explode(.)",
		expected:       []string{"D0, P[4], (!!map)::r: 10\ny: 2\nx: 1\n"},
	},
	{
		skipDoc:        true,
		description:    "Duplicate keys",
		subdescription: "outside merge anchor",
		document:       `{a: 1, a: 2}`,
		expression:     `explode(.)`,
		expected: []string{
			"D0, P[], (!!map)::{a: 1, a: 2}\n",
		},
	},
}

var anchorOperatorScenarios = []expressionScenario{
	{
		skipDoc:        true,
		description:    "merge anchor after existing keys",
		subdescription: "legacy: overrides existing keys",
		document:       explodeWhenKeysExistDocument,
		expression:     "explode(.)",
		expected:       []string{explodeWhenKeysExistLegacy},
	},
	{
		skipDoc:       true,
		description:   "merge anchor not map",
		document:      "a: &a\n  - 0\nc:\n  <<: [*a]\n",
		expectedError: "can only use merge anchors with maps (!!map) or sequences (!!seq) of maps, but got sequence containing !!seq",
		expression:    "explode(.)",
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
	//TODO The following 2 tests warn about overwriting [3].r not being to spec while they shouldn't
	{
		description:    "Override",
		subdescription: "see https://yaml.org/type/merge.html",
		document:       specDocument + "- << : [ *BIG, *LEFT, *SMALL ]\n  x: 1\n",
		expression:     ".[4] | explode(.)",
		expected:       []string{"D0, P[4], (!!map)::r: 10\ny: 2\nx: 1\n"},
	},
	// Correctly warns about overwriting [4].x
	{
		description: "Override with local key",
		subdescription: "like https://yaml.org/type/merge.html, but with x: 1 before the merge key. " +
			"This is legacy behavior, see --yaml-fix-merge-anchor-to-spec",
		document:   specDocument + "- x: 1\n  << : [ *BIG, *LEFT, *SMALL ]\n",
		expression: ".[4] | explode(.)",
		expected:   []string{"D0, P[4], (!!map)::x: 0\nr: 10\ny: 2\n"},
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
			"D0, P[foobarList], (!!map)::{b: bar_b, thing: foo_thing, a: foo_a, c: foobarList_c}\n",
			"D0, P[foobar], (!!map)::{c: foo_c, a: foo_a, thing: foobar_thing}\n",
		},
	},
	{
		skipDoc:    true,
		document:   mergeDocSample,
		expression: `.foo* | explode(explode(.)) | (. style="flow")`,
		expected: []string{
			"D0, P[foo], (!!map)::{a: foo_a, thing: foo_thing, c: foo_c}\n",
			"D0, P[foobarList], (!!map)::{b: bar_b, thing: foo_thing, a: foo_a, c: foobarList_c}\n",
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
	},
	{
		description:    "Dereference and update a field",
		subdescription: "Use explode with multiply to dereference an object",
		document:       simpleArrayRef,
		expression:     `.thingOne |= explode(.) * {"value": false}`,
		expected:       []string{expectedUpdatedArrayRef},
	},
	{
		skipDoc:     true,
		description: "Merge anchor with inline map",
		document:    `{<<: {a: 42}}`,
		expression:  `explode(.)`,
		expected: []string{
			"D0, P[], (!!map)::{a: 42}\n",
		},
	},
	{
		skipDoc:     true,
		description: "Merge anchor with sequence with inline map",
		document:    `{<<: [{a: 42}]}`,
		expression:  `explode(.)`,
		expected: []string{
			"D0, P[], (!!map)::{a: 42}\n",
		},
	},
	{
		skipDoc:     true,
		description: "Merge anchor with aliased sequence with inline map",
		document:    `{s: &s [{a: 42}], m: {<<: *s}}`,
		expression:  `.m | explode(.)`,
		expected: []string{
			"D0, P[m], (!!map)::{a: 42}\n",
		},
	},
	{
		skipDoc:        true,
		description:    "Exploding merge anchor should not explode neighbors",
		subdescription: "b must not be exploded, as `r: *a` will become invalid",
		document:       `{b: &b {a: &a 42}, r: *a, c: {<<: *b}}`,
		expression:     `explode(.c)`,
		expected: []string{
			"D0, P[], (!!map)::{b: &b {a: &a 42}, r: *a, c: {a: &a 42}}\n",
		},
	},
	{
		skipDoc:        true,
		description:    "Exploding sequence merge anchor should not explode neighbors",
		subdescription: "b must not be exploded, as `r: *a` will become invalid",
		document:       `{b: &b {a: &a 42}, r: *a, c: {<<: [*b]}}`,
		expression:     `explode(.c)`,
		expected: []string{
			"D0, P[], (!!map)::{b: &b {a: &a 42}, r: *a, c: {a: &a 42}}\n",
		},
	},
	{
		skipDoc:        true,
		description:    "Exploding inline merge anchor",
		subdescription: "`<<` map must be exploded, otherwise `c: *b` will become invalid",
		document:       `{a: {b: &b 42}, <<: {c: *b}}`,
		expression:     `explode(.)`,
		expected: []string{
			"D0, P[], (!!map)::{a: {b: 42}, c: 42}\n",
		},
	},
	{
		skipDoc:        true,
		description:    "Exploding inline merge anchor in sequence",
		subdescription: "`<<` map must be exploded, otherwise `c: *b` will become invalid",
		document:       `{a: {b: &b 42}, <<: [{c: *b}]}`,
		expression:     `explode(.)`,
		expected: []string{
			"D0, P[], (!!map)::{a: {b: 42}, c: 42}\n",
		},
	},
	{
		skipDoc:        true,
		description:    "Duplicate keys",
		subdescription: "outside merge anchor",
		document:       `{a: 1, a: 2}`,
		expression:     `explode(.)`,
		expected: []string{
			"D0, P[], (!!map)::{a: 1, a: 2}\n",
		},
	},
}

func TestAnchorAliasOperatorScenarios(t *testing.T) {
	for _, tt := range anchorOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "anchor-and-alias-operators", anchorOperatorScenarios)
}

func TestAnchorAliasOperatorAlignedToSpecScenarios(t *testing.T) {
	ConfiguredYamlPreferences.FixMergeAnchorToSpec = true
	for _, tt := range fixedAnchorOperatorScenarios {
		testScenario(t, &tt)
	}
	ConfiguredYamlPreferences.FixMergeAnchorToSpec = false
}
