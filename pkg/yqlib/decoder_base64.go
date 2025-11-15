//go:build !yq_nobase64

package yqlib

import (
	"bytes"
	"encoding/base64"
	"io"
	"strings"
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

func (dec *base64Decoder) Init(reader io.Reader) error {
	// Read all data from the reader and strip leading/trailing whitespace
	// This is necessary because base64 decoding needs to see the complete input
	// to handle padding correctly, and we need to strip whitespace before decoding.
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(reader); err != nil {
		return err
	}

	// Strip leading and trailing whitespace
	stripped := strings.TrimSpace(buf.String())

	// Add padding if needed (base64 strings should be a multiple of 4 characters)
	padLen := len(stripped) % 4
	if padLen > 0 {
		stripped += strings.Repeat("=", 4-padLen)
	}

	// Create a new reader from the stripped and padded data
	dec.reader = strings.NewReader(stripped)
	dec.readAnything = false
	dec.finished = false
	return nil
}

func (dec *base64Decoder) Decode() (*CandidateNode, error) {
	if dec.finished {
		return nil, io.EOF
	}
	base64Reader := base64.NewDecoder(&dec.encoding, dec.reader)
	buf := new(bytes.Buffer)

	if _, err := buf.ReadFrom(base64Reader); err != nil {
		return nil, err
	}
	if buf.Len() == 0 {
		dec.finished = true

		// if we've read _only_ an empty string, lets return that
		// otherwise if we've already read some bytes, and now we get
		// an empty string, then we are done.
		if dec.readAnything {
			return nil, io.EOF
		}
	}
	dec.readAnything = true
	return createStringScalarNode(buf.String()), nil
}
