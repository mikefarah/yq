package yqlib

import (
	"testing"
)

var doc1 = `list:
# Hi this is a comment.
# Hello this is another comment.
- "abc"`

var doc2 = `list2:
# This is yet another comment.
# Indeed this is yet another comment.
- "123"`

var docExpected = `D0, P[], (!!map)::list:
    # Hi this is a comment.
    # Hello this is another comment.
    - "abc"
list2:
    # This is yet another comment.
    # Indeed this is yet another comment.
    - "123"
`

var mergeArrayWithAnchors = `sample:
- &a
- <<: *a
`

var mergeArraysObjectKeysText = `It's a complex command, the trickyness comes from needing to have the right context in the expressions.
First we save the second array into a variable '$two' which lets us reference it later.
We then need to update the first array. We will use the relative update (|=) because we need to update relative to the current element of the array in the LHS in the RHS expression. 
We set the current element of the first array as $cur. Now we multiply (merge) $cur with the matching entry in $two, by passing $two through a select filter.
`

var docWithHeader = `
# here

a: apple
`

var nodeWithHeader = `
# here
a: apple
`

var docNoComments = `
b: banana
`

var docWithFooter = `
a: apple

# footer
`

var nodeWithFooter = `
a: apple
# footer`

var document = `
a: &cat {name: cat}
b: {name: dog}
c: 
  <<: *cat
`

var multiplyOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		document:   mergeArrayWithAnchors,
		expression: `. * .`,
		expected: []string{
			"D0, P[], (!!map)::sample:\n    - &a\n    - !!merge <<: *a\n",
		},
	},
	{
		skipDoc:    true,
		document:   docWithHeader,
		document2:  docNoComments,
		expression: `select(fi == 0) * select(fi == 1)`,
		expected: []string{
			"D0, P[], (!!map)::# here\na: apple\nb: banana\n",
		},
	},
	{
		skipDoc:    true,
		document:   nodeWithHeader,
		document2:  docNoComments,
		expression: `select(fi == 0) * select(fi == 1)`,
		expected: []string{
			"D0, P[], (!!map)::# here\na: apple\nb: banana\n",
		},
	},
	{
		skipDoc:    true,
		document:   docNoComments,
		document2:  docWithHeader,
		expression: `select(fi == 0) * select(fi == 1)`,
		expected: []string{
			"D0, P[], (!!map)::# here\nb: banana\na: apple\n",
		},
	},
	{
		skipDoc:    true,
		document:   docNoComments,
		document2:  nodeWithHeader,
		expression: `select(fi == 0) * select(fi == 1)`,
		expected: []string{
			"D0, P[], (!!map)::b: banana\n# here\na: apple\n",
		},
	},
	{
		skipDoc:    true,
		document:   docWithFooter,
		document2:  docNoComments,
		expression: `select(fi == 0) * select(fi == 1)`,
		expected: []string{
			"D0, P[], (!!map)::a: apple\nb: banana\n\n# footer\n",
		},
	},
	{
		skipDoc:    true,
		document:   nodeWithFooter,
		document2:  docNoComments,
		expression: `select(fi == 0) * select(fi == 1)`,
		expected: []string{ // not sure why there's an extra newline *shrug*
			"D0, P[], (!!map)::a: apple\n# footer\n\nb: banana\n",
		},
	},
	{
		skipDoc:    true,
		document:   docNoComments,
		document2:  docWithFooter,
		expression: `select(fi == 0) * select(fi == 1)`,
		expected: []string{
			"D0, P[], (!!map)::b: banana\na: apple\n\n# footer\n",
		},
	},
	{
		description: "Multiply integers",
		expression:  `3 * 4`,
		expected: []string{
			"D0, P[], (!!int)::12\n",
		},
	},
	{
		skipDoc:    true,
		document:   doc1,
		document2:  doc2,
		expression: `select(fi == 0) * select(fi == 1)`,
		expected: []string{
			docExpected,
		},
	},
	{
		skipDoc:    true,
		expression: `.x =  {"things": "whatever"} * {}`,
		expected: []string{
			"D0, P[], ()::x:\n    things: whatever\n",
		},
	},
	{
		skipDoc:    true,
		expression: `.x = {} * {"things": "whatever"}`,
		expected: []string{
			"D0, P[], ()::x:\n    things: whatever\n",
		},
	},
	{
		skipDoc:    true,
		expression: `3 * 4.5`,
		expected: []string{
			"D0, P[], (!!float)::13.5\n",
		},
	},
	{
		skipDoc:    true,
		expression: `4.5 * 3`,
		expected: []string{
			"D0, P[], (!!float)::13.5\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {also: [1]}, b: {also: me}}`,
		expression: `. * {"a" : .b}`,
		expected: []string{
			"D0, P[], (!!map)::{a: {also: me}, b: {also: me}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   "# b\nb:\n  # a\n  a: cat",
		expression: "{} * .",
		expected: []string{
			"D0, P[], (!!map)::# b\nb:\n    # a\n    a: cat\n",
		},
	},
	{
		skipDoc:    true,
		document:   "# b\nb:\n  # a\n  a: cat",
		expression: ". * {}",
		expected: []string{
			"D0, P[], (!!map)::# b\nb:\n    # a\n    a: cat\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: &a { b: &b { c: &c cat } } }`,
		expression: `{} * .`,
		expected: []string{
			"D0, P[], (!!map)::{a: &a {b: &b {c: &c cat}}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: 2, b: 5}`,
		document2:  `{a: 3, b: 10}`,
		expression: `.a * .b`,
		expected: []string{
			"D0, P[a], (!!int)::10\n",
			"D0, P[a], (!!int)::20\n",
			"D0, P[a], (!!int)::15\n",
			"D0, P[a], (!!int)::30\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: 2}`,
		document2:  `{b: 10}`,
		expression: `select(fi ==0) * select(fi==1)`,
		expected: []string{
			"D0, P[], (!!map)::{a: 2, b: 10}\n",
		},
	},
	{
		skipDoc:    true,
		expression: `{} * {"cat":"dog"}`,
		expected: []string{
			"D0, P[], (!!map)::cat: dog\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {also: me}, b: {also: [1]}}`,
		expression: `. * {"a":.b}`,
		expected: []string{
			"D0, P[], (!!map)::{a: {also: [1]}, b: {also: [1]}}\n",
		},
	},
	{
		description: "Merge objects together, returning merged result only",
		document:    `{a: {field: me, fieldA: cat}, b: {field: {g: wizz}, fieldB: dog}}`,
		expression:  `.a * .b`,
		expected: []string{
			"D0, P[a], (!!map)::{field: {g: wizz}, fieldA: cat, fieldB: dog}\n",
		},
	},
	{
		description: "Merge objects together, returning parent object",
		document:    `{a: {field: me, fieldA: cat}, b: {field: {g: wizz}, fieldB: dog}}`,
		expression:  `. * {"a":.b}`,
		expected: []string{
			"D0, P[], (!!map)::{a: {field: {g: wizz}, fieldA: cat, fieldB: dog}, b: {field: {g: wizz}, fieldB: dog}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {also: {g: wizz}}, b: {also: me}}`,
		expression: `. * {"a":.b}`,
		expected: []string{
			"D0, P[], (!!map)::{a: {also: me}, b: {also: me}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {also: {g: wizz}}, b: {also: [1]}}`,
		expression: `. * {"a":.b}`,
		expected: []string{
			"D0, P[], (!!map)::{a: {also: [1]}, b: {also: [1]}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {also: [1]}, b: {also: {g: wizz}}}`,
		expression: `. * {"a":.b}`,
		expected: []string{
			"D0, P[], (!!map)::{a: {also: {g: wizz}}, b: {also: {g: wizz}}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {things: great}, b: {also: me}}`,
		expression: `. * {"a": .b}`,
		expected: []string{
			"D0, P[], (!!map)::{a: {things: great, also: me}, b: {also: me}}\n",
		},
	},
	{
		description:           "Merge keeps style of LHS",
		dontFormatInputForDoc: true,
		document:              "a: {things: great}\nb:\n  also: \"me\"",
		expression:            `. * {"a":.b}`,
		expected: []string{
			"D0, P[], (!!map)::a: {things: great, also: \"me\"}\nb:\n    also: \"me\"\n",
		},
	},
	{
		description: "Merge arrays",
		document:    `{a: [1,2,3], b: [3,4,5]}`,
		expression:  `. * {"a":.b}`,
		expected: []string{
			"D0, P[], (!!map)::{a: [3, 4, 5], b: [3, 4, 5]}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: [1], b: [2]}`,
		expression: `.a *+ .b`,
		expected: []string{
			"D0, P[a], (!!seq)::[1, 2]\n",
		},
	},
	{
		description: "Merge, only existing fields",
		document:    `{a: {thing: one, cat: frog}, b: {missing: two, thing: two}}`,
		expression:  `.a *? .b`,
		expected: []string{
			"D0, P[a], (!!map)::{thing: two, cat: frog}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: [{thing: one}], b: [{missing: two, thing: two}]}`,
		expression: `.a *?d .b`,
		expected: []string{
			"D0, P[a], (!!seq)::[{thing: two}]\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {array: [1]}, b: {}}`,
		expression: `.b *+ .a`,
		expected: []string{
			"D0, P[b], (!!map)::{array: [1]}\n",
		},
	},
	{
		description: "Merge, appending arrays",
		document:    `{a: {array: [1, 2, animal: dog], value: coconut}, b: {array: [3, 4, animal: cat], value: banana}}`,
		expression:  `.a *+ .b`,
		expected: []string{
			"D0, P[a], (!!map)::{array: [1, 2, {animal: dog}, 3, 4, {animal: cat}], value: banana}\n",
		},
	},
	{
		description: "Merge, only existing fields, appending arrays",
		document:    `{a: {thing: [1,2]}, b: {thing: [3,4], another: [1]}}`,
		expression:  `.a *?+ .b`,
		expected: []string{
			"D0, P[a], (!!map)::{thing: [1, 2, 3, 4]}\n",
		},
	},
	{
		description:    "Merge, deeply merging arrays",
		subdescription: "Merging arrays deeply means arrays are merge like objects, with indexes as their key. In this case, we merge the first item in the array, and do nothing with the second.",
		document:       `{a: [{name: fred, age: 12}, {name: bob, age: 32}], b: [{name: fred, age: 34}]}`,
		expression:     `.a *d .b`,
		expected: []string{
			"D0, P[a], (!!seq)::[{name: fred, age: 34}, {name: bob, age: 32}]\n",
		},
	},
	{
		description:    "Merge arrays of objects together, matching on a key",
		subdescription: mergeArraysObjectKeysText,
		document:       `[{a: apple, b: appleB}, {a: kiwi, b: kiwiB}, {a: banana, b: bananaB}]`,
		document2:      `[{a: banana, c: bananaC}, {a: apple, b: appleB2}, {a: dingo, c: dingoC}]`,
		expression:     `(select(fi==1) | .[]) as $two | select(fi==0) | .[] |= (. as $cur |  $cur * ($two | select(.a == $cur.a)))`,
		expected: []string{
			"D0, P[], (doc)::[{a: apple, b: appleB2}, {a: kiwi, b: kiwiB}, {a: banana, b: bananaB, c: bananaC}]\n",
		},
	},
	{
		description: "Merge to prefix an element",
		document:    `{a: cat, b: dog}`,
		expression:  `. * {"a": {"c": .a}}`,
		expected: []string{
			"D0, P[], (!!map)::{a: {c: cat}, b: dog}\n",
		},
	},
	{
		description: "Merge with simple aliases",
		document:    `{a: &cat {c: frog}, b: {f: *cat}, c: {g: thongs}}`,
		expression:  `.c * .b`,
		expected: []string{
			"D0, P[c], (!!map)::{g: thongs, f: *cat}\n",
		},
	},
	{
		description: "Merge copies anchor names",
		document:    `{a: {c: &cat frog}, b: {f: *cat}, c: {g: thongs}}`,
		expression:  `.c * .a`,
		expected: []string{
			"D0, P[c], (!!map)::{g: thongs, c: &cat frog}\n",
		},
	},
	{
		description: "Merge with merge anchors",
		document:    mergeDocSample,
		expression:  `.foobar * .foobarList`,
		expected: []string{
			"D0, P[foobar], (!!map)::c: foobarList_c\n!!merge <<: [*foo, *bar]\nthing: foobar_thing\nb: foobarList_b\n",
		},
	},
	{
		skipDoc:    true,
		document:   document,
		expression: `.b * .c`,
		expected: []string{
			"D0, P[b], (!!map)::{name: dog, <<: *cat}\n",
		},
	},
}

func TestMultiplyOperatorScenarios(t *testing.T) {
	for _, tt := range multiplyOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Multiply (Merge)", multiplyOperatorScenarios)
}
