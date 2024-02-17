package yqlib

import (
	"encoding/csv"
	"fmt"
	"io"
)

type csvEncoder struct {
	separator rune
}

func NewCsvEncoder(prefs CsvPreferences) Encoder {
	return &csvEncoder{separator: prefs.Separator}
}

func (e *csvEncoder) CanHandleAliases() bool {
	return false
}

func (e *csvEncoder) PrintDocumentSeparator(_ io.Writer) error {
	return nil
}

func (e *csvEncoder) PrintLeadingContent(_ io.Writer, _ string) error {
	return nil
}

func (e *csvEncoder) encodeRow(csvWriter *csv.Writer, contents []*CandidateNode) error {
	stringValues := make([]string, len(contents))

	for i, child := range contents {

		if child.Kind != ScalarNode {
			return fmt.Errorf("csv encoding only works for arrays of scalars (string/numbers/booleans), child[%v] is a %v", i, child.Tag)
		}
		stringValues[i] = child.Value
	}
	return csvWriter.Write(stringValues)
}

func (e *csvEncoder) encodeArrays(csvWriter *csv.Writer, content []*CandidateNode) error {
	for i, child := range content {

		if child.Kind != SequenceNode {
			return fmt.Errorf("csv encoding only works for arrays of scalars (string/numbers/booleans), child[%v] is a %v", i, child.Tag)
		}
		err := e.encodeRow(csvWriter, child.Content)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *csvEncoder) extractHeader(child *CandidateNode) ([]*CandidateNode, error) {
	if child.Kind != MappingNode {
		return nil, fmt.Errorf("csv object encoding only works for arrays of flat objects (string key => string/numbers/boolean value), child[0] is a %v", child.Tag)
	}
	mapKeys := getMapKeys(child)
	return mapKeys.Content, nil
}

func (e *csvEncoder) createChildRow(child *CandidateNode, headers []*CandidateNode) []*CandidateNode {
	childRow := make([]*CandidateNode, 0)
	for _, header := range headers {
		keyIndex := findKeyInMap(child, header)
		value := createScalarNode(nil, "")
		if keyIndex != -1 {
			value = child.Content[keyIndex+1]
		}
		childRow = append(childRow, value)
	}
	return childRow

}

func (e *csvEncoder) encodeObjects(csvWriter *csv.Writer, content []*CandidateNode) error {
	headers, err := e.extractHeader(content[0])
	if err != nil {
		return nil
	}

	err = e.encodeRow(csvWriter, headers)
	if err != nil {
		return nil
	}

	for i, child := range content {
		if child.Kind != MappingNode {
			return fmt.Errorf("csv object encoding only works for arrays of flat objects (string key => string/numbers/boolean value), child[%v] is a %v", i, child.Tag)
		}
		row := e.createChildRow(child, headers)
		err = e.encodeRow(csvWriter, row)
		if err != nil {
			return err
		}

	}
	return nil
}

func (e *csvEncoder) Encode(writer io.Writer, node *CandidateNode) error {
	if node.Kind == ScalarNode {
		return writeString(writer, node.Value+"\n")
	}

	csvWriter := csv.NewWriter(writer)
	csvWriter.Comma = e.separator

	// node must be a sequence
	if node.Kind != SequenceNode {
		return fmt.Errorf("csv encoding only works for arrays, got: %v", node.Tag)
	} else if len(node.Content) == 0 {
		return nil
	}
	if node.Content[0].Kind == ScalarNode {
		return e.encodeRow(csvWriter, node.Content)
	}

	if node.Content[0].Kind == MappingNode {
		return e.encodeObjects(csvWriter, node.Content)
	}

	return e.encodeArrays(csvWriter, node.Content)

}
