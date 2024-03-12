package xml

import "github.com/mikefarah/yq/v4/pkg/yqlib"

type XmlPreferences struct {
	Indent          int
	AttributePrefix string
	ContentName     string
	StrictMode      bool
	KeepNamespace   bool
	UseRawToken     bool
	ProcInstPrefix  string
	DirectiveName   string
	SkipProcInst    bool
	SkipDirectives  bool
}

func NewDefaultXmlPreferences() XmlPreferences {
	return XmlPreferences{
		Indent:          2,
		AttributePrefix: "+@",
		ContentName:     "+content",
		StrictMode:      false,
		KeepNamespace:   true,
		UseRawToken:     true,
		ProcInstPrefix:  "+p_",
		DirectiveName:   "+directive",
		SkipProcInst:    false,
		SkipDirectives:  false,
	}
}

func (p *XmlPreferences) Copy() XmlPreferences {
	return XmlPreferences{
		Indent:          p.Indent,
		AttributePrefix: p.AttributePrefix,
		ContentName:     p.ContentName,
		StrictMode:      p.StrictMode,
		KeepNamespace:   p.KeepNamespace,
		UseRawToken:     p.UseRawToken,
		ProcInstPrefix:  p.ProcInstPrefix,
		DirectiveName:   p.DirectiveName,
		SkipProcInst:    p.SkipProcInst,
		SkipDirectives:  p.SkipDirectives,
	}
}

var ConfiguredXMLPreferences = NewDefaultXmlPreferences()

var XMLFormat = &yqlib.Format{"xml", []string{"x"},
	func() yqlib.Encoder { return NewXMLEncoder(ConfiguredXMLPreferences) },
	func() yqlib.Decoder { return NewXMLDecoder(ConfiguredXMLPreferences) },
}

var xmlYqRules = []*yqlib.ParticipleYqRule{
	{"XMLEncodeWithIndent", `to_?xml\([0-9]+\)`, encodeParseIndent(XMLFormat), 0},
	{"XmlDecode", `from_?xml|@xmld`, decodeOp(XMLFormat), 0},
	{"XMLEncode", `to_?xml`, encodeWithIndent(XMLFormat, 2), 0},
	{"XMLEncodeNoIndent", `@xml`, encodeWithIndent(XMLFormat, 0), 0},
	{"LoadXML", `load_?xml|xml_?load`, loadOp(NewXMLDecoder(ConfiguredXMLPreferences), false), 0},
}

func RegisterXmlFormat() {
	yqlib.RegisterFormat(XMLFormat)
	yqlib.RegisterRules(xmlYqRules)

}
