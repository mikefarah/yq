package treeops

import (
	"testing"

	"github.com/mikefarah/yq/v3/test"
)

var tokeniserTests = []struct {
	path           string
	expectedTokens []interface{}
}{ // TODO: Ensure ALL documented examples have tests! sheesh

	{"a OR (b OR c)", append(make([]interface{}, 0), "a", "OR", "(", "b", "OR", "c", ")")},
	{"a .- (b OR c)", append(make([]interface{}, 0), "a", " .- ", "(", "b", "OR", "c", ")")},
	{"(animal==3)", append(make([]interface{}, 0), "(", "animal", "==", int64(3), ")")},
	{"(animal==f3)", append(make([]interface{}, 0), "(", "animal", "==", "f3", ")")},
	{"apples.BANANAS", append(make([]interface{}, 0), "apples", ".", "BANANAS")},
	{"appl*.BANA*", append(make([]interface{}, 0), "appl*", ".", "BANA*")},
	{"a.b.**", append(make([]interface{}, 0), "a", ".", "b", ".", "**")},
	{"a.\"=\".frog", append(make([]interface{}, 0), "a", ".", "=", ".", "frog")},
	{"a.b.*", append(make([]interface{}, 0), "a", ".", "b", ".", "*")},
	{"a.b.thin*", append(make([]interface{}, 0), "a", ".", "b", ".", "thin*")},
	{"a.b[0]", append(make([]interface{}, 0), "a", ".", "b", ".", int64(0))},
	{"a.b.[0]", append(make([]interface{}, 0), "a", ".", "b", ".", int64(0))},
	{"a.b[*]", append(make([]interface{}, 0), "a", ".", "b", ".", "[*]")},
	{"a.b.[*]", append(make([]interface{}, 0), "a", ".", "b", ".", "[*]")},
	{"a.b[+]", append(make([]interface{}, 0), "a", ".", "b", ".", "[+]")},
	{"a.b.[+]", append(make([]interface{}, 0), "a", ".", "b", ".", "[+]")},
	{"a.b[-12]", append(make([]interface{}, 0), "a", ".", "b", ".", int64(-12))},
	{"a.b.0", append(make([]interface{}, 0), "a", ".", "b", ".", int64(0))},
	// {"a.b.-12", append(make([]interface{}, 0), "a", ".", "b", ".", int64(-12))},
	{"a", append(make([]interface{}, 0), "a")},
	{"\"a.b\".c", append(make([]interface{}, 0), "a.b", ".", "c")},
	{`b."foo.bar"`, append(make([]interface{}, 0), "b", ".", "foo.bar")},
	{"animals(.==cat)", append(make([]interface{}, 0), "animals", ".", "(", ".==", "cat", ")")},
	{"animals.(.==cat)", append(make([]interface{}, 0), "animals", ".", "(", ".==", "cat", ")")},
	{"animals(. == cat)", append(make([]interface{}, 0), "animals", ".", "(", ". == ", "cat", ")")},
	{"animals(.==c*)", append(make([]interface{}, 0), "animals", ".", "(", ".==", "c*", ")")},
	{"animals(a.b==c*)", append(make([]interface{}, 0), "animals", ".", "(", "a", ".", "b", "==", "c*", ")")},
	{"animals.(a.b==c*)", append(make([]interface{}, 0), "animals", ".", "(", "a", ".", "b", "==", "c*", ")")},
	{"(a.b==c*).animals", append(make([]interface{}, 0), "(", "a", ".", "b", "==", "c*", ")", ".", "animals")},
	{"(a.b==c*)animals", append(make([]interface{}, 0), "(", "a", ".", "b", "==", "c*", ")", ".", "animals")},
	{"[1].a.d", append(make([]interface{}, 0), int64(1), ".", "a", ".", "d")},
	{"[1]a.d", append(make([]interface{}, 0), int64(1), ".", "a", ".", "d")},
	{"a[0]c", append(make([]interface{}, 0), "a", ".", int64(0), ".", "c")},
	{"a.[0].c", append(make([]interface{}, 0), "a", ".", int64(0), ".", "c")},
	{"[0]", append(make([]interface{}, 0), int64(0))},
	{"0", append(make([]interface{}, 0), int64(0))},
	{"a.b[+]c", append(make([]interface{}, 0), "a", ".", "b", ".", "[+]", ".", "c")},
	{"a.cool(s.d.f == cool)", append(make([]interface{}, 0), "a", ".", "cool", ".", "(", "s", ".", "d", ".", "f", " == ", "cool", ")")},
	{"a.cool.(s.d.f==cool OR t.b.h==frog).caterpillar", append(make([]interface{}, 0), "a", ".", "cool", ".", "(", "s", ".", "d", ".", "f", "==", "cool", "OR", "t", ".", "b", ".", "h", "==", "frog", ")", ".", "caterpillar")},
	{"a.cool(s.d.f==cool and t.b.h==frog)*", append(make([]interface{}, 0), "a", ".", "cool", ".", "(", "s", ".", "d", ".", "f", "==", "cool", "and", "t", ".", "b", ".", "h", "==", "frog", ")", ".", "*")},
	{"a.cool(s.d.f==cool and t.b.h==frog).th*", append(make([]interface{}, 0), "a", ".", "cool", ".", "(", "s", ".", "d", ".", "f", "==", "cool", "and", "t", ".", "b", ".", "h", "==", "frog", ")", ".", "th*")},
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
