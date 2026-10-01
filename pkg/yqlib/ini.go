package yqlib

type INIPreferences struct {
	ColorsEnabled           bool
	PreserveSurroundedQuote bool
	KeyValueDelimiters      string
}

func NewDefaultINIPreferences() INIPreferences {
	return INIPreferences{
		ColorsEnabled:           false,
		PreserveSurroundedQuote: false,
		KeyValueDelimiters:      "=:",
	}
}

func (p *INIPreferences) Copy() INIPreferences {
	return INIPreferences{
		ColorsEnabled:           p.ColorsEnabled,
		PreserveSurroundedQuote: p.PreserveSurroundedQuote,
		KeyValueDelimiters:      p.KeyValueDelimiters,
	}
}

var ConfiguredINIPreferences = NewDefaultINIPreferences()
