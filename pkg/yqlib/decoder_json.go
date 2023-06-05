//go:build !yq_nojson

package yqlib

import (
	"io"

	"github.com/goccy/go-json"
)

type jsonDecoder struct {
	decoder json.Decoder
}

func NewJSONDecoder() Decoder {
	return &jsonDecoder{}
}

func (dec *jsonDecoder) Init(reader io.Reader) error {
	dec.decoder = *json.NewDecoder(reader)
	return nil
}

func (dec *jsonDecoder) Decode() (*CandidateNode, error) {

	var dataBucket CandidateNode
	err := dec.decoder.Decode(&dataBucket)
	if err != nil {
		return nil, err
	}

	return &dataBucket, nil
}
