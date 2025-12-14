package yqlib

type TomlPreferences struct {
	ColorsEnabled bool
}

func NewDefaultTomlPreferences() TomlPreferences {
	return TomlPreferences{ColorsEnabled: false}
}

func (p *TomlPreferences) Copy() TomlPreferences {
	return TomlPreferences{ColorsEnabled: p.ColorsEnabled}
}

var ConfiguredTomlPreferences = NewDefaultTomlPreferences()
