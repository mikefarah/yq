package yqlib

import (
	"fmt"
	"strings"
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

var mergeArraysObjectKeysText = `
This is a fairly complex expression - you can use it as is by providing the environment variables as seen in the example below.

It merges in the array provided in the second file into the first - matching on equal keys.

Explanation:

The approach, at a high level, is to reduce into a merged map (keyed by the unique key)
and then convert that back into an array.

First the expression will create a map from the arrays keyed by the idPath, the unique field we want to merge by.
The reduce operator is merging '({}; . * $item )', so array elements with the matching key will be merged together.

Next, we convert the map back to an array, using reduce again, concatenating all the map values together.

Finally, we set the result of the merged array back into the first doc.

Thanks Kev from [stackoverflow](https://stackoverflow.com/a/70109529/1168223)
`

var mergeExpression = `
(
  (( (eval(strenv(originalPath)) + eval(strenv(otherPath)))  | .[] | {(eval(strenv(idPath))):  .}) as $item ireduce ({}; . * $item )) as $uniqueMap
  | ( $uniqueMap  | to_entries | .[]) as $item ireduce([]; . + $item.value)
) as $mergedArray
| select(fi == 0) | (eval(strenv(originalPath))) = $mergedArray
`

var docWithHeader = `# here

a: apple
`

var nodeWithHeader = `node:
  # here
  a: apple
`

var docNoComments = `b: banana
`

var docWithFooter = `a: apple

# footer
`

var nodeWithFooter = `a: apple
# footer`

var document = `a: &cat {name: cat}
b: {name: dog}
c: 
  <<: *cat
`

var mergeWithGlobA = `
"**cat": things,
"meow**cat": stuff
`

var mergeWithGlobB = `
"**cat": newThings,
`

var multiplyOperatorScenarios = []expressionScenario{
	{
		description: "multiple should be readonly",
		skipDoc:     true,
		document:    "",
		expression:  ".x |= (root | (.a * .b))",
		expected: []string{
			"D0, P[], ()::x: null\n",
		},
	},
	{
		description: "glob keys are treated as literals when merging",
		skipDoc:     true,
		document:    mergeWithGlobA,
		document2:   mergeWithGlobB,
		expression:  `select(fi == 0) * select(fi == 1)`,
		expected: []string{
			"D0, P[], (!!map)::\n\"**cat\": newThings,\n\"meow**cat\": stuff\n",
		},
	},
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
		document:   `[[c], [b]]`,
		expression: `.[] | . *+ ["a"]`,
		expected: []string{
			"D0, P[0], (!!seq)::[c, a]\n",
			"D0, P[1], (!!seq)::[b, a]\n",
		},
	},
	{
		skipDoc:    true,
		document:   docWithHeader,
		document2:  docNoComments,
		expression: `select(fi == 0) * select(fi == 1)`,
		expected: []string{
			"D0, P[], (!!map)::# here\n\na: apple\nb: banana\n",
		},
	},
	{
		skipDoc:    true,
		document:   nodeWithHeader,
		document2:  docNoComments,
		expression: `(select(fi == 0) | .node) * select(fi == 1)`,
		expected: []string{
			"D0, P[node], (!!map)::# here\na: apple\nb: banana\n",
		},
	},
	{
		skipDoc:    true,
		document:   docNoComments,
		document2:  docWithHeader,
		expression: `select(fi == 0) * select(fi == 1)`,
		expected: []string{
			"D0, P[], (!!map)::# here\n\nb: banana\na: apple\n",
		},
	},
	{
		skipDoc:    true,
		document:   docNoComments,
		document2:  nodeWithHeader,
		expression: `select(fi == 0) * (select(fi == 1) | .node)`,
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
			"D0, P[], (!!map)::a: apple\nb: banana\n# footer\n",
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
		description: "doc 2 has footer",
		skipDoc:     true,
		document:    docNoComments,
		document2:   docWithFooter,
		expression:  `select(fi == 0) * select(fi == 1)`,
		expected: []string{
			"D0, P[], (!!map)::b: banana\na: apple\n# footer\n",
		},
	},
	{
		description: "Multiply integers",
		document:    "a: 3\nb: 4",
		expression:  `.a *= .b`,
		expected: []string{
			"D0, P[], (!!map)::a: 12\nb: 4\n",
		},
	},
	{
		description: "Multiply string node X int",
		document:    docNoComments,
		expression:  ".b * 4",
		expected: []string{
			fmt.Sprintf("D0, P[b], (!!str)::%s\n", strings.Repeat("banana", 4)),
		},
	},
	{
		description: "Multiply int X string node",
		document:    docNoComments,
		expression:  "4 * .b",
		expected: []string{
			fmt.Sprintf("D0, P[], (!!str)::%s\n", strings.Repeat("banana", 4)),
		},
	},
	{
		description: "Multiply string X int node",
		document: `n: 4
`,
		expression: `"banana" * .n`,
		expected: []string{
			fmt.Sprintf("D0, P[], (!!str)::%s\n", strings.Repeat("banana", 4)),
		},
	},
	{
		description:   "Multiply string X by negative int",
		skipDoc:       true,
		document:      `n: -4`,
		expression:    `"banana" * .n`,
		expectedError: "cannot repeat string by a negative number (-4)",
	},
	{
		description: "Multiply string X by more than 100 million",
		// very large string.repeats causes a panic
		skipDoc:       true,
		document:      `n: 100000001`,
		expression:    `"banana" * .n`,
		expectedError: "cannot repeat string by more than 100 million (100000001)",
	},
	{
		description: "Multiply int node X string",
		document: `n: 4
`,
		expression: `.n * "banana"`,
		expected: []string{
			fmt.Sprintf("D0, P[n], (!!str)::%s\n", strings.Repeat("banana", 4)),
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
			"D0, P[], (!!map)::a: &a\n    b: &b\n        c: &c cat\n",
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
		description: "Merge, only new fields",
		document:    `{a: {thing: one, cat: frog}, b: {missing: two, thing: two}}`,
		expression:  `.a *n .b`,
		expected: []string{
			"D0, P[a], (!!map)::{thing: one, cat: frog, missing: two}\n",
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
		document:   `{a: [{thing: one}], b: [{missing: two, thing: two}]}`,
		expression: `.a *nd .b`,
		expected: []string{
			"D0, P[a], (!!seq)::[{thing: one, missing: two}]\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {array: [1]}, b: {}}`,
		expression: `.b *+ .a`,
		expected: []string{
			"D0, P[b], (!!map)::array: [1]\n",
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
		description:    "Merge, only new fields, appending arrays",
		subdescription: "Append (+) with (n) has no effect.",
		skipDoc:        true,
		document:       `{a: {thing: [1,2]}, b: {thing: [3,4], another: [1]}}`,
		expression:     `.a *n+ .b`,
		expected: []string{
			"D0, P[a], (!!map)::{thing: [1, 2], another: [1]}\n",
		},
	},
	{
		description:    "Merge, deeply merging arrays",
		subdescription: "Merging arrays deeply means arrays are merged like objects, with indices as their key. In this case, we merge the first item in the array and do nothing with the second.",
		document:       `{a: [{name: fred, age: 12}, {name: bob, age: 32}], b: [{name: fred, age: 34}]}`,
		expression:     `.a *d .b`,
		expected: []string{
			"D0, P[a], (!!seq)::[{name: fred, age: 34}, {name: bob, age: 32}]\n",
		},
	},
	{
		description:          "Merge arrays of objects together, matching on a key",
		subdescription:       mergeArraysObjectKeysText,
		document:             `{myArray: [{a: apple, b: appleB}, {a: kiwi, b: kiwiB}, {a: banana, b: bananaB}], something: else}`,
		document2:            `newArray: [{a: banana, c: bananaC}, {a: apple, b: appleB2}, {a: dingo, c: dingoC}]`,
		environmentVariables: map[string]string{"originalPath": ".myArray", "otherPath": ".newArray", "idPath": ".a"},
		expression:           mergeExpression,
		expected: []string{
			"D0, P[], (!!map)::{myArray: [{a: apple, b: appleB2}, {a: kiwi, b: kiwiB}, {a: banana, b: bananaB, c: bananaC}, {a: dingo, c: dingoC}], something: else}\n",
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
	{
		description:    "Custom types: that are really numbers",
		subdescription: "When custom tags are encountered, yq will try to decode the underlying type.",
		document:       "a: !horse 2\nb: !goat 3",
		expression:     ".a = .a * .b",
		expected: []string{
			"D0, P[], (!!map)::a: !horse 6\nb: !goat 3\n",
		},
	},
	{
		skipDoc:     true,
		description: "Custom types: that are really numbers",
		document:    "a: !horse 2.5\nb: !goat 3.5",
		expression:  ".a = .a * .b",
		expected: []string{
			"D0, P[], (!!map)::a: !horse 8.75\nb: !goat 3.5\n",
		},
	},
	{
		skipDoc:     true,
		description: "Custom types: that are really numbers",
		document:    "a: 2\nb: !goat 3.5",
		expression:  ".a = .a * .b",
		expected: []string{
			"D0, P[], (!!map)::a: !!float 7\nb: !goat 3.5\n",
		},
	},
	{
		skipDoc:     true,
		description: "Custom types: that are really arrays",
		document:    "a: !horse [1,2]\nb: !goat [3]",
		expression:  ".a = .a * .b",
		expected: []string{
			"D0, P[], (!!map)::a: !horse [3]\nb: !goat [3]\n",
		},
	},
	{
		description:    "Custom types: that are really maps",
		subdescription: "Custom tags will be maintained.",
		document:       "a: !horse {cat: meow}\nb: !goat {dog: woof}",
		expression:     ".a = .a * .b",
		expected: []string{
			"D0, P[], (!!map)::a: !horse {cat: meow, dog: woof}\nb: !goat {dog: woof}\n",
		},
	},
	{
		description:    "Custom types: clobber tags",
		subdescription: "Use the `c` option to clobber custom tags. Note that the second tag is now used.",
		document:       "a: !horse {cat: meow}\nb: !goat {dog: woof}",
		expression:     ".a *=c .b",
		expected: []string{
			"D0, P[], (!!map)::a: !goat {cat: meow, dog: woof}\nb: !goat {dog: woof}\n",
		},
	},
	{
		skipDoc:        true,
		description:    "Custom types: clobber tags - *=",
		subdescription: "Use the `c` option to clobber custom tags - on both the `=` and `*` operator. Note that the second tag is now used.",
		document:       "a: !horse {cat: meow}\nb: !goat {dog: woof}",
		expression:     ".a =c .a *c .b",
		expected: []string{
			"D0, P[], (!!map)::a: !goat {cat: meow, dog: woof}\nb: !goat {dog: woof}\n",
		},
	},
	{
		skipDoc:        true,
		description:    "Custom types: dont clobber tags - *=",
		subdescription: "Use the `c` option to clobber custom tags - on both the `=` and `*` operator. Note that the second tag is now used.",
		document:       "a: !horse {cat: meow}\nb: !goat {dog: woof}",
		expression:     ".a *= .b",
		expected: []string{
			"D0, P[], (!!map)::a: !horse {cat: meow, dog: woof}\nb: !goat {dog: woof}\n",
		},
	},
	{
		skipDoc:     true,
		description: "Custom types: that are really maps",
		document:    "a: {cat: !horse meow}\nb: {cat: 5}",
		expression:  ".a = .a * .b",
		expected: []string{
			"D0, P[], (!!map)::a: {cat: !horse 5}\nb: {cat: 5}\n",
		},
	},
	{
		skipDoc:     true,
		description: "Relative merge, new fields only",
		document:    "a: {a: original}\n",
		expression:  `.a *=n load("../../examples/thing.yml")`,
		expected: []string{
			"D0, P[], (!!map)::a: {a: original, b: cool.}\n",
		},
	},
	{
		skipDoc:     true,
		description: "Relative merge",
		document:    "a: {a: original}\n",
		expression:  `.a *= load("../../examples/thing.yml")`,
		expected: []string{
			"D0, P[], (!!map)::a: {a: apple is included, b: cool.}\n",
		},
	},
	{
		description: "Merging a null with a map",
		expression:  `null * {"some": "thing"}`,
		expected: []string{
			"D0, P[], (!!map)::some: thing\n",
		},
	},
	{
		description: "Merging a map with null",
		expression:  `{"some": "thing"} * null`,
		expected: []string{
			"D0, P[], (!!map)::some: thing\n",
		},
	},
	{
		description: "Merging a null with an array",
		expression:  `null * ["some"]`,
		expected: []string{
			"D0, P[], (!!seq)::- some\n",
		},
	},
	{
		description: "Merging an array with null",
		expression:  `["some"] * null`,
		expected: []string{
			"D0, P[], (!!seq)::- some\n",
		},
	},
	{
		skipDoc:    true,
		expression: `null * null`,
		expected: []string{
			"D0, P[], (!!null)::null\n",
		},
	},
}

func TestMultiplyOperatorScenarios(t *testing.T) {
	for _, tt := range multiplyOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "multiply-merge", multiplyOperatorScenarios)
}
