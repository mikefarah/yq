package yqlib

import (
	"testing"
)

var stringsOperatorScenarios = []expressionScenario{
	// {
	// 	description: "Join strings",
	// 	document:    `[cat, meow, 1, null, true]`,
	// 	expression:  `join("; ")`,
	// 	expected: []string{
	// 		"D0, P[], (!!str)::cat; meow; 1; ; true\n",
	// 	},
	// },
	// {
	// 	description: "Match string",
	// 	document:    `cat`,
	// 	expression:  `match("at")`,
	// 	expected: []string{
	// 		"D0, P[], ()::string: at\noffset: 1\nlength: 2\ncaptures: []\n",
	// 	},
	// },
	// {
	// 	description: "Match string, case insensitive",
	// 	document:    `cAt`,
	// 	expression:  `match("(?i)at")`,
	// 	expected: []string{
	// 		"D0, P[], ()::string: At\noffset: 1\nlength: 2\ncaptures: []\n",
	// 	},
	// },
	// {
	// 	description: "Match with capture groups",
	// 	document:    `a cat`,
	// 	expression:  `match("c(.t)")`,
	// 	expected: []string{
	// 		"D0, P[], ()::string: cat\noffset: 2\nlength: 3\ncaptures:\n    - string: at\n      offset: 3\n      length: 2\n",
	// 	},
	// },
	{
		description: "Match with names capture groups",
		document:    `a cat`,
		expression:  `match("c(?P<cool>.t)")`,
		expected: []string{
			"D0, P[], ()::string: cat\noffset: 2\nlength: 3\ncaptures:\n    - string: at\n      offset: 3\n      length: 2\n      name: cool\n",
		},
	},
	// {
	// 	description:    "Substitute / Replace string",
	// 	subdescription: "This uses golang regex, described [here](https://github.com/google/re2/wiki/Syntax)\nNote the use of `|=` to run in context of the current string value.",
	// 	document:       `a: dogs are great`,
	// 	expression:     `.a |= sub("dogs", "cats")`,
	// 	expected: []string{
	// 		"D0, P[], (doc)::a: cats are great\n",
	// 	},
	// },
	// {
	// 	description:    "Substitute / Replace string with regex",
	// 	subdescription: "This uses golang regex, described [here](https://github.com/google/re2/wiki/Syntax)\nNote the use of `|=` to run in context of the current string value.",
	// 	document:       "a: cat\nb: heat",
	// 	expression:     `.[] |= sub("(a)", "${1}r")`,
	// 	expected: []string{
	// 		"D0, P[], (doc)::a: cart\nb: heart\n",
	// 	},
	// },
	// {
	// 	description: "Split strings",
	// 	document:    `"cat; meow; 1; ; true"`,
	// 	expression:  `split("; ")`,
	// 	expected: []string{
	// 		"D0, P[], (!!seq)::- cat\n- meow\n- \"1\"\n- \"\"\n- \"true\"\n",
	// 	},
	// },
	// {
	// 	description: "Split strings one match",
	// 	document:    `"word"`,
	// 	expression:  `split("; ")`,
	// 	expected: []string{
	// 		"D0, P[], (!!seq)::- word\n",
	// 	},
	// },
	// {
	// 	skipDoc:    true,
	// 	document:   `""`,
	// 	expression: `split("; ")`,
	// 	expected: []string{
	// 		"D0, P[], (!!seq)::[]\n", // dont actually want this, just not to error
	// 	},
	// },
	// {
	// 	skipDoc:    true,
	// 	expression: `split("; ")`,
	// 	expected:   []string{},
	// },
}

func TestStringsOperatorScenarios(t *testing.T) {
	for _, tt := range stringsOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "String Operators", stringsOperatorScenarios)
}
