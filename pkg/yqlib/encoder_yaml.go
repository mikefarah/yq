package yqlib

import (
	"bytes"
	"io"
	"strings"

	"go.yaml.in/yaml/v4"
)

type yamlEncoder struct {
	prefs YamlPreferences
}

func NewYamlEncoder(prefs YamlPreferences) Encoder {
	return &yamlEncoder{prefs}
}

func (ye *yamlEncoder) CanHandleAliases() bool {
	return true
}

func (ye *yamlEncoder) PrintDocumentSeparator(writer io.Writer) error {
	return PrintYAMLDocumentSeparator(writer, ye.prefs.PrintDocSeparators)
}

func (ye *yamlEncoder) PrintLeadingContent(writer io.Writer, content string) error {
	return PrintYAMLLeadingContent(writer, content, ye.prefs.PrintDocSeparators, ye.prefs.ColorsEnabled)
}

func (ye *yamlEncoder) Encode(writer io.Writer, node *CandidateNode) error {
	log.Debugf("encoderYaml - going to print %v", NodeToString(node))
	// Detect line ending style from LeadingContent
	lineEnding := "\n"
	if strings.Contains(node.LeadingContent, "\r\n") {
		lineEnding = "\r\n"
	}
	if node.Kind == ScalarNode && ye.prefs.UnwrapScalar && !bareStringNeedsQuoting(node) {
		valueToPrint := node.Value
		if node.LeadingContent == "" || valueToPrint != "" {
			valueToPrint = valueToPrint + lineEnding
		}
		return writeString(writer, valueToPrint)
	}

	destination := writer
	tempBuffer := bytes.NewBuffer(nil)
	if ye.prefs.ColorsEnabled {
		destination = tempBuffer
	}

	var encoder = yaml.NewEncoder(destination)

	encoder.SetIndent(ye.prefs.Indent)
	if ye.prefs.CompactSequenceIndent {
		encoder.CompactSeqIndent()
	}

	target, err := node.MarshalYAML()

	if err != nil {
		return err
	}

	trailingContent := target.FootComment
	target.FootComment = ""

	if err := encoder.Encode(target); err != nil {
		return err
	}

	if err := ye.PrintLeadingContent(destination, trailingContent); err != nil {
		return err
	}

	if ye.prefs.ColorsEnabled {
		return colorizeAndPrint(tempBuffer.Bytes(), writer)
	}
	return nil
}

// bareStringNeedsQuoting reports whether a top-level string scalar would be
// structurally reinterpreted if emitted as an unquoted bare value. The
// unwrap-scalar fast-path writes node.Value verbatim, which silently turns a
// string like "this: should really work" into a mapping on the next read, or
// "- item" into a sequence. When this returns true the caller falls through
// to the full yaml encoder, which applies the quoting style required by the
// YAML spec. Scalar-to-scalar reinterpretations (e.g. "123" parsing as an int
// tag) are not covered here: they preserve the node shape and are handled by
// callers that care about explicit tag preservation.
func bareStringNeedsQuoting(node *CandidateNode) bool {
	if node.Tag != "!!str" {
		return false
	}
	var parsed yaml.Node
	if err := yaml.Unmarshal([]byte(node.Value), &parsed); err != nil {
		// Unparseable bare form (e.g. control characters): leave it on the
		// fast-path so callers that check for those characters still see them.
		return false
	}
	if parsed.Kind != yaml.DocumentNode || len(parsed.Content) != 1 {
		return false
	}
	return parsed.Content[0].Kind != yaml.ScalarNode
}
