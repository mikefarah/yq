package yqlib

import (
	"testing"
)

var documentToPrune = `
parentA: bob
parentB:
  child1: i am child1
  child2: i am child2
parentC:
  child1: me child1
  child2: me child2
`

var pathOperatorScenarios = []expressionScenario{
	{
		description: "Map path",
		document:    `{a: {b: cat}}`,
		expression:  `.a.b | path`,
		expected: []string{
			"D0, P[a b], (!!seq)::- a\n- b\n",
		},
	},
	{
		skipDoc: true,
		document: `a:
  b:
    c:
    - 0
    - 1
    - 2
    - 3`,
		expression: `.a.b.c.[]`,
		expected: []string{
			"D0, P[a b c 0], (!!int)::0\n",
			"D0, P[a b c 1], (!!int)::1\n",
			"D0, P[a b c 2], (!!int)::2\n",
			"D0, P[a b c 3], (!!int)::3\n",
		},
	},
	{
		description: "Get map key",
		document:    `{a: {b: cat}}`,
		expression:  `.a.b | path | .[-1]`,
		expected: []string{
			"D0, P[a b 1], (!!str)::b\n",
		},
	},
	{
		description: "Array path",
		document:    `{a: [cat, dog]}`,
		expression:  `.a.[] | select(. == "dog") | path`,
		expected: []string{
			"D0, P[a 1], (!!seq)::- a\n- 1\n",
		},
	},
	{
		description: "Get array index",
		document:    `{a: [cat, dog]}`,
		expression:  `.a.[] | select(. == "dog") | path | .[-1]`,
		expected: []string{
			"D0, P[a 1 1], (!!int)::1\n",
		},
	},
	{
		description: "Print path and value",
		document:    `{a: [cat, dog, frog]}`,
		expression:  `.a[] | select(. == "*og") | [{"path":path, "value":.}]`,
		expected: []string{
			"D0, P[a 1], (!!seq)::- path:\n    - a\n    - 1\n  value: dog\n",
			"D0, P[a 2], (!!seq)::- path:\n    - a\n    - 2\n  value: frog\n",
		},
	},
	{
		description: "Set path",
		document:    `{a: {b: cat}}`,
		expression:  `setpath(["a", "b"]; "things")`,
		expected: []string{
			"D0, P[], (!!map)::{a: {b: things}}\n",
		},
	},
	{
		description: "Set on empty document",
		expression:  `setpath(["a", "b"]; "things")`,
		expected: []string{
			"D0, P[], ()::a:\n    b: things\n",
		},
	},
	{
		description:    "Set path to prune deep paths",
		subdescription: "Like pick but recursive. This uses `ireduce` to deeply set the selected paths into an empty object.",
		document:       documentToPrune,
		expression:     "(.parentB.child2, .parentC.child1) as $i\n  ireduce({}; setpath($i | path; $i))",
		expected: []string{
			"D0, P[], (!!map)::parentB:\n    child2: i am child2\nparentC:\n    child1: me child1\n",
		},
	},
	{
		description: "Set array path",
		document:    `a: [cat, frog]`,
		expression:  `setpath(["a", 0]; "things")`,
		expected: []string{
			"D0, P[], (!!map)::a: [things, frog]\n",
		},
	},
	{
		description: "Set array path empty",
		expression:  `setpath(["a", 0]; "things")`,
		expected: []string{
			"D0, P[], ()::a:\n    - things\n",
		},
	},
	{
		description:    "Delete path",
		subdescription: "Notice delpaths takes an _array_ of paths.",
		document:       `{a: {b: cat, c: dog, d: frog}}`,
		expression:     `delpaths([["a", "c"], ["a", "d"]])`,
		expected: []string{
			"D0, P[], (!!map)::{a: {b: cat}}\n",
		},
	},
	{
		description: "Delete array path",
		document:    `a: [cat, frog]`,
		expression:  `delpaths([["a", 0]])`,
		expected: []string{
			"D0, P[], (!!map)::a: [frog]\n",
		},
	},
	{
		description:    "Delete - wrong parameter",
		subdescription: "delpaths does not work with a single path array",
		document:       `a: [cat, frog]`,
		expression:     `delpaths(["a", 0])`,
		expectedError:  "DELPATHS: expected entry [0] to be a sequence, but its a !!str. Note that delpaths takes an array of path arrays, e.g. [[\"a\", \"b\"]]",
	},
}

func TestPathOperatorsScenarios(t *testing.T) {
	for _, tt := range pathOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "path", pathOperatorScenarios)
}
