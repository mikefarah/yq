package yqlib

type PropertiesPreferences struct {
	UnwrapScalar      bool
	KeyValueSeparator string
	UseArrayBrackets  bool
}

func NewDefaultPropertiesPreferences() PropertiesPreferences {
	return PropertiesPreferences{
		UnwrapScalar:      true,
		KeyValueSeparator: " = ",
		UseArrayBrackets:  false,
	}
}

func (p *PropertiesPreferences) Copy() PropertiesPreferences {
	return PropertiesPreferences{
		UnwrapScalar:      p.UnwrapScalar,
		KeyValueSeparator: p.KeyValueSeparator,
		UseArrayBrackets:  p.UseArrayBrackets,
	}
}

var ConfiguredPropertiesPreferences = NewDefaultPropertiesPreferences()
