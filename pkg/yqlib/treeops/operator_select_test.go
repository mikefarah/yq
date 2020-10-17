package treeops

import (
	"testing"
)

var selectOperatorScenarios = []expressionScenario{
	{
		document:   `[cat,goat,dog]`,
		expression: `.[] | select(. == "*at")`,
		expected: []string{
			"D0, P[0], (!!str)::cat\n",
			"D0, P[1], (!!str)::goat\n",
		},
	}, {
		document:   `[hot, fot, dog]`,
		expression: `.[] | select(. == "*at")`,
		expected:   []string{},
	}, {
		document:   `a: [cat,goat,dog]`,
		expression: `.a[] | select(. == "*at")`,
		expected: []string{
			"D0, P[a 0], (!!str)::cat\n",
			"D0, P[a 1], (!!str)::goat\n"},
	}, {
		document:   `a: { things: cat, bob: goat, horse: dog }`,
		expression: `.a[] | select(. == "*at")`,
		expected: []string{
			"D0, P[a things], (!!str)::cat\n",
			"D0, P[a bob], (!!str)::goat\n"},
	}, {
		document:   `a: { things: {include: true}, notMe: {include: false}, andMe: {include: fold} }`,
		expression: `.a[] | select(.include)`,
		expected: []string{
			"D0, P[a things], (!!map)::{include: true}\n",
			"D0, P[a andMe], (!!map)::{include: fold}\n",
		},
	},
}

func TestSelectOperatorScenarios(t *testing.T) {
	for _, tt := range selectOperatorScenarios {
		testScenario(t, &tt)
	}
}
