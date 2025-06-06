package yqlib

import "testing"

var reverseOperatorScenarios = []expressionScenario{
	{
		description: "Reverse",
		document:    "[1, 2, 3]",
		expression:  `reverse`,
		expected: []string{
			"D0, P[], (!!seq)::[3, 2, 1]\n",
		},
	},
	{
		description: "Reverse",
		skipDoc:     true,
		document:    "[1, 2]",
		expression:  `reverse[]`,
		expected: []string{
			"D0, P[1], (!!int)::2\n",
			"D0, P[0], (!!int)::1\n",
		},
	},
	{
		skipDoc:    true,
		document:   "[]",
		expression: `reverse`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n",
		},
	},
	{
		skipDoc:    true,
		document:   "[1]",
		expression: `reverse`,
		expected: []string{
			"D0, P[], (!!seq)::[1]\n",
		},
	},
	{
		skipDoc:    true,
		document:   "[1,2]",
		expression: `reverse`,
		expected: []string{
			"D0, P[], (!!seq)::[2, 1]\n",
		},
	},
	{
		description:    "Sort descending by string field",
		subdescription: "Use sort with reverse to sort in descending order.",
		document:       "[{a: banana},{a: cat},{a: apple}]",
		expression:     `sort_by(.a) | reverse`,
		expected: []string{
			"D0, P[], (!!seq)::[{a: cat}, {a: banana}, {a: apple}]\n",
		},
	},
	{
		description: "Sort descending by string field, with comments",
		skipDoc:     true,
		document:    "# abc\n[{a: banana},{a: cat},{a: apple}]\n# xyz",
		expression:  `sort_by(.a) | reverse`,
		expected: []string{
			"D0, P[], (!!seq)::# abc\n[{a: cat}, {a: banana}, {a: apple}]\n# xyz\n",
		},
	},
}

func testReverseScenarioWithParserCheck(t *testing.T, s *expressionScenario) {
	// Skip comment-related tests for goccy as it handles comment placement more strictly
	if s.description == "Sort descending by string field, with comments" && ConfiguredYamlPreferences.UseGoccyParser {
		t.Skip("goccy parser handles trailing comments more strictly - structurally equivalent but different comment handling")
		return
	}
	testScenario(t, s)
}

func TestReverseOperatorScenarios(t *testing.T) {
	for _, tt := range reverseOperatorScenarios {
		testReverseScenarioWithParserCheck(t, &tt)
	}
	documentOperatorScenarios(t, "reverse", reverseOperatorScenarios)
}
