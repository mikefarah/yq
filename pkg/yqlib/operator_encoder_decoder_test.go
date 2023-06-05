package yqlib

import (
	"testing"
)

var prefix = "D0, P[], (!!map)::a:\n    cool:\n        bob: dylan\n"

var encoderDecoderOperatorScenarios = []expressionScenario{
	{
		requiresFormat: "json",
		description:    "Encode value as json string",
		document:       `{a: {cool: "thing"}}`,
		expression:     `.b = (.a | to_json)`,
		expected: []string{
			`D0, P[], (!!map)::{a: {cool: "thing"}, b: "{\n  \"cool\": \"thing\"\n}\n"}
`,
		},
	},
	{
		requiresFormat: "json",
		description:    "Encode value as json string, on one line",
		subdescription: "Pass in a 0 indent to print json on a single line.",
		document:       `{a: {cool: "thing"}}`,
		expression:     `.b = (.a | to_json(0))`,
		expected: []string{
			`D0, P[], (!!map)::{a: {cool: "thing"}, b: '{"cool":"thing"}'}
`,
		},
	},
	{
		requiresFormat: "json",
		description:    "Encode value as json string, on one line shorthand",
		subdescription: "Pass in a 0 indent to print json on a single line.",
		document:       `{a: {cool: "thing"}}`,
		expression:     `.b = (.a | @json)`,
		expected: []string{
			`D0, P[], (!!map)::{a: {cool: "thing"}, b: '{"cool":"thing"}'}
`,
		},
	},
	{
		requiresFormat: "json",
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
			`D0, P[], (!!map)::{a: {cool: "thing"}, b: "cool = thing\n"}
`,
		},
	},
	{
		description: "Encode value as props string",
		document:    `{a: {cool: "thing"}}`,
		expression:  `.b = (.a | @props)`,
		expected: []string{
			`D0, P[], (!!map)::{a: {cool: "thing"}, b: "cool = thing\n"}
`,
		},
	},
	{
		description: "Decode props encoded string",
		document:    `a: "cats=great\ndogs=cool as well"`,
		expression:  `.a |= @propsd`,
		expected: []string{
			"D0, P[], (!!map)::a:\n    cats: great\n    dogs: cool as well\n",
		},
	},
	{
		description: "Decode csv encoded string",
		document:    `a: "cats,dogs\ngreat,cool as well"`,
		expression:  `.a |= @csvd`,
		expected: []string{
			"D0, P[], (!!map)::a:\n    - cats: great\n      dogs: cool as well\n",
		},
	},
	{
		description: "Decode tsv encoded string",
		document:    `a: "cats	dogs\ngreat	cool as well"`,
		expression:  `.a |= @tsvd`,
		expected: []string{
			"D0, P[], (!!map)::a:\n    - cats: great\n      dogs: cool as well\n",
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
			`D0, P[], (!!map)::{a: {cool: "thing"}, b: "{cool: \"thing\"}\n"}
`,
		},
	},
	{
		description: "Decode a yaml encoded string",
		document:    `a: "foo: bar"`,
		expression:  `.b = (.a | from_yaml)`,
		expected: []string{
			"D0, P[], (!!map)::a: \"foo: bar\"\nb:\n    foo: bar\n",
		},
	},
	{
		description:           "Update a multiline encoded yaml string",
		dontFormatInputForDoc: true,
		document:              "a: |\n  foo: bar\n  baz: dog\n",
		expression:            `.a |= (from_yaml | .foo = "cat" | to_yaml)`,
		expected: []string{
			"D0, P[], (!!map)::a: |\n    foo: cat\n    baz: dog\n",
		},
	},
	{
		skipDoc:               true,
		dontFormatInputForDoc: true,
		document:              "a: |-\n  foo: bar\n  baz: dog\n",
		expression:            `.a |= (from_yaml | .foo = "cat" | to_yaml)`,
		expected: []string{
			"D0, P[], (!!map)::a: |-\n    foo: cat\n    baz: dog\n",
		},
	},
	{
		description:           "Update a single line encoded yaml string",
		dontFormatInputForDoc: true,
		document:              "a: 'foo: bar'",
		expression:            `.a |= (from_yaml | .foo = "cat" | to_yaml)`,
		expected: []string{
			"D0, P[], (!!map)::a: 'foo: cat'\n",
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
		description:    "Encode array of arrays as tsv string",
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
			"D0, P[], (!!map)::a: \"foo:\\n  a: frog\"\n",
		},
	},
	{
		requiresFormat: "xml",
		description:    "Encode value as xml string",
		document:       `{a: {cool: {foo: "bar", +@id: hi}}}`,
		expression:     `.a | to_xml`,
		expected: []string{
			"D0, P[a], (!!str)::<cool id=\"hi\">\n  <foo>bar</foo>\n</cool>\n\n",
		},
	},
	{
		requiresFormat: "xml",
		description:    "Encode value as xml string on a single line",
		document:       `{a: {cool: {foo: "bar", +@id: hi}}}`,
		expression:     `.a | @xml`,
		expected: []string{
			"D0, P[a], (!!str)::<cool id=\"hi\"><foo>bar</foo></cool>\n\n",
		},
	},
	{
		requiresFormat: "xml",
		description:    "Encode value as xml string with custom indentation",
		document:       `{a: {cool: {foo: "bar", +@id: hi}}}`,
		expression:     `{"cat": .a | to_xml(1)}`,
		expected: []string{
			"D0, P[], (!!map)::cat: |\n    <cool id=\"hi\">\n     <foo>bar</foo>\n    </cool>\n",
		},
	},
	{
		requiresFormat: "xml",
		description:    "Decode a xml encoded string",
		document:       `a: "<foo>bar</foo>"`,
		expression:     `.b = (.a | from_xml)`,
		expected: []string{
			"D0, P[], (!!map)::a: \"<foo>bar</foo>\"\nb:\n    foo: bar\n",
		},
	},
	{
		description: "Encode a string to base64",
		document:    "coolData: a special string",
		expression:  ".coolData | @base64",
		expected: []string{
			"D0, P[coolData], (!!str)::YSBzcGVjaWFsIHN0cmluZw==\n",
		},
	},
	{
		description:    "Encode a yaml document to base64",
		subdescription: "Pipe through @yaml first to convert to a string, then use @base64 to encode it.",
		document:       "a: apple",
		expression:     "@yaml | @base64",
		expected: []string{
			"D0, P[], (!!str)::YTogYXBwbGUK\n",
		},
	},
	{
		description: "Encode a string to uri",
		document:    "coolData: this has & special () characters *",
		expression:  ".coolData | @uri",
		expected: []string{
			"D0, P[coolData], (!!str)::this+has+%26+special+%28%29+characters+%2A\n",
		},
	},
	{
		description: "Decode a URI to a string",
		document:    "this+has+%26+special+%28%29+characters+%2A",
		expression:  "@urid",
		expected: []string{
			"D0, P[], (!!str)::this has & special () characters *\n",
		},
	},
	{
		description:    "Encode a string to sh",
		subdescription: "Sh/Bash friendly string",
		document:       "coolData: strings with spaces and a 'quote'",
		expression:     ".coolData | @sh",
		expected: []string{
			"D0, P[coolData], (!!str)::strings' with spaces and a '\\'quote\\'\n",
		},
	},
	{
		description:    "Encode a string to sh",
		subdescription: "Watch out for stray '' (empty strings)",
		document:       "coolData: \"'starts, contains more '' and ends with a quote'\"",
		expression:     ".coolData | @sh",
		expected: []string{
			"D0, P[coolData], (!!str)::\\'starts,' contains more '\\'\\'' and ends with a quote'\\'\n",
		},
		skipDoc: true,
	},
	{
		description:    "Decode a base64 encoded string",
		subdescription: "Decoded data is assumed to be a string.",
		document:       "coolData: V29ya3Mgd2l0aCBVVEYtMTYg8J+Yig==",
		expression:     ".coolData | @base64d",
		expected: []string{
			"D0, P[coolData], (!!str)::Works with UTF-16 😊\n",
		},
	},
	{
		description:    "Decode a base64 encoded yaml document",
		subdescription: "Pipe through `from_yaml` to parse the decoded base64 string as a yaml document.",
		document:       "coolData: YTogYXBwbGUK",
		expression:     ".coolData |= (@base64d | from_yaml)",
		expected: []string{
			"D0, P[], (!!map)::coolData:\n    a: apple\n",
		},
	},
	{
		description: "empty base64 decode",
		skipDoc:     true,
		expression:  `"" | @base64d`,
		expected: []string{
			"D0, P[], (!!str)::\n",
		},
	},
	{
		description: "base64 missing padding test",
		skipDoc:     true,
		expression:  `"Y2F0cw" | @base64d`,
		expected: []string{
			"D0, P[], (!!str)::cats\n",
		},
	},
	{
		description: "base64 missing padding test",
		skipDoc:     true,
		expression:  `"cats" | @base64 | @base64d`,
		expected: []string{
			"D0, P[], (!!str)::cats\n",
		},
	},
	{
		requiresFormat: "xml",
		description:    "empty xml decode",
		skipDoc:        true,
		expression:     `"" | @xmld`,
		expected: []string{
			"D0, P[], (!!null)::\n",
		},
	},
}

func TestEncoderDecoderOperatorScenarios(t *testing.T) {
	for _, tt := range encoderDecoderOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "encode-decode", encoderDecoderOperatorScenarios)
}
