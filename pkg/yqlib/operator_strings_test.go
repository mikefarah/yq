package yqlib

import (
	"testing"
)

var stringsOperatorScenarios = []expressionScenario{
	{
		description: "Interpolation",
		document:    "value: things\nanother: stuff",
		expression:  `.message = "I like \(.value) and \(.another)"`,
		expected: []string{
			"D0, P[], (!!map)::value: things\nanother: stuff\nmessage: I like things and stuff\n",
		},
	},
	{
		description: "Interpolation - not a string",
		document:    `value: {an: apple}`,
		expression:  `.message = "I like \(.value)"`,
		expected: []string{
			"D0, P[], (!!map)::value: {an: apple}\nmessage: 'I like {an: apple}'\n",
		},
	},
	{
		skipDoc:     true,
		description: "Interpolation - just escape",
		expression:  `"\\"`,
		expected: []string{
			"D0, P[], (!!str)::\\\n",
		},
	},
	{
		skipDoc:     true,
		description: "Interpolation - nested",
		document:    `value: things`,
		expression:  `"Hi \( (.value) )"`,
		expected: []string{
			"D0, P[], (!!str)::Hi things\n",
		},
	},
	{
		skipDoc:     true,
		description: "Interpolation - don't",
		document:    `value: things`,
		expression:  `"Hi (.value)"`,
		expected: []string{
			"D0, P[], (!!str)::Hi (.value)\n",
		},
	},
	{
		skipDoc:     true,
		description: "Interpolation - don't!",
		document:    `value: things`,
		expression:  `"Hi \\(.value)"`,
		expected: []string{
			"D0, P[], (!!str)::Hi \\(.value)\n",
		},
	},
	{
		skipDoc:     true,
		description: "Interpolation - random close bracket",
		document:    `value: things`,
		expression:  `"Hi )"`,
		expected: []string{
			"D0, P[], (!!str)::Hi )\n",
		},
	},
	{
		skipDoc:       true,
		description:   "Interpolation - unclosed interpolation string",
		document:      `value: things`,
		expression:    `"Hi \("`,
		expectedError: "unclosed interpolation string \\(",
	},
	{
		skipDoc:       true,
		description:   "Interpolation - unclosed interpolation string due to escape",
		document:      `value: things`,
		expression:    `"Hi \(\)"`,
		expectedError: "unclosed interpolation string \\(",
	},
	{
		description:    "To up (upper) case",
		subdescription: "Works with unicode characters",
		document:       `água`,
		expression:     "upcase",
		expected: []string{
			"D0, P[], (!!str)::ÁGUA\n",
		},
	},
	{
		skipDoc:    true,
		document:   `!camel água`,
		expression: "upcase",
		expected: []string{
			"D0, P[], (!camel)::ÁGUA\n",
		},
	},
	{
		description:    "To down (lower) case",
		subdescription: "Works with unicode characters",
		document:       `ÁgUA`,
		expression:     "downcase",
		expected: []string{
			"D0, P[], (!!str)::água\n",
		},
	},
	{
		skipDoc:    true,
		document:   `!camel ÁgUA`,
		expression: "downcase",
		expected: []string{
			"D0, P[], (!camel)::água\n",
		},
	},
	{
		description: "Join strings",
		document:    `[cat, meow, 1, null, true]`,
		expression:  `join("; ")`,
		expected: []string{
			"D0, P[], (!!str)::cat; meow; 1; ; true\n",
		},
	},
	{
		description: "Trim strings",
		document:    `[" cat", "dog ", " cow cow ", horse]`,
		expression:  `.[] | trim`,
		expected: []string{
			"D0, P[0], (!!str)::cat\n",
			"D0, P[1], (!!str)::dog\n",
			"D0, P[2], (!!str)::cow cow\n",
			"D0, P[3], (!!str)::horse\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[!horse cat, !goat meow, !frog 1, null, true]`,
		expression: `join("; ")`,
		expected: []string{
			"D0, P[], (!!str)::cat; meow; 1; ; true\n",
		},
	},
	{
		description: "Match string",
		document:    `foo bar foo`,
		expression:  `match("foo")`,
		expected: []string{
			"D0, P[], (!!map)::string: foo\noffset: 0\nlength: 3\ncaptures: []\n",
		},
	},
	{
		skipDoc:    true,
		document:   `!horse foo bar foo`,
		expression: `match("foo")`,
		expected: []string{
			"D0, P[], (!!map)::string: foo\noffset: 0\nlength: 3\ncaptures: []\n",
		},
	},
	{
		description: "Match string, case insensitive",
		document:    `foo bar FOO`,
		expression:  `[match("(?i)foo"; "g")]`,
		expected: []string{
			"D0, P[], (!!seq)::- string: foo\n  offset: 0\n  length: 3\n  captures: []\n- string: FOO\n  offset: 8\n  length: 3\n  captures: []\n",
		},
	},
	{
		description: "Match with global capture group",
		document:    `abc abc`,
		expression:  `[match("(ab)(c)"; "g")]`,
		expected: []string{
			"D0, P[], (!!seq)::- string: abc\n  offset: 0\n  length: 3\n  captures:\n    - string: ab\n      offset: 0\n      length: 2\n    - string: c\n      offset: 2\n      length: 1\n- string: abc\n  offset: 4\n  length: 3\n  captures:\n    - string: ab\n      offset: 4\n      length: 2\n    - string: c\n      offset: 6\n      length: 1\n",
		},
	},
	{
		description: "Match with named capture groups",
		document:    `foo bar foo foo  foo`,
		expression:  `[match("foo (?P<bar123>bar)? foo"; "g")]`,
		expected: []string{
			"D0, P[], (!!seq)::- string: foo bar foo\n  offset: 0\n  length: 11\n  captures:\n    - string: bar\n      offset: 4\n      length: 3\n      name: bar123\n- string: foo  foo\n  offset: 12\n  length: 8\n  captures:\n    - string: null\n      offset: -1\n      length: 0\n      name: bar123\n",
		},
	},
	{
		description: "Capture named groups into a map",
		document:    `xyzzy-14`,
		expression:  `capture("(?P<a>[a-z]+)-(?P<n>[0-9]+)")`,
		expected: []string{
			"D0, P[], (!!map)::a: xyzzy\nn: \"14\"\n",
		},
	},
	{
		skipDoc:    true,
		document:   `!horse xyzzy-14`,
		expression: `capture("(?P<a>[a-z]+)-(?P<n>[0-9]+)")`,
		expected: []string{
			"D0, P[], (!!map)::a: xyzzy\nn: \"14\"\n",
		},
	},
	{
		skipDoc:     true,
		description: "Capture named groups into a map, with null",
		document:    `xyzzy-14`,
		expression:  `capture("(?P<a>[a-z]+)-(?P<n>[0-9]+)(?P<bar123>bar)?")`,
		expected: []string{
			"D0, P[], (!!map)::a: xyzzy\nn: \"14\"\nbar123: null\n",
		},
	},
	{
		description: "Match without global flag",
		document:    `cat cat`,
		expression:  `match("cat")`,
		expected: []string{
			"D0, P[], (!!map)::string: cat\noffset: 0\nlength: 3\ncaptures: []\n",
		},
	},
	{
		description: "Match with global flag",
		document:    `cat cat`,
		expression:  `[match("cat"; "g")]`,
		expected: []string{
			"D0, P[], (!!seq)::- string: cat\n  offset: 0\n  length: 3\n  captures: []\n- string: cat\n  offset: 4\n  length: 3\n  captures: []\n",
		},
	},
	{
		skipDoc:    true,
		document:   `!horse cat cat`,
		expression: `[match("cat"; "g")]`,
		expected: []string{
			"D0, P[], (!!seq)::- string: cat\n  offset: 0\n  length: 3\n  captures: []\n- string: cat\n  offset: 4\n  length: 3\n  captures: []\n",
		},
	},
	{
		skipDoc:     true,
		description: "No match",
		document:    `dog`,
		expression:  `match("cat"; "g")`,
		expected:    []string{},
	},
	{
		skipDoc:     true,
		description: "No match",
		expression:  `"dog" | match("cat", "g")`,
		expected:    []string{},
	},
	{
		skipDoc:     true,
		description: "No match",
		expression:  `"dog" | match("cat")`,
		expected:    []string{},
	},
	{
		description:    "Test using regex",
		subdescription: "Like jq's equivalent, this works like match but only returns true/false instead of full match details",
		document:       `["cat", "dog"]`,
		expression:     `.[] | test("at")`,
		expected: []string{
			"D0, P[0], (!!bool)::true\n",
			"D0, P[1], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[!horse "cat", !cat "dog"]`,
		expression: `.[] | test("at")`,
		expected: []string{
			"D0, P[0], (!!bool)::true\n",
			"D0, P[1], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `["cat*", "cat*", "cat"]`,
		expression: `.[] | test("cat\*")`,
		expected: []string{
			"D0, P[0], (!!bool)::true\n",
			"D0, P[1], (!!bool)::true\n",
			"D0, P[2], (!!bool)::false\n",
		},
	},
	{
		description:    "Substitute / Replace string",
		subdescription: "This uses Golang's regex, described [here](https://github.com/google/re2/wiki/Syntax).\nNote the use of `|=` to run in context of the current string value.",
		document:       `a: dogs are great`,
		expression:     `.a |= sub("dogs", "cats")`,
		expected: []string{
			"D0, P[], (!!map)::a: cats are great\n",
		},
	},
	{
		description:    "Substitute / Replace string with regex",
		subdescription: "This uses Golang's regex, described [here](https://github.com/google/re2/wiki/Syntax).\nNote the use of `|=` to run in context of the current string value.",
		document:       "a: cat\nb: heat",
		expression:     `.[] |= sub("(a)", "${1}r")`,
		expected: []string{
			"D0, P[], (!!map)::a: cart\nb: heart\n",
		},
	},
	{
		description:    "Custom types: that are really strings",
		subdescription: "When custom tags are encountered, yq will try to decode the underlying type.",
		document:       "a: !horse cat\nb: !goat heat",
		expression:     `.[] |= sub("(a)", "${1}r")`,
		expected: []string{
			"D0, P[], (!!map)::a: !horse cart\nb: !goat heart\n",
		},
	},
	{
		description: "Split strings",
		document:    `"cat; meow; 1; ; true"`,
		expression:  `split("; ")`,
		expected: []string{
			"D0, P[], (!!seq)::- cat\n- meow\n- \"1\"\n- \"\"\n- \"true\"\n",
		},
	},
	{
		description: "Split strings one match",
		document:    `"word"`,
		expression:  `split("; ")`,
		expected: []string{
			"D0, P[], (!!seq)::- word\n",
		},
	},
	{
		skipDoc:    true,
		document:   `!horse "word"`,
		expression: `split("; ")`,
		expected: []string{
			"D0, P[], (!!seq)::- word\n",
		},
	},
	{
		skipDoc:    true,
		document:   `""`,
		expression: `split("; ")`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n", // dont actually want this, just not to error
		},
	},
	{
		skipDoc:    true,
		expression: `split("; ")`,
		expected:   []string{},
	},
	{
		description:    "To string",
		subdescription: "Note that you may want to force `yq` to leave scalar values wrapped by passing in `--unwrapScalar=false` or `-r=f`",
		document:       `[1, true, null, ~, cat, {an: object}, [array, 2]]`,
		expression:     ".[] |= to_string",
		expected: []string{
			"D0, P[], (!!seq)::[\"1\", \"true\", \"null\", \"~\", cat, \"{an: object}\", \"[array, 2]\"]\n",
		},
	},
}

func TestStringsOperatorScenarios(t *testing.T) {
	for _, tt := range stringsOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "string-operators", stringsOperatorScenarios)
}
