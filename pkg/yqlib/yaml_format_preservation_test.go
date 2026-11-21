package yqlib

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
	yaml "go.yaml.in/yaml/v4"
)

// decodePreserving decodes a single yaml document with formatting
// preservation enabled (as `yq -i` does) and returns the candidate node.
func decodePreserving(t *testing.T, document string) *CandidateNode {
	t.Helper()
	prefs := NewDefaultYamlPreferences()
	prefs.PreserveFormatting = true
	decoder := NewYamlDecoder(prefs)
	inputs, err := ReadDocuments(strings.NewReader(document), decoder)
	if err != nil {
		t.Fatal(err)
	}
	if inputs.Len() != 1 {
		t.Fatalf("expected 1 document, got %v", inputs.Len())
	}
	return inputs.Front().Value.(*CandidateNode)
}

// encodePreserving encodes the node the way the printer + yaml encoder do
// when formatting preservation is enabled (as `yq -i` does).
func encodePreserving(t *testing.T, node *CandidateNode) string {
	t.Helper()
	prefs := NewDefaultYamlPreferences()
	prefs.PreserveFormatting = true
	prefs.UnwrapScalar = false
	encoder := NewYamlEncoder(prefs)
	var buf bytes.Buffer
	if err := encoder.PrintLeadingContent(&buf, node.LeadingContent); err != nil {
		t.Fatal(err)
	}
	if err := encoder.Encode(&buf, node); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

// setMappingValue replaces the value of a mapping key, mimicking what an
// assignment expression like `.key = value` does to the tree.
func setMappingValue(t *testing.T, root *CandidateNode, key string, value string) {
	t.Helper()
	if root.Kind != MappingNode {
		t.Fatalf("expected mapping root, got %v", root.Kind)
	}
	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value == key {
			root.Content[i+1].Value = value
			root.Content[i+1].Style = 0
			if _, err := strconv.Atoi(value); err == nil {
				root.Content[i+1].Tag = "!!int"
			} else if value == "true" || value == "false" {
				root.Content[i+1].Tag = "!!bool"
			} else {
				root.Content[i+1].Tag = "!!str"
			}
			return
		}
	}
	t.Fatalf("key %q not found", key)
}

func TestPreserveFormattingIssue465(t *testing.T) {
	// the exact example from https://github.com/mikefarah/yq/issues/465
	document := `# Log level for operator
logLevel: 3

# Annotations applied to operator deployment
annotations: {}

# Annotations passed to operator pod(s).
podAnnotations: {}

# Node labels for pod assignment
nodeSelector:  # +doc-gen:break
  beta.kubernetes.io/os: linux
  beta.kubernetes.io/arch: amd64
`
	expected := `# Log level for operator
logLevel: 4

# Annotations applied to operator deployment
annotations: {}

# Annotations passed to operator pod(s).
podAnnotations: {}

# Node labels for pod assignment
nodeSelector:  # +doc-gen:break
  beta.kubernetes.io/os: linux
  beta.kubernetes.io/arch: amd64
`
	node := decodePreserving(t, document)
	setMappingValue(t, node, "logLevel", "4")
	test.AssertResult(t, expected, encodePreserving(t, node))
}

func TestPreserveFormattingIndentationAndQuotes(t *testing.T) {
	document := `# top comment

parent:
    # nested comment
    child1: 'single quoted'

    child2: "double quoted"

    list:
        - one

        - two
`
	expected := `# top comment

parent:
    # nested comment
    child1: changed

    child2: "double quoted"

    list:
        - one

        - two
`
	node := decodePreserving(t, document)
	parent := node.Content[1]
	setMappingValue(t, parent, "child1", "changed")
	test.AssertResult(t, expected, encodePreserving(t, node))
}

func TestPreserveFormattingLiteralBlockAndFootComment(t *testing.T) {
	document := `literal: |
  line1
  line2

after: true

foot: 1
# foot comment of foot

last: 2
`
	expected := `literal: |
  line1
  line2

after: false

foot: 1
# foot comment of foot

last: 2
`
	node := decodePreserving(t, document)
	setMappingValue(t, node, "after", "false")
	test.AssertResult(t, expected, encodePreserving(t, node))
}

func TestPreserveFormattingCRLF(t *testing.T) {
	document := "a: 1\r\n\r\n# comment\r\nb: 2\r\n"
	expected := "a: 10\r\n\r\n# comment\r\nb: 2\r\n"
	node := decodePreserving(t, document)
	setMappingValue(t, node, "a", "10")
	test.AssertResult(t, expected, encodePreserving(t, node))
}

func TestPreserveFormattingSequenceEntries(t *testing.T) {
	document := `seq:
  - a

  - b
  - c
last: 1
`
	expected := `seq:
  - a

  - B
  - c
last: 1
`
	node := decodePreserving(t, document)
	seq := node.Content[1]
	if seq.Kind != SequenceNode {
		t.Fatalf("expected sequence, got %v", seq.Kind)
	}
	seq.Content[1].Value = "B"
	test.AssertResult(t, expected, encodePreserving(t, node))
}

func TestPreserveFormattingDeletedNodeCollapsesBlanks(t *testing.T) {
	document := `keep: 1

# comment on gone
gone: 2

after: 3
`
	node := decodePreserving(t, document)
	// mimic del(.gone): remove the key and value from the root mapping
	content := make([]*CandidateNode, 0, len(node.Content))
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value != "gone" {
			content = append(content, node.Content[i], node.Content[i+1])
		}
	}
	node.Content = content
	// the two blank lines that surrounded the deleted entry are kept, but the
	// comment and the entry itself are gone
	expected := `keep: 1


after: 3
`
	test.AssertResult(t, expected, encodePreserving(t, node))
}

func TestWithoutPreserveFormattingNormalizes(t *testing.T) {
	// control test: without PreserveFormatting the output stays normalised
	document := `# comment
logLevel: 3

annotations: {}
`
	expected := `# comment
logLevel: 4
annotations: {}
`
	prefs := NewDefaultYamlPreferences()
	decoder := NewYamlDecoder(prefs)
	inputs, err := ReadDocuments(strings.NewReader(document), decoder)
	if err != nil {
		t.Fatal(err)
	}
	node := inputs.Front().Value.(*CandidateNode)
	setMappingValue(t, node, "logLevel", "4")
	encoder := NewYamlEncoder(prefs)
	var buf bytes.Buffer
	if err := encoder.PrintLeadingContent(&buf, node.LeadingContent); err != nil {
		t.Fatal(err)
	}
	if err := encoder.Encode(&buf, node); err != nil {
		t.Fatal(err)
	}
	test.AssertResult(t, expected, buf.String())
}

func TestDetectOriginalIndent(t *testing.T) {
	mkdoc := func(doc string) *yaml.Node {
		var node yaml.Node
		loader, err := yaml.NewLoader(strings.NewReader(doc))
		if err != nil {
			t.Fatal(err)
		}
		if err := loader.Load(&node); err != nil {
			t.Fatal(err)
		}
		return &node
	}
	if indent, ok := detectOriginalIndent(mkdoc("a:\n  b: 1\n").Content[0]); !ok || indent != 2 {
		t.Errorf("expected 2, got %v (%v)", indent, ok)
	}
	if indent, ok := detectOriginalIndent(mkdoc("a:\n    b: 1\n").Content[0]); !ok || indent != 4 {
		t.Errorf("expected 4, got %v (%v)", indent, ok)
	}
	if indent, ok := detectOriginalIndent(mkdoc("a:\n        b: 1\n").Content[0]); !ok || indent != 8 {
		t.Errorf("expected 8, got %v (%v)", indent, ok)
	}
	if _, ok := detectOriginalIndent(mkdoc("a: 1\nb: 2\n").Content[0]); ok {
		t.Error("expected no indent detection for flat doc")
	}
	// flow collections do not count
	if _, ok := detectOriginalIndent(mkdoc("a: {b: 1}\n").Content[0]); ok {
		t.Error("expected no indent detection for flow doc")
	}
}

func TestInlineCommentGap(t *testing.T) {
	test.AssertResult(t, 2, inlineCommentGap("key: value  # comment"))
	test.AssertResult(t, 1, inlineCommentGap("key: value # comment"))
	test.AssertResult(t, 4, inlineCommentGap("key: value    # comment"))
	test.AssertResult(t, 1, inlineCommentGap("no comment here"))
	test.AssertResult(t, 1, inlineCommentGap("#only comment"))
}

func TestCountBlankSourceLines(t *testing.T) {
	lines := splitSourceLines("a: 1\n\n\n# c\nb: 2\n")
	test.AssertResult(t, 2, countBlankSourceLines(lines, 2, 3))
	test.AssertResult(t, 0, countBlankSourceLines(lines, 4, 4))
	test.AssertResult(t, 1, countBlankSourceLines(lines, 2, 2))
	// out of range lines are ignored; lines 2,3 and 6 are blank
	test.AssertResult(t, 3, countBlankSourceLines(lines, 2, 100))
}
