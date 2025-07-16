//go:build yq_noprops

package yqlib

func NewPropertiesDecoder() Decoder {
	return nil
}

func NewPropertiesEncoder(prefs PropertiesPreferences) Encoder {
	return nil
}
