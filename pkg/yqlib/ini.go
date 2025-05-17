package yqlib

type INIPreferences struct {
	ColorsEnabled bool
}

func NewDefaultINIPreferences() INIPreferences {
	return INIPreferences{
		ColorsEnabled: false,
	}
}

func (p *INIPreferences) Copy() INIPreferences {
	return INIPreferences{
		ColorsEnabled: p.ColorsEnabled,
	}
}

var ConfiguredINIPreferences = NewDefaultINIPreferences()
