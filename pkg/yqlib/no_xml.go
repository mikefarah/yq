//go:build yq_noxml

package yqlib

func NewXMLDecoder(prefs XmlPreferences) Decoder {
	return nil
}

func NewXMLEncoder(prefs XmlPreferences) Encoder {
	return nil
}
