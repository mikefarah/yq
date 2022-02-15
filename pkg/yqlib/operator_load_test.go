package yqlib

import (
	"testing"
)

var loadScenarios = []expressionScenario{
	{
		description: "Simple example",
		document:    `{myFile: "../../examples/thing.yml"}`,
		expression:  `load(.myFile)`,
		expected: []string{
			"D0, P[], (doc)::a: apple is included\nb: cool.\n",
		},
	},
	{
		description:    "Replace node with referenced file",
		subdescription: "Note that you can modify the filename in the load operator if needed.",
		document:       `{something: {file: "thing.yml"}}`,
		expression:     `.something |= load("../../examples/" + .file)`,
		expected: []string{
			"D0, P[], (doc)::{something: {a: apple is included, b: cool.}}\n",
		},
	},
	{
		description:    "Replace _all_ nodes with referenced file",
		subdescription: "Recursively match all the nodes (`..`) and then filter the ones that have a 'file' attribute. ",
		document:       `{something: {file: "thing.yml"}, over: {here: [{file: "thing.yml"}]}}`,
		expression:     `(.. | select(has("file"))) |= load("../../examples/" + .file)`,
		expected: []string{
			"D0, P[], (!!map)::{something: {a: apple is included, b: cool.}, over: {here: [{a: apple is included, b: cool.}]}}\n",
		},
	},
	{
		description:    "Replace node with referenced file as string",
		subdescription: "This will work for any text based file",
		document:       `{something: {file: "thing.yml"}}`,
		expression:     `.something |= load_str("../../examples/" + .file)`,
		expected: []string{
			"D0, P[], (doc)::{something: \"a: apple is included\\nb: cool.\"}\n",
		},
	},
	{
		description: "Load from XML",
		document:    "cool: things",
		expression:  `.more_stuff = load_xml("../../examples/small.xml")`,
		expected: []string{
			"D0, P[], (doc)::cool: things\nmore_stuff:\n    this: is some xml\n",
		},
	},
	{
		description: "Load from Properties",
		document:    "cool: things",
		expression:  `.more_stuff = load_props("../../examples/small.properties")`,
		expected: []string{
			"D0, P[], (doc)::cool: things\nmore_stuff:\n    this:\n        is: a properties file\n",
		},
	},
	{
		description:    "Merge from properties",
		subdescription: "This can be used as a convenient way to update a yaml document",
		document:       "this:\n  is: from yaml\n  cool: ay\n",
		expression:     `. *= load_props("../../examples/small.properties")`,
		expected: []string{
			"D0, P[], (!!map)::this:\n    is: a properties file\n    cool: ay\n",
		},
	},
}

func TestLoadScenarios(t *testing.T) {
	for _, tt := range loadScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "load", loadScenarios)
}
