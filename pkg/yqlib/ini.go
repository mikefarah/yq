package yqlib

type INIPreferences struct {
	ColorsEnabled           bool
	PreserveSurroundedQuote bool
}

func NewDefaultINIPreferences() INIPreferences {
	return INIPreferences{
		ColorsEnabled:           false,
		PreserveSurroundedQuote: false,
	}
}

func (p *INIPreferences) Copy() INIPreferences {
	return INIPreferences{
		ColorsEnabled:           p.ColorsEnabled,
		PreserveSurroundedQuote: p.PreserveSurroundedQuote,
	}
}

var ConfiguredINIPreferences = NewDefaultINIPreferences()
