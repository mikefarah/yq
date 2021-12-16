package yqlib

import (
	"io"

	yaml "gopkg.in/yaml.v3"
)

type Decoder interface {
	Decode(node *yaml.Node) error
}

type yamlDecoder struct {
	decoder *yaml.Decoder
}

func NewYamlDecoder(reader io.Reader) Decoder {
	return &yamlDecoder{decoder: yaml.NewDecoder(reader)}
}

func (dec *yamlDecoder) Decode(rootYamlNode *yaml.Node) error {
	return dec.decoder.Decode(rootYamlNode)
}
