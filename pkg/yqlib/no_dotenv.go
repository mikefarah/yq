//go:build yq_nodotenv

package yqlib

func NewDotenvDecoder() Decoder {
	return nil
}

func NewDotenvEncoder() Encoder {
	return nil
}
