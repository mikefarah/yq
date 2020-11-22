package yqlib

import (
	"testing"
)

var commentOperatorScenarios = []expressionScenario{
	{
		description: "Set line comment",
		document:    `a: cat`,
		expression:  `.a lineComment="single"`,
		expected: []string{
			"D0, P[], (doc)::a: cat # single\n",
		},
	},
	{
		description: "Set head comment",
		document:    `a: cat`,
		expression:  `. headComment="single"`,
		expected: []string{
			"D0, P[], (doc)::# single\n\na: cat\n",
		},
	},
	{
		description: "Set foot comment, using an expression",
		document:    `a: cat`,
		expression:  `. footComment=.a`,
		expected: []string{
			"D0, P[], (doc)::a: cat\n\n# cat\n",
		},
	},
	{
		description: "Remove comment",
		document:    "a: cat # comment\nb: dog # leave this",
		expression:  `.a lineComment=""`,
		expected: []string{
			"D0, P[], (doc)::a: cat\nb: dog # leave this\n",
		},
	},
	{
		description: "Remove all comments",
		document:    "# hi\n\na: cat # comment\n\n# great\n",
		expression:  `.. comments=""`,
		expected: []string{
			"D0, P[], (!!map)::a: cat\n",
		},
	},
	{
		description: "Get line comment",
		document:    "# welcome!\n\na: cat # meow\n\n# have a great day",
		expression:  `.a | lineComment`,
		expected: []string{
			"D0, P[a], (!!str)::meow\n",
		},
	},
	{
		description: "Get head comment",
		document:    "# welcome!\n\na: cat # meow\n\n# have a great day",
		expression:  `. | headComment`,
		expected: []string{
			"D0, P[], (!!str)::welcome!\n",
		},
	},
	{
		description: "Get foot comment",
		document:    "# welcome!\n\na: cat # meow\n\n# have a great day",
		expression:  `. | footComment`,
		expected: []string{
			"D0, P[], (!!str)::have a great day\n",
		},
	},
}

func TestCommentOperatorScenarios(t *testing.T) {
	for _, tt := range commentOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Comment Operators", commentOperatorScenarios)
}
