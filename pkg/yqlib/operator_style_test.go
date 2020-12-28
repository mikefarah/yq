package yqlib

import (
	"testing"
)

var styleOperatorScenarios = []expressionScenario{
	{
		description: "Set tagged style",
		document:    `{a: cat, b: 5, c: 3.2, e: true}`,
		expression:  `.. style="tagged"`,
		expected: []string{
			"D0, P[], (!!map)::!!map\na: !!str cat\nb: !!int 5\nc: !!float 3.2\ne: !!bool true\n",
		},
	},
	{
		description: "Set double quote style",
		document:    `{a: cat, b: 5, c: 3.2, e: true}`,
		expression:  `.. style="double"`,
		expected: []string{
			"D0, P[], (!!map)::a: \"cat\"\nb: \"5\"\nc: \"3.2\"\ne: \"true\"\n",
		},
	},
	{
		description: "Set double quote style on map keys too",
		document:    `{a: cat, b: 5, c: 3.2, e: true}`,
		expression:  `... style="double"`,
		expected: []string{
			"D0, P[], (!!map)::\"a\": \"cat\"\n\"b\": \"5\"\n\"c\": \"3.2\"\n\"e\": \"true\"\n",
		},
	},
	{
		skipDoc:    true,
		document:   "bing: &foo frog\na:\n  c: cat\n  <<: [*foo]",
		expression: `(... | select(tag=="!!str")) style="single"`,
		expected: []string{
			"D0, P[], (!!map)::'bing': &foo 'frog'\n'a':\n    'c': 'cat'\n    !!merge <<: [*foo]\n",
		},
	},
	{
		description: "Set single quote style",
		document:    `{a: cat, b: 5, c: 3.2, e: true}`,
		expression:  `.. style="single"`,
		expected: []string{
			"D0, P[], (!!map)::a: 'cat'\nb: '5'\nc: '3.2'\ne: 'true'\n",
		},
	},
	{
		description: "Set literal quote style",
		document:    `{a: cat, b: 5, c: 3.2, e: true}`,
		expression:  `.. style="literal"`,
		expected: []string{
			`D0, P[], (!!map)::a: |-
    cat
b: |-
    5
c: |-
    3.2
e: |-
    true
`,
		},
	},
	{
		description: "Set folded quote style",
		document:    `{a: cat, b: 5, c: 3.2, e: true}`,
		expression:  `.. style="folded"`,
		expected: []string{
			`D0, P[], (!!map)::a: >-
    cat
b: >-
    5
c: >-
    3.2
e: >-
    true
`,
		},
	},
	{
		description: "Set flow quote style",
		document:    `{a: cat, b: 5, c: 3.2, e: true}`,
		expression:  `.. style="flow"`,
		expected: []string{
			"D0, P[], (!!map)::{a: cat, b: 5, c: 3.2, e: true}\n",
		},
	},
	{
		description:    "Pretty print",
		subdescription: "Set empty (default) quote style, note the usage of `...` to match keys too.",
		document:       `{a: cat, "b": 5, 'c': 3.2, "e": true}`,
		expression:     `... style=""`,
		expected: []string{
			"D0, P[], (!!map)::a: cat\nb: 5\nc: 3.2\ne: true\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: cat, b: double}`,
		expression: `.a style=.b`,
		expected: []string{
			"D0, P[], (doc)::{a: \"cat\", b: double}\n",
		},
	},
	{
		description:           "Read style",
		document:              `{a: "cat", b: 'thing'}`,
		dontFormatInputForDoc: true,
		expression:            `.. | style`,
		expected: []string{
			"D0, P[], (!!str)::flow\n",
			"D0, P[a], (!!str)::double\n",
			"D0, P[b], (!!str)::single\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: cat`,
		expression: `.. | style`,
		expected: []string{
			"D0, P[], (!!str)::\"\"\n",
			"D0, P[a], (!!str)::\"\"\n",
		},
	},
}

func TestStyleOperatorScenarios(t *testing.T) {
	for _, tt := range styleOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Style", styleOperatorScenarios)
}
