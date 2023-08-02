package yqlib

import (
	"io"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type luaEncoder struct {
	docPrefix string
	docSuffix string
	unquoted  bool
	escape    *strings.Replacer
}

func (le *luaEncoder) CanHandleAliases() bool {
	return false
}

func NewLuaEncoder(prefs LuaPreferences) Encoder {
	escape := strings.NewReplacer(
		"\000", "\\000",
		"\001", "\\001",
		"\002", "\\002",
		"\003", "\\003",
		"\004", "\\004",
		"\005", "\\005",
		"\006", "\\006",
		"\007", "\\a",
		"\010", "\\b",
		"\011", "\\t",
		"\012", "\\n",
		"\013", "\\v",
		"\014", "\\f",
		"\015", "\\r",
		"\016", "\\014",
		"\017", "\\015",
		"\020", "\\016",
		"\021", "\\017",
		"\022", "\\018",
		"\023", "\\019",
		"\024", "\\020",
		"\025", "\\021",
		"\026", "\\022",
		"\027", "\\023",
		"\030", "\\024",
		"\031", "\\025",
		"\032", "\\026",
		"\033", "\\027",
		"\034", "\\028",
		"\035", "\\029",
		"\036", "\\030",
		"\037", "\\031",
		"\"", "\\\"",
		"'", "\\'",
		"\\", "\\\\",
		"\177", "\\127",
	)
	unescape := strings.NewReplacer(
		"\\'", "'",
		"\\\"", "\"",
		"\\n", "\n",
		"\\r", "\r",
		"\\t", "\t",
		"\\\\", "\\",
	)
	return &luaEncoder{unescape.Replace(prefs.DocPrefix), unescape.Replace(prefs.DocSuffix), prefs.UnquotedKeys, escape}
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

func needsQuoting(s string) bool {
	// known keywords as of Lua 5.4
	switch s {
	case "do", "and", "else", "break",
		"if", "end", "goto", "false",
		"in", "for", "then", "local",
		"or", "nil", "true", "until",
		"elseif", "function", "not",
		"repeat", "return", "while":
		return true
	}
	// [%a_][%w_]*
	for i, c := range s {
		if i == 0 {
			if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '_') {
				return true
			}
		} else {
			if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '_') {
				return true
			}
		}
	}
	return false
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
		} else if le.unquoted && child.Tag == "!!str" && !needsQuoting(child.Value) {
			err = writeString(writer, child.Value+"=")
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
		panic("Lua encoder NYI -- " + node.ShortTag())
	}
}

func (le *luaEncoder) Encode(writer io.Writer, node *yaml.Node) error {
	return le.encodeAny(writer, node)
}
