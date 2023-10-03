//go:build yq_nolua

package yqlib

func NewLuaEncoder(prefs LuaPreferences) Encoder {
	return nil
}

func NewLuaDecoder(prefs LuaPreferences) Decoder {
	return nil
}
