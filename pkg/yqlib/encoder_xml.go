package yqlib

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type xmlEncoder struct {
	xmlEncoder      *xml.Encoder
	attributePrefix string
}

func NewXmlEncoder(writer io.Writer, indent int, attributePrefix string) Encoder {
	encoder := xml.NewEncoder(writer)
	var indentString = ""

	for index := 0; index < indent; index++ {
		indentString = indentString + " "
	}
	encoder.Indent("", indentString)
	return &xmlEncoder{encoder, attributePrefix}
}
func (e *xmlEncoder) Encode(node *yaml.Node) error {
	switch node.Kind {
	case yaml.MappingNode:
		return e.encodeTopLevelMap(node)
	case yaml.DocumentNode:
		return e.Encode(unwrapDoc(node))
	case yaml.ScalarNode:
		var charData xml.CharData = []byte(node.Value)
		return e.xmlEncoder.EncodeToken(charData)
	}
	return fmt.Errorf("unsupported type %v", node.Tag)
}

func (e *xmlEncoder) encodeTopLevelMap(node *yaml.Node) error {
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]

		start := xml.StartElement{Name: xml.Name{Local: key.Value}}
		err := e.doEncode(value, start)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *xmlEncoder) doEncode(node *yaml.Node, start xml.StartElement) error {
	switch node.Kind {
	case yaml.MappingNode:
		return e.encodeMap(node, start)
	case yaml.SequenceNode:
		return e.encodeArray(node, start)
	case yaml.ScalarNode:
		err := e.xmlEncoder.EncodeToken(start)
		if err != nil {
			return err
		}

		var charData xml.CharData = []byte(node.Value)
		err = e.xmlEncoder.EncodeToken(charData)

		if err != nil {
			return err
		}
		return e.xmlEncoder.EncodeToken(start.End())
	}
	return fmt.Errorf("unsupported type %v", node.Tag)
}

func (e *xmlEncoder) encodeArray(node *yaml.Node, start xml.StartElement) error {
	for i := 0; i < len(node.Content); i++ {
		value := node.Content[i]
		err := e.doEncode(value, start.Copy())
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *xmlEncoder) encodeMap(node *yaml.Node, start xml.StartElement) error {

	//first find all the attributes and put them on the start token
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]

		if strings.HasPrefix(key.Value, e.attributePrefix) {
			if value.Kind == yaml.ScalarNode {
				attributeName := strings.Replace(key.Value, e.attributePrefix, "", 1)
				start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: attributeName}, Value: value.Value})
			} else {
				return fmt.Errorf("cannot use %v as attribute, only scalars are supported", value.Tag)
			}
		}
	}

	err := e.xmlEncoder.EncodeToken(start)
	if err != nil {
		return err
	}

	//now we encode non attribute tokens
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]

		if !strings.HasPrefix(key.Value, e.attributePrefix) {
			start := xml.StartElement{Name: xml.Name{Local: key.Value}}
			err := e.doEncode(value, start)
			if err != nil {
				return err
			}
		}
	}

	return e.xmlEncoder.EncodeToken(start.End())
}
