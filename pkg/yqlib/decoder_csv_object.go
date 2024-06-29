package yqlib

import (
	"encoding/csv"
	"errors"
	"io"

	"github.com/dimchansky/utfbom"
)

type csvObjectDecoder struct {
	prefs    CsvPreferences
	reader   csv.Reader
	finished bool
}

func NewCSVObjectDecoder(prefs CsvPreferences) Decoder {
	return &csvObjectDecoder{prefs: prefs}
}

func (dec *csvObjectDecoder) Init(reader io.Reader) error {
	cleanReader, enc := utfbom.Skip(reader)
	log.Debugf("Detected encoding: %s\n", enc)
	dec.reader = *csv.NewReader(cleanReader)
	dec.reader.Comma = dec.prefs.Separator
	dec.finished = false
	return nil
}

func (dec *csvObjectDecoder) convertToNode(content string) *CandidateNode {
	node, err := parseSnippet(content)
	// if we're not auto-parsing, then we wont put in parsed objects or arrays
	// but we still parse scalars
	if err != nil || (!dec.prefs.AutoParse && (node.Kind != ScalarNode || node.Value != content)) {
		return createScalarNode(content, content)
	}
	return node
}

func (dec *csvObjectDecoder) createObject(headerRow []string, contentRow []string) *CandidateNode {
	objectNode := &CandidateNode{Kind: MappingNode, Tag: "!!map"}

	for i, header := range headerRow {
		objectNode.AddKeyValueChild(createScalarNode(header, header), dec.convertToNode(contentRow[i]))
	}
	return objectNode
}

func (dec *csvObjectDecoder) Decode() (*CandidateNode, error) {
	if dec.finished {
		return nil, io.EOF
	}
	headerRow, err := dec.reader.Read()
	log.Debugf(": headerRow%v", headerRow)
	if err != nil {
		return nil, err
	}

	rootArray := &CandidateNode{Kind: SequenceNode, Tag: "!!seq"}

	contentRow, err := dec.reader.Read()

	for err == nil && len(contentRow) > 0 {
		log.Debugf("Adding contentRow: %v", contentRow)
		rootArray.AddChild(dec.createObject(headerRow, contentRow))
		contentRow, err = dec.reader.Read()
		log.Debugf("Read next contentRow: %v, %v", contentRow, err)
	}
	if !errors.Is(err, io.EOF) {
		return nil, err
	}

	return rootArray, nil
}
