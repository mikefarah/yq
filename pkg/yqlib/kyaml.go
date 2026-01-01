//go:build !yq_nokyaml

package yqlib

type KYamlPreferences struct {
	Indent             int
	ColorsEnabled      bool
	PrintDocSeparators bool
	UnwrapScalar       bool
}

func NewDefaultKYamlPreferences() KYamlPreferences {
	return KYamlPreferences{
		Indent:             2,
		ColorsEnabled:      false,
		PrintDocSeparators: true,
		UnwrapScalar:       true,
	}
}

func (p *KYamlPreferences) Copy() KYamlPreferences {
	return KYamlPreferences{
		Indent:             p.Indent,
		ColorsEnabled:      p.ColorsEnabled,
		PrintDocSeparators: p.PrintDocSeparators,
		UnwrapScalar:       p.UnwrapScalar,
	}
}

var ConfiguredKYamlPreferences = NewDefaultKYamlPreferences()
