package yqlib

import (
	"encoding/json"
	"io"

	yaml "gopkg.in/yaml.v3"
)

type jsonEncoder struct {
	indentString string
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

func NewJONEncoder(indent int) Encoder {
	var indentString = ""

	for index := 0; index < indent; index++ {
		indentString = indentString + " "
	}

	return &jsonEncoder{indentString}
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
	var encoder = json.NewEncoder(writer)
	encoder.SetEscapeHTML(false) // do not escape html chars e.g. &, <, >
	encoder.SetIndent("", je.indentString)

	var dataBucket orderedMap
	// firstly, convert all map keys to strings
	mapKeysToStrings(node)
	errorDecoding := node.Decode(&dataBucket)
	if errorDecoding != nil {
		return errorDecoding
	}
	return encoder.Encode(dataBucket)
}
