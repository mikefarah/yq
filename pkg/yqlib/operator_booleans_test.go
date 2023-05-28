package yqlib

import (
	"testing"
)

var booleanOperatorScenarios = []expressionScenario{
	{
		description: "`or` example",
		expression:  `true or false`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		description:    "\"yes\" and \"no\" are strings",
		subdescription: "In the yaml 1.2 standard, support for yes/no as booleans was dropped - they are now considered strings. See '10.2.1.2. Boolean' in https://yaml.org/spec/1.2.2/",
		document:       `[yes, no]`,
		expression:     `.[] | tag`,
		expected: []string{
			"D0, P[0], (!!str)::!!str\n",
			"D0, P[1], (!!str)::!!str\n",
		},
	},
	{
		skipDoc:    true,
		document:   "b: hi",
		expression: `.a or .c`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "b: false",
		expression: `.b or .c`,
		expected: []string{
			"D0, P[b], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "b: hi",
		expression: `select(.a or .b)`,
		expected: []string{
			"D0, P[], (doc)::b: hi\n",
		},
	},
	{
		skipDoc:    true,
		document:   "b: hi",
		expression: `select((.a and .b) | not)`,
		expected: []string{
			"D0, P[], (doc)::b: hi\n",
		},
	},
	{
		skipDoc:     true,
		description: "And should not run 2nd arg if first is false",
		expression:  `false and test(3)`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		skipDoc:     true,
		description: "Or should not run 2nd arg if first is true",
		expression:  `true or test(3)`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		description: "`and` example",
		expression:  `true and false`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		document:    "[{a: bird, b: dog}, {a: frog, b: bird}, {a: cat, b: fly}]",
		description: "Matching nodes with select, equals and or",
		expression:  `[.[] | select(.a == "cat" or .b == "dog")]`,
		expected: []string{
			"D0, P[], (!!seq)::- {a: bird, b: dog}\n- {a: cat, b: fly}\n",
		},
	},
	{
		description: "`any` returns true if any boolean in a given array is true",
		document:    `[false, true]`,
		expression:  "any",
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		description: "`any` returns false for an empty array",
		document:    `[]`,
		expression:  "any",
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		description: "`any_c` returns true if any element in the array is true for the given condition.",
		document:    "a: [rad, awesome]\nb: [meh, whatever]",
		expression:  `.[] |= any_c(. == "awesome")`,
		expected: []string{
			"D0, P[], (doc)::a: true\nb: false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[{pet: cat}]`,
		expression: `any_c(.name == "harry") as $c | .`,
		expected: []string{
			"D0, P[], (doc)::[{pet: cat}]\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[{pet: cat}]`,
		expression: `any_c(.name == "harry") as $c | $c`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[{pet: cat}]`,
		expression: `all_c(.name == "harry") as $c | $c`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[false, false]`,
		expression: "any",
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		description: "`all` returns true if all booleans in a given array are true",
		document:    `[true, true]`,
		expression:  "all",
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[false, true]`,
		expression: "all",
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		description: "`all` returns true for an empty array",
		document:    `[]`,
		expression:  "all",
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		description: "`all_c` returns true if all elements in the array are true for the given condition.",
		document:    "a: [rad, awesome]\nb: [meh, 12]",
		expression:  `.[] |= all_c(tag == "!!str")`,
		expected: []string{
			"D0, P[], (doc)::a: true\nb: false\n",
		},
	},
	{
		skipDoc:    true,
		expression: `false or false`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: true, b: false}`,
		expression: `.[] or (false, true)`,
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
			"D0, P[b], (!!bool)::false\n",
			"D0, P[b], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: true, b: false}`,
		expression: `.[] and (false, true)`,
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
			"D0, P[a], (!!bool)::true\n",
			"D0, P[b], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{}`,
		expression: `(.a.b or .c) as $x | .`,
		expected: []string{
			"D0, P[], (doc)::{}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{}`,
		expression: `(.a.b and .c) as $x | .`,
		expected: []string{
			"D0, P[], (doc)::{}\n",
		},
	},
	{
		description: "Not true is false",
		expression:  `true | not`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		description: "Not false is true",
		expression:  `false | not`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		description: "String values considered to be true",
		expression:  `"cat" | not`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		description: "Empty string value considered to be true",
		expression:  `"" | not`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		description: "Numbers are considered to be true",
		expression:  `1 | not`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		description: "Zero is considered to be true",
		expression:  `0 | not`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{

		description: "Null is considered to be false",
		expression:  `~ | not`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
}

func TestBooleanOperatorScenarios(t *testing.T) {
	for _, tt := range booleanOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "boolean-operators", booleanOperatorScenarios)
}
