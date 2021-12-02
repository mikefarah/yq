package yqlib

import (
	"encoding/csv"
	"fmt"
	"io"

	yaml "gopkg.in/yaml.v3"
)

type csvEncoder struct {
	destination csv.Writer
}

func NewCsvEncoder(destination io.Writer, separator rune) Encoder {
	csvWriter := *csv.NewWriter(destination)
	csvWriter.Comma = separator
	return &csvEncoder{csvWriter}
}

func (e *csvEncoder) encodeRow(contents []*yaml.Node) error {
	stringValues := make([]string, len(contents))

	for i, child := range contents {

		if child.Kind != yaml.ScalarNode {
			return fmt.Errorf("csv encoding only works for arrays of scalars (string/numbers/booleans), child[%v] is a %v", i, child.Tag)
		}
		stringValues[i] = child.Value
	}
	return e.destination.Write(stringValues)
}

func (e *csvEncoder) Encode(originalNode *yaml.Node) error {
	// node must be a sequence
	node := unwrapDoc(originalNode)
	if node.Kind != yaml.SequenceNode {
		return fmt.Errorf("csv encoding only works for arrays, got: %v", node.Tag)
	} else if len(node.Content) == 0 {
		return nil
	}
	if node.Content[0].Kind == yaml.ScalarNode {
		return e.encodeRow(node.Content)
	}

	for i, child := range node.Content {

		if child.Kind != yaml.SequenceNode {
			return fmt.Errorf("csv encoding only works for arrays of scalars (string/numbers/booleans), child[%v] is a %v", i, child.Tag)
		}
		err := e.encodeRow(child.Content)
		if err != nil {
			return err
		}
	}
	return nil
}
