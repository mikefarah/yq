package yqlib

import (
	"testing"
)

var valueOperatorScenarios = []expressionScenario{
	{
		document:   ``,
		expression: `1`,
		expected: []string{
			"D0, P[], (!!int)::1\n",
		},
	},
	{
		document:   ``,
		expression: `-1`,
		expected: []string{
			"D0, P[], (!!int)::-1\n",
		},
	}, {
		document:   ``,
		expression: `1.2`,
		expected: []string{
			"D0, P[], (!!float)::1.2\n",
		},
	}, {
		document:   ``,
		expression: `-5.2e11`,
		expected: []string{
			"D0, P[], (!!float)::-5.2e11\n",
		},
	}, {
		document:   ``,
		expression: `5e-10`,
		expected: []string{
			"D0, P[], (!!float)::5e-10\n",
		},
	},
	{
		document:   ``,
		expression: `"cat"`,
		expected: []string{
			"D0, P[], (!!str)::cat\n",
		},
	},
	{
		document:   ``,
		expression: `"frog jumps"`,
		expected: []string{
			"D0, P[], (!!str)::frog jumps\n",
		},
	},
	{
		document:   ``,
		expression: `"1.3"`,
		expected: []string{
			"D0, P[], (!!str)::\"1.3\"\n",
		},
	}, {
		document:   ``,
		expression: `"true"`,
		expected: []string{
			"D0, P[], (!!str)::\"true\"\n",
		},
	}, {
		document:   ``,
		expression: `true`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	}, {
		document:   ``,
		expression: `false`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		document:   ``,
		expression: `Null`,
		expected: []string{
			"D0, P[], (!!null)::Null\n",
		},
	},
	{
		document:   ``,
		expression: `~`,
		expected: []string{
			"D0, P[], (!!null)::~\n",
		},
	},
}

func TestValueOperatorScenarios(t *testing.T) {
	for _, tt := range valueOperatorScenarios {
		testScenario(t, &tt)
	}
}
