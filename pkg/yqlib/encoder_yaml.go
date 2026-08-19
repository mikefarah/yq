package yqlib

import (
	"bytes"
	"cmp"
	"fmt"
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
	if node.Kind == ScalarNode && ye.prefs.UnwrapScalar {
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

	dumper, err := yaml.NewDumper(destination,
		yaml.WithV3Defaults(),
		yaml.WithIndent(cmp.Or(ye.prefs.Indent, 2)),
		yaml.WithCompactSeqIndent(ye.prefs.CompactSequenceIndent),
		yaml.WithLineWidth(-1),
	)
	if err != nil {
		return fmt.Errorf("configure YAML encoding %#v: %v", ye.prefs, err)
	}

	target, err := node.MarshalYAML()

	if err != nil {
		return err
	}

	trailingContent := target.FootComment
	target.FootComment = ""

	if err := dumper.Dump(target); err != nil {
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
