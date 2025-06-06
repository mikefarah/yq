package yqlib

import (
	"testing"
)

var styleOperatorScenarios = []expressionScenario{
	{
		description: "Update and set style of a particular node (simple)",
		document:    `a: {b: thing, c: something}`,
		expression:  `.a.b = "new" | .a.b style="double"`,
		expected: []string{
			"D0, P[], (!!map)::a: {b: \"new\", c: something}\n",
		},
	},
	{
		description: "Update and set style of a particular node using path variables",
		document:    `a: {b: thing, c: something}`,
		expression:  `with(.a.b ; . = "new" | . style="double")`,
		expected: []string{
			"D0, P[], (!!map)::a: {b: \"new\", c: something}\n",
		},
	},
	{
		description: "Set tagged style",
		document:    `{a: cat, b: 5, c: 3.2, e: true, f: [1,2,3], g: { something: cool}}`,
		expression:  `.. style="tagged"`,
		expected: []string{
			"D0, P[], (!!map)::!!map\na: !!str cat\nb: !!int 5\nc: !!float 3.2\ne: !!bool true\nf: !!seq\n    - !!int 1\n    - !!int 2\n    - !!int 3\ng: !!map\n    something: !!str cool\n",
		},
	},
	{
		description: "Set double quote style",
		document:    `{a: cat, b: 5, c: 3.2, e: true, f: [1,2,3], g: { something: cool}}`,
		expression:  `.. style="double"`,
		expected: []string{
			"D0, P[], (!!map)::a: \"cat\"\nb: \"5\"\nc: \"3.2\"\ne: \"true\"\nf:\n    - \"1\"\n    - \"2\"\n    - \"3\"\ng:\n    something: \"cool\"\n",
		},
	},
	{
		description: "Set double quote style on map keys too",
		document:    `{a: cat, b: 5, c: 3.2, e: true, f: [1,2,3], g: { something: cool}}`,
		expression:  `... style="double"`,
		expected: []string{
			"D0, P[], (!!map)::\"a\": \"cat\"\n\"b\": \"5\"\n\"c\": \"3.2\"\n\"e\": \"true\"\n\"f\":\n    - \"1\"\n    - \"2\"\n    - \"3\"\n\"g\":\n    \"something\": \"cool\"\n",
		},
	},
	{
		skipDoc:    true,
		document:   "bing: &foo {x: z}\na:\n  c: cat\n  <<: [*foo]",
		expression: `(... | select(tag=="!!str")) style="single"`,
		expected: []string{
			"D0, P[], (!!map)::'bing': &foo {'x': 'z'}\n'a':\n    'c': 'cat'\n    !!merge <<: [*foo]\n",
		},
	},
	{
		description: "Set single quote style",
		document:    `{a: cat, b: 5, c: 3.2, e: true, f: [1,2,3], g: { something: cool}}`,
		expression:  `.. style="single"`,
		expected: []string{
			"D0, P[], (!!map)::a: 'cat'\nb: '5'\nc: '3.2'\ne: 'true'\nf:\n    - '1'\n    - '2'\n    - '3'\ng:\n    something: 'cool'\n",
		},
	},
	{
		description: "Set literal quote style",
		document:    `{a: cat, b: 5, c: 3.2, e: true, f: [1,2,3], g: { something: cool}}`,
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
f:
    - |-
      1
    - |-
      2
    - |-
      3
g:
    something: |-
        cool
`,
		},
	},
	{
		description: "Set folded quote style",
		document:    `{a: cat, b: 5, c: 3.2, e: true, f: [1,2,3], g: { something: cool}}`,
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
f:
    - >-
      1
    - >-
      2
    - >-
      3
g:
    something: >-
        cool
`,
		},
	},
	{
		description: "Set flow quote style",
		document:    `{a: cat, b: 5, c: 3.2, e: true, f: [1,2,3], g: { something: cool}}`,
		expression:  `.. style="flow"`,
		expected: []string{
			"D0, P[], (!!map)::{a: cat, b: 5, c: 3.2, e: true, f: [1, 2, 3], g: {something: cool}}\n",
		},
	},
	{
		description:           "Reset style - or pretty print",
		subdescription:        "Set empty (default) quote style, note the usage of `...` to match keys too. Note that there is a `--prettyPrint/-P` short flag for this.",
		dontFormatInputForDoc: true,
		document:              `{a: cat, "b": 5, 'c': 3.2, "e": true,  f: [1,2,3], "g": { something: "cool"} }`,
		expression:            `... style=""`,
		expected: []string{
			"D0, P[], (!!map)::a: cat\nb: 5\nc: 3.2\ne: true\nf:\n    - 1\n    - 2\n    - 3\ng:\n    something: cool\n",
		},
	},
	{
		description: "Set style relatively with assign-update",
		document:    `{a: single, b: double}`,
		expression:  `.[] style |= .`,
		expected: []string{
			"D0, P[], (!!map)::{a: 'single', b: \"double\"}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: cat, b: double}`,
		expression: `.a style=.b`,
		expected: []string{
			"D0, P[], (!!map)::{a: \"cat\", b: double}\n",
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
			"D0, P[], (!!str)::\n",
			"D0, P[a], (!!str)::\n",
		},
	},
}

func testStyleScenarioWithParserCheck(t *testing.T, s *expressionScenario) {
	// Skip tests where goccy correctly rejects invalid YAML at parse time
	if ConfiguredYamlPreferences.UseGoccyParser {
		// Check if the document contains merge anchor with sequence (invalid YAML)
		if s.document == "bing: &foo frog\na:\n  c: cat\n  <<: [*foo]" {
			t.Skip("goccy parser correctly rejects merge anchors with sequences at parse time")
			return
		}
	}
	testScenario(t, s)
}

func TestStyleOperatorScenarios(t *testing.T) {
	for _, tt := range styleOperatorScenarios {
		testStyleScenarioWithParserCheck(t, &tt)
	}
	documentOperatorScenarios(t, "style", styleOperatorScenarios)
}
