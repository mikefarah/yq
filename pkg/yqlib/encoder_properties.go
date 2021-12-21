package yqlib

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/magiconair/properties"
	yaml "gopkg.in/yaml.v3"
)

type propertiesEncoder struct {
}

func NewPropertiesEncoder() Encoder {
	return &propertiesEncoder{}
}

func (pe *propertiesEncoder) CanHandleAliases() bool {
	return false
}

func (pe *propertiesEncoder) PrintDocumentSeparator(writer io.Writer) error {
	return nil
}

func (pe *propertiesEncoder) PrintLeadingContent(writer io.Writer, content string) error {
	reader := bufio.NewReader(strings.NewReader(content))
	for {

		readline, errReading := reader.ReadString('\n')
		if errReading != nil && !errors.Is(errReading, io.EOF) {
			return errReading
		}
		if strings.Contains(readline, "$yqDocSeperator$") {

			if err := pe.PrintDocumentSeparator(writer); err != nil {
				return err
			}

		} else {
			if err := writeString(writer, readline); err != nil {
				return err
			}
		}

		if errors.Is(errReading, io.EOF) {
			if readline != "" {
				// the last comment we read didn't have a new line, put one in
				if err := writeString(writer, "\n"); err != nil {
					return err
				}
			}
			break
		}
	}
	return nil
}

func (pe *propertiesEncoder) Encode(writer io.Writer, node *yaml.Node) error {
	mapKeysToStrings(node)
	p := properties.NewProperties()
	err := pe.doEncode(p, node, "")
	if err != nil {
		return err
	}

	_, err = p.WriteComment(writer, "#", properties.UTF8)
	return err
}

func (pe *propertiesEncoder) doEncode(p *properties.Properties, node *yaml.Node, path string) error {
	p.SetComment(path, headAndLineComment(node))
	switch node.Kind {
	case yaml.ScalarNode:
		_, _, err := p.Set(path, node.Value)
		return err
	case yaml.DocumentNode:
		return pe.doEncode(p, node.Content[0], path)
	case yaml.SequenceNode:
		return pe.encodeArray(p, node.Content, path)
	case yaml.MappingNode:
		return pe.encodeMap(p, node.Content, path)
	case yaml.AliasNode:
		return pe.doEncode(p, node.Alias, path)
	default:
		return fmt.Errorf("Unsupported node %v", node.Tag)
	}
}

func (pe *propertiesEncoder) appendPath(path string, key interface{}) string {
	if path == "" {
		return fmt.Sprintf("%v", key)
	}
	return fmt.Sprintf("%v.%v", path, key)
}

func (pe *propertiesEncoder) encodeArray(p *properties.Properties, kids []*yaml.Node, path string) error {
	for index, child := range kids {
		err := pe.doEncode(p, child, pe.appendPath(path, index))
		if err != nil {
			return err
		}
	}
	return nil
}

func (pe *propertiesEncoder) encodeMap(p *properties.Properties, kids []*yaml.Node, path string) error {
	for index := 0; index < len(kids); index = index + 2 {
		key := kids[index]
		value := kids[index+1]
		err := pe.doEncode(p, value, pe.appendPath(path, key.Value))
		if err != nil {
			return err
		}
	}
	return nil
}
