//go:build yq_nocsv

package yqlib

func NewCSVObjectDecoder(prefs CsvPreferences) Decoder {
	return nil
}

func NewCsvEncoder(prefs CsvPreferences) Encoder {
	return nil
}
