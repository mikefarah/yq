package yqlib

import (
	"testing"
)

var commentOperatorScenarios = []expressionScenario{
	{
		description: "Set line comment",
		document:    `a: cat`,
		expression:  `.a line_comment="single"`,
		expected: []string{
			"D0, P[], (doc)::a: cat # single\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: cat\nb: dog",
		expression: `.a line_comment=.b`,
		expected: []string{
			"D0, P[], (doc)::a: cat # dog\nb: dog\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: cat\n---\na: dog",
		expression: `.a line_comment |= documentIndex`,
		expected: []string{
			"D0, P[], (doc)::a: cat # 0\n",
			"D1, P[], (doc)::a: dog # 1\n",
		},
	},
	{
		description: "Use update assign to perform relative updates",
		document:    "a: cat\nb: dog",
		expression:  `.. line_comment |= .`,
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
		expression:  `. head_comment="single"`,
		expected: []string{
			"D0, P[], (doc)::# single\n\na: cat\n",
		},
	},
	{
		description: "Set foot comment, using an expression",
		document:    `a: cat`,
		expression:  `. foot_comment=.a`,
		expected: []string{
			"D0, P[], (doc)::a: cat\n# cat\n",
		},
	},
	{
		skipDoc:     true,
		description: "Set foot comment, using an expression",
		document:    "a: cat\n\n# hi",
		expression:  `. foot_comment=""`,
		expected: []string{
			"D0, P[], (doc)::a: cat\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: cat`,
		expression: `. foot_comment=.b.d`,
		expected: []string{
			"D0, P[], (doc)::a: cat\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: cat`,
		expression: `. foot_comment|=.b.d`,
		expected: []string{
			"D0, P[], (doc)::a: cat\n",
		},
	},
	{
		description: "Remove comment",
		document:    "a: cat # comment\nb: dog # leave this",
		expression:  `.a line_comment=""`,
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
		expression:  `.a | line_comment`,
		expected: []string{
			"D0, P[a], (!!str)::meow\n",
		},
	},
	{
		description:           "Get head comment",
		dontFormatInputForDoc: true,
		document:              "# welcome!\n\na: cat # meow\n\n# have a great day",
		expression:            `. | head_comment`,
		expected: []string{
			"D0, P[], (!!str)::welcome!\n",
		},
	},
	{
		skipDoc:     true,
		description: "strip trailing comment recurse all",
		document:    "a: cat\n\n# haha",
		expression:  `... comments= ""`,
		expected: []string{
			"D0, P[], (!!map)::a: cat\n",
		},
	},
	{
		skipDoc:     true,
		description: "strip trailing comment recurse values",
		document:    "a: cat\n\n# haha",
		expression:  `.. comments= ""`,
		expected: []string{
			"D0, P[], (!!map)::a: cat\n",
		},
	},
	{
		description:           "Head comment with document split",
		dontFormatInputForDoc: true,
		document:              "# welcome!\n---\n# bob\na: cat # meow\n\n# have a great day",
		expression:            `head_comment`,
		expected: []string{
			"D0, P[], (!!str)::welcome!\nbob\n",
		},
	},
	{
		description:           "Get foot comment",
		dontFormatInputForDoc: true,
		document:              "# welcome!\n\na: cat # meow\n\n# have a great day\n# no really",
		expression:            `. | foot_comment`,
		expected: []string{
			"D0, P[], (!!str)::have a great day\nno really\n",
		},
	},
}

func TestCommentOperatorScenarios(t *testing.T) {
	for _, tt := range commentOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "comment-operators", commentOperatorScenarios)
}
