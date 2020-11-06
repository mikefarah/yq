package yqlib

import (
	"testing"
)

var commentOperatorScenarios = []expressionScenario{
	{
		description: "Add line comment",
		document:    `a: cat`,
		expression:  `.a lineComment="single"`,
		expected: []string{
			"D0, P[], (doc)::a: cat # single\n",
		},
	},
	{
		description: "Add head comment",
		document:    `a: cat`,
		expression:  `. headComment="single"`,
		expected: []string{
			"D0, P[], (doc)::# single\n\na: cat\n",
		},
	},
	{
		description: "Add foot comment, using an expression",
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
}

func TestCommentOperatorScenarios(t *testing.T) {
	for _, tt := range commentOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Comments Operator", commentOperatorScenarios)
}
