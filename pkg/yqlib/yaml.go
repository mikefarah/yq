package yqlib

type YamlPreferences struct {
	LeadingContentPreProcessing bool
	PrintDocSeparators          bool
	UnwrapScalar                bool
}

func NewDefaultYamlPreferences() YamlPreferences {
	return YamlPreferences{
		LeadingContentPreProcessing: true,
		PrintDocSeparators:          true,
		UnwrapScalar:                true,
	}
}

var ConfiguredYamlPreferences = NewDefaultYamlPreferences()
