package yqlib

import (
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var pathTests = []struct {
	path            string
	expectedTokens  []interface{}
	expectedPostFix []interface{}
}{
	{
		`[]`,
		append(make([]interface{}, 0), "[", "]"),
		append(make([]interface{}, 0), "EMPTY", "COLLECT", "SHORT_PIPE"),
	},
	{
		`.[]`,
		append(make([]interface{}, 0), "TRAVERSE_ARRAY", "[", "]"),
		append(make([]interface{}, 0), "EMPTY", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY"),
	},
	{
		`.a[]`,
		append(make([]interface{}, 0), "a", "SHORT_PIPE", "TRAVERSE_ARRAY", "[", "]"),
		append(make([]interface{}, 0), "a", "EMPTY", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY", "SHORT_PIPE"),
	},
	{
		`.a.[]`,
		append(make([]interface{}, 0), "a", "SHORT_PIPE", "TRAVERSE_ARRAY", "[", "]"),
		append(make([]interface{}, 0), "a", "EMPTY", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY", "SHORT_PIPE"),
	},
	{
		`.a[0]`,
		append(make([]interface{}, 0), "a", "SHORT_PIPE", "TRAVERSE_ARRAY", "[", "0 (int64)", "]"),
		append(make([]interface{}, 0), "a", "0 (int64)", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY", "SHORT_PIPE"),
	},
	{
		`.a.[0]`,
		append(make([]interface{}, 0), "a", "SHORT_PIPE", "TRAVERSE_ARRAY", "[", "0 (int64)", "]"),
		append(make([]interface{}, 0), "a", "0 (int64)", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY", "SHORT_PIPE"),
	},
	{
		`.a[].c`,
		append(make([]interface{}, 0), "a", "SHORT_PIPE", "TRAVERSE_ARRAY", "[", "]", "SHORT_PIPE", "c"),
		append(make([]interface{}, 0), "a", "EMPTY", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY", "SHORT_PIPE", "c", "SHORT_PIPE"),
	},
	{
		`[3]`,
		append(make([]interface{}, 0), "[", "3 (int64)", "]"),
		append(make([]interface{}, 0), "3 (int64)", "COLLECT", "SHORT_PIPE"),
	},
	{
		`.a | .[].b == "apple"`,
		append(make([]interface{}, 0), "a", "PIPE", "TRAVERSE_ARRAY", "[", "]", "SHORT_PIPE", "b", "EQUALS", "apple (string)"),
		append(make([]interface{}, 0), "a", "EMPTY", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY", "b", "SHORT_PIPE", "apple (string)", "EQUALS", "PIPE"),
	},
	{
		`(.a | .[].b) == "apple"`,
		append(make([]interface{}, 0), "(", "a", "PIPE", "TRAVERSE_ARRAY", "[", "]", "SHORT_PIPE", "b", ")", "EQUALS", "apple (string)"),
		append(make([]interface{}, 0), "a", "EMPTY", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY", "b", "SHORT_PIPE", "PIPE", "apple (string)", "EQUALS"),
	},
	{
		`.[] | select(. == "*at")`,
		append(make([]interface{}, 0), "TRAVERSE_ARRAY", "[", "]", "PIPE", "SELECT", "(", "SELF", "EQUALS", "*at (string)", ")"),
		append(make([]interface{}, 0), "EMPTY", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY", "SELF", "*at (string)", "EQUALS", "SELECT", "PIPE"),
	},
	{
		`[true]`,
		append(make([]interface{}, 0), "[", "true (bool)", "]"),
		append(make([]interface{}, 0), "true (bool)", "COLLECT", "SHORT_PIPE"),
	},
	{
		`[true, false]`,
		append(make([]interface{}, 0), "[", "true (bool)", "UNION", "false (bool)", "]"),
		append(make([]interface{}, 0), "true (bool)", "false (bool)", "UNION", "COLLECT", "SHORT_PIPE"),
	},
	{
		`"mike": .a`,
		append(make([]interface{}, 0), "mike (string)", "CREATE_MAP", "a"),
		append(make([]interface{}, 0), "mike (string)", "a", "CREATE_MAP"),
	},
	{
		`.a: "mike"`,
		append(make([]interface{}, 0), "a", "CREATE_MAP", "mike (string)"),
		append(make([]interface{}, 0), "a", "mike (string)", "CREATE_MAP"),
	},
	{
		`{"mike": .a}`,
		append(make([]interface{}, 0), "{", "mike (string)", "CREATE_MAP", "a", "}"),
		append(make([]interface{}, 0), "mike (string)", "a", "CREATE_MAP", "COLLECT_OBJECT", "SHORT_PIPE"),
	},
	{
		`{.a: "mike"}`,
		append(make([]interface{}, 0), "{", "a", "CREATE_MAP", "mike (string)", "}"),
		append(make([]interface{}, 0), "a", "mike (string)", "CREATE_MAP", "COLLECT_OBJECT", "SHORT_PIPE"),
	},
	{
		`{.a: .c, .b.[]: .f.g.[]}`,
		append(make([]interface{}, 0), "{", "a", "CREATE_MAP", "c", "UNION", "b", "SHORT_PIPE", "TRAVERSE_ARRAY", "[", "]", "CREATE_MAP", "f", "SHORT_PIPE", "g", "SHORT_PIPE", "TRAVERSE_ARRAY", "[", "]", "}"),
		append(make([]interface{}, 0), "a", "c", "CREATE_MAP", "b", "EMPTY", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY", "SHORT_PIPE", "f", "g", "SHORT_PIPE", "EMPTY", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY", "SHORT_PIPE", "CREATE_MAP", "UNION", "COLLECT_OBJECT", "SHORT_PIPE"),
	},
	{
		`explode(.a.b)`,
		append(make([]interface{}, 0), "EXPLODE", "(", "a", "SHORT_PIPE", "b", ")"),
		append(make([]interface{}, 0), "a", "b", "SHORT_PIPE", "EXPLODE"),
	},
	{
		`.a.b style="folded"`,
		append(make([]interface{}, 0), "a", "SHORT_PIPE", "b", "ASSIGN_STYLE", "folded (string)"),
		append(make([]interface{}, 0), "a", "b", "SHORT_PIPE", "folded (string)", "ASSIGN_STYLE"),
	},
	{
		`tag == "str"`,
		append(make([]interface{}, 0), "GET_TAG", "EQUALS", "str (string)"),
		append(make([]interface{}, 0), "GET_TAG", "str (string)", "EQUALS"),
	},
	{
		`. tag= "str"`,
		append(make([]interface{}, 0), "SELF", "ASSIGN_TAG", "str (string)"),
		append(make([]interface{}, 0), "SELF", "str (string)", "ASSIGN_TAG"),
	},
	{
		`lineComment == "str"`,
		append(make([]interface{}, 0), "GET_COMMENT", "EQUALS", "str (string)"),
		append(make([]interface{}, 0), "GET_COMMENT", "str (string)", "EQUALS"),
	},
	{
		`. lineComment= "str"`,
		append(make([]interface{}, 0), "SELF", "ASSIGN_COMMENT", "str (string)"),
		append(make([]interface{}, 0), "SELF", "str (string)", "ASSIGN_COMMENT"),
	},
	{
		`. lineComment |= "str"`,
		append(make([]interface{}, 0), "SELF", "ASSIGN_COMMENT", "str (string)"),
		append(make([]interface{}, 0), "SELF", "str (string)", "ASSIGN_COMMENT"),
	},
	{
		`.a.b tag="!!str"`,
		append(make([]interface{}, 0), "a", "SHORT_PIPE", "b", "ASSIGN_TAG", "!!str (string)"),
		append(make([]interface{}, 0), "a", "b", "SHORT_PIPE", "!!str (string)", "ASSIGN_TAG"),
	},
	{
		`""`,
		append(make([]interface{}, 0), " (string)"),
		append(make([]interface{}, 0), " (string)"),
	},
	{
		`.foo* | (. style="flow")`,
		append(make([]interface{}, 0), "foo*", "PIPE", "(", "SELF", "ASSIGN_STYLE", "flow (string)", ")"),
		append(make([]interface{}, 0), "foo*", "SELF", "flow (string)", "ASSIGN_STYLE", "PIPE"),
	},
	{
		`{}`,
		append(make([]interface{}, 0), "{", "}"),
		append(make([]interface{}, 0), "EMPTY", "COLLECT_OBJECT", "SHORT_PIPE"),
	},
}

var tokeniser = newExpressionTokeniser()
var postFixer = newExpressionPostFixer()

func TestPathParsing(t *testing.T) {
	for _, tt := range pathTests {
		tokens, err := tokeniser.Tokenise(tt.path)
		if err != nil {
			t.Error(tt.path, err)
		}
		var tokenValues []interface{}
		for _, token := range tokens {
			tokenValues = append(tokenValues, token.toString())
		}
		test.AssertResultComplexWithContext(t, tt.expectedTokens, tokenValues, fmt.Sprintf("tokenise: %v", tt.path))

		results, errorP := postFixer.ConvertToPostfix(tokens)

		var readableResults []interface{}
		for _, token := range results {
			readableResults = append(readableResults, token.toString())
		}

		if errorP != nil {
			t.Error(tt.path, err)
		}

		test.AssertResultComplexWithContext(t, tt.expectedPostFix, readableResults, fmt.Sprintf("postfix: %v", tt.path))

	}
}
