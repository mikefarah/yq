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
	documentScenarios(t, "Contains", containsOperatorScenarios)
}
