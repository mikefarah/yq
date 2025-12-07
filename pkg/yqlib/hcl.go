package yqlib

type HclPreferences struct {
	ColorsEnabled bool
}

func NewDefaultHclPreferences() HclPreferences {
	return HclPreferences{ColorsEnabled: false}
}

func (p *HclPreferences) Copy() HclPreferences {
	return HclPreferences{ColorsEnabled: p.ColorsEnabled}
}

var ConfiguredHclPreferences = NewDefaultHclPreferences()
