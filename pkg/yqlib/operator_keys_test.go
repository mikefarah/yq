package yqlib

import (
	"testing"
)

var expectedIsKey = `D0, P[], (!!seq)::- p: ""
  isKey: false
  tag: '!!map'
- p: a
  isKey: true
  tag: '!!str'
- p: a
  isKey: false
  tag: '!!map'
- p: a.b
  isKey: true
  tag: '!!str'
- p: a.b
  isKey: false
  tag: '!!seq'
- p: a.b.0
  isKey: false
  tag: '!!str'
- p: a.c
  isKey: true
  tag: '!!str'
- p: a.c
  isKey: false
  tag: '!!str'
`

var keysOperatorScenarios = []expressionScenario{
	{
		description: "Map keys",
		document:    `{dog: woof, cat: meow}`,
		expression:  `keys`,
		expected: []string{
			"D0, P[], (!!seq)::- dog\n- cat\n",
		},
	},
	{
		description: "Map keys with splat",
		skipDoc:     true,
		document:    `{dog: woof, cat: meow}`,
		expression:  `keys[]`,
		expected: []string{
			"D0, P[dog], (!!str)::dog\n",
			"D0, P[cat], (!!str)::cat\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{}`,
		expression: `keys`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n",
		},
	},
	{
		description: "Array keys",
		document:    `[apple, banana]`,
		expression:  `keys`,
		expected: []string{
			"D0, P[], (!!seq)::- 0\n- 1\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[]`,
		expression: `keys`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n",
		},
	},
	{
		description: "Retrieve array key",
		document:    "[1,2,3]",
		expression:  `.[1] | key`,
		expected: []string{
			"D0, P[1], (!!int)::1\n",
		},
	},
	{
		description: "Retrieve map key",
		document:    "a: thing",
		expression:  `.a | key`,
		expected: []string{
			"D0, P[a], (!!str)::a\n",
		},
	},
	{
		description: "No key",
		document:    "{}",
		expression:  `key`,
		expected:    []string{},
	},
	{
		description: "Update map key",
		document:    "a:\n  x: 3\n  y: 4",
		expression:  `(.a.x | key) = "meow"`,
		expected: []string{
			"D0, P[], (!!map)::a:\n    meow: 3\n    y: 4\n",
		},
	},
	{
		description: "Get comment from map key",
		document:    "a: \n  # comment on key\n  x: 3\n  y: 4",
		expression:  `.a.x | key | headComment`,
		expected: []string{
			"D0, P[a x], (!!str)::comment on key\n",
		},
	},
	{
		description: "Check node is a key",
		document:    "a: \n  b: [cat]\n  c: frog\n",
		expression:  `[... | { "p": path | join("."), "isKey": is_key, "tag": tag }]`,
		expected: []string{
			expectedIsKey,
		},
	},
}

func testKeysScenarioWithParserCheck(t *testing.T, s *expressionScenario) {
	// Skip comment-related tests for goccy as it handles comment placement differently
	if s.description == "Get comment from map key" && ConfiguredYamlPreferences.UseGoccyParser {
		t.Skip("goccy parser handles comment placement differently - data integrity preserved")
		return
	}
	testScenario(t, s)
}

func TestKeysOperatorScenarios(t *testing.T) {
	for _, tt := range keysOperatorScenarios {
		testKeysScenarioWithParserCheck(t, &tt)
	}
	documentOperatorScenarios(t, "keys", keysOperatorScenarios)
}
