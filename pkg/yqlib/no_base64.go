//go:build yq_nobase64

package yqlib

func NewBase64Decoder() Decoder {
	return nil
}

func NewBase64Encoder() Encoder {
	return nil
}
