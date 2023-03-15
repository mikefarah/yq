//go:build yq_nojson

package yqlib

func NewJSONDecoder() Decoder {
	return nil
}

func NewJSONEncoder(indent int, colorise bool, unwrapScalar bool) Encoder {
	return nil
}
