package yqlib

type YamlPreferences struct {
	Indent                      int
	ColorsEnabled               bool
	LeadingContentPreProcessing bool
	PrintDocSeparators          bool
	UnwrapScalar                bool
	EvaluateTogether            bool
	FixMergeAnchorToSpec        bool
	CompactSequenceIndent       bool
}

func NewDefaultYamlPreferences() YamlPreferences {
	return YamlPreferences{
		Indent:                      2,
		ColorsEnabled:               false,
		LeadingContentPreProcessing: true,
		PrintDocSeparators:          true,
		UnwrapScalar:                true,
		EvaluateTogether:            false,
		FixMergeAnchorToSpec:        false,
		CompactSequenceIndent:       false,
	}
}

func (p *YamlPreferences) Copy() YamlPreferences {
	return YamlPreferences{
		Indent:                      p.Indent,
		ColorsEnabled:               p.ColorsEnabled,
		LeadingContentPreProcessing: p.LeadingContentPreProcessing,
		PrintDocSeparators:          p.PrintDocSeparators,
		UnwrapScalar:                p.UnwrapScalar,
		EvaluateTogether:            p.EvaluateTogether,
		FixMergeAnchorToSpec:        p.FixMergeAnchorToSpec,
		CompactSequenceIndent:       p.CompactSequenceIndent,
	}
}

var ConfiguredYamlPreferences = NewDefaultYamlPreferences()
