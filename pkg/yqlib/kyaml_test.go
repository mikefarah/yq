//go:build !yq_nokyaml

package yqlib

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string {
	return ansiRe.ReplaceAllString(s, "")
}

var kyamlFormatScenarios = []formatScenario{
	{
		description:    "Encode kyaml: plain string scalar",
		subdescription: "Strings are always double-quoted in KYaml output.",
		scenarioType:   "encode",
		indent:         2,
		input:          "cat\n",
		expected:       "\"cat\"\n",
	},
	{
		description:  "encode plain int scalar",
		scenarioType: "encode",
		indent:       2,
		input:        "12\n",
		expected:     "12\n",
		skipDoc:      true,
	},
	{
		description:  "encode plain bool scalar",
		scenarioType: "encode",
		indent:       2,
		input:        "true\n",
		expected:     "true\n",
		skipDoc:      true,
	},
	{
		description:  "encode plain null scalar",
		scenarioType: "encode",
		indent:       2,
		input:        "null\n",
		expected:     "null\n",
		skipDoc:      true,
	},
	{
		description:  "encode flow mapping and sequence",
		scenarioType: "encode",
		indent:       2,
		input:        "a: b\nc:\n  - d\n",
		expected: "{\n" +
			"  a: \"b\",\n" +
			"  c: [\n" +
			"    \"d\",\n" +
			"  ],\n" +
			"}\n",
	},
	{
		description:  "encode non-string scalars",
		scenarioType: "encode",
		indent:       2,
		input: "a: 12\n" +
			"b: true\n" +
			"c: null\n" +
			"d: \"true\"\n",
		expected: "{\n" +
			"  a: 12,\n" +
			"  b: true,\n" +
			"  c: null,\n" +
			"  d: \"true\",\n" +
			"}\n",
	},
	{
		description:  "quote non-identifier keys",
		scenarioType: "encode",
		indent:       2,
		input:        "\"1a\": b\n\"has space\": c\n",
		expected: "{\n" +
			"  \"1a\": \"b\",\n" +
			"  \"has space\": \"c\",\n" +
			"}\n",
	},
	{
		description:  "escape quoted strings",
		scenarioType: "encode",
		indent:       2,
		input:        "a: \"line1\\nline2\\t\\\"q\\\"\"\n",
		expected: "{\n" +
			"  a: \"line1\\nline2\\t\\\"q\\\"\",\n" +
			"}\n",
	},
	{
		description:  "preserve comments when encoding",
		scenarioType: "encode",
		indent:       2,
		input: "# leading\n" +
			"a: 1 # a line\n" +
			"# head b\n" +
			"b: 2\n" +
			"c:\n" +
			"  # head d\n" +
			"  - d # d line\n" +
			"  - e\n" +
			"# trailing\n",
		expected: "# leading\n" +
			"{\n" +
			"  a: 1, # a line\n" +
			"  # head b\n" +
			"  b: 2,\n" +
			"  c: [\n" +
			"    # head d\n" +
			"    \"d\", # d line\n" +
			"    \"e\",\n" +
			"  ],\n" +
			"  # trailing\n" +
			"}\n",
	},
	{
		description:    "Encode kyaml: anchors and aliases",
		subdescription: "KYaml output does not support anchors/aliases; they are expanded to concrete values.",
		scenarioType:   "encode",
		indent:         2,
		input: "base: &base\n" +
			"  a: b\n" +
			"copy: *base\n",
		expected: "{\n" +
			"  base: {\n" +
			"    a: \"b\",\n" +
			"  },\n" +
			"  copy: {\n" +
			"    a: \"b\",\n" +
			"  },\n" +
			"}\n",
	},
	{
		description:    "Encode kyaml: yaml to kyaml shows formatting differences",
		subdescription: "KYaml uses flow-style collections (braces/brackets) and explicit commas.",
		scenarioType:   "encode",
		indent:         2,
		input: "person:\n" +
			"  name: John\n" +
			"  pets:\n" +
			"    - cat\n" +
			"    - dog\n",
		expected: "{\n" +
			"  person: {\n" +
			"    name: \"John\",\n" +
			"    pets: [\n" +
			"      \"cat\",\n" +
			"      \"dog\",\n" +
			"    ],\n" +
			"  },\n" +
			"}\n",
	},
	{
		description:    "Encode kyaml: nested lists of objects",
		subdescription: "Lists and objects can be nested arbitrarily; KYaml always uses flow-style collections.",
		scenarioType:   "encode",
		indent:         2,
		input: "- name: a\n" +
			"  items:\n" +
			"    - id: 1\n" +
			"      tags:\n" +
			"        - k: x\n" +
			"          v: y\n" +
			"        - k: x2\n" +
			"          v: y2\n" +
			"    - id: 2\n" +
			"      tags:\n" +
			"        - k: z\n" +
			"          v: w\n",
		expected: "[\n" +
			"  {\n" +
			"    name: \"a\",\n" +
			"    items: [\n" +
			"      {\n" +
			"        id: 1,\n" +
			"        tags: [\n" +
			"          {\n" +
			"            k: \"x\",\n" +
			"            v: \"y\",\n" +
			"          },\n" +
			"          {\n" +
			"            k: \"x2\",\n" +
			"            v: \"y2\",\n" +
			"          },\n" +
			"        ],\n" +
			"      },\n" +
			"      {\n" +
			"        id: 2,\n" +
			"        tags: [\n" +
			"          {\n" +
			"            k: \"z\",\n" +
			"            v: \"w\",\n" +
			"          },\n" +
			"        ],\n" +
			"      },\n" +
			"    ],\n" +
			"  },\n" +
			"]\n",
	},
}

func testKYamlScenario(t *testing.T, s formatScenario) {
	prefs := ConfiguredKYamlPreferences.Copy()
	prefs.Indent = s.indent
	prefs.UnwrapScalar = false

	switch s.scenarioType {
	case "encode":
		test.AssertResultWithContext(
			t,
			s.expected,
			mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewKYamlEncoder(prefs)),
			s.description,
		)
	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentKYamlScenario(_ *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(formatScenario)
	if s.skipDoc {
		return
	}

	switch s.scenarioType {
	case "encode":
		documentKYamlEncodeScenario(w, s)
	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentKYamlEncodeScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.yml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```yaml\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")

	expression := s.expression
	if expression == "" {
		expression = "."
	}

	if s.indent == 2 {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -o=kyaml '%v' sample.yml\n```\n", expression))
	} else {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -o=kyaml -I=%v '%v' sample.yml\n```\n", s.indent, expression))
	}

	writeOrPanic(w, "will output\n")

	prefs := ConfiguredKYamlPreferences.Copy()
	prefs.Indent = s.indent
	prefs.UnwrapScalar = false

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewKYamlEncoder(prefs))))
}

func TestKYamlFormatScenarios(t *testing.T) {
	for _, s := range kyamlFormatScenarios {
		testKYamlScenario(t, s)
	}

	genericScenarios := make([]interface{}, len(kyamlFormatScenarios))
	for i, s := range kyamlFormatScenarios {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "kyaml", genericScenarios, documentKYamlScenario)
}

func TestKYamlEncoderPrintDocumentSeparator(t *testing.T) {
	t.Run("enabled", func(t *testing.T) {
		prefs := NewDefaultKYamlPreferences()
		prefs.PrintDocSeparators = true

		var buf bytes.Buffer
		err := NewKYamlEncoder(prefs).PrintDocumentSeparator(&buf)
		if err != nil {
			t.Fatal(err)
		}
		if buf.String() != "---\n" {
			t.Fatalf("expected doc separator, got %q", buf.String())
		}
	})

	t.Run("disabled", func(t *testing.T) {
		prefs := NewDefaultKYamlPreferences()
		prefs.PrintDocSeparators = false

		var buf bytes.Buffer
		err := NewKYamlEncoder(prefs).PrintDocumentSeparator(&buf)
		if err != nil {
			t.Fatal(err)
		}
		if buf.String() != "" {
			t.Fatalf("expected no output, got %q", buf.String())
		}
	})
}

func TestKYamlEncoderEncodeUnwrapScalar(t *testing.T) {
	prefs := NewDefaultKYamlPreferences()
	prefs.UnwrapScalar = true

	var buf bytes.Buffer
	err := NewKYamlEncoder(prefs).Encode(&buf, &CandidateNode{
		Kind:  ScalarNode,
		Tag:   "!!str",
		Value: "cat",
	})
	if err != nil {
		t.Fatal(err)
	}
	if buf.String() != "cat\n" {
		t.Fatalf("expected unwrapped scalar, got %q", buf.String())
	}
}

func TestKYamlEncoderEncodeColorsEnabled(t *testing.T) {
	prefs := NewDefaultKYamlPreferences()
	prefs.UnwrapScalar = false
	prefs.ColorsEnabled = true

	var buf bytes.Buffer
	err := NewKYamlEncoder(prefs).Encode(&buf, &CandidateNode{
		Kind: MappingNode,
		Content: []*CandidateNode{
			{Kind: ScalarNode, Tag: "!!str", Value: "a"},
			{Kind: ScalarNode, Tag: "!!str", Value: "b"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	out := stripANSI(buf.String())
	if !strings.Contains(out, "a:") || !strings.Contains(out, "\"b\"") {
		t.Fatalf("expected colourised output to contain rendered tokens, got %q", out)
	}
}

func TestKYamlEncoderWriteNodeAliasAndUnknown(t *testing.T) {
	ke := NewKYamlEncoder(NewDefaultKYamlPreferences()).(*kyamlEncoder)

	t.Run("alias_nil", func(t *testing.T) {
		var buf bytes.Buffer
		err := ke.writeNode(&buf, &CandidateNode{Kind: AliasNode}, 0)
		if err != nil {
			t.Fatal(err)
		}
		if buf.String() != "null" {
			t.Fatalf("expected null for nil alias, got %q", buf.String())
		}
	})

	t.Run("alias_value", func(t *testing.T) {
		var buf bytes.Buffer
		err := ke.writeNode(&buf, &CandidateNode{
			Kind:  AliasNode,
			Alias: &CandidateNode{Kind: ScalarNode, Tag: "!!int", Value: "12"},
		}, 0)
		if err != nil {
			t.Fatal(err)
		}
		if buf.String() != "12" {
			t.Fatalf("expected dereferenced alias value, got %q", buf.String())
		}
	})

	t.Run("unknown_kind", func(t *testing.T) {
		var buf bytes.Buffer
		err := ke.writeNode(&buf, &CandidateNode{Kind: Kind(12345)}, 0)
		if err != nil {
			t.Fatal(err)
		}
		if buf.String() != "null" {
			t.Fatalf("expected null for unknown kind, got %q", buf.String())
		}
	})
}

func TestKYamlEncoderEmptyCollections(t *testing.T) {
	ke := NewKYamlEncoder(NewDefaultKYamlPreferences()).(*kyamlEncoder)

	t.Run("empty_mapping", func(t *testing.T) {
		var buf bytes.Buffer
		err := ke.writeNode(&buf, &CandidateNode{Kind: MappingNode}, 0)
		if err != nil {
			t.Fatal(err)
		}
		if buf.String() != "{}" {
			t.Fatalf("expected empty mapping, got %q", buf.String())
		}
	})

	t.Run("empty_sequence", func(t *testing.T) {
		var buf bytes.Buffer
		err := ke.writeNode(&buf, &CandidateNode{Kind: SequenceNode}, 0)
		if err != nil {
			t.Fatal(err)
		}
		if buf.String() != "[]" {
			t.Fatalf("expected empty sequence, got %q", buf.String())
		}
	})
}

func TestKYamlEncoderScalarFallbackAndEscaping(t *testing.T) {
	ke := NewKYamlEncoder(NewDefaultKYamlPreferences()).(*kyamlEncoder)

	t.Run("unknown_tag_falls_back_to_string", func(t *testing.T) {
		var buf bytes.Buffer
		err := ke.writeNode(&buf, &CandidateNode{Kind: ScalarNode, Tag: "!!timestamp", Value: "2020-01-01T00:00:00Z"}, 0)
		if err != nil {
			t.Fatal(err)
		}
		if buf.String() != "\"2020-01-01T00:00:00Z\"" {
			t.Fatalf("expected quoted fallback, got %q", buf.String())
		}
	})

	t.Run("escape_double_quoted", func(t *testing.T) {
		got := escapeDoubleQuotedString("a\\b\"c\n\r\t" + string(rune(0x01)))
		want := "a\\\\b\\\"c\\n\\r\\t\\u0001"
		if got != want {
			t.Fatalf("expected %q, got %q", want, got)
		}
	})

	t.Run("valid_bare_key", func(t *testing.T) {
		if isValidKYamlBareKey("") {
			t.Fatalf("expected empty string to be invalid")
		}
		if isValidKYamlBareKey("1a") {
			t.Fatalf("expected leading digit to be invalid")
		}
		if !isValidKYamlBareKey("a_b-2") {
			t.Fatalf("expected identifier-like key to be valid")
		}
	})
}

func TestKYamlEncoderCommentsInMapping(t *testing.T) {
	prefs := NewDefaultKYamlPreferences()
	prefs.UnwrapScalar = false
	ke := NewKYamlEncoder(prefs).(*kyamlEncoder)

	var buf bytes.Buffer
	err := ke.writeNode(&buf, &CandidateNode{
		Kind: MappingNode,
		Content: []*CandidateNode{
			{
				Kind:        ScalarNode,
				Tag:         "!!str",
				Value:       "a",
				HeadComment: "key head",
				LineComment: "key line",
				FootComment: "key foot",
			},
			{
				Kind:        ScalarNode,
				Tag:         "!!str",
				Value:       "b",
				HeadComment: "value head",
			},
		},
	}, 0)
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "# key head\n") {
		t.Fatalf("expected key head comment, got %q", out)
	}
	if !strings.Contains(out, "# value head\n") {
		t.Fatalf("expected value head comment, got %q", out)
	}
	if !strings.Contains(out, ", # key line\n") {
		t.Fatalf("expected inline key comment fallback, got %q", out)
	}
	if !strings.Contains(out, "# key foot\n") {
		t.Fatalf("expected foot comment fallback, got %q", out)
	}
}

func TestKYamlEncoderCommentBlockAndInlineComment(t *testing.T) {
	ke := NewKYamlEncoder(NewDefaultKYamlPreferences()).(*kyamlEncoder)

	t.Run("comment_block_prefixing_and_crlf", func(t *testing.T) {
		var buf bytes.Buffer
		err := ke.writeCommentBlock(&buf, "line1\r\n\r\n# already\r\nline2", 2)
		if err != nil {
			t.Fatal(err)
		}
		want := "  # line1\n  # already\n  # line2\n"
		if buf.String() != want {
			t.Fatalf("expected %q, got %q", want, buf.String())
		}
	})

	t.Run("inline_comment_prefix_and_first_line_only", func(t *testing.T) {
		var buf bytes.Buffer
		err := ke.writeInlineComment(&buf, "hello\r\nsecond line")
		if err != nil {
			t.Fatal(err)
		}
		if buf.String() != " # hello" {
			t.Fatalf("expected %q, got %q", " # hello", buf.String())
		}
	})

	t.Run("inline_comment_already_prefixed", func(t *testing.T) {
		var buf bytes.Buffer
		err := ke.writeInlineComment(&buf, "# hello")
		if err != nil {
			t.Fatal(err)
		}
		if buf.String() != " # hello" {
			t.Fatalf("expected %q, got %q", " # hello", buf.String())
		}
	})
}
