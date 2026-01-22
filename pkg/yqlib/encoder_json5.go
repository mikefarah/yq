//go:build !yq_nojson5

package yqlib

import (
	"bytes"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
)

type json5Encoder struct {
	prefs        JsonPreferences
	indentString string
}

func NewJSON5Encoder(prefs JsonPreferences) Encoder {
	indentString := ""
	for i := 0; i < prefs.Indent; i++ {
		indentString += " "
	}
	return &json5Encoder{prefs: prefs, indentString: indentString}
}

func (je *json5Encoder) CanHandleAliases() bool {
	return false
}

func (je *json5Encoder) PrintDocumentSeparator(_ io.Writer) error {
	return nil
}

func (je *json5Encoder) PrintLeadingContent(_ io.Writer, _ string) error {
	return nil
}

func (je *json5Encoder) Encode(writer io.Writer, node *CandidateNode) error {
	if node.Kind == ScalarNode && je.prefs.UnwrapScalar {
		return writeString(writer, node.Value+"\n")
	}

	destination := writer
	tempBuffer := bytes.NewBuffer(nil)
	if je.prefs.ColorsEnabled {
		destination = tempBuffer
	}

	if err := writeJSON5CommentBlock(destination, node.HeadComment, je.indentString, 0); err != nil {
		return err
	}
	if err := encodeJSON5Node(destination, node, je.indentString, 0); err != nil {
		return err
	}
	if err := writeJSON5InlineComment(destination, node.LineComment, true); err != nil {
		return err
	}
	if err := writeString(destination, "\n"); err != nil {
		return err
	}
	if err := writeJSON5CommentBlock(destination, node.FootComment, je.indentString, 0); err != nil {
		return err
	}

	if je.prefs.ColorsEnabled {
		return colorizeAndPrint(tempBuffer.Bytes(), writer)
	}
	return nil
}

func encodeJSON5Node(writer io.Writer, node *CandidateNode, indentString string, depth int) error {
	if node == nil {
		return writeString(writer, "null")
	}

	switch node.Kind {
	case AliasNode:
		return encodeJSON5Node(writer, node.Alias, indentString, depth)
	case ScalarNode:
		return writeString(writer, json5ScalarString(node))
	case SequenceNode:
		return encodeJSON5Sequence(writer, node, indentString, depth)
	case MappingNode:
		return encodeJSON5Mapping(writer, node, indentString, depth)
	default:
		return writeString(writer, "null")
	}
}

func encodeJSON5Sequence(writer io.Writer, node *CandidateNode, indentString string, depth int) error {
	if err := writeString(writer, "["); err != nil {
		return err
	}
	if len(node.Content) == 0 {
		return writeString(writer, "]")
	}

	pretty := indentString != ""
	if pretty {
		if err := writeString(writer, "\n"); err != nil {
			return err
		}
	}

	for i, child := range node.Content {
		if pretty {
			if err := writeJSON5CommentBlock(writer, child.HeadComment, indentString, depth+1); err != nil {
				return err
			}
			if err := writeString(writer, strings.Repeat(indentString, depth+1)); err != nil {
				return err
			}
		}
		if err := encodeJSON5Node(writer, child, indentString, depth+1); err != nil {
			return err
		}
		if err := writeJSON5InlineComment(writer, child.LineComment, true); err != nil {
			return err
		}
		if i != len(node.Content)-1 {
			if err := writeString(writer, ","); err != nil {
				return err
			}
		}
		if pretty {
			if err := writeString(writer, "\n"); err != nil {
				return err
			}
		}
	}

	if pretty {
		if err := writeString(writer, strings.Repeat(indentString, depth)); err != nil {
			return err
		}
	}
	return writeString(writer, "]")
}

func encodeJSON5Mapping(writer io.Writer, node *CandidateNode, indentString string, depth int) error {
	if err := writeString(writer, "{"); err != nil {
		return err
	}
	if len(node.Content) == 0 {
		return writeString(writer, "}")
	}

	pretty := indentString != ""
	if pretty {
		if err := writeString(writer, "\n"); err != nil {
			return err
		}
	}

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		keyBytes, err := json.Marshal(keyNode.Value)
		if err != nil {
			return err
		}

		if pretty {
			if err := writeJSON5CommentBlock(writer, keyNode.HeadComment, indentString, depth+1); err != nil {
				return err
			}
			if err := writeString(writer, strings.Repeat(indentString, depth+1)); err != nil {
				return err
			}
		}
		if err := writeString(writer, string(keyBytes)); err != nil {
			return err
		}
		if err := writeJSON5InlineComment(writer, keyNode.LineComment, true); err != nil {
			return err
		}
		if err := writeString(writer, ":"); err != nil {
			return err
		}
		if pretty && strings.TrimSpace(valueNode.HeadComment) != "" && strings.Contains(valueNode.HeadComment, "\n") {
			if err := writeString(writer, "\n"); err != nil {
				return err
			}
			if err := writeJSON5CommentBlock(writer, valueNode.HeadComment, indentString, depth+1); err != nil {
				return err
			}
			if err := writeString(writer, strings.Repeat(indentString, depth+1)); err != nil {
				return err
			}
		} else {
			if pretty {
				if err := writeString(writer, " "); err != nil {
					return err
				}
			}
			if err := writeJSON5InlineComment(writer, valueNode.HeadComment, false); err != nil {
				return err
			}
			if strings.TrimSpace(valueNode.HeadComment) != "" {
				if err := writeString(writer, " "); err != nil {
					return err
				}
			}
		}

		if err := encodeJSON5Node(writer, valueNode, indentString, depth+1); err != nil {
			return err
		}
		if err := writeJSON5InlineComment(writer, valueNode.LineComment, true); err != nil {
			return err
		}
		if i != len(node.Content)-2 {
			if err := writeString(writer, ","); err != nil {
				return err
			}
		}
		if pretty {
			if err := writeString(writer, "\n"); err != nil {
				return err
			}
		}
	}

	if pretty {
		if err := writeString(writer, strings.Repeat(indentString, depth)); err != nil {
			return err
		}
	}
	return writeString(writer, "}")
}

func json5ScalarString(node *CandidateNode) string {
	tag := node.guessTagFromCustomType()

	switch tag {
	case "!!null":
		return "null"
	case "!!bool":
		if isTruthyNode(node) {
			return "true"
		}
		return "false"
	case "!!int":
		value, err := node.GetValueRep()
		if err != nil {
			return "null"
		}
		intBytes, err := json.Marshal(value)
		if err != nil {
			return "null"
		}
		return string(intBytes)
	case "!!float":
		return json5FloatString(node.Value)
	case "!!str":
		stringBytes, err := json.Marshal(node.Value)
		if err != nil {
			return "null"
		}
		return string(stringBytes)
	default:
		stringBytes, err := json.Marshal(node.Value)
		if err != nil {
			return "null"
		}
		return string(stringBytes)
	}
}

func json5FloatString(value string) string {
	switch strings.ToLower(value) {
	case ".inf", "+.inf", "inf", "+inf", "+infinity", "infinity":
		return "Infinity"
	case "-.inf", "-inf", "-infinity":
		return "-Infinity"
	case ".nan", "nan":
		return "NaN"
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return "null"
	}

	if math.IsNaN(parsed) {
		return "NaN"
	}
	if math.IsInf(parsed, 1) {
		return "Infinity"
	}
	if math.IsInf(parsed, -1) {
		return "-Infinity"
	}

	floatBytes, err := json.Marshal(parsed)
	if err != nil {
		return "null"
	}
	return string(floatBytes)
}

func writeJSON5CommentBlock(writer io.Writer, comment string, indentString string, depth int) error {
	comment = strings.ReplaceAll(comment, "\r\n", "\n")
	comment = strings.TrimSpace(comment)
	if comment == "" {
		return nil
	}

	indent := strings.Repeat(indentString, depth)
	lines := strings.Split(comment, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = strings.TrimSpace(strings.TrimPrefix(line, "#"))
		line = strings.TrimSpace(strings.TrimPrefix(line, "//"))
		if err := writeString(writer, indent); err != nil {
			return err
		}
		if err := writeString(writer, "// "+line+"\n"); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON5InlineComment(writer io.Writer, comment string, forceLeadingSpace bool) error {
	comment = strings.ReplaceAll(comment, "\r\n", "\n")
	comment = strings.TrimSpace(comment)
	if comment == "" {
		return nil
	}
	comment = strings.TrimSpace(strings.TrimPrefix(comment, "#"))
	comment = strings.TrimSpace(strings.TrimPrefix(comment, "//"))
	comment = strings.ReplaceAll(comment, "\n", " ")
	comment = strings.Join(strings.Fields(comment), " ")

	prefix := ""
	if forceLeadingSpace {
		prefix = " "
	}
	return writeString(writer, prefix+"/* "+comment+" */")
}
