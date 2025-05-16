package yqlib

type INIPreferences struct {
	Indent        int
	ColorsEnabled bool
}

func NewDefaultINIPreferences() INIPreferences {
	return INIPreferences{
		Indent:        2,
		ColorsEnabled: false,
	}
}

func (p *INIPreferences) Copy() INIPreferences {
	return INIPreferences{
		Indent:        p.Indent,
		ColorsEnabled: p.ColorsEnabled,
	}
}

var ConfiguredINIPreferences = NewDefaultINIPreferences()
