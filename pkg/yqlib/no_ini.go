//go:build yq_noini

package yqlib

func NewINIDecoder(prefs INIPreferences) Decoder {
	return nil
}

func NewINIEncoder() Encoder {
	return nil
}
