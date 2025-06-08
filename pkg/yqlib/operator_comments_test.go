package yqlib

import (
	"testing"
)

var expectedWhereIsMyCommentMapKey = `D0, P[], (!!seq)::- p: ""
  isKey: false
  hc: ""
  lc: ""
  fc: ""
- p: hello
  isKey: true
  hc: ""
  lc: hello-world-comment
  fc: ""
- p: hello
  isKey: false
  hc: ""
  lc: ""
  fc: ""
- p: hello.message
  isKey: true
  hc: ""
  lc: ""
  fc: ""
- p: hello.message
  isKey: false
  hc: ""
  lc: ""
  fc: ""
`

var expectedWhereIsMyCommentArray = `D0, P[], (!!seq)::- p: ""
  isKey: false
  hc: ""
  lc: ""
  fc: ""
- p: name
  isKey: true
  hc: ""
  lc: ""
  fc: ""
- p: name
  isKey: false
  hc: ""
  lc: ""
  fc: ""
- p: name.0
  isKey: false
  hc: under-name-comment
  lc: ""
  fc: ""
`

var commentOperatorScenarios = []expressionScenario{
	{
		description:    "Set line comment",
		subdescription: "Set the comment on the key node for more reliability (see below).",
		document:       `a: cat`,
		expression:     `.a line_comment="single"`,
		expected: []string{
			"D0, P[], (!!map)::a: cat # single\n",
		},
		skipForGoccy: true,
	},
	{
		description:    "Set line comment of a maps/arrays",
		subdescription: "For maps and arrays, you need to set the line comment on the _key_ node. This will also work for scalars.",
		document:       "a:\n  b: things",
		expression:     `(.a | key) line_comment="single"`,
		expected: []string{
			"D0, P[], (!!map)::a: # single\n    b: things\n",
		},
		skipForGoccy: true,
	},
	{
		skipDoc:    true,
		document:   "a: cat\nb: dog",
		expression: `.a line_comment=.b`,
		expected: []string{
			"D0, P[], (!!map)::a: cat # dog\nb: dog\n",
		},
		skipForGoccy: true,
	},
	{
		skipDoc:    true,
		document:   "a: cat\n---\na: dog",
		expression: `.a line_comment |= documentIndex`,
		expected: []string{
			"D0, P[], (!!map)::a: cat # 0\n",
			"D1, P[], (!!map)::a: dog # 1\n",
		},
		skipForGoccy: true,
	},
	{
		description: "Use update assign to perform relative updates",
		document:    "a: cat\nb: dog",
		expression:  `.. line_comment |= .`,
		expected: []string{
			"D0, P[], (!!map)::a: cat # cat\nb: dog # dog\n",
		},
		skipForGoccy: true,
	},
	{
		skipDoc:    true,
		document:   "a: cat\nb: dog",
		expression: `.. comments |= .`,
		expected: []string{
			"D0, P[], (!!map)::a: cat # cat\n# cat\n\n# cat\nb: dog # dog\n# dog\n\n# dog\n",
		},
		skipForGoccy: true,
	},
	{
		description:    "Where is the comment - map key example",
		subdescription: "The underlying yaml parser can assign comments in a document to surprising nodes. Use an expression like this to find where you comment is. 'p' indicates the path, 'isKey' is if the node is a map key (as opposed to a map value).\nFrom this, you can see the 'hello-world-comment' is actually on the 'hello' key",
		document:       "hello: # hello-world-comment\n  message: world",
		expression:     `[... | {"p": path | join("."), "isKey": is_key, "hc": headComment, "lc": lineComment, "fc": footComment}]`,
		expected: []string{
			expectedWhereIsMyCommentMapKey,
		},
		skipForGoccy: true,
	},
	{
		description:    "Retrieve comment - map key example",
		subdescription: "From the previous example, we know that the comment is on the 'hello' _key_ as a lineComment",
		document:       "hello: # hello-world-comment\n  message: world",
		expression:     `.hello | key | line_comment`,
		expected: []string{
			"D0, P[hello], (!!str)::hello-world-comment\n",
		},
		skipForGoccy: true,
	},
	{
		description:    "Where is the comment - array example",
		subdescription: "The underlying yaml parser can assign comments in a document to surprising nodes. Use an expression like this to find where you comment is. 'p' indicates the path, 'isKey' is if the node is a map key (as opposed to a map value).\nFrom this, you can see the 'under-name-comment' is actually on the first child",
		document:       "name:\n  # under-name-comment\n  - first-array-child",
		expression:     `[... | {"p": path | join("."), "isKey": is_key, "hc": headComment, "lc": lineComment, "fc": footComment}]`,
		expected: []string{
			expectedWhereIsMyCommentArray,
		},
		skipForGoccy: true,
	},
	{
		description:    "Retrieve comment - array example",
		subdescription: "From the previous example, we know that the comment is on the first child as a headComment",
		document:       "name:\n  # under-name-comment\n  - first-array-child",
		expression:     `.name[0] | headComment`,
		expected: []string{
			"D0, P[name 0], (!!str)::under-name-comment\n",
		},
		skipForGoccy: true,
	},
	{
		description: "Set head comment",
		document:    `a: cat`,
		expression:  `. head_comment="single"`,
		expected: []string{
			"D0, P[], (!!map)::# single\na: cat\n",
		},
		skipForGoccy: true,
	},
	{
		description: "Set head comment of a map entry",
		document:    "f: foo\na:\n  b: cat",
		expression:  `(.a | key) head_comment="single"`,
		expected: []string{
			"D0, P[], (!!map)::f: foo\n# single\na:\n    b: cat\n",
		},
		skipForGoccy: true,
	},
	{
		description: "Set foot comment, using an expression",
		document:    `a: cat`,
		expression:  `. foot_comment=.a`,
		expected: []string{
			"D0, P[], (!!map)::a: cat\n# cat\n",
		},
		skipForGoccy: true,
	},
	{
		skipDoc:     true,
		description: "Set foot comment, using an expression",
		document:    "a: cat\n\n# hi",
		expression:  `. foot_comment=""`,
		expected: []string{
			"D0, P[], (!!map)::a: cat\n",
		},
		skipForGoccy: true,
	},
	{
		skipDoc:    true,
		document:   `a: cat`,
		expression: `. foot_comment=.b.d`,
		expected: []string{
			"D0, P[], (!!map)::a: cat\n",
		},
		skipForGoccy: true,
	},
	{
		skipDoc:    true,
		document:   `a: cat`,
		expression: `. foot_comment|=.b.d`,
		expected: []string{
			"D0, P[], (!!map)::a: cat\n",
		},
		skipForGoccy: true,
	},
	{
		description: "Remove comment",
		document:    "a: cat # comment\nb: dog # leave this",
		expression:  `.a line_comment=""`,
		expected: []string{
			"D0, P[], (!!map)::a: cat\nb: dog # leave this\n",
		},
		skipForGoccy: true,
	},
	{
		description:    "Remove (strip) all comments",
		subdescription: "Note the use of `...` to ensure key nodes are included.",
		document:       "# hi\n\na: cat # comment\n\n# great\n\nb: # key comment",
		expression:     `... comments=""`,
		expected: []string{
			"D0, P[], (!!map)::a: cat\nb:\n",
		},
		skipForGoccy: true,
	},
	{
		description: "Get line comment",
		document:    "# welcome!\n\na: cat # meow\n\n# have a great day",
		expression:  `.a | line_comment`,
		expected: []string{
			"D0, P[a], (!!str)::meow\n",
		},
		skipForGoccy: true,
	},
	{
		description:           "Get head comment",
		dontFormatInputForDoc: true,
		document:              "# welcome!\n\na: cat # meow\n\n# have a great day",
		expression:            `. | head_comment`,
		expected: []string{
			"D0, P[], (!!str)::welcome!\n\n",
		},
		skipForGoccy: true,
	},
	{
		skipDoc:     true,
		description: "strip trailing comment recurse all",
		document:    "a: cat\n\n# haha",
		expression:  `... comments= ""`,
		expected: []string{
			"D0, P[], (!!map)::a: cat\n",
		},
		skipForGoccy: true,
	},
	{
		skipDoc:     true,
		description: "strip trailing comment recurse values",
		document:    "a: cat\n\n# haha",
		expression:  `.. comments= ""`,
		expected: []string{
			"D0, P[], (!!map)::a: cat\n",
		},
		skipForGoccy: true,
	},
	{
		description:           "Head comment with document split",
		dontFormatInputForDoc: true,
		document:              "# welcome!\n---\n# bob\na: cat # meow\n\n# have a great day",
		expression:            `head_comment`,
		expected: []string{
			"D0, P[], (!!str)::welcome!\nbob\n",
		},
		skipForGoccy: true,
	},
	{
		description:           "Get foot comment",
		dontFormatInputForDoc: true,
		document:              "# welcome!\n\na: cat # meow\n\n# have a great day\n# no really",
		expression:            `. | foot_comment`,
		expected: []string{
			"D0, P[], (!!str)::have a great day\nno really\n",
		},
		skipForGoccy: true,
	},
	{
		description: "leading spaces",
		skipDoc:     true,
		document:    " # hi",
		expression:  `.`,
		expected: []string{
			"D0, P[], (!!null):: # hi\n",
		},
		skipForGoccy: true,
	},
	{
		description: "string spaces",
		skipDoc:     true,
		document:    "# hi\ncat\n",
		expression:  `.`,
		expected: []string{
			"D0, P[], (!!str)::# hi\ncat\n",
		},
		skipForGoccy: true,
	},
	{
		description: "leading spaces with new line",
		skipDoc:     true,
		document:    " # hi\n",
		expression:  `.`,
		expected: []string{
			"D0, P[], (!!null):: # hi\n",
		},
		skipForGoccy: true,
	},
	{
		description: "directive",
		skipDoc:     true,
		document:    "%YAML 1.1\n# hi\n",
		expression:  `.`,
		expected: []string{
			"D0, P[], (!!null)::%YAML 1.1\n# hi\n",
		},
		skipForGoccy: true,
	},
}

// Goccy-specific comment operator scenarios - these validate the actual behaviour of Goccy parser
var goccyCommentOperatorScenarios = []expressionScenario{
	{
		description:    "Goccy: Set line comment - basic functionality",
		subdescription: "Goccy parser maintains comment functionality with potentially different placement",
		document:       `a: cat`,
		expression:     `.a line_comment="single"`,
		expected: []string{
			"D0, P[], (!!map)::a: cat # single\n",
		},
		skipForYamlV3: true,
	},
	{
		description:    "Goccy: Get line comment - basic functionality",
		subdescription: "Goccy parser can read existing line comments",
		document:       "a: cat # meow",
		expression:     `.a | line_comment`,
		expected: []string{
			"D0, P[a], (!!str)::meow\n",
		},
		skipForYamlV3: true,
	},
	{
		description:    "Goccy: Set head comment - basic functionality",
		subdescription: "Goccy parser can set head comments",
		document:       `a: cat`,
		expression:     `. head_comment="header"`,
		expected: []string{
			"D0, P[], (!!map)::# header\na: cat\n",
		},
		skipForYamlV3: true,
	},
	{
		description:    "Goccy: Get head comment - basic functionality",
		subdescription: "Goccy parser can read existing head comments",
		document:       "# welcome!\na: cat",
		expression:     `. | head_comment`,
		expected: []string{
			"D0, P[], (!!str)::welcome!\n",
		},
		skipForYamlV3: true,
	},
	{
		description:    "Goccy: Set foot comment - basic functionality",
		subdescription: "Goccy parser can set foot comments",
		document:       `a: cat`,
		expression:     `. foot_comment="footer"`,
		expected: []string{
			"D0, P[], (!!map)::a: cat\n# footer\n",
		},
		skipForYamlV3: true,
	},
	{
		description:    "Goccy: Remove line comment",
		subdescription: "Goccy parser can remove comments",
		document:       "a: cat # remove_me\nb: dog # keep_me",
		expression:     `.a line_comment=""`,
		expected: []string{
			"D0, P[], (!!map)::a: cat\nb: dog # keep_me\n",
		},
		skipForYamlV3: true,
	},
	{
		description:    "Goccy: Comment preservation during data operations",
		subdescription: "Goccy parser preserves structural integrity while handling comments",
		document:       "# header\na: cat # inline\nb: dog\n# footer",
		expression:     `.c = "new"`,
		expected: []string{
			"D0, P[], (!!map)::# header\na: cat # inline\nb: dog\nc: new\n# footer\n",
		},
		skipForYamlV3: true,
	},
}

func TestCommentOperatorScenarios(t *testing.T) {
	for _, tt := range commentOperatorScenarios {
		testScenario(t, &tt)
	}
	for _, tt := range goccyCommentOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "comment-operators", commentOperatorScenarios)
}
