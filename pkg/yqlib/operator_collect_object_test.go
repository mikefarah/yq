package yqlib

import (
	"testing"
)

var collectObjectOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		expression: `{"name": "mike"} | .name`,
		expected: []string{
			"D0, P[name], (!!str)::mike\n",
		},
	},
	{
		skipDoc:       true,
		expression:    `{"c": "a", "b", "d"}`,
		expectedError: "CollectObject: mismatching node sizes; are you creating a map with mismatching key value pairs?",
	},
	{
		skipDoc:    true,
		expression: `{"person": {"names": ["mike"]}} | .person.names[0]`,
		expected: []string{
			"D0, P[person names 0], (!!str)::mike\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[{name: cat}, {name: dog}]`,
		expression: `.[] | {.name: "great"}`,
		expected: []string{
			"D0, P[], (!!map)::cat: great\n",
			"D0, P[], (!!map)::dog: great\n",
		},
	},
	{
		description: "collect splat",
		skipDoc:     true,
		document:    `[{name: cat}, {name: dog}]`,
		expression:  `.[] | {.name: "great"}[]`,
		expected: []string{
			"D0, P[cat], (!!str)::great\n",
			"D0, P[dog], (!!str)::great\n",
		},
	},
	{
		skipDoc:    true,
		expression: `({} + {}) | (.b = 3)`,
		expected: []string{
			"D0, P[], (!!map)::b: 3\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: []",
		expression: `.a += [{"key": "att2", "value": "val2"}]`,
		expected: []string{
			"D0, P[], (!!map)::a:\n    - key: att2\n      value: val2\n",
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
		skipDoc:     true,
		description: "Two documents",
		document:    "{name: Mike}\n",
		document2:   "{name: Bob}\n",
		expression:  `{"wrap": .}`,
		expected: []string{
			"D0, P[], (!!map)::wrap: {name: Mike}\n",
			"D0, P[], (!!map)::wrap: {name: Bob}\n",
		},
	},
	{
		skipDoc:     true,
		description: "two embedded documents",
		document:    "{name: Mike}\n---\n{name: Bob}",
		expression:  `{"wrap": .}`,
		expected: []string{
			"D0, P[], (!!map)::wrap: {name: Mike}\n",
			"D1, P[], (!!map)::wrap: {name: Bob}\n",
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
			"D1, P[], (!!map)::Rosey: monkey\n",
			"D1, P[], (!!map)::Rosey: sheep\n",
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
		document:   "name: Mike\npets:\n  cows:\n    - apl\n    - bba",
		document2:  "name: Rosey\npets:\n  sheep:\n    - frog\n    - meow",
		expression: `{"a":.name, "b":.pets}`,
		expected: []string{
			"D0, P[], (!!map)::a: Mike\nb:\n    cows:\n        - apl\n        - bba\n",
			"D0, P[], (!!map)::a: Rosey\nb:\n    sheep:\n        - frog\n        - meow\n",
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
		expression: `{"wrap": "frog", "bing": "bong"}`,
		expected: []string{
			"D0, P[], (!!map)::wrap: frog\nbing: bong\n",
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
	{
		description: "Creating yaml from scratch with multiple objects",
		expression:  `(.a.b = "foo") | (.d.e = "bar")`,
		expected: []string{
			"D0, P[], ()::a:\n    b: foo\nd:\n    e: bar\n",
		},
	},
}

func TestCollectObjectOperatorScenarios(t *testing.T) {
	for _, tt := range collectObjectOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "create-collect-into-object", collectObjectOperatorScenarios)
}
