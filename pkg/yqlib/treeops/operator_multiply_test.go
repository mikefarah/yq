package treeops

import (
	"testing"
)

var multiplyOperatorScenarios = []expressionScenario{
	{
		document:   `{a: {also: [1]}, b: {also: me}}`,
		expression: `.a * .b`,
		expected: []string{
			"D0, P[], (!!map)::{a: {also: me}, b: {also: me}}\n",
		},
	}, {
		document:   `{a: {also: me}, b: {also: [1]}}`,
		expression: `.a * .b`,
		expected: []string{
			"D0, P[], (!!map)::{a: {also: [1]}, b: {also: [1]}}\n",
		},
	}, {
		document:   `{a: {also: me}, b: {also: {g: wizz}}}`,
		expression: `.a * .b`,
		expected: []string{
			"D0, P[], (!!map)::{a: {also: {g: wizz}}, b: {also: {g: wizz}}}\n",
		},
	}, {
		document:   `{a: {also: {g: wizz}}, b: {also: me}}`,
		expression: `.a * .b`,
		expected: []string{
			"D0, P[], (!!map)::{a: {also: me}, b: {also: me}}\n",
		},
	}, {
		document:   `{a: {also: {g: wizz}}, b: {also: [1]}}`,
		expression: `.a * .b`,
		expected: []string{
			"D0, P[], (!!map)::{a: {also: [1]}, b: {also: [1]}}\n",
		},
	}, {
		document:   `{a: {also: [1]}, b: {also: {g: wizz}}}`,
		expression: `.a * .b`,
		expected: []string{
			"D0, P[], (!!map)::{a: {also: {g: wizz}}, b: {also: {g: wizz}}}\n",
		},
	}, {
		document:   `{a: {things: great}, b: {also: me}}`,
		expression: `.a * .b`,
		expected: []string{
			"D0, P[], (!!map)::{a: {things: great, also: me}, b: {also: me}}\n",
		},
	}, {
		document: `a: {things: great}
b:
  also: "me"
`,
		expression: `.a * .b`,
		expected: []string{
			`D0, P[], (!!map)::a:
    things: great
    also: "me"
b:
    also: "me"
`,
		},
	},
}

func TestMultiplyOperatorScenarios(t *testing.T) {
	for _, tt := range multiplyOperatorScenarios {
		testScenario(t, &tt)
	}
}
