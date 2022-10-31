package yqlib

import "testing"

var containsOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		expression: `null | contains(~)`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		expression: `3 | contains(3)`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		expression: `3 | contains(32)`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		description:    "Array contains array",
		subdescription: "Array is equal or subset of",
		document:       `["foobar", "foobaz", "blarp"]`,
		expression:     `contains(["baz", "bar"])`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		description:    "Array has a subset array",
		subdescription: "Subtract the superset array from the subset, if there's anything left, it's not a subset",
		document:       `["foobar", "foobaz", "blarp"]`,
		expression:     `["baz", "bar"] - . | length == 0`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		expression: `["dog", "cat", "giraffe"] | contains(["camel"])`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		description: "Object included in array",
		document:    `{"foo": 12, "bar":[1,2,{"barp":12, "blip":13}]}`,
		expression:  `contains({"bar": [{"barp": 12}]})`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		description: "Object not included in array",
		document:    `{"foo": 12, "bar":[1,2,{"barp":12, "blip":13}]}`,
		expression:  `contains({"foo": 12, "bar": [{"barp": 15}]})`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		description: "String contains substring",
		document:    `"foobar"`,
		expression:  `contains("bar")`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		description: "String equals string",
		document:    `"meow"`,
		expression:  `contains("meow")`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
}

func TestContainsOperatorScenarios(t *testing.T) {
	for _, tt := range containsOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "contains", containsOperatorScenarios)
}
