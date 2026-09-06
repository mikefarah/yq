package yqlib

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"go.yaml.in/yaml/v4"
)

type yamlEncoder struct {
	prefs YamlPreferences
}

// scalarWouldNotRoundTrip reports whether the given string scalar would be
// re-parsed as something other than a string (an int, float, bool, null,
// timestamp, or a non-scalar node) if emitted without quotes, breaking
// roundtrip safety of the output document.
func scalarWouldNotRoundTrip(node *CandidateNode) bool {
	if node.Tag != "!!str" || node.Value == "" {
		return false
	}
	// Control characters (e.g. NUL) cannot appear in plain scalars; emitting
	// them quoted would bypass the NUL-separated output safety check, which
	// relies on the plain write path erroring out.
	if strings.ContainsFunc(node.Value, func(r rune) bool { return r < 0x20 && r != '\t' }) {
		return false
	}
	decoder := NewYamlDecoder(YamlPreferences{})
	if err := decoder.Init(bytes.NewReader([]byte(node.Value))); err != nil {
		return true
	}
	reencoded, err := decoder.Decode()
	if err != nil || reencoded == nil {
		return true
	}
	return reencoded.Tag != "!!str" || reencoded.Kind != ScalarNode || reencoded.Value != node.Value
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
	if node.Kind == ScalarNode && ye.prefs.UnwrapScalar && !scalarWouldNotRoundTrip(node) {
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

	indent := ye.prefs.Indent
	if indent < 2 {
		indent = 2
	} else if indent > 9 {
		indent = 9
	}

	dumper, err := yaml.NewDumper(destination,
		yaml.WithV3Defaults(),
		yaml.WithIndent(indent),
		yaml.WithCompactSeqIndent(ye.prefs.CompactSequenceIndent),
		yaml.WithLineWidth(-1),
	)
	if err != nil {
		return fmt.Errorf("configure YAML encoding: %w", err)
	}

	target, err := node.MarshalYAML()
	if err != nil {
		_ = dumper.Close()
		return err
	}

	trailingContent := target.FootComment
	target.FootComment = ""

	err = dumper.Dump(target)
	if closeErr := dumper.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
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
