//go:build !yq_nojson5

package yqlib

import (
	"bytes"
	"testing"
)

func TestJSON5EncoderPrintDocumentSeparatorIsNoop(t *testing.T) {
	prefs := ConfiguredJSONPreferences.Copy()
	var buf bytes.Buffer

	if err := NewJSON5Encoder(prefs).PrintDocumentSeparator(&buf); err != nil {
		t.Fatalf("PrintDocumentSeparator returned error: %v", err)
	}
	if got := buf.String(); got != "" {
		t.Fatalf("expected no output, got %q", got)
	}
}

func TestJSON5EncoderUnwrapScalar(t *testing.T) {
	prefs := ConfiguredJSONPreferences.Copy()
	prefs.UnwrapScalar = true

	node := createStringScalarNode("hi")

	var buf bytes.Buffer
	if err := NewJSON5Encoder(prefs).Encode(&buf, node); err != nil {
		t.Fatalf("Encode returned error: %v", err)
	}

	if got, want := buf.String(), "hi\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestEncodeJSON5NodeNilAndUnknownKind(t *testing.T) {
	var buf bytes.Buffer
	if err := encodeJSON5Node(&buf, nil, "  ", 0); err != nil {
		t.Fatalf("encodeJSON5Node(nil) returned error: %v", err)
	}
	if got, want := buf.String(), "null"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	buf.Reset()
	if err := encodeJSON5Node(&buf, &CandidateNode{}, "  ", 0); err != nil {
		t.Fatalf("encodeJSON5Node(unknown kind) returned error: %v", err)
	}
	if got, want := buf.String(), "null"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestJSON5ScalarAndFloatFormatting(t *testing.T) {
	if got, want := json5FloatString(".inf"), "Infinity"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if got, want := json5FloatString("1e309"), "null"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if got, want := json5FloatString(".nan"), "NaN"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if got, want := json5FloatString("definitely-not-a-float"), "null"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	badInt := &CandidateNode{Kind: ScalarNode, Tag: "!!int", Value: "nope"}
	if got, want := json5ScalarString(badInt), "null"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	unknownTag := &CandidateNode{Kind: ScalarNode, Tag: "!!unknown", Value: "hi"}
	if got, want := json5ScalarString(unknownTag), "\"hi\""; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestJSON5EncoderEncodeSequenceMappingAndComments(t *testing.T) {
	prefs := ConfiguredJSONPreferences.Copy()
	prefs.Indent = 2
	prefs.UnwrapScalar = false
	prefs.ColorsEnabled = false

	aliasTarget := createScalarNode(3, "3")
	seq := &CandidateNode{
		Kind:        SequenceNode,
		HeadComment: "# root head\n\n// root second\n",
		LineComment: "# root inline",
		FootComment: "# foot",
		Content: []*CandidateNode{
			{
				Kind:        ScalarNode,
				Tag:         "!!int",
				Value:       "1",
				HeadComment: "# first",
				LineComment: "// lc1",
			},
			{
				Kind:        AliasNode,
				Alias:       aliasTarget,
				LineComment: "# alias",
			},
			{
				Kind:        MappingNode,
				HeadComment: "# map head",
				LineComment: "# map lc",
				Content: []*CandidateNode{
					{
						Kind:        ScalarNode,
						Tag:         "!!str",
						Value:       "a",
						HeadComment: "# key a",
						LineComment: "# key inline",
						IsMapKey:    true,
					},
					{
						Kind:        ScalarNode,
						Tag:         "!!float",
						Value:       ".inf",
						HeadComment: "# before a line1\n# before a line2",
						LineComment: "# value inline",
					},
					{
						Kind:     ScalarNode,
						Tag:      "!!str",
						Value:    "b",
						IsMapKey: true,
					},
					{
						Kind:        ScalarNode,
						Tag:         "!!bool",
						Value:       "true",
						HeadComment: "# before b",
					},
				},
			},
			createStringScalarNode("hi"),
		},
	}

	var buf bytes.Buffer
	if err := NewJSON5Encoder(prefs).Encode(&buf, seq); err != nil {
		t.Fatalf("Encode returned error: %v", err)
	}

	const expected = `// root head
// root second
[
  // first
  1 /* lc1 */,
  3 /* alias */,
  // map head
  {
    // key a
    "a" /* key inline */:
    // before a line1
    // before a line2
    Infinity /* value inline */,
    "b": /* before b */ true
  } /* map lc */,
  "hi"
] /* root inline */
// foot
`
	if got := buf.String(); got != expected {
		t.Fatalf("unexpected output:\n%s", got)
	}
}
