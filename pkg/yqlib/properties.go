package yqlib

type PropertiesPreferences struct {
	KeyValueSeparator string
}

func NewDefaultPropertiesPreferences() PropertiesPreferences {
	return PropertiesPreferences{
		KeyValueSeparator: " = ",
	}
}

var ConfiguredPropertiesPreferences = NewDefaultPropertiesPreferences()
