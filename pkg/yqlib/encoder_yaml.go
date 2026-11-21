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

	// When preserving the original formatting (e.g. yq -i), the output is
	// buffered so it can be reconciled with the original source layout
	// (blank lines, indentation, spacing before inline comments) afterwards.
	preserve := ye.prefs.PreserveFormatting
	var sourceLines []string
	if preserve {
		sourceLines = node.GetSourceLines()
	}

	destination := writer
	tempBuffer := bytes.NewBuffer(nil)
	if ye.prefs.ColorsEnabled || preserve {
		destination = tempBuffer
	}

	indent := ye.prefs.Indent
	if indent < 2 {
		indent = 2
	} else if indent > 9 {
		indent = 9
	}

	target, err := node.MarshalYAML()
	if err != nil {
		return err
	}

	if preserve {
		// reuse the indentation of the original document when detectable
		if detected, ok := detectOriginalIndent(target); ok {
			indent = detected
		}
		// re-insert the blank lines that existed between nodes originally
		restoreBlankLines(target, sourceLines)
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

	if preserve {
		processed := postProcessPreservedYAML(tempBuffer.Bytes(), sourceLines, target)
		if ye.prefs.ColorsEnabled {
			return colorizeAndPrint(processed, writer)
		}
		_, err := writer.Write(processed)
		return err
	}

	if ye.prefs.ColorsEnabled {
		return colorizeAndPrint(tempBuffer.Bytes(), writer)
	}
	return nil
}
