package yqlib

type CsvPreferences struct {
	Separator rune
	AutoParse bool
}

func NewDefaultCsvPreferences() CsvPreferences {
	return CsvPreferences{
		Separator: ',',
		AutoParse: true,
	}
}

func NewDefaultTsvPreferences() CsvPreferences {
	return CsvPreferences{
		Separator: '\t',
		AutoParse: true,
	}
}

var ConfiguredCsvPreferences = NewDefaultCsvPreferences()
var ConfiguredTsvPreferences = NewDefaultTsvPreferences()
