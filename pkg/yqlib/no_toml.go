//go:build yq_notoml

package yqlib

func NewTomlDecoder() Decoder {
	return nil
}

func NewTomlEncoder() Encoder {
	return nil
}

func NewTomlEncoderWithPrefs(prefs TomlPreferences) Encoder {
	return nil
}
