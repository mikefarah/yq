//go:build yq_nojson

package yqlib

func NewJSONDecoder() Decoder {
	return nil
}

func NewJSONEncoder(prefs JsonPreferences) Encoder {
	return nil
}
