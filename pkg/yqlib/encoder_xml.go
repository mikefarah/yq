package yqlib

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

var XmlPreferences = xmlPreferences{AttributePrefix: "+", ContentName: "+content"}

type xmlEncoder struct {
	attributePrefix string
	contentName     string
	indentString    string
}

func NewXmlEncoder(indent int, attributePrefix string, contentName string) Encoder {
	var indentString = ""

	for index := 0; index < indent; index++ {
		indentString = indentString + " "
	}
	return &xmlEncoder{attributePrefix, contentName, indentString}
}

func (e *xmlEncoder) CanHandleAliases() bool {
	return false
}

func (e *xmlEncoder) PrintDocumentSeparator(writer io.Writer) error {
	return nil
}

func (e *xmlEncoder) PrintLeadingContent(writer io.Writer, content string) error {
	return nil
}

func (e *xmlEncoder) Encode(writer io.Writer, node *yaml.Node) error {
	encoder := xml.NewEncoder(writer)
	encoder.Indent("", e.indentString)

	switch node.Kind {
	case yaml.MappingNode:
		err := e.encodeTopLevelMap(encoder, node)
		if err != nil {
			return err
		}
	case yaml.DocumentNode:
		err := e.encodeComment(encoder, headAndLineComment(node))
		if err != nil {
			return err
		}
		// this used to call encode...
		err = e.encodeTopLevelMap(encoder, unwrapDoc(node))
		if err != nil {
			return err
		}
		err = e.encodeComment(encoder, footComment(node))
		if err != nil {
			return err
		}
	case yaml.ScalarNode:
		var charData xml.CharData = []byte(node.Value)
		err := encoder.EncodeToken(charData)
		if err != nil {
			return err
		}
		return encoder.Flush()
	default:
		return fmt.Errorf("unsupported type %v", node.Tag)
	}
	var charData xml.CharData = []byte("\n")
	return encoder.EncodeToken(charData)

}

func (e *xmlEncoder) encodeTopLevelMap(encoder *xml.Encoder, node *yaml.Node) error {
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]

		start := xml.StartElement{Name: xml.Name{Local: key.Value}}
		err := e.encodeComment(encoder, headAndLineComment(key))
		if err != nil {
			return err
		}

		err = e.doEncode(encoder, value, start)
		if err != nil {
			return err
		}
		err = e.encodeComment(encoder, footComment(key))
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *xmlEncoder) encodeStart(encoder *xml.Encoder, node *yaml.Node, start xml.StartElement) error {
	err := encoder.EncodeToken(start)
	if err != nil {
		return err
	}
	return e.encodeComment(encoder, headAndLineComment(node))
}

func (e *xmlEncoder) encodeEnd(encoder *xml.Encoder, node *yaml.Node, start xml.StartElement) error {
	err := encoder.EncodeToken(start.End())
	if err != nil {
		return err
	}
	return e.encodeComment(encoder, footComment(node))
}

func (e *xmlEncoder) doEncode(encoder *xml.Encoder, node *yaml.Node, start xml.StartElement) error {
	switch node.Kind {
	case yaml.MappingNode:
		return e.encodeMap(encoder, node, start)
	case yaml.SequenceNode:
		return e.encodeArray(encoder, node, start)
	case yaml.ScalarNode:
		err := e.encodeStart(encoder, node, start)
		if err != nil {
			return err
		}

		var charData xml.CharData = []byte(node.Value)
		err = encoder.EncodeToken(charData)
		if err != nil {
			return err
		}

		return e.encodeEnd(encoder, node, start)
	}
	return fmt.Errorf("unsupported type %v", node.Tag)
}

func (e *xmlEncoder) encodeComment(encoder *xml.Encoder, commentStr string) error {
	if commentStr != "" {
		var comment xml.Comment = []byte(commentStr)
		err := encoder.EncodeToken(comment)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *xmlEncoder) encodeArray(encoder *xml.Encoder, node *yaml.Node, start xml.StartElement) error {
	for i := 0; i < len(node.Content); i++ {
		value := node.Content[i]
		err := e.doEncode(encoder, value, start.Copy())
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *xmlEncoder) encodeMap(encoder *xml.Encoder, node *yaml.Node, start xml.StartElement) error {

	//first find all the attributes and put them on the start token
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]

		if strings.HasPrefix(key.Value, e.attributePrefix) && key.Value != e.contentName {
			if value.Kind == yaml.ScalarNode {
				attributeName := strings.Replace(key.Value, e.attributePrefix, "", 1)
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

		if !strings.HasPrefix(key.Value, e.attributePrefix) && key.Value != e.contentName {
			start := xml.StartElement{Name: xml.Name{Local: key.Value}}
			err := e.doEncode(encoder, value, start)
			if err != nil {
				return err
			}
		} else if key.Value == e.contentName {
			// directly encode the contents
			var charData xml.CharData = []byte(value.Value)
			err = encoder.EncodeToken(charData)
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
