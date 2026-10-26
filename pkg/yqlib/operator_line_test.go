package yqlib

import (
	"testing"
)

var lineOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		document:   "---\n\nfoo:\n  bar: 42\n",
		expression: `.foo.bar | line`,
		expected:   []string{"D0, P[foo bar], (!!int)::4\n"},
	},
	{
		skipDoc:    true,
		document:   "# comment\nfoo:\n  bar: 42\n",
		expression: `.foo.bar | line`,
		expected:   []string{"D0, P[foo bar], (!!int)::3\n"},
	},
	{
		skipDoc:    true,
		document:   "%YAML 1.1\r\n---\r\n\r\nfoo: 42\r\n",
		expression: `.foo | key | line`,
		expected:   []string{"D0, P[foo], (!!int)::4\n"},
	},
	{
		skipDoc:    true,
		document:   "# header\nfoo: 1\n---\nfoo: 2\n",
		expression: `.foo | line`,
		expected:   []string{"D0, P[foo], (!!int)::2\n", "D1, P[foo], (!!int)::4\n"},
	},
	{
		skipDoc:    true,
		document:   "# header\nfoo: 1\n",
		document2:  "foo: 2\n",
		expression: `.foo | line`,
		expected:   []string{"D0, P[foo], (!!int)::2\n", "D0, P[foo], (!!int)::1\n"},
	},
	{
		skipDoc:    true,
		document:   "--- {foo: 42}",
		expression: `.foo | line`,
		expected:   []string{"D0, P[foo], (!!int)::1\n"},
	},
	{
		skipDoc:    true,
		document:   "--- # header\nfoo: 42\n",
		expression: `.foo | line`,
		expected:   []string{"D0, P[foo], (!!int)::2\n"},
	},
	{
		description: "Returns line of _value_ node",
		document:    "a: cat\nb:\n   c: cat",
		expression:  `.b | line`,
		expected: []string{
			"D0, P[b], (!!int)::3\n",
		},
	},
	{
		description:    "Returns line of _key_ node",
		subdescription: "Pipe through the key operator to get the line of the key",
		document:       "a: cat\nb:\n   c: cat",
		expression:     `.b | key | line`,
		expected: []string{
			"D0, P[b], (!!int)::2\n",
		},
	},
	{
		description: "First line is 1",
		document:    "a: cat",
		expression:  `.a | line`,
		expected: []string{
			"D0, P[a], (!!int)::1\n",
		},
	},
	{
		description: "No line data is 0",
		expression:  `{"a": "new entry"} | line`,
		expected: []string{
			"D0, P[], (!!int)::0\n",
		},
	},
}

func TestLineOperatorScenarios(t *testing.T) {
	for _, tt := range lineOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "line", lineOperatorScenarios)
}
