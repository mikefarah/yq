//go:build yq_nojson5

package yqlib

func NewJSON5Decoder() Decoder {
	return nil
}

func NewJSON5Encoder(prefs JsonPreferences) Encoder {
	return nil
}
