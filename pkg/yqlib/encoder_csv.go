package yqlib

import (
	"encoding/csv"
	"fmt"
	"io"

	yaml "gopkg.in/yaml.v3"
)

type csvEncoder struct {
	separator rune
}

func NewCsvEncoder(separator rune) Encoder {
	return &csvEncoder{separator}
}

func (e *csvEncoder) CanHandleAliases() bool {
	return false
}

func (e *csvEncoder) PrintDocumentSeparator(writer io.Writer) error {
	return nil
}

func (e *csvEncoder) PrintLeadingContent(writer io.Writer, content string) error {
	return nil
}

func (e *csvEncoder) encodeRow(csvWriter *csv.Writer, contents []*yaml.Node) error {
	stringValues := make([]string, len(contents))

	for i, child := range contents {

		if child.Kind != yaml.ScalarNode {
			return fmt.Errorf("csv encoding only works for arrays of scalars (string/numbers/booleans), child[%v] is a %v", i, child.Tag)
		}
		stringValues[i] = child.Value
	}
	return csvWriter.Write(stringValues)
}

func (e *csvEncoder) Encode(writer io.Writer, originalNode *yaml.Node) error {
	csvWriter := csv.NewWriter(writer)
	csvWriter.Comma = e.separator

	// node must be a sequence
	node := unwrapDoc(originalNode)
	if node.Kind != yaml.SequenceNode {
		return fmt.Errorf("csv encoding only works for arrays, got: %v", node.Tag)
	} else if len(node.Content) == 0 {
		return nil
	}
	if node.Content[0].Kind == yaml.ScalarNode {
		return e.encodeRow(csvWriter, node.Content)
	}

	for i, child := range node.Content {

		if child.Kind != yaml.SequenceNode {
			return fmt.Errorf("csv encoding only works for arrays of scalars (string/numbers/booleans), child[%v] is a %v", i, child.Tag)
		}
		err := e.encodeRow(csvWriter, child.Content)
		if err != nil {
			return err
		}
	}
	return nil
}
