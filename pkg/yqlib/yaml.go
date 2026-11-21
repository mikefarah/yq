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
	// PreserveFormatting keeps the original document layout (blank lines,
	// indentation, comment spacing) of untouched nodes when re-encoding,
	// e.g. when writing a file in place with -i. It relies on the line/column
	// information captured when the document was decoded.
	PreserveFormatting bool
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
		PreserveFormatting:          false,
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
		PreserveFormatting:          p.PreserveFormatting,
	}
}

var ConfiguredYamlPreferences = NewDefaultYamlPreferences()
