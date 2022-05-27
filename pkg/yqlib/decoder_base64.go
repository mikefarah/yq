package yqlib

import (
	"bytes"
	"encoding/base64"
	"io"

	yaml "gopkg.in/yaml.v3"
)

type base64Decoder struct {
	reader       io.Reader
	finished     bool
	readAnything bool
	encoding     base64.Encoding
}

func NewBase64Decoder() Decoder {
	return &base64Decoder{finished: false, encoding: *base64.StdEncoding}
}

func (dec *base64Decoder) Init(reader io.Reader) {
	dec.reader = reader
	dec.readAnything = false
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

		// if we've read _only_ an empty string, lets return that
		// otherwise if we've already read some bytes, and now we get
		// an empty string, then we are done.
		if dec.readAnything {
			return io.EOF
		}
	}
	dec.readAnything = true
	rootYamlNode.Kind = yaml.ScalarNode
	rootYamlNode.Tag = "!!str"
	rootYamlNode.Value = buf.String()
	return nil
}
