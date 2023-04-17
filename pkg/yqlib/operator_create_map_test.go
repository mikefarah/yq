package yqlib

import (
	"testing"
)

var createMapOperatorScenarios = []expressionScenario{
	{
		document:   ``,
		expression: `"frog": "jumps"`,
		expected: []string{
			"D0, P[], (!!seq)::- [{frog: jumps}]\n",
		},
	},
	{
		skipDoc:     true,
		description: "sets key properly",
		expression:  `("frog": "jumps") | .[0][0] | .frog`,
		expected: []string{
			"D0, P[frog], (!!str)::jumps\n",
		},
	},
	{
		skipDoc:     true,
		description: "sets key properly on map",
		expression:  `{"frog": "jumps"} | .frog`,
		expected: []string{
			"D0, P[frog], (!!str)::jumps\n",
		},
	},
	{
		document:   `{name: Mike, pets: [cat, dog]}`,
		expression: `(.name: .pets.[]) | .[0][0] | ..`,
		expected: []string{
			"D0, P[], (!!map)::Mike: cat\n",
			"D0, P[Mike], (!!str)::cat\n",
		},
	},
	{
		description: "check path of nested child",
		document:    "pets:\n  cows: value",
		expression:  `("b":.pets) | .[0][0] | .b.cows`,
		expected: []string{
			"D0, P[b cows], (!!str)::value\n",
		},
	},
	{
		document:   `{name: Mike, age: 32}`,
		expression: `.name: .age`,
		expected: []string{
			"D0, P[], (!!seq)::- [{Mike: 32}]\n",
		},
	},
	{
		document:   `{name: Mike, pets: [cat, dog]}`,
		expression: `.name: .pets.[]`,
		expected: []string{
			"D0, P[], (!!seq)::- [{Mike: cat}, {Mike: dog}]\n",
		},
	},
	{
		document:   `{name: Mike, pets: [cat, dog], food: [hotdog, burger]}`,
		expression: `.name: .pets.[], "f":.food.[]`,
		expected: []string{
			"D0, P[], (!!seq)::- [{Mike: cat}, {Mike: dog}]\n",
			"D0, P[], (!!seq)::- [{f: hotdog}, {f: burger}]\n",
		},
	},
	{
		document:   "{name: Mike, pets: [cat, dog], food: [hotdog, burger]}\n---\n{name: Fred, pets: [mouse], food: [pizza, onion, apple]}",
		expression: `.name: .pets.[], "f":.food.[]`,
		expected: []string{
			"D0, P[], (!!seq)::- [{Mike: cat}, {Mike: dog}]\n- [{Fred: mouse}]\n",
			"D0, P[], (!!seq)::- [{f: hotdog}, {f: burger}]\n- [{f: pizza}, {f: onion}, {f: apple}]\n",
		},
	},
	{
		document:   `{name: Mike, pets: {cows: [apl, bba]}}`,
		expression: `"a":.name, "b":.pets`,
		expected: []string{
			"D0, P[], (!!seq)::- [{a: Mike}]\n",
			"D0, P[], (!!seq)::- [{b: {cows: [apl, bba]}}]\n",
		},
	},
	{
		document:   `{name: Mike}`,
		expression: `"wrap": .`,
		expected: []string{
			"D0, P[], (!!seq)::- [{wrap: {name: Mike}}]\n",
		},
	},
	{
		document:   "{name: Mike}\n---\n{name: Bob}",
		expression: `"wrap": .`,
		expected: []string{
			"D0, P[], (!!seq)::- [{wrap: {name: Mike}}]\n- [{wrap: {name: Bob}}]\n",
		},
	},
	{
		document:   "{name: Mike}\n---\n{name: Bob}",
		expression: `"wrap": ., .name: "great"`,
		expected: []string{
			"D0, P[], (!!seq)::- [{wrap: {name: Mike}}]\n- [{wrap: {name: Bob}}]\n",
			"D0, P[], (!!seq)::- [{Mike: great}]\n- [{Bob: great}]\n",
		},
	},
}

func TestCreateMapOperatorScenarios(t *testing.T) {
	for _, tt := range createMapOperatorScenarios {
		testScenario(t, &tt)
	}
}
