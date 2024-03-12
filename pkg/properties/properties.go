package properties

import "github.com/mikefarah/yq/v4/pkg/yqlib"

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

var PropertiesFormat = &yqlib.Format{"props", []string{"p", "properties"},
	func() yqlib.Encoder { return NewPropertiesEncoder(ConfiguredPropertiesPreferences) },
	func() yqlib.Decoder { return NewPropertiesDecoder() },
}

var propertyYqRules = []*yqlib.ParticipleYqRule{
	{"PropertiesDecode", `from_?props|@propsd`, decodeOp(PropertiesFormat), 0},
	{"PropsEncode", `to_?props|@props`, encodeWithIndent(PropertiesFormat, 2), 0},
	{"LoadProperties", `load_?props`, loadOp(NewPropertiesDecoder(), false), 0},
}

func RegisterPropertiesFormat() {
	yqlib.RegisterFormat(PropertiesFormat)
	yqlib.RegisterRules(propertyYqRules)

}

var ConfiguredPropertiesPreferences = NewDefaultPropertiesPreferences()
