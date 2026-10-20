package yqlib

import (
	"testing"
)

var unionOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		document:   "{}",
		expression: `(.a, .b.c) as $x | .`,
		expected: []string{
			"D0, P[], (!!map)::{}\n",
		},
	},
	{
		skipDoc:     true,
		description: "clone test",
		expression:  `"abc" as $a | [$a, "cat"]`,
		expected: []string{
			"D0, P[], (!!seq)::- abc\n- cat\n",
		},
	},
	{
		skipDoc:    true,
		expression: `(.foo = "bar"), (.toe = "jam")`,
		expected: []string{
			"D0, P[], ()::foo: bar\ntoe: jam\n",
		},
	},
	{
		skipDoc:     true,
		description: "Three assignments modify one document",
		expression:  `.a.b.c = "foo", .a.b.d = "bar", .a.b.e = "baz"`,
		expected: []string{
			"D0, P[], ()::a:\n    b:\n        c: foo\n        d: bar\n        e: baz\n",
		},
	},
	{
		skipDoc:     true,
		description: "Four assignments modify one document",
		expression:  `.a = "one", .b = "two", .c = "three", .d = "four"`,
		expected: []string{
			"D0, P[], ()::a: one\nb: two\nc: three\nd: four\n",
		},
	},
	{
		skipDoc:     true,
		description: "Right-nested assignments modify one document",
		expression:  `.a = "one", (.b = "two", .c = "three")`,
		expected: []string{
			"D0, P[], ()::a: one\nb: two\nc: three\n",
		},
	},
	{
		skipDoc:     true,
		description: "Left-nested assignments modify one document",
		expression:  `(.a = "one", .b = "two"), .c = "three"`,
		expected: []string{
			"D0, P[], ()::a: one\nb: two\nc: three\n",
		},
	},
	{
		skipDoc:     true,
		description: "Grouped assignments modify one document",
		expression:  `(.a = "one", .b = "two"), (.c = "three", .d = "four")`,
		expected: []string{
			"D0, P[], ()::a: one\nb: two\nc: three\nd: four\n",
		},
	},
	{
		skipDoc:     true,
		description: "Assignments preserve an existing input document",
		document:    "a: old\nb: old\n",
		expression:  `.a = "new", .b = "new", .c = "new"`,
		expected: []string{
			"D0, P[], (!!map)::a: new\nb: new\nc: new\n",
		},
	},
	{
		skipDoc:     true,
		description: "Pipe assignments continue producing one document",
		expression:  `.a = "one" | .b = "two" | .c = "three"`,
		expected: []string{
			"D0, P[], ()::a: one\nb: two\nc: three\n",
		},
	},
	{
		description: "Combine scalars",
		expression:  `1, true, "cat"`,
		expected: []string{
			"D0, P[], (!!int)::1\n",
			"D0, P[], (!!bool)::true\n",
			"D0, P[], (!!str)::cat\n",
		},
	},
	{
		description: "Combine selected paths",
		document:    `{a: fieldA, b: fieldB, c: fieldC}`,
		expression:  `.a, .c`,
		expected: []string{
			"D0, P[a], (!!str)::fieldA\n",
			"D0, P[c], (!!str)::fieldC\n",
		},
	},
	{
		skipDoc:     true,
		description: "Identical scalars remain separate union results",
		expression:  `1, 1`,
		expected: []string{
			"D0, P[], (!!int)::1\n",
			"D0, P[], (!!int)::1\n",
		},
	},
	{
		skipDoc:     true,
		description: "Identical selected paths remain separate union results",
		document:    `{a: cat}`,
		expression:  `.a, .a`,
		expected: []string{
			"D0, P[a], (!!str)::cat\n",
			"D0, P[a], (!!str)::cat\n",
		},
	},
	{
		skipDoc:     true,
		description: "Modified document and scalar remain separate union results",
		expression:  `(.a = "foo"), "other"`,
		expected: []string{
			"D0, P[], ()::a: foo\n",
			"D0, P[], (!!str)::other\n",
		},
	},
}

func TestUnionOperatorScenarios(t *testing.T) {
	for _, tt := range unionOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "union", unionOperatorScenarios)
}
