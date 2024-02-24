package yqlib

type JsonPreferences struct {
	Indent        int
	ColorsEnabled bool
	UnwrapScalar  bool
}

func NewDefaultJsonPreferences() JsonPreferences {
	return JsonPreferences{
		Indent:        2,
		ColorsEnabled: true,
		UnwrapScalar:  true,
	}
}

func (p *JsonPreferences) Copy() JsonPreferences {
	return JsonPreferences{
		Indent:        p.Indent,
		ColorsEnabled: p.ColorsEnabled,
		UnwrapScalar:  p.UnwrapScalar,
	}
}

var ConfiguredJSONPreferences = NewDefaultJsonPreferences()
