package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v3/test"
)

var tokeniserTests = []struct {
	path           string
	expectedTokens []interface{}
}{ // TODO: Ensure ALL documented examples have tests! sheesh

	{"apples.BANANAS", append(make([]interface{}, 0), "apples", "BANANAS")},
	{"a.b.**", append(make([]interface{}, 0), "a", "b", "**")},
	{"a.\"=\".frog", append(make([]interface{}, 0), "a", "=", "frog")},
	{"a.b.*", append(make([]interface{}, 0), "a", "b", "*")},
	{"a.b.thin*", append(make([]interface{}, 0), "a", "b", "thin*")},
	{"a.b[0]", append(make([]interface{}, 0), "a", "b", int64(0))},
	{"a.b[*]", append(make([]interface{}, 0), "a", "b", "[*]")},
	{"a.b[-12]", append(make([]interface{}, 0), "a", "b", int64(-12))},
	{"a.b.0", append(make([]interface{}, 0), "a", "b", int64(0))},
	{"a.b.d[+]", append(make([]interface{}, 0), "a", "b", "d", "[+]")},
	{"a", append(make([]interface{}, 0), "a")},
	{"\"a.b\".c", append(make([]interface{}, 0), "a.b", "c")},
	{`b."foo.bar"`, append(make([]interface{}, 0), "b", "foo.bar")},
	{"animals(.==cat)", append(make([]interface{}, 0), "animals", "(", "==", "cat", ")")}, // TODO validate this dot is not a join?
	{"animals(.==c*)", append(make([]interface{}, 0), "animals", "(", "==", "c*", ")")},   // TODO validate this dot is not a join?
	{"[1].a.d", append(make([]interface{}, 0), int64(1), "a", "d")},
	{"a[0].c", append(make([]interface{}, 0), "a", int64(0), "c")},
	{"[0]", append(make([]interface{}, 0), int64(0))},
	{"a.cool(s.d.f==cool)", append(make([]interface{}, 0), "a", "cool", "(", "s", "d", "f", "==", "cool", ")")},
	{"a.cool(s.d.f==cool OR t.b.h==frog).caterpillar", append(make([]interface{}, 0), "a", "cool", "(", "s", "d", "f", "==", "cool", "OR", "t", "b", "h", "==", "frog", ")", "caterpillar")},
	{"a.cool(s.d.f==cool and t.b.h==frog)*", append(make([]interface{}, 0), "a", "cool", "(", "s", "d", "f", "==", "cool", "and", "t", "b", "h", "==", "frog", ")", "*")},
	{"a.cool(s.d.f==cool and t.b.h==frog).th*", append(make([]interface{}, 0), "a", "cool", "(", "s", "d", "f", "==", "cool", "and", "t", "b", "h", "==", "frog", ")", "th*")},
}

var tokeniser = NewPathTokeniser()

func TestTokeniser(t *testing.T) {
	for _, tt := range tokeniserTests {
		tokens, err := tokeniser.Tokenise(tt.path)
		if err != nil {
			t.Error(tt.path, err)
		}
		var tokenValues []interface{}
		for _, token := range tokens {
			tokenValues = append(tokenValues, token.Value)
		}
		test.AssertResultComplex(t, tt.expectedTokens, tokenValues)
	}
}
