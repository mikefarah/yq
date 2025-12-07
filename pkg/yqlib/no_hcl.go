//go:build yq_nohcl

package yqlib

func NewHclDecoder() Decoder {
	return nil
}

func NewHclEncoder() Encoder {
	return nil
}
