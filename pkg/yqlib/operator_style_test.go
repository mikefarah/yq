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
		skipForGoccy: true,
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
	// Flow map test scenarios demonstrating parser differences
	{
		description:    "Basic flow map display - legacy-v3 behaviour",
		subdescription: "legacy-v3 preserves original flow map syntax in output",
		document:       `{a: 1, b: 2}`,
		expression:     `.`,
		expected: []string{
			"D0, P[], (!!map)::{a: 1, b: 2}\n",
		},
		skipForGoccy: true,
	},
	{
		description:    "Basic flow map display - goccy behaviour",
		subdescription: "Goccy normalizes flow maps to block style in output",
		document:       `{a: 1, b: 2}`,
		expression:     `.`,
		expected: []string{
			"D0, P[], (!!map)::a: 1\nb: 2\n",
		},
		skipForYamlV3: true,
	},
	{
		description:    "Flow maps in arrays - legacy-v3 behaviour",
		subdescription: "legacy-v3 maintains flow syntax for maps within arrays",
		document:       `items: [{name: item1, value: 100}, {name: item2, value: 200}]`,
		expression:     `.`,
		expected: []string{
			"D0, P[], (!!map)::items: [{name: item1, value: 100}, {name: item2, value: 200}]\n",
		},
		skipForGoccy: true,
	},
	{
		description:    "Flow maps in arrays - goccy behaviour",
		subdescription: "Goccy converts flow maps in arrays to block style",
		document:       `items: [{name: item1, value: 100}, {name: item2, value: 200}]`,
		expression:     `.`,
		expected: []string{
			"D0, P[], (!!map)::items:\n- name: item1\n  value: 100\n- name: item2\n  value: 200\n",
		},
		skipForYamlV3: true,
	},
	{
		description:    "Nested flow maps - legacy-v3 behaviour",
		subdescription: "legacy-v3 preserves nested flow map structure",
		document:       `config: {database: {host: localhost, port: 5432}, cache: {enabled: true}}`,
		expression:     `.`,
		expected: []string{
			"D0, P[], (!!map)::config: {database: {host: localhost, port: 5432}, cache: {enabled: true}}\n",
		},
		skipForGoccy: true,
	},
	{
		description:    "Nested flow maps - goccy behaviour",
		subdescription: "Goccy normalizes nested flow maps to block style",
		document:       `config: {database: {host: localhost, port: 5432}, cache: {enabled: true}}`,
		expression:     `.`,
		expected: []string{
			"D0, P[], (!!map)::config:\n  cache:\n    enabled: true\n  database:\n    host: localhost\n    port: 5432\n",
		},
		skipForYamlV3: true,
	},
	{
		description:    "Mixed flow and block styles - legacy-v3 behaviour",
		subdescription: "legacy-v3 preserves original mix of flow and block styles",
		document:       "data:\n  users: [{id: 1, name: \"Alice\"}, {id: 2, name: \"Bob\"}]\n  settings: {theme: dark, lang: en}",
		expression:     `.`,
		expected: []string{
			"D0, P[], (!!map)::data:\n    users: [{id: 1, name: \"Alice\"}, {id: 2, name: \"Bob\"}]\n    settings: {theme: dark, lang: en}\n",
		},
		skipForGoccy: true,
	},
	{
		description:    "Mixed flow and block styles - goccy behaviour",
		subdescription: "Goccy normalizes all to block style regardless of original format",
		document:       "data:\n  users: [{id: 1, name: \"Alice\"}, {id: 2, name: \"Bob\"}]\n  settings: {theme: dark, lang: en}",
		expression:     `.`,
		expected: []string{
			"D0, P[], (!!map)::data:\n  settings:\n    lang: en\n    theme: dark\n  users:\n  - id: 1\n    name: Alice\n  - id: 2\n    name: Bob\n",
		},
		skipForYamlV3: true,
	},
	{
		description:    "Empty flow maps - legacy-v3 behaviour",
		subdescription: "legacy-v3 preserves empty flow map syntax",
		document:       `{empty: {}, data: {key: value}}`,
		expression:     `.`,
		expected: []string{
			"D0, P[], (!!map)::{empty: {}, data: {key: value}}\n",
		},
		skipForGoccy: true,
	},
	{
		description:    "Empty flow maps - goccy behaviour",
		subdescription: "Goccy normalizes empty flow maps to block style",
		document:       `{empty: {}, data: {key: value}}`,
		expression:     `.`,
		expected: []string{
			"D0, P[], (!!map)::data:\n    key: value\nempty: {}\n",
		},
		skipForYamlV3: true,
	},
	{
		description:    "Flow map content access - both parsers",
		subdescription: "Both parsers access flow map content identically despite display differences",
		document:       `{a: 1, b: 2, c: 3}`,
		expression:     `.b`,
		expected: []string{
			"D0, P[b], (!!int)::2\n",
		},
	},
	{
		description:    "Complex flow map navigation - both parsers",
		subdescription: "Both parsers navigate complex flow maps identically",
		document:       `config: {database: {host: localhost, port: 5432}, cache: {enabled: true}}`,
		expression:     `.config.database.host`,
		expected: []string{
			"D0, P[config database host], (!!str)::localhost\n",
		},
	},
	{
		description:    "Flow map style preservation - forced flow style",
		subdescription: "Both parsers can output in flow style when explicitly requested",
		document:       `a: 1\nb: 2`,
		expression:     `. style="flow"`,
		expected: []string{
			"D0, P[], (!!map)::{a: 1, b: 2}\n",
		},
	},
}

func TestStyleOperatorScenarios(t *testing.T) {
	for _, tt := range styleOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "style", styleOperatorScenarios)
}
