package yqlib

import (
	"io"

	yaml "gopkg.in/yaml.v3"
)

type yamlDecoder struct {
	decoder yaml.Decoder
}

func NewYamlDecoder() Decoder {
	return &yamlDecoder{}
}

func (dec *yamlDecoder) Init(reader io.Reader) {
	dec.decoder = *yaml.NewDecoder(reader)
}

func (dec *yamlDecoder) Decode(rootYamlNode *yaml.Node) error {
	return dec.decoder.Decode(rootYamlNode)
}
