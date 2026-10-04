package yqlib

import (
	"bytes"
	"strings"
	"testing"
)

// TestYamlEncoderUnwrapScalarRoundtripSafety verifies that a top-level string
// scalar whose unquoted form would re-parse as a non-scalar node (map or
// sequence) is emitted quoted even when UnwrapScalar is enabled. Safe plain
// strings continue to round-trip through the existing fast-path. See #2608.
func TestYamlEncoderUnwrapScalarRoundtripSafety(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		wantBare bool // true: output equals value+"\n"; false: output must differ
	}{
		{name: "colon_parses_as_map", value: "this: should really work"},
		{name: "dash_parses_as_seq", value: "- item"},
		{name: "multiline_maplike", value: "a: a\nb: b"},
		{name: "safe_plain_string", value: "hello world", wantBare: true},
		{name: "safe_identifier", value: "cat", wantBare: true},
		{name: "safe_digits_preserved", value: "123", wantBare: true},
		{name: "safe_null_word_preserved", value: "null", wantBare: true},
		{name: "safe_tag_shorthand_preserved", value: "!!int", wantBare: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			prefs := NewDefaultYamlPreferences()
			prefs.UnwrapScalar = true

			var buf bytes.Buffer
			err := NewYamlEncoder(prefs).Encode(&buf, &CandidateNode{
				Kind:  ScalarNode,
				Tag:   "!!str",
				Value: tc.value,
			})
			if err != nil {
				t.Fatalf("encode failed: %v", err)
			}
			got := buf.String()

			if tc.wantBare {
				if got != tc.value+"\n" {
					t.Fatalf("expected bare %q, got %q", tc.value+"\n", got)
				}
				return
			}

			// Ambiguous input: must not be emitted as the bare value.
			if got == tc.value+"\n" {
				t.Fatalf("value %q was emitted bare; expected quoted form", tc.value)
			}

			// The output must round-trip back to a string scalar with the
			// same value, proving structural roundtrip safety.
			decoder := NewYamlDecoder(NewDefaultYamlPreferences())
			nodes, err := readDocuments(strings.NewReader(got), "test.yaml", 0, decoder)
			if err != nil {
				t.Fatalf("decode of %q failed: %v", got, err)
			}
			if nodes.Len() != 1 {
				t.Fatalf("expected one document, got %d", nodes.Len())
			}
			candidate := nodes.Front().Value.(*CandidateNode)
			// readDocuments wraps the document; descend to the scalar.
			scalar := candidate
			for scalar.Kind != ScalarNode && len(scalar.Content) == 1 {
				scalar = scalar.Content[0]
			}
			if scalar.Kind != ScalarNode {
				t.Fatalf("round-tripped node is not a scalar: kind=%v value=%q", scalar.Kind, scalar.Value)
			}
			if scalar.Tag != "!!str" {
				t.Fatalf("round-tripped tag is %q, want !!str", scalar.Tag)
			}
			if scalar.Value != tc.value {
				t.Fatalf("round-tripped value is %q, want %q", scalar.Value, tc.value)
			}
		})
	}
}
