//go:build yq_noini

package yqlib

func NewINIDecoder() Decoder {
	return nil
}

func NewINIEncoder() Encoder {
	return nil
}
