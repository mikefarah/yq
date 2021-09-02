package yqlib

import (
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var variableWithNewLine = `"cat
"`

var pathTests = []struct {
	path            string
	expectedTokens  []interface{}
	expectedPostFix []interface{}
}{
	{
		"0x12",
		append(make([]interface{}, 0), "18 (int64)"),
		append(make([]interface{}, 0), "18 (int64)"),
	},
	{
		"0X12",
		append(make([]interface{}, 0), "18 (int64)"),
		append(make([]interface{}, 0), "18 (int64)"),
	},
	{
		".a\n",
		append(make([]interface{}, 0), "a"),
		append(make([]interface{}, 0), "a"),
	},
	{
		variableWithNewLine,
		append(make([]interface{}, 0), "cat\n (string)"),
		append(make([]interface{}, 0), "cat\n (string)"),
	},
	{
		`.[0]`,
		append(make([]interface{}, 0), "SELF", "TRAVERSE_ARRAY", "[", "0 (int64)", "]"),
		append(make([]interface{}, 0), "SELF", "0 (int64)", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY"),
	},
	{
		`.[0][1]`,
		append(make([]interface{}, 0), "SELF", "TRAVERSE_ARRAY", "[", "0 (int64)", "]", "TRAVERSE_ARRAY", "[", "1 (int64)", "]"),
		append(make([]interface{}, 0), "SELF", "0 (int64)", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY", "1 (int64)", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY"),
	},
	{
		`"\""`,
		append(make([]interface{}, 0), "\" (string)"),
		append(make([]interface{}, 0), "\" (string)"),
	},
	{
		`[]|join(".")`,
		append(make([]interface{}, 0), "[", "EMPTY", "]", "PIPE", "JOIN", "(", ". (string)", ")"),
		append(make([]interface{}, 0), "EMPTY", "COLLECT", "SHORT_PIPE", ". (string)", "JOIN", "PIPE"),
	},
	{
		`{"cool": .b or .c}`,
		append(make([]interface{}, 0), "{", "cool (string)", "CREATE_MAP", "b", "OR", "c", "}"),
		append(make([]interface{}, 0), "cool (string)", "b", "c", "OR", "CREATE_MAP", "COLLECT_OBJECT", "SHORT_PIPE"),
	},
	{
		`{"cool": []|join(".")}`,
		append(make([]interface{}, 0), "{", "cool (string)", "CREATE_MAP", "[", "EMPTY", "]", "PIPE", "JOIN", "(", ". (string)", ")", "}"),
		append(make([]interface{}, 0), "cool (string)", "EMPTY", "COLLECT", "SHORT_PIPE", ". (string)", "JOIN", "PIPE", "CREATE_MAP", "COLLECT_OBJECT", "SHORT_PIPE"),
	},
	{
		`.a as $item ireduce (0; . + $item)`, // note - add code to shuffle reduce to this position for postfix
		append(make([]interface{}, 0), "a", "ASSIGN_VARIABLE", "GET_VARIABLE", "REDUCE", "(", "0 (int64)", "BLOCK", "SELF", "ADD", "GET_VARIABLE", ")"),
		append(make([]interface{}, 0), "a", "GET_VARIABLE", "ASSIGN_VARIABLE", "0 (int64)", "SELF", "GET_VARIABLE", "ADD", "BLOCK", "REDUCE"),
	},
	{
		`.a | .b | .c`,
		append(make([]interface{}, 0), "a", "PIPE", "b", "PIPE", "c"),
		append(make([]interface{}, 0), "a", "b", "c", "PIPE", "PIPE"),
	},
	{
		`[]`,
		append(make([]interface{}, 0), "[", "EMPTY", "]"),
		append(make([]interface{}, 0), "EMPTY", "COLLECT", "SHORT_PIPE"),
	},
	{
		`{}`,
		append(make([]interface{}, 0), "{", "EMPTY", "}"),
		append(make([]interface{}, 0), "EMPTY", "COLLECT_OBJECT", "SHORT_PIPE"),
	},
	{
		`[{}]`,
		append(make([]interface{}, 0), "[", "{", "EMPTY", "}", "]"),
		append(make([]interface{}, 0), "EMPTY", "COLLECT_OBJECT", "SHORT_PIPE", "COLLECT", "SHORT_PIPE"),
	},
	{
		`.realnames as $names | $names["anon"]`,
		append(make([]interface{}, 0), "realnames", "ASSIGN_VARIABLE", "GET_VARIABLE", "PIPE", "GET_VARIABLE", "TRAVERSE_ARRAY", "[", "anon (string)", "]"),
		append(make([]interface{}, 0), "realnames", "GET_VARIABLE", "ASSIGN_VARIABLE", "GET_VARIABLE", "anon (string)", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY", "PIPE"),
	},
	{
		`.b[.a]`,
		append(make([]interface{}, 0), "b", "TRAVERSE_ARRAY", "[", "a", "]"),
		append(make([]interface{}, 0), "b", "a", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY"),
	},
	{
		`.b[.a]?`,
		append(make([]interface{}, 0), "b", "TRAVERSE_ARRAY", "[", "a", "]"),
		append(make([]interface{}, 0), "b", "a", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY"),
	},
	{
		`.[]`,
		append(make([]interface{}, 0), "SELF", "TRAVERSE_ARRAY", "[", "EMPTY", "]"),
		append(make([]interface{}, 0), "SELF", "EMPTY", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY"),
	},
	{
		`.a[]`,
		append(make([]interface{}, 0), "a", "TRAVERSE_ARRAY", "[", "EMPTY", "]"),
		append(make([]interface{}, 0), "a", "EMPTY", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY"),
	},
	{
		`.a[]?`,
		append(make([]interface{}, 0), "a", "TRAVERSE_ARRAY", "[", "EMPTY", "]"),
		append(make([]interface{}, 0), "a", "EMPTY", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY"),
	},
	{
		`.a.[]`,
		append(make([]interface{}, 0), "a", "TRAVERSE_ARRAY", "[", "EMPTY", "]"),
		append(make([]interface{}, 0), "a", "EMPTY", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY"),
	},
	{
		`.a[0]`,
		append(make([]interface{}, 0), "a", "TRAVERSE_ARRAY", "[", "0 (int64)", "]"),
		append(make([]interface{}, 0), "a", "0 (int64)", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY"),
	},
	{
		`.a[0]?`,
		append(make([]interface{}, 0), "a", "TRAVERSE_ARRAY", "[", "0 (int64)", "]"),
		append(make([]interface{}, 0), "a", "0 (int64)", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY"),
	},
	{
		`.a.[0]`,
		append(make([]interface{}, 0), "a", "TRAVERSE_ARRAY", "[", "0 (int64)", "]"),
		append(make([]interface{}, 0), "a", "0 (int64)", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY"),
	},
	{
		`.a[].c`,
		append(make([]interface{}, 0), "a", "TRAVERSE_ARRAY", "[", "EMPTY", "]", "SHORT_PIPE", "c"),
		append(make([]interface{}, 0), "a", "EMPTY", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY", "c", "SHORT_PIPE"),
	},
	{
		`[3]`,
		append(make([]interface{}, 0), "[", "3 (int64)", "]"),
		append(make([]interface{}, 0), "3 (int64)", "COLLECT", "SHORT_PIPE"),
	},
	{
		`.key.array + .key.array2`,
		append(make([]interface{}, 0), "key", "SHORT_PIPE", "array", "ADD", "key", "SHORT_PIPE", "array2"),
		append(make([]interface{}, 0), "key", "array", "SHORT_PIPE", "key", "array2", "SHORT_PIPE", "ADD"),
	},
	{
		`.key.array * .key.array2`,
		append(make([]interface{}, 0), "key", "SHORT_PIPE", "array", "MULTIPLY", "key", "SHORT_PIPE", "array2"),
		append(make([]interface{}, 0), "key", "array", "SHORT_PIPE", "key", "array2", "SHORT_PIPE", "MULTIPLY"),
	},
	{
		`.key.array // .key.array2`,
		append(make([]interface{}, 0), "key", "SHORT_PIPE", "array", "ALTERNATIVE", "key", "SHORT_PIPE", "array2"),
		append(make([]interface{}, 0), "key", "array", "SHORT_PIPE", "key", "array2", "SHORT_PIPE", "ALTERNATIVE"),
	},
	{
		`.a | .[].b == "apple"`,
		append(make([]interface{}, 0), "a", "PIPE", "SELF", "TRAVERSE_ARRAY", "[", "EMPTY", "]", "SHORT_PIPE", "b", "EQUALS", "apple (string)"),
		append(make([]interface{}, 0), "a", "SELF", "EMPTY", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY", "b", "SHORT_PIPE", "apple (string)", "EQUALS", "PIPE"),
	},
	{
		`(.a | .[].b) == "apple"`,
		append(make([]interface{}, 0), "(", "a", "PIPE", "SELF", "TRAVERSE_ARRAY", "[", "EMPTY", "]", "SHORT_PIPE", "b", ")", "EQUALS", "apple (string)"),
		append(make([]interface{}, 0), "a", "SELF", "EMPTY", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY", "b", "SHORT_PIPE", "PIPE", "apple (string)", "EQUALS"),
	},
	{
		`.[] | select(. == "*at")`,
		append(make([]interface{}, 0), "SELF", "TRAVERSE_ARRAY", "[", "EMPTY", "]", "PIPE", "SELECT", "(", "SELF", "EQUALS", "*at (string)", ")"),
		append(make([]interface{}, 0), "SELF", "EMPTY", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY", "SELF", "*at (string)", "EQUALS", "SELECT", "PIPE"),
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
		`{.a: .c, .b.[]: .f.g[]}`,
		append(make([]interface{}, 0), "{", "a", "CREATE_MAP", "c", "UNION", "b", "TRAVERSE_ARRAY", "[", "EMPTY", "]", "CREATE_MAP", "f", "SHORT_PIPE", "g", "TRAVERSE_ARRAY", "[", "EMPTY", "]", "}"),
		append(make([]interface{}, 0), "a", "c", "CREATE_MAP", "b", "EMPTY", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY", "f", "g", "EMPTY", "COLLECT", "SHORT_PIPE", "TRAVERSE_ARRAY", "SHORT_PIPE", "CREATE_MAP", "UNION", "COLLECT_OBJECT", "SHORT_PIPE"),
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
			tokenValues = append(tokenValues, token.toString(false))
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
