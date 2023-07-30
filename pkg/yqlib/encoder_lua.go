package yqlib

import (
	"io"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type luaEncoder struct {
	docPrefix string
	docSuffix string
	escape    *strings.Replacer
}

func (le *luaEncoder) CanHandleAliases() bool {
	return false
}

func NewLuaEncoder() Encoder {
	escape := strings.NewReplacer("\n", "\\n", "\"", "\\\"", "\\", "\\\\")
	return &luaEncoder{"return ", ";\n", escape}
}

func (le *luaEncoder) PrintDocumentSeparator(writer io.Writer) error {
	return nil
}

func (le *luaEncoder) PrintLeadingContent(writer io.Writer, content string) error {
	return nil
}

func (le *luaEncoder) encodeString(writer io.Writer, node *yaml.Node) error {
	return writeString(writer, "\""+le.escape.Replace(node.Value)+"\"")
}

func (le *luaEncoder) encodeArray(writer io.Writer, node *yaml.Node) error {
	err := writeString(writer, "{")
	if err != nil {
		return err
	}
	for _, child := range node.Content {
		err := le.Encode(writer, child)
		if err != nil {
			return err
		}
		err = writeString(writer, ",")
		if err != nil {
			return err
		}
	}
	return writeString(writer, "}")
}

func (le *luaEncoder) encodeMap(writer io.Writer, node *yaml.Node) error {
	err := writeString(writer, "{")
	if err != nil {
		return err
	}
	for i, child := range node.Content {
		if (i % 2) == 1 {
			// value
			err = le.Encode(writer, child)
			if err != nil {
				return err
			}
			err = writeString(writer, ";")
			if err != nil {
				return err
			}
		} else {
			// key
			err := writeString(writer, "[")
			if err != nil {
				return err
			}
			err = le.encodeAny(writer, child)
			if err != nil {
				return err
			}
			err = writeString(writer, "]=")
			if err != nil {
				return err
			}
		}
	}
	return writeString(writer, "}")
}

func (le *luaEncoder) encodeAny(writer io.Writer, node *yaml.Node) error {
	switch node.Kind {
	case yaml.SequenceNode:
		return le.encodeArray(writer, node)
	case yaml.MappingNode:
		return le.encodeMap(writer, node)
	case yaml.ScalarNode:
		switch node.Tag {
		case "!!str":
			return le.encodeString(writer, node)
		case "!!null":
			return writeString(writer, "nil")
		default:
			return writeString(writer, node.Value)
		}
	case yaml.DocumentNode:
		err := writeString(writer, le.docPrefix)
		if err != nil {
			return err
		}
		err = le.encodeAny(writer, node.Content[0])
		if err != nil {
			return err
		}
		return writeString(writer, le.docSuffix)
	default:
		return writeString(writer, "nil --[[ encoder NYI -- "+node.ShortTag()+" ]]")
	}
}

func (le *luaEncoder) Encode(writer io.Writer, node *yaml.Node) error {
	return le.encodeAny(writer, node)
}
