package yqlib

import "testing"

var withOperatorScenarios = []expressionScenario{
	{
		description: "Update and style",
		document:    `a: {deeply: {nested: value}}`,
		expression:  `with(.a.deeply.nested; . = "newValue" | . style="single")`,
		expected: []string{
			"D0, P[], (doc)::a: {deeply: {nested: 'newValue'}}\n",
		},
	},
	{
		description: "Update multiple deeply nested properties",
		document:    `a: {deeply: {nested: value, other: thing}}`,
		expression:  `with(.a.deeply; .nested = "newValue" | .other= "newThing")`,
		expected: []string{
			"D0, P[], (doc)::a: {deeply: {nested: newValue, other: newThing}}\n",
		},
	},
	{
		description:    "Update array elements relatively",
		subdescription: "The second expression runs with each element of the array as it's contextual root. This allows you to make updates relative to the element.",
		document:       `myArray: [{a: apple},{a: banana}]`,
		expression:     `with(.myArray[]; .b = .a + " yum")`,
		expected: []string{
			"D0, P[], (doc)::myArray: [{a: apple, b: apple yum}, {a: banana, b: banana yum}]\n",
		},
	},
	{
		description:    "Update array elements relatively +=",
		skipDoc:        true,
		subdescription: "The second expression runs with each element of the array as it's contextual root. This allows you to make updates relative to the element.",
		document:       `myArray: [{a: apple},{a: banana}]`,
		expression:     `with(.myArray[]; .a += .a)`,
		expected: []string{
			"D0, P[], (doc)::myArray: [{a: appleapple}, {a: bananabanana}]\n",
		},
	},
}

func TestWithOperatorScenarios(t *testing.T) {
	for _, tt := range withOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "with", withOperatorScenarios)
}
