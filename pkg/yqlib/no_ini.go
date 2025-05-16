//go:build yq_noini

package yqlib

func NewINIDecoder() Decoder {
	return nil
}

func NewINIEncoder(indent int) Encoder {
	return nil
}
