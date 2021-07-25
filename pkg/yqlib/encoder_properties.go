package yqlib

import (
	"fmt"
	"io"

	"github.com/magiconair/properties"
	yaml "gopkg.in/yaml.v3"
)

type propertiesEncoder struct {
	destination io.Writer
}

func NewPropertiesEncoder(destination io.Writer) Encoder {
	return &propertiesEncoder{destination}
}

func (pe *propertiesEncoder) Encode(node *yaml.Node) error {
	mapKeysToStrings(node)
	p := properties.NewProperties()
	err := pe.doEncode(p, node, "")
	if err != nil {
		return err
	}

	_, err = p.WriteComment(pe.destination, "#", properties.UTF8)
	return err
}

func (pe *propertiesEncoder) doEncode(p *properties.Properties, node *yaml.Node, path string) error {
	p.SetComment(path, node.HeadComment+node.LineComment)
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
