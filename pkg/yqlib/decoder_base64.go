package yqlib

import (
	"bytes"
	"encoding/base64"
	"io"

	yaml "gopkg.in/yaml.v3"
)

type base64Decoder struct {
	reader   io.Reader
	finished bool
	encoding base64.Encoding
}

func NewBase64Decoder() Decoder {
	return &base64Decoder{finished: false, encoding: *base64.StdEncoding}
}

func (dec *base64Decoder) Init(reader io.Reader) {
	dec.reader = reader
	dec.finished = false
}

func (dec *base64Decoder) Decode(rootYamlNode *yaml.Node) error {
	if dec.finished {
		return io.EOF
	}
	base64Reader := base64.NewDecoder(&dec.encoding, dec.reader)
	buf := new(bytes.Buffer)

	if _, err := buf.ReadFrom(base64Reader); err != nil {
		return err
	}
	if buf.Len() == 0 {
		dec.finished = true
		return io.EOF
	}
	rootYamlNode.Kind = yaml.ScalarNode
	rootYamlNode.Tag = "!!str"
	rootYamlNode.Value = buf.String()
	return nil
}
