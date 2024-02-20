package yqlib

type PropertiesPreferences struct {
	KeyValueSeparator string
	UseArrayBrackets  bool
}

func NewDefaultPropertiesPreferences() PropertiesPreferences {
	return PropertiesPreferences{
		KeyValueSeparator: " = ",
		UseArrayBrackets:  false,
	}
}

var ConfiguredPropertiesPreferences = NewDefaultPropertiesPreferences()
