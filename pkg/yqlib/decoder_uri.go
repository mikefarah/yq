package yqlib

import (
	"bytes"
	"io"
	"net/url"
)

type uriDecoder struct {
	reader       io.Reader
	finished     bool
	readAnything bool
}

func NewUriDecoder() Decoder {
	return &uriDecoder{finished: false}
}

func (dec *uriDecoder) Init(reader io.Reader) error {
	dec.reader = reader
	dec.readAnything = false
	dec.finished = false
	return nil
}

func (dec *uriDecoder) Decode() (*CandidateNode, error) {
	if dec.finished {
		return nil, io.EOF
	}

	buf := new(bytes.Buffer)

	if _, err := buf.ReadFrom(dec.reader); err != nil {
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
	newValue, err := url.QueryUnescape(buf.String())
	if err != nil {
		return nil, err
	}
	dec.readAnything = true
	return createStringScalarNode(newValue), nil
}
