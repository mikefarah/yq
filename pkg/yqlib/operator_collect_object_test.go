package yqlib

import (
	"testing"
)

var collectObjectOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		document:   "a: []",
		expression: `.a += [{"key": "att2", "value": "val2"}]`,
		expected: []string{
			"D0, P[], (doc)::a: [{key: att2, value: val2}]\n",
		},
	},
	{
		skipDoc:    true,
		document:   "",
		expression: `.a += {"key": "att2", "value": "val2"}`,
		expected: []string{
			"D0, P[], ()::a:\n    key: att2\n    value: val2\n",
		},
	},
	{
		skipDoc:    true,
		document:   "",
		expression: `.a += [0]`,
		expected: []string{
			"D0, P[], ()::a:\n    - 0\n",
		},
	},
	{
		description: `Collect empty object`,
		document:    ``,
		expression:  `{}`,
		expected: []string{
			"D0, P[], (!!map)::{}\n",
		},
	},
	{
		description: `Wrap (prefix) existing object`,
		document:    "{name: Mike}\n",
		expression:  `{"wrap": .}`,
		expected: []string{
			"D0, P[], (!!map)::wrap: {name: Mike}\n",
		},
	},
	{
		skipDoc:    true,
		document:   "{name: Mike}\n---\n{name: Bob}",
		expression: `{"wrap": .}`,
		expected: []string{
			"D0, P[], (!!map)::wrap: {name: Mike}\n",
			"D0, P[], (!!map)::wrap: {name: Bob}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{name: Mike, age: 32}`,
		expression: `{.name: .age}`,
		expected: []string{
			"D0, P[], (!!map)::Mike: 32\n",
		},
	},
	{
		description: `Using splat to create multiple objects`,
		document:    `{name: Mike, pets: [cat, dog]}`,
		expression:  `{.name: .pets.[]}`,
		expected: []string{
			"D0, P[], (!!map)::Mike: cat\n",
			"D0, P[], (!!map)::Mike: dog\n",
		},
	},
	{
		description:           `Working with multiple documents`,
		dontFormatInputForDoc: false,
		document:              "{name: Mike, pets: [cat, dog]}\n---\n{name: Rosey, pets: [monkey, sheep]}",
		expression:            `{.name: .pets.[]}`,
		expected: []string{
			"D0, P[], (!!map)::Mike: cat\n",
			"D0, P[], (!!map)::Mike: dog\n",
			"D0, P[], (!!map)::Rosey: monkey\n",
			"D0, P[], (!!map)::Rosey: sheep\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{name: Mike, pets: [cat, dog], food: [hotdog, burger]}`,
		expression: `{.name: .pets.[], "f":.food.[]}`,
		expected: []string{
			"D0, P[], (!!map)::Mike: cat\nf: hotdog\n",
			"D0, P[], (!!map)::Mike: cat\nf: burger\n",
			"D0, P[], (!!map)::Mike: dog\nf: hotdog\n",
			"D0, P[], (!!map)::Mike: dog\nf: burger\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{name: Mike, pets: {cows: [apl, bba]}}`,
		expression: `{"a":.name, "b":.pets}`,
		expected: []string{
			`D0, P[], (!!map)::a: Mike
b: {cows: [apl, bba]}
`,
		},
	},
	{
		description: "Creating yaml from scratch",
		document:    ``,
		expression:  `{"wrap": "frog"}`,
		expected: []string{
			"D0, P[], (!!map)::wrap: frog\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{name: Mike}`,
		expression: `{"wrap": .}`,
		expected: []string{
			"D0, P[], (!!map)::wrap: {name: Mike}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{name: Mike}`,
		expression: `{"wrap": {"further": .}} | (.. style= "flow")`,
		expected: []string{
			"D0, P[], (!!map)::{wrap: {further: {name: Mike}}}\n",
		},
	},
}

func TestCollectObjectOperatorScenarios(t *testing.T) {
	for _, tt := range collectObjectOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Create, Collect into Object", collectObjectOperatorScenarios)
}
