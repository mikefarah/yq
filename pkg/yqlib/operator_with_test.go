package yqlib

import "testing"

var withOperatorScenarios = []expressionScenario{
	{
		description: "Update and style",
		document:    `a: {deeply: {nested: value}}`,
		expression:  `with(.a.deeply.nested ; . = "newValue" | . style="single")`,
		expected: []string{
			"D0, P[], (doc)::a: {deeply: {nested: 'newValue'}}\n",
		},
	},
	{
		description: "Update multiple deeply nested properties",
		document:    `a: {deeply: {nested: value, other: thing}}`,
		expression:  `with(.a.deeply ; .nested = "newValue" | .other= "newThing")`,
		expected: []string{
			"D0, P[], (doc)::a: {deeply: {nested: newValue, other: newThing}}\n",
		},
	},
	{
		description: "Update array elements relatively",
		document:    `myArray: [{a: apple},{a: banana}]`,
		expression:  `with(.myArray[] ; .b = .a + " yum")`,
		expected: []string{
			"D0, P[], (doc)::myArray: [{a: apple, b: apple yum}, {a: banana, b: banana yum}]\n",
		},
	},
}

func TestWithOperatorScenarios(t *testing.T) {
	for _, tt := range withOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "With", withOperatorScenarios)
}
