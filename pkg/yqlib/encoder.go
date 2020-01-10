package yqlib

import (
	"encoding/json"
	"io"

	yaml "gopkg.in/yaml.v3"
)

type Encoder interface {
	Encode(node *yaml.Node) error
}

type yamlEncoder struct {
	encoder *yaml.Encoder
}

func NewYamlEncoder(destination io.Writer) Encoder {
	var encoder = yaml.NewEncoder(destination)
	encoder.SetIndent(2)
	return &yamlEncoder{encoder}
}

func (ye *yamlEncoder) Encode(node *yaml.Node) error {
	return ye.encoder.Encode(node)
}

type jsonEncoder struct {
	encoder *json.Encoder
}

func NewJsonEncoder(destination io.Writer) Encoder {
	var encoder = json.NewEncoder(destination)
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
