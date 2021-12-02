package yqlib

import (
	"testing"
)

var prefix = "D0, P[], (doc)::a:\n    cool:\n        bob: dylan\n"

var encoderDecoderOperatorScenarios = []expressionScenario{
	{
		description: "Encode value as json string",
		document:    `{a: {cool: "thing"}}`,
		expression:  `.b = (.a | to_json)`,
		expected: []string{
			`D0, P[], (doc)::{a: {cool: "thing"}, b: "{\n  \"cool\": \"thing\"\n}\n"}
`,
		},
	},
	{
		description:    "Encode value as json string, on one line",
		subdescription: "Pass in a 0 indent to print json on a single line.",
		document:       `{a: {cool: "thing"}}`,
		expression:     `.b = (.a | to_json(0))`,
		expected: []string{
			`D0, P[], (doc)::{a: {cool: "thing"}, b: '{"cool":"thing"}'}
`,
		},
	},
	{
		description:    "Encode value as json string, on one line shorthand",
		subdescription: "Pass in a 0 indent to print json on a single line.",
		document:       `{a: {cool: "thing"}}`,
		expression:     `.b = (.a | @json)`,
		expected: []string{
			`D0, P[], (doc)::{a: {cool: "thing"}, b: '{"cool":"thing"}'}
`,
		},
	},
	{
		description:    "Decode a json encoded string",
		subdescription: "Keep in mind JSON is a subset of YAML. If you want idiomatic yaml, pipe through the style operator to clear out the JSON styling.",
		document:       `a: '{"cool":"thing"}'`,
		expression:     `.a | from_json | ... style=""`,
		expected: []string{
			"D0, P[a], (!!map)::cool: thing\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {cool: "thing"}}`,
		expression: `.b = (.a | to_props)`,
		expected: []string{
			`D0, P[], (doc)::{a: {cool: "thing"}, b: "cool = thing\n"}
`,
		},
	},
	{
		description: "Encode value as props string",
		document:    `{a: {cool: "thing"}}`,
		expression:  `.b = (.a | @props)`,
		expected: []string{
			`D0, P[], (doc)::{a: {cool: "thing"}, b: "cool = thing\n"}
`,
		},
	},
	{
		skipDoc:    true,
		document:   "a:\n  cool:\n    bob: dylan",
		expression: `.b = (.a | @yaml)`,
		expected: []string{
			prefix + "b: |\n    cool:\n      bob: dylan\n",
		},
	},
	{
		description:    "Encode value as yaml string",
		subdescription: "Indent defaults to 2",
		document:       "a:\n  cool:\n    bob: dylan",
		expression:     `.b = (.a | to_yaml)`,
		expected: []string{
			prefix + "b: |\n    cool:\n      bob: dylan\n",
		},
	},
	{
		description:    "Encode value as yaml string, with custom indentation",
		subdescription: "You can specify the indentation level as the first parameter.",
		document:       "a:\n  cool:\n    bob: dylan",
		expression:     `.b = (.a | to_yaml(8))`,
		expected: []string{
			prefix + "b: |\n    cool:\n            bob: dylan\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {cool: "thing"}}`,
		expression: `.b = (.a | to_yaml)`,
		expected: []string{
			`D0, P[], (doc)::{a: {cool: "thing"}, b: "{cool: \"thing\"}\n"}
`,
		},
	},
	{
		description: "Decode a yaml encoded string",
		document:    `a: "foo: bar"`,
		expression:  `.b = (.a | from_yaml)`,
		expected: []string{
			"D0, P[], (doc)::a: \"foo: bar\"\nb:\n    foo: bar\n",
		},
	},
	{
		description:           "Update a multiline encoded yaml string",
		dontFormatInputForDoc: true,
		document:              "a: |\n  foo: bar\n  baz: dog\n",
		expression:            `.a |= (from_yaml | .foo = "cat" | to_yaml)`,
		expected: []string{
			"D0, P[], (doc)::a: |\n    foo: cat\n    baz: dog\n",
		},
	},
	{
		skipDoc:               true,
		dontFormatInputForDoc: true,
		document:              "a: |-\n  foo: bar\n  baz: dog\n",
		expression:            `.a |= (from_yaml | .foo = "cat" | to_yaml)`,
		expected: []string{
			"D0, P[], (doc)::a: |-\n    foo: cat\n    baz: dog\n",
		},
	},
	{
		description:           "Update a single line encoded yaml string",
		dontFormatInputForDoc: true,
		document:              "a: 'foo: bar'",
		expression:            `.a |= (from_yaml | .foo = "cat" | to_yaml)`,
		expected: []string{
			"D0, P[], (doc)::a: 'foo: cat'\n",
		},
	},
	{
		description:    "Encode array of scalars as csv string",
		subdescription: "Scalars are strings, numbers and booleans.",
		document:       `[cat, "thing1,thing2", true, 3.40]`,
		expression:     `@csv`,
		expected: []string{
			"D0, P[], (!!str)::cat,\"thing1,thing2\",true,3.40\n",
		},
	},
	{
		description: "Encode array of arrays as csv string",
		document:    `[[cat, "thing1,thing2", true, 3.40], [dog, thing3, false, 12]]`,
		expression:  `@csv`,
		expected: []string{
			"D0, P[], (!!str)::cat,\"thing1,thing2\",true,3.40\ndog,thing3,false,12\n",
		},
	},
	{
		description:    "Encode array of array scalars as tsv string",
		subdescription: "Scalars are strings, numbers and booleans.",
		document:       `[[cat, "thing1,thing2", true, 3.40], [dog, thing3, false, 12]]`,
		expression:     `@tsv`,
		expected: []string{
			"D0, P[], (!!str)::cat\tthing1,thing2\ttrue\t3.40\ndog\tthing3\tfalse\t12\n",
		},
	},
	{
		skipDoc:               true,
		dontFormatInputForDoc: true,
		document:              "a: \"foo: bar\"",
		expression:            `.a |= (from_yaml | .foo = {"a": "frog"} | to_yaml)`,
		expected: []string{
			"D0, P[], (doc)::a: \"foo:\\n  a: frog\"\n",
		},
	},
}

func TestEncoderDecoderOperatorScenarios(t *testing.T) {
	for _, tt := range encoderDecoderOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "encode-decode", encoderDecoderOperatorScenarios)
}
