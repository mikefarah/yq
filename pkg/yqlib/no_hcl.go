//go:build yq_nohcl

package yqlib

func NewHclDecoder() Decoder {
	return nil
}

func NewHclEncoder(_ HclPreferences) Encoder {
	return nil
}
