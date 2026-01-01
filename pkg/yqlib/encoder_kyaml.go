//go:build !yq_nokyaml

package yqlib

import (
	"bytes"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type kyamlEncoder struct {
	prefs KYamlPreferences
}

func NewKYamlEncoder(prefs KYamlPreferences) Encoder {
	return &kyamlEncoder{prefs: prefs}
}

func (ke *kyamlEncoder) CanHandleAliases() bool {
	// KYAML is a restricted subset; avoid emitting anchors/aliases.
	return false
}

func (ke *kyamlEncoder) PrintDocumentSeparator(writer io.Writer) error {
	return PrintYAMLDocumentSeparator(writer, ke.prefs.PrintDocSeparators)
}

func (ke *kyamlEncoder) PrintLeadingContent(writer io.Writer, content string) error {
	return PrintYAMLLeadingContent(writer, content, ke.prefs.PrintDocSeparators, ke.prefs.ColorsEnabled)
}

func (ke *kyamlEncoder) Encode(writer io.Writer, node *CandidateNode) error {
	log.Debug("encoderKYaml - going to print %v", NodeToString(node))
	if node.Kind == ScalarNode && ke.prefs.UnwrapScalar {
		return writeString(writer, node.Value+"\n")
	}

	destination := writer
	tempBuffer := bytes.NewBuffer(nil)
	if ke.prefs.ColorsEnabled {
		destination = tempBuffer
	}

	// Mirror the YAML encoder behaviour: trailing comments on the document root
	// are stored in FootComment and need to be printed after the document.
	trailingContent := node.FootComment

	if err := ke.writeCommentBlock(destination, node.HeadComment, 0); err != nil {
		return err
	}
	if err := ke.writeNode(destination, node, 0); err != nil {
		return err
	}
	if err := ke.writeInlineComment(destination, node.LineComment); err != nil {
		return err
	}
	if err := writeString(destination, "\n"); err != nil {
		return err
	}
	if err := ke.PrintLeadingContent(destination, trailingContent); err != nil {
		return err
	}

	if ke.prefs.ColorsEnabled {
		return colorizeAndPrint(tempBuffer.Bytes(), writer)
	}
	return nil
}

func (ke *kyamlEncoder) writeNode(writer io.Writer, node *CandidateNode, indent int) error {
	switch node.Kind {
	case MappingNode:
		return ke.writeMapping(writer, node, indent)
	case SequenceNode:
		return ke.writeSequence(writer, node, indent)
	case ScalarNode:
		return writeString(writer, ke.formatScalar(node))
	case AliasNode:
		// Should have been exploded by the printer, but handle defensively.
		if node.Alias == nil {
			return writeString(writer, "null")
		}
		return ke.writeNode(writer, node.Alias, indent)
	default:
		return writeString(writer, "null")
	}
}

func (ke *kyamlEncoder) writeMapping(writer io.Writer, node *CandidateNode, indent int) error {
	if len(node.Content) == 0 {
		return writeString(writer, "{}")
	}
	if err := writeString(writer, "{\n"); err != nil {
		return err
	}

	for i := 0; i+1 < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		entryIndent := indent + ke.prefs.Indent
		if err := ke.writeCommentBlock(writer, keyNode.HeadComment, entryIndent); err != nil {
			return err
		}
		if valueNode.HeadComment != "" && valueNode.HeadComment != keyNode.HeadComment {
			if err := ke.writeCommentBlock(writer, valueNode.HeadComment, entryIndent); err != nil {
				return err
			}
		}

		if err := ke.writeIndent(writer, entryIndent); err != nil {
			return err
		}
		if err := writeString(writer, ke.formatKey(keyNode)); err != nil {
			return err
		}
		if err := writeString(writer, ": "); err != nil {
			return err
		}
		if err := ke.writeNode(writer, valueNode, entryIndent); err != nil {
			return err
		}

		// Always emit a trailing comma; KYAML encourages explicit separators,
		// and this ensures all quoted strings have a trailing `",` as requested.
		if err := writeString(writer, ","); err != nil {
			return err
		}
		inline := valueNode.LineComment
		if inline == "" {
			inline = keyNode.LineComment
		}
		if err := ke.writeInlineComment(writer, inline); err != nil {
			return err
		}
		if err := writeString(writer, "\n"); err != nil {
			return err
		}

		foot := valueNode.FootComment
		if foot == "" {
			foot = keyNode.FootComment
		}
		if err := ke.writeCommentBlock(writer, foot, entryIndent); err != nil {
			return err
		}
	}

	if err := ke.writeIndent(writer, indent); err != nil {
		return err
	}
	return writeString(writer, "}")
}

func (ke *kyamlEncoder) writeSequence(writer io.Writer, node *CandidateNode, indent int) error {
	if len(node.Content) == 0 {
		return writeString(writer, "[]")
	}
	if err := writeString(writer, "[\n"); err != nil {
		return err
	}

	for _, child := range node.Content {
		itemIndent := indent + ke.prefs.Indent
		if err := ke.writeCommentBlock(writer, child.HeadComment, itemIndent); err != nil {
			return err
		}
		if err := ke.writeIndent(writer, itemIndent); err != nil {
			return err
		}
		if err := ke.writeNode(writer, child, itemIndent); err != nil {
			return err
		}
		if err := writeString(writer, ","); err != nil {
			return err
		}
		if err := ke.writeInlineComment(writer, child.LineComment); err != nil {
			return err
		}
		if err := writeString(writer, "\n"); err != nil {
			return err
		}
		if err := ke.writeCommentBlock(writer, child.FootComment, itemIndent); err != nil {
			return err
		}
	}

	if err := ke.writeIndent(writer, indent); err != nil {
		return err
	}
	return writeString(writer, "]")
}

func (ke *kyamlEncoder) writeIndent(writer io.Writer, indent int) error {
	if indent <= 0 {
		return nil
	}
	return writeString(writer, strings.Repeat(" ", indent))
}

func (ke *kyamlEncoder) formatKey(keyNode *CandidateNode) string {
	// KYAML examples use bare keys. Quote keys only when needed.
	key := keyNode.Value
	if isValidKYamlBareKey(key) {
		return key
	}
	return `"` + escapeDoubleQuotedString(key) + `"`
}

func (ke *kyamlEncoder) formatScalar(node *CandidateNode) string {
	switch node.Tag {
	case "!!null":
		return "null"
	case "!!bool":
		return strings.ToLower(node.Value)
	case "!!int", "!!float":
		return node.Value
	case "!!str":
		return `"` + escapeDoubleQuotedString(node.Value) + `"`
	default:
		// Fall back to a string representation to avoid implicit typing surprises.
		return `"` + escapeDoubleQuotedString(node.Value) + `"`
	}
}

var kyamlBareKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_-]*$`)

func isValidKYamlBareKey(s string) bool {
	// Conservative: require an identifier-like key; otherwise quote.
	if s == "" {
		return false
	}
	return kyamlBareKeyRe.MatchString(s)
}

func escapeDoubleQuotedString(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 2)

	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString(`\\`)
		case '"':
			b.WriteString(`\"`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			if r < 0x20 {
				// YAML double-quoted strings support \uXXXX escapes.
				b.WriteString(`\u`)
				hex := "0000" + strings.ToUpper(strconv.FormatInt(int64(r), 16))
				b.WriteString(hex[len(hex)-4:])
			} else {
				b.WriteRune(r)
			}
		}
	}
	return b.String()
}

func (ke *kyamlEncoder) writeCommentBlock(writer io.Writer, comment string, indent int) error {
	if strings.TrimSpace(comment) == "" {
		return nil
	}

	lines := strings.Split(strings.ReplaceAll(comment, "\r\n", "\n"), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if err := ke.writeIndent(writer, indent); err != nil {
			return err
		}

		toWrite := line
		if !commentLineRe.MatchString(toWrite) {
			toWrite = "# " + toWrite
		}
		if err := writeString(writer, toWrite); err != nil {
			return err
		}
		if err := writeString(writer, "\n"); err != nil {
			return err
		}
	}
	return nil
}

func (ke *kyamlEncoder) writeInlineComment(writer io.Writer, comment string) error {
	comment = strings.TrimSpace(strings.ReplaceAll(comment, "\r\n", "\n"))
	if comment == "" {
		return nil
	}

	lines := strings.Split(comment, "\n")
	first := strings.TrimSpace(lines[0])
	if first == "" {
		return nil
	}

	if !strings.HasPrefix(first, "#") {
		first = "# " + first
	}

	if err := writeString(writer, " "); err != nil {
		return err
	}
	return writeString(writer, first)
}
