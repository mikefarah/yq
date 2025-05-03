//go:build !yq_nolua

package yqlib

import (
	"fmt"
	"io"
	"strings"
)

type luaEncoder struct {
	docPrefix string
	docSuffix string
	indent    int
	indentStr string
	unquoted  bool
	globals   bool
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
	return &luaEncoder{unescape.Replace(prefs.DocPrefix), unescape.Replace(prefs.DocSuffix), 0, "\t", prefs.UnquotedKeys, prefs.Globals, escape}
}

func (le *luaEncoder) PrintDocumentSeparator(_ io.Writer) error {
	return nil
}

func (le *luaEncoder) PrintLeadingContent(_ io.Writer, _ string) error {
	return nil
}

func (le *luaEncoder) encodeString(writer io.Writer, node *CandidateNode) error {
	quote := "\""
	switch node.Style {
	case LiteralStyle, FoldedStyle, FlowStyle:
		for i := 0; i < 10; i++ {
			if !strings.Contains(node.Value, "]"+strings.Repeat("=", i)+"]") {
				err := writeString(writer, "["+strings.Repeat("=", i)+"[\n")
				if err != nil {
					return err
				}
				err = writeString(writer, node.Value)
				if err != nil {
					return err
				}
				return writeString(writer, "]"+strings.Repeat("=", i)+"]")
			}
		}
	case SingleQuotedStyle:
		quote = "'"

		// fallthrough to regular ol' string
	}
	return writeString(writer, quote+le.escape.Replace(node.Value)+quote)
}

func (le *luaEncoder) writeIndent(writer io.Writer) error {
	if le.indentStr == "" {
		return nil
	}
	err := writeString(writer, "\n")
	if err != nil {
		return err
	}
	return writeString(writer, strings.Repeat(le.indentStr, le.indent))
}

func (le *luaEncoder) encodeArray(writer io.Writer, node *CandidateNode) error {
	err := writeString(writer, "{")
	if err != nil {
		return err
	}
	le.indent++
	for _, child := range node.Content {
		err = le.writeIndent(writer)
		if err != nil {
			return err
		}
		err := le.encodeAny(writer, child)
		if err != nil {
			return err
		}
		err = writeString(writer, ",")
		if err != nil {
			return err
		}
		if child.LineComment != "" {
			sansPrefix, _ := strings.CutPrefix(child.LineComment, "#")
			err = writeString(writer, " --"+sansPrefix)
			if err != nil {
				return err
			}
		}
	}
	le.indent--
	if len(node.Content) != 0 {
		err = le.writeIndent(writer)
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
			// keeping for legacy reasons, upgraded linter
			//nolint:staticcheck
			if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '_') {
				return true
			}
		} else {
			// keeping for legacy reasons, upgraded linter
			//nolint:staticcheck
			if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '_') {
				return true
			}
		}
	}
	return false
}

func (le *luaEncoder) encodeMap(writer io.Writer, node *CandidateNode, global bool) error {
	if !global {
		err := writeString(writer, "{")
		if err != nil {
			return err
		}
		le.indent++
	}
	for i, child := range node.Content {
		if (i % 2) == 1 {
			// value
			err := le.encodeAny(writer, child)
			if err != nil {
				return err
			}
			err = writeString(writer, ";")
			if err != nil {
				return err
			}
		} else {
			// key
			if !global || i > 0 {
				err := le.writeIndent(writer)
				if err != nil {
					return err
				}
			}
			if (le.unquoted || global) && child.Tag == "!!str" && !needsQuoting(child.Value) {
				err := writeString(writer, child.Value+" = ")
				if err != nil {
					return err
				}
			} else {
				if global {
					// This only works in Lua 5.2+
					err := writeString(writer, "_ENV")
					if err != nil {
						return err
					}
				}
				err := writeString(writer, "[")
				if err != nil {
					return err
				}
				err = le.encodeAny(writer, child)
				if err != nil {
					return err
				}
				err = writeString(writer, "] = ")
				if err != nil {
					return err
				}
			}
		}
		if child.LineComment != "" {
			sansPrefix, _ := strings.CutPrefix(child.LineComment, "#")
			err := writeString(writer, strings.Repeat(" ", i%2)+"--"+sansPrefix)
			if err != nil {
				return err
			}
			if (i % 2) == 0 {
				// newline and indent after comments on keys
				err = le.writeIndent(writer)
				if err != nil {
					return err
				}
			}
		}
	}
	if global {
		return writeString(writer, "\n")
	}
	le.indent--
	if len(node.Content) != 0 {
		err := le.writeIndent(writer)
		if err != nil {
			return err
		}
	}
	return writeString(writer, "}")
}

func (le *luaEncoder) encodeAny(writer io.Writer, node *CandidateNode) error {
	switch node.Kind {
	case SequenceNode:
		return le.encodeArray(writer, node)
	case MappingNode:
		return le.encodeMap(writer, node, false)
	case ScalarNode:
		switch node.Tag {
		case "!!str":
			return le.encodeString(writer, node)
		case "!!null":
			// TODO reject invalid use as a table key
			return writeString(writer, "nil")
		case "!!bool":
			// Yaml 1.2 has case variation e.g. True, FALSE etc but Lua only has
			// lower case
			return writeString(writer, strings.ToLower(node.Value))
		case "!!int":
			if strings.HasPrefix(node.Value, "0o") {
				_, octalValue, err := parseInt64(node.Value)
				if err != nil {
					return err
				}
				return writeString(writer, fmt.Sprintf("%d", octalValue))
			}
			return writeString(writer, strings.ToLower(node.Value))
		case "!!float":
			switch strings.ToLower(node.Value) {
			case ".inf", "+.inf":
				return writeString(writer, "(1/0)")
			case "-.inf":
				return writeString(writer, "(-1/0)")
			case ".nan":
				return writeString(writer, "(0/0)")
			default:
				return writeString(writer, node.Value)
			}
		default:
			return fmt.Errorf("lua encoder NYI -- %s", node.Tag)
		}
	default:
		return fmt.Errorf("lua encoder NYI -- %s", node.Tag)
	}
}

func (le *luaEncoder) encodeTopLevel(writer io.Writer, node *CandidateNode) error {
	err := writeString(writer, le.docPrefix)
	if err != nil {
		return err
	}
	err = le.encodeAny(writer, node)
	if err != nil {
		return err
	}
	return writeString(writer, le.docSuffix)
}

func (le *luaEncoder) Encode(writer io.Writer, node *CandidateNode) error {

	if le.globals {
		if node.Kind != MappingNode {
			return fmt.Errorf("--lua-global requires a top level MappingNode")
		}
		return le.encodeMap(writer, node, true)
	}
	return le.encodeTopLevel(writer, node)
}
