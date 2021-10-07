package yqlib

import (
	"testing"
)

var selectOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		document:   `cat`,
		expression: `select(false, true)`,
		expected: []string{
			"D0, P[], (doc)::cat\n",
		},
	},
	{
		skipDoc:    true,
		document:   `cat`,
		expression: `select(true, false)`,
		expected: []string{
			"D0, P[], (doc)::cat\n",
		},
	},
	{
		skipDoc:    true,
		document:   `cat`,
		expression: `select(false)`,
		expected:   []string{},
	},
	{
		description: "Select elements from array",
		document:    `[cat,goat,dog]`,
		expression:  `.[] | select(. == "*at")`,
		expected: []string{
			"D0, P[0], (!!str)::cat\n",
			"D0, P[1], (!!str)::goat\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: hello",
		document2:  "b: world",
		expression: `select(.a == "hello" or .b == "world")`,
		expected: []string{
			"D0, P[], (doc)::a: hello\n",
			"D0, P[], (doc)::b: world\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[{animal: cat, legs: {cool: true}}, {animal: fish}]`,
		expression: `(.[] | select(.legs.cool == true).canWalk) = true | (.[] | .alive.things) = "yes"`,
		expected: []string{
			"D0, P[], (doc)::[{animal: cat, legs: {cool: true}, canWalk: true, alive: {things: yes}}, {animal: fish, alive: {things: yes}}]\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[hot, fot, dog]`,
		expression: `.[] | select(. == "*at")`,
		expected:   []string{},
	},
	{
		skipDoc:    true,
		document:   `a: [cat,goat,dog]`,
		expression: `.a.[] | select(. == "*at")`,
		expected: []string{
			"D0, P[a 0], (!!str)::cat\n",
			"D0, P[a 1], (!!str)::goat\n"},
	},
	{
		description: "Select and update matching values in map",
		document:    `a: { things: cat, bob: goat, horse: dog }`,
		expression:  `(.a.[] | select(. == "*at")) |= "rabbit"`,
		expected: []string{
			"D0, P[], (doc)::a: {things: rabbit, bob: rabbit, horse: dog}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: { things: {include: true}, notMe: {include: false}, andMe: {include: fold} }`,
		expression: `.a.[] | select(.include)`,
		expected: []string{
			"D0, P[a things], (!!map)::{include: true}\n",
			"D0, P[a andMe], (!!map)::{include: fold}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[cat,~,dog]`,
		expression: `.[] | select(. == ~)`,
		expected: []string{
			"D0, P[1], (!!null)::~\n",
		},
	},
}

func TestSelectOperatorScenarios(t *testing.T) {
	for _, tt := range selectOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Select", selectOperatorScenarios)
}
