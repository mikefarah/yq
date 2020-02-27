package yqlib

import (
	"bytes"
	"encoding/json"
	"io"

	yaml "gopkg.in/yaml.v3"
)

type Encoder interface {
	Encode(node *yaml.Node) error
}

type yamlEncoder struct {
	destination io.Writer
	indent      int
	colorise    bool
	firstDoc    bool
}

func NewYamlEncoder(destination io.Writer, indent int, colorise bool) Encoder {
	if indent < 0 {
		indent = 0
	}
	return &yamlEncoder{destination, indent, colorise, true}
}

func (ye *yamlEncoder) Encode(node *yaml.Node) error {

	destination := ye.destination
	tempBuffer := bytes.NewBuffer(nil)
	if ye.colorise {
		destination = tempBuffer
	}

	var encoder = yaml.NewEncoder(destination)

	encoder.SetIndent(ye.indent)
	// TODO: work out if the first doc had a separator or not.
	if ye.firstDoc {
		ye.firstDoc = false
	} else if _, err := destination.Write([]byte("---\n")); err != nil {
		return err
	}

	if err := encoder.Encode(node); err != nil {
		return err
	}

	if ye.colorise {
		return ColorizeAndPrint(tempBuffer.Bytes(), ye.destination)
	}
	return nil
}

type jsonEncoder struct {
	encoder *json.Encoder
}

func NewJsonEncoder(destination io.Writer, prettyPrint bool, indent int) Encoder {
	var encoder = json.NewEncoder(destination)
	var indentString = ""

	for index := 0; index < indent; index++ {
		indentString = indentString + " "
	}
	if prettyPrint {
		encoder.SetIndent("", indentString)
	}
	return &jsonEncoder{encoder}
}

func (je *jsonEncoder) Encode(node *yaml.Node) error {
	var dataBucket interface{}
	errorDecoding := node.Decode(&dataBucket)
	if errorDecoding != nil {
		return errorDecoding
	}
	return je.encoder.Encode(dataBucket)
}
