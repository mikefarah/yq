package yqlib

import (
	"bytes"
	"io"

	"github.com/goccy/go-json"
	yaml "gopkg.in/yaml.v3"
)

type jsonEncoder struct {
	indentString string
	colorise     bool
}

func mapKeysToStrings(node *yaml.Node) {

	if node.Kind == yaml.MappingNode {
		for index, child := range node.Content {
			if index%2 == 0 { // its a map key
				child.Tag = "!!str"
			}
		}
	}

	for _, child := range node.Content {
		mapKeysToStrings(child)
	}
}

func NewJSONEncoder(indent int, colorise bool) Encoder {
	var indentString = ""

	for index := 0; index < indent; index++ {
		indentString = indentString + " "
	}

	return &jsonEncoder{indentString, colorise}
}

func (je *jsonEncoder) CanHandleAliases() bool {
	return false
}

func (je *jsonEncoder) PrintDocumentSeparator(writer io.Writer) error {
	return nil
}

func (je *jsonEncoder) PrintLeadingContent(writer io.Writer, content string) error {
	return nil
}

func (je *jsonEncoder) Encode(writer io.Writer, node *yaml.Node) error {

	destination := writer
	tempBuffer := bytes.NewBuffer(nil)
	if je.colorise {
		destination = tempBuffer
	}

	var encoder = json.NewEncoder(destination)
	encoder.SetEscapeHTML(false) // do not escape html chars e.g. &, <, >
	encoder.SetIndent("", je.indentString)

	var dataBucket orderedMap
	// firstly, convert all map keys to strings
	mapKeysToStrings(node)
	errorDecoding := node.Decode(&dataBucket)
	if errorDecoding != nil {
		return errorDecoding
	}
	err := encoder.Encode(dataBucket)
	if err != nil {
		return err
	}
	if je.colorise {
		return colorizeAndPrint(tempBuffer.Bytes(), writer)
	}
	return nil
}
