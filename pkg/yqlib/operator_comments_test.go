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
		skipDoc:    true,
		document:   "a: cat\nb: dog",
		expression: `.a lineComment=.b`,
		expected: []string{
			"D0, P[], (doc)::a: cat # dog\nb: dog\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: cat\n---\na: dog",
		expression: `.a lineComment |= documentIndex`,
		expected: []string{
			"D0, P[], (doc)::a: cat # 0\n",
			"D1, P[], (doc)::a: dog # 1\n",
		},
	},
	{
		description: "Use update assign to perform relative updates",
		document:    "a: cat\nb: dog",
		expression:  `.. lineComment |= .`,
		expected: []string{
			"D0, P[], (!!map)::a: cat # cat\nb: dog # dog\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: cat\nb: dog",
		expression: `.. comments |= .`,
		expected: []string{
			"D0, P[], (!!map)::a: cat # cat\n# cat\n\n# cat\nb: dog # dog\n# dog\n\n# dog\n",
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
		skipDoc:    true,
		document:   `a: cat`,
		expression: `. footComment=.b.d`,
		expected: []string{
			"D0, P[], (doc)::a: cat\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: cat`,
		expression: `. footComment|=.b.d`,
		expected: []string{
			"D0, P[], (doc)::a: cat\n",
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
		description:    "Remove (strip) all comments",
		subdescription: "Note the use of `...` to ensure key nodes are included.",
		document:       "# hi\n\na: cat # comment\n\n# great\n\nb: # key comment",
		expression:     `... comments=""`,
		expected: []string{
			"D0, P[], (!!map)::a: cat\nb:\n",
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
		description:           "Get head comment",
		dontFormatInputForDoc: true,
		document:              "# welcome!\n\na: cat # meow\n\n# have a great day",
		expression:            `. | headComment`,
		expected: []string{
			"D0, P[], (!!str)::welcome!\n",
		},
	},
	{
		description:           "Get foot comment",
		dontFormatInputForDoc: true,
		document:              "# welcome!\n\na: cat # meow\n\n# have a great day",
		expression:            `. | footComment`,
		expected: []string{
			"D0, P[], (!!str)::have a great day\n",
		},
	},
}

func TestCommentOperatorScenarios(t *testing.T) {
	for _, tt := range commentOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "comment-operators", commentOperatorScenarios)
}
