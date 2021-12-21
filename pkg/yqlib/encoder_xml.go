package yqlib

import (
	"encoding/xml"
	"fmt"
	"io"

	yaml "gopkg.in/yaml.v3"
)

type xmlEncoder struct {
	xmlEncoder *xml.Encoder
}

func NewXmlEncoder(writer io.Writer, indent int) Encoder {
	encoder := xml.NewEncoder(writer)
	var indentString = ""

	for index := 0; index < indent; index++ {
		indentString = indentString + " "
	}
	encoder.Indent("", indentString)
	return &xmlEncoder{encoder}
}
func (e *xmlEncoder) Encode(node *yaml.Node) error {
	switch node.Kind {
	case yaml.MappingNode:
		return e.encodeMap(node)
	case yaml.DocumentNode:
		return e.Encode(unwrapDoc(node))
	case yaml.ScalarNode:
		var charData xml.CharData = []byte(node.Value)
		return e.xmlEncoder.EncodeToken(charData)
	}
	return fmt.Errorf("unsupported type %v", node.Tag)
}

func (e *xmlEncoder) doEncode(node *yaml.Node, start xml.StartElement) error {
	switch node.Kind {
	case yaml.MappingNode:
		err := e.xmlEncoder.EncodeToken(start)
		if err != nil {
			return err
		}
		err = e.encodeMap(node)
		if err != nil {
			return err
		}
		return e.xmlEncoder.EncodeToken(start.End())
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

func (e *xmlEncoder) encodeMap(node *yaml.Node) error {
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
