//go:build !yq_noxml

package yqlib

import (
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type xmlEncoder struct {
	indentString   string
	writer         io.Writer
	prefs          XmlPreferences
	leadingContent string
}

func NewXMLEncoder(prefs XmlPreferences) Encoder {
	var indentString = ""

	for index := 0; index < prefs.Indent; index++ {
		indentString = indentString + " "
	}
	return &xmlEncoder{indentString, nil, prefs, ""}
}

func (e *xmlEncoder) CanHandleAliases() bool {
	return false
}

func (e *xmlEncoder) PrintDocumentSeparator(_ io.Writer) error {
	return nil
}

func (e *xmlEncoder) PrintLeadingContent(_ io.Writer, content string) error {
	e.leadingContent = content
	return nil
}

func (e *xmlEncoder) Encode(writer io.Writer, node *CandidateNode) error {
	encoder := xml.NewEncoder(writer)
	// hack so we can manually add newlines to procInst and directives
	e.writer = writer
	encoder.Indent("", e.indentString)
	var newLine xml.CharData = []byte("\n")

	if node.Tag == "!!map" {
		// make sure <?xml .. ?> processing instructions are encoded first
		for i := 0; i < len(node.Content); i += 2 {
			key := node.Content[i]
			value := node.Content[i+1]

			if key.Value == (e.prefs.ProcInstPrefix + "xml") {
				name := strings.Replace(key.Value, e.prefs.ProcInstPrefix, "", 1)
				procInst := xml.ProcInst{Target: name, Inst: []byte(value.Value)}
				if err := encoder.EncodeToken(procInst); err != nil {
					return err
				}
				if _, err := e.writer.Write([]byte("\n")); err != nil {
					log.Warningf("Unable to write newline, skipping: %v", err)
				}
			}
		}
	}

	if e.leadingContent != "" {

		// remove first and last newlines if present
		err := e.encodeComment(encoder, e.leadingContent)
		if err != nil {
			return err
		}
		err = encoder.EncodeToken(newLine)
		if err != nil {
			return err
		}
	}

	switch node.Kind {
	case MappingNode:
		err := e.encodeTopLevelMap(encoder, node)
		if err != nil {
			return err
		}
	case ScalarNode:
		var charData xml.CharData = []byte(node.Value)
		err := encoder.EncodeToken(charData)
		if err != nil {
			return err
		}
		return encoder.Flush()
	default:
		return fmt.Errorf("cannot encode %v to XML - only maps can be encoded", node.Tag)
	}

	return encoder.EncodeToken(newLine)

}

func (e *xmlEncoder) encodeTopLevelMap(encoder *xml.Encoder, node *CandidateNode) error {
	err := e.encodeComment(encoder, headAndLineComment(node))
	if err != nil {
		return err
	}
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]

		start := xml.StartElement{Name: xml.Name{Local: key.Value}}
		log.Debugf("comments of key %v", key.Value)
		err := e.encodeComment(encoder, headAndLineComment(key))
		if err != nil {
			return err
		}
		if headAndLineComment(key) != "" {
			var newLine xml.CharData = []byte("\n")
			err = encoder.EncodeToken(newLine)
			if err != nil {
				return err
			}
		}

		if key.Value == (e.prefs.ProcInstPrefix + "xml") { //nolint
			// dont double process these.
		} else if strings.HasPrefix(key.Value, e.prefs.ProcInstPrefix) {
			name := strings.Replace(key.Value, e.prefs.ProcInstPrefix, "", 1)
			procInst := xml.ProcInst{Target: name, Inst: []byte(value.Value)}
			if err := encoder.EncodeToken(procInst); err != nil {
				return err
			}
			if _, err := e.writer.Write([]byte("\n")); err != nil {
				log.Warningf("Unable to write newline, skipping: %v", err)
			}
		} else if key.Value == e.prefs.DirectiveName {
			var directive xml.Directive = []byte(value.Value)
			if err := encoder.EncodeToken(directive); err != nil {
				return err
			}
			if _, err := e.writer.Write([]byte("\n")); err != nil {
				log.Warningf("Unable to write newline, skipping: %v", err)
			}
		} else {

			log.Debugf("recursing")

			if err := validateXMLName(key.Value, "element"); err != nil {
				return err
			}
			err = e.doEncode(encoder, value, start)
			if err != nil {
				return err
			}
		}
		err = e.encodeComment(encoder, footComment(key))
		if err != nil {
			return err
		}
	}
	return e.encodeComment(encoder, footComment(node))
}

func (e *xmlEncoder) encodeStart(encoder *xml.Encoder, node *CandidateNode, start xml.StartElement) error {
	err := encoder.EncodeToken(start)
	if err != nil {
		return err
	}
	return e.encodeComment(encoder, headComment(node))
}

func (e *xmlEncoder) encodeEnd(encoder *xml.Encoder, node *CandidateNode, start xml.StartElement) error {
	err := encoder.EncodeToken(start.End())
	if err != nil {
		return err
	}
	return e.encodeComment(encoder, footComment(node))
}

func (e *xmlEncoder) doEncode(encoder *xml.Encoder, node *CandidateNode, start xml.StartElement) error {
	switch node.Kind {
	case MappingNode:
		return e.encodeMap(encoder, node, start)
	case SequenceNode:
		return e.encodeArray(encoder, node, start)
	case ScalarNode:
		err := e.encodeStart(encoder, node, start)
		if err != nil {
			return err
		}

		var charData xml.CharData = []byte(node.Value)
		err = encoder.EncodeToken(charData)
		if err != nil {
			return err
		}

		if err = e.encodeComment(encoder, lineComment(node)); err != nil {
			return err
		}

		return e.encodeEnd(encoder, node, start)
	}
	return fmt.Errorf("unsupported type %v", node.Tag)
}

var xmlEncodeMultilineCommentRegex = regexp.MustCompile(`(^|\n) *# ?(.*)`)
var xmlEncodeSingleLineCommentRegex = regexp.MustCompile(`^\s*#(.*)\n?`)
var chompRegexp = regexp.MustCompile(`\n$`)

func (e *xmlEncoder) encodeComment(encoder *xml.Encoder, commentStr string) error {
	if commentStr != "" {
		log.Debugf("got comment [%v]", commentStr)
		// multi line string
		if len(commentStr) > 2 && strings.Contains(commentStr[1:len(commentStr)-1], "\n") {
			commentStr = chompRegexp.ReplaceAllString(commentStr, "")
			log.Debugf("chompRegexp [%v]", commentStr)
			commentStr = xmlEncodeMultilineCommentRegex.ReplaceAllString(commentStr, "$1$2")
			log.Debugf("processed multiline [%v]", commentStr)
			// if the first line is non blank, add a space
			if commentStr[0] != '\n' && commentStr[0] != ' ' {
				commentStr = " " + commentStr
			}

		} else {
			commentStr = xmlEncodeSingleLineCommentRegex.ReplaceAllString(commentStr, "$1")
		}

		if !strings.HasSuffix(commentStr, " ") && !strings.HasSuffix(commentStr, "\n") {
			commentStr = commentStr + " "
			log.Debugf("added suffix [%v]", commentStr)
		}
		log.Debugf("encoding comment [%v]", commentStr)

		var comment xml.Comment = []byte(commentStr)
		err := encoder.EncodeToken(comment)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *xmlEncoder) encodeArray(encoder *xml.Encoder, node *CandidateNode, start xml.StartElement) error {

	if err := e.encodeComment(encoder, headAndLineComment(node)); err != nil {
		return err
	}

	for i := 0; i < len(node.Content); i++ {
		value := node.Content[i]
		if err := e.doEncode(encoder, value, start.Copy()); err != nil {
			return err
		}
	}
	return e.encodeComment(encoder, footComment(node))
}

// isValidXMLName reports whether name matches the XML 1.0 Name production.
// The Go XML encoder writes element and attribute names verbatim, so a name
// that is not valid XML produces output that yq (and any other XML parser)
// cannot read back.
func isValidXMLName(name string) bool {
	if name == "" {
		return false
	}
	for index, r := range name {
		if index == 0 {
			if !isXMLNameStartChar(r) {
				return false
			}
		} else if !isXMLNameChar(r) {
			return false
		}
	}
	return true
}

func isXMLNameStartChar(r rune) bool {
	return r == ':' || r == '_' ||
		(r >= 'A' && r <= 'Z') ||
		(r >= 'a' && r <= 'z') ||
		(r >= 0xC0 && r <= 0xD6) ||
		(r >= 0xD8 && r <= 0xF6) ||
		(r >= 0xF8 && r <= 0x2FF) ||
		(r >= 0x370 && r <= 0x37D) ||
		(r >= 0x37F && r <= 0x1FFF) ||
		(r >= 0x200C && r <= 0x200D) ||
		(r >= 0x2070 && r <= 0x218F) ||
		(r >= 0x2C00 && r <= 0x2FEF) ||
		(r >= 0x3001 && r <= 0xD7FF) ||
		(r >= 0xF900 && r <= 0xFDCF) ||
		(r >= 0xFDF0 && r <= 0xFFFD) ||
		(r >= 0x10000 && r <= 0xEFFFF)
}

func isXMLNameChar(r rune) bool {
	return isXMLNameStartChar(r) ||
		r == '-' || r == '.' || r == 0xB7 ||
		(r >= '0' && r <= '9') ||
		(r >= 0x300 && r <= 0x36F) ||
		(r >= 0x203F && r <= 0x2040)
}

func validateXMLName(name string, nameType string) error {
	if !isValidXMLName(name) {
		return fmt.Errorf("cannot encode key '%v' to XML: invalid %v name", name, nameType)
	}
	return nil
}

func (e *xmlEncoder) isAttribute(name string) bool {
	return strings.HasPrefix(name, e.prefs.AttributePrefix) &&
		name != e.prefs.ContentName &&
		name != e.prefs.DirectiveName &&
		!strings.HasPrefix(name, e.prefs.ProcInstPrefix)
}

func (e *xmlEncoder) encodeMap(encoder *xml.Encoder, node *CandidateNode, start xml.StartElement) error {
	log.Debug("its a map")

	//first find all the attributes and put them on the start token
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]

		if e.isAttribute(key.Value) {
			if value.Kind == ScalarNode {
				attributeName := strings.Replace(key.Value, e.prefs.AttributePrefix, "", 1)
				if err := validateXMLName(attributeName, "attribute"); err != nil {
					return err
				}
				start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: attributeName}, Value: value.Value})
			} else {
				return fmt.Errorf("cannot use %v as attribute, only scalars are supported", value.Tag)
			}
		}
	}

	err := e.encodeStart(encoder, node, start)
	if err != nil {
		return err
	}

	//now we encode non attribute tokens
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]

		err := e.encodeComment(encoder, headAndLineComment(key))
		if err != nil {
			return err
		}
		if strings.HasPrefix(key.Value, e.prefs.ProcInstPrefix) {
			name := strings.Replace(key.Value, e.prefs.ProcInstPrefix, "", 1)
			procInst := xml.ProcInst{Target: name, Inst: []byte(value.Value)}
			if err := encoder.EncodeToken(procInst); err != nil {
				return err
			}
		} else if key.Value == e.prefs.DirectiveName {
			var directive xml.Directive = []byte(value.Value)
			if err := encoder.EncodeToken(directive); err != nil {
				return err
			}
		} else if key.Value == e.prefs.ContentName {
			// directly encode the contents
			err = e.encodeComment(encoder, headAndLineComment(value))
			if err != nil {
				return err
			}
			var charData xml.CharData = []byte(value.Value)
			err = encoder.EncodeToken(charData)
			if err != nil {
				return err
			}
			err = e.encodeComment(encoder, footComment(value))
			if err != nil {
				return err
			}
		} else if !e.isAttribute(key.Value) {
			if err := validateXMLName(key.Value, "element"); err != nil {
				return err
			}
			start := xml.StartElement{Name: xml.Name{Local: key.Value}}
			err := e.doEncode(encoder, value, start)
			if err != nil {
				return err
			}
		}
		err = e.encodeComment(encoder, footComment(key))
		if err != nil {
			return err
		}
	}

	return e.encodeEnd(encoder, node, start)
}
