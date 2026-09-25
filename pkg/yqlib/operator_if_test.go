package yqlib

import (
	"testing"
)

var ifOperatorScenarios = []expressionScenario{
	{
		description:    "Basic if-then-else",
		subdescription: "The condition is evaluated against each input; `null` and `false` are falsy, everything else is truthy.",
		document:       `[1, 2, 3]`,
		expression:     `.[] |= if . > 1 then "big" else "small" end`,
		expected: []string{
			"D0, P[], (!!seq)::[small, big, big]\n",
		},
	},
	{
		description: "elif",
		document:    `[1, 2, 3]`,
		expression:  `.[] | if . == 1 then "one" elif . == 2 then "two" else "many" end`,
		expected: []string{
			"D0, P[], (!!str)::one\n",
			"D0, P[], (!!str)::two\n",
			"D0, P[], (!!str)::many\n",
		},
	},
	{
		description:    "else is optional",
		subdescription: "Like jq, if there is no `else` and the condition is falsy, the input is returned unchanged.",
		document:       `[1, 2, 3]`,
		expression:     `.[] |= if . == 2 then "two" end`,
		expected: []string{
			"D0, P[], (!!seq)::[1, two, 3]\n",
		},
	},
	{
		description:    "Update matching branch",
		subdescription: "The branches can be used to choose which path to update.",
		document:       `{enabled: true, a: 1, b: 2}`,
		expression:     `(if .enabled then .a else .b end) = 10`,
		expected: []string{
			"D0, P[], (!!map)::{enabled: true, a: 10, b: 2}\n",
		},
	},
	{
		description:    "Condition with multiple results",
		subdescription: "Like jq, each result of the condition produces a result.",
		expression:     `if true, false then "yes" else "no" end`,
		expected: []string{
			"D0, P[], (!!str)::yes\n",
			"D0, P[], (!!str)::no\n",
		},
	},
	{
		description: "Keywords can still be used as keys",
		document:    `{if: {then: {else: {end: 1}}}}`,
		expression:  `if .if.then.else.end == 1 then "OK" else "OOPS" end`,
		expected: []string{
			"D0, P[], (!!str)::OK\n",
		},
	},
	{
		description:    "Missing keys are falsy",
		subdescription: "A condition that returns no results, such as a missing key, is treated as `null` - consistent with `and`, `or` and `//`.",
		document:       `{a: 1}`,
		expression:     `if .b then "has b" else "no b" end`,
		expected: []string{
			"D0, P[], (!!str)::no b\n",
		},
	},
	{
		skipDoc:     true,
		description: "missing key in condition does not create the key",
		document:    `{a: 1}`,
		expression:  `.c = (if .b then "yes" else "no" end)`,
		expected: []string{
			"D0, P[], (!!map)::{a: 1, c: no}\n",
		},
	},
	{
		skipDoc:     true,
		description: "missing else on falsy input returns input",
		document:    `a: 1`,
		expression:  `if .b then "yes" end`,
		expected: []string{
			"D0, P[], (!!map)::a: 1\n",
		},
	},
	{
		skipDoc:     true,
		description: "elif without else",
		document:    `[1, 2, 3]`,
		expression:  `.[] | if . == 1 then "one" elif . == 2 then "two" end`,
		expected: []string{
			"D0, P[], (!!str)::one\n",
			"D0, P[], (!!str)::two\n",
			"D0, P[2], (!!int)::3\n",
		},
	},
	{
		skipDoc:     true,
		description: "pipes, unions and boolean ops within the condition and branches",
		document:    `{a: {b: 5}}`,
		expression:  `if (.a | .b > 3) and .a.b < 10 then .a.b | . * 2, 1 else 0 end`,
		expected: []string{
			"D0, P[a b], (!!int)::10\n",
			"D0, P[], (!!int)::1\n",
		},
	},
	{
		skipDoc:     true,
		description: "nested if",
		document:    `{a: 1}`,
		expression:  `if (if .a then false else true end) then "x" else if .a == 1 then "y" else "z" end end`,
		expected: []string{
			"D0, P[], (!!str)::y\n",
		},
	},
	{
		skipDoc:     true,
		description: "traverse after end",
		document:    `{a: {b: 5}}`,
		expression:  `if . then .a end.b`,
		expected: []string{
			"D0, P[a b], (!!int)::5\n",
		},
	},
	{
		skipDoc:     true,
		description: "if within collect, object and arithmetic",
		document:    `{a: 1}`,
		expression:  `[if .a then "x" else "y" end + "!", {"k": if .b then 1 else 2 end}]`,
		expected: []string{
			"D0, P[], (!!seq)::- x!\n- k: 2\n",
		},
	},
	{
		skipDoc:     true,
		description: "variables",
		document:    `{a: 5}`,
		expression:  `.a as $x | if $x == 5 then "five" else "other" end`,
		expected: []string{
			"D0, P[], (!!str)::five\n",
		},
	},
	{
		skipDoc:              true,
		description:          "env",
		environmentVariables: map[string]string{"myenv": "fred"},
		expression:           `if strenv(myenv) == "fred" then "is-fred" else "is-not-fred" end`,
		expected: []string{
			"D0, P[], (!!str)::is-fred\n",
		},
	},
}

func TestIfOperatorScenarios(t *testing.T) {
	for _, tt := range ifOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "if-then-else", ifOperatorScenarios)
}
