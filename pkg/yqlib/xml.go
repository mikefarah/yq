package yqlib

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
