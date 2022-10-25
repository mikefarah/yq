package yqlib

type YamlPreferences struct {
	LeadingContentPreProcessing bool
	PrintDocSeparators          bool
	UnwrapScalar                bool
	EvaluateTogether            bool
}

func NewDefaultYamlPreferences() YamlPreferences {
	return YamlPreferences{
		LeadingContentPreProcessing: true,
		PrintDocSeparators:          true,
		UnwrapScalar:                true,
		EvaluateTogether:            false,
	}
}

var ConfiguredYamlPreferences = NewDefaultYamlPreferences()
