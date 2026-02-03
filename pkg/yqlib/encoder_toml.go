//go:build !yq_notoml

package yqlib

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
)

type tomlEncoder struct {
	wroteRootAttr bool // Track if we wrote root-level attributes before tables
	prefs         TomlPreferences
}

func NewTomlEncoder() Encoder {
	return NewTomlEncoderWithPrefs(ConfiguredTomlPreferences)
}

func NewTomlEncoderWithPrefs(prefs TomlPreferences) Encoder {
	return &tomlEncoder{prefs: prefs}
}

func (te *tomlEncoder) Encode(writer io.Writer, node *CandidateNode) error {
	if node.Kind != MappingNode {
		// For standalone selections, TOML tests expect raw value for scalars
		if node.Kind == ScalarNode {
			return writeString(writer, node.Value+"\n")
		}
		return fmt.Errorf("TOML encoder expects a mapping at the root level")
	}

	// Encode to a buffer first if colors are enabled
	var buf bytes.Buffer
	var targetWriter io.Writer
	targetWriter = writer
	if te.prefs.ColorsEnabled {
		targetWriter = &buf
	}

	// Encode a root mapping as a sequence of attributes, tables, and arrays of tables
	if err := te.encodeRootMapping(targetWriter, node); err != nil {
		return err
	}

	if te.prefs.ColorsEnabled {
		colourised := te.colorizeToml(buf.Bytes())
		_, err := writer.Write(colourised)
		return err
	}

	return nil
}

func (te *tomlEncoder) PrintDocumentSeparator(_ io.Writer) error {
	return nil
}

func (te *tomlEncoder) PrintLeadingContent(_ io.Writer, _ string) error {
	return nil
}

func (te *tomlEncoder) CanHandleAliases() bool {
	return false
}

// ---- helpers ----

func (te *tomlEncoder) writeComment(w io.Writer, comment string) error {
	if comment == "" {
		return nil
	}
	lines := strings.Split(comment, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "#") {
			line = "# " + line
		}
		if _, err := w.Write([]byte(line + "\n")); err != nil {
			return err
		}
	}
	return nil
}

func (te *tomlEncoder) formatScalar(node *CandidateNode) string {
	switch node.Tag {
	case "!!str":
		// Quote strings per TOML spec
		return fmt.Sprintf("%q", node.Value)
	case "!!bool", "!!int", "!!float":
		return node.Value
	case "!!null":
		// TOML does not have null; encode as empty string
		return `""`
	default:
		return node.Value
	}
}

func (te *tomlEncoder) encodeRootMapping(w io.Writer, node *CandidateNode) error {
	te.wroteRootAttr = false // Reset state

	// Write root head comment if present (at the very beginning, no leading blank line)
	if node.HeadComment != "" {
		if err := te.writeComment(w, node.HeadComment); err != nil {
			return err
		}
	}

	// Preserve existing order by iterating Content
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valNode := node.Content[i+1]
		if err := te.encodeTopLevelEntry(w, []string{keyNode.Value}, valNode); err != nil {
			return err
		}
	}
	return nil
}

// encodeTopLevelEntry encodes a key/value at the root, dispatching to attribute, table, or array-of-tables
func (te *tomlEncoder) encodeTopLevelEntry(w io.Writer, path []string, node *CandidateNode) error {
	if len(path) == 0 {
		return fmt.Errorf("cannot encode TOML entry with empty path")
	}

	switch node.Kind {
	case ScalarNode:
		// key = value
		return te.writeAttribute(w, path[len(path)-1], node)
	case SequenceNode:
		// Empty arrays should be encoded as [] attributes
		if len(node.Content) == 0 {
			return te.writeArrayAttribute(w, path[len(path)-1], node)
		}

		// If all items are mappings => array of tables; else => array attribute
		allMaps := true
		for _, it := range node.Content {
			if it.Kind != MappingNode {
				allMaps = false
				break
			}
		}
		if allMaps {
			key := path[len(path)-1]
			for _, it := range node.Content {
				// [[key]] then body
				if _, err := w.Write([]byte("[[" + key + "]]\n")); err != nil {
					return err
				}
				if err := te.encodeMappingBodyWithPath(w, []string{key}, it); err != nil {
					return err
				}
			}
			return nil
		}
		// Regular array attribute
		return te.writeArrayAttribute(w, path[len(path)-1], node)
	case MappingNode:
		// Inline table if not EncodeSeparate, else emit separate tables/arrays of tables for children under this path
		if !node.EncodeSeparate {
			// If children contain mappings or arrays of mappings, prefer separate sections
			if te.hasEncodeSeparateChild(node) || te.hasStructuralChildren(node) {
				return te.encodeSeparateMapping(w, path, node)
			}
			return te.writeInlineTableAttribute(w, path[len(path)-1], node)
		}
		return te.encodeSeparateMapping(w, path, node)
	default:
		return fmt.Errorf("unsupported node kind for TOML: %v", node.Kind)
	}
}

func (te *tomlEncoder) writeAttribute(w io.Writer, key string, value *CandidateNode) error {
	te.wroteRootAttr = true // Mark that we wrote a root attribute

	// Write head comment before the attribute
	if err := te.writeComment(w, value.HeadComment); err != nil {
		return err
	}

	// Write the attribute
	line := key + " = " + te.formatScalar(value)

	// Add line comment if present
	if value.LineComment != "" {
		lineComment := strings.TrimSpace(value.LineComment)
		if !strings.HasPrefix(lineComment, "#") {
			lineComment = "# " + lineComment
		}
		line += "  " + lineComment
	}

	_, err := w.Write([]byte(line + "\n"))
	return err
}

func (te *tomlEncoder) writeArrayAttribute(w io.Writer, key string, seq *CandidateNode) error {
	te.wroteRootAttr = true // Mark that we wrote a root attribute

	// Write head comment before the array
	if err := te.writeComment(w, seq.HeadComment); err != nil {
		return err
	}

	// Handle empty arrays
	if len(seq.Content) == 0 {
		line := key + " = []"
		if seq.LineComment != "" {
			lineComment := strings.TrimSpace(seq.LineComment)
			if !strings.HasPrefix(lineComment, "#") {
				lineComment = "# " + lineComment
			}
			line += "  " + lineComment
		}
		_, err := w.Write([]byte(line + "\n"))
		return err
	}

	// Check if any array elements have head comments - if so, use multiline format
	hasElementComments := false
	for _, it := range seq.Content {
		if it.HeadComment != "" {
			hasElementComments = true
			break
		}
	}

	if hasElementComments {
		// Write multiline array format with comments
		if _, err := w.Write([]byte(key + " = [\n")); err != nil {
			return err
		}

		for i, it := range seq.Content {
			// Write head comment for this element
			if it.HeadComment != "" {
				commentLines := strings.Split(it.HeadComment, "\n")
				for _, commentLine := range commentLines {
					if strings.TrimSpace(commentLine) != "" {
						if !strings.HasPrefix(strings.TrimSpace(commentLine), "#") {
							commentLine = "# " + commentLine
						}
						if _, err := w.Write([]byte("  " + commentLine + "\n")); err != nil {
							return err
						}
					}
				}
			}

			// Write the element value
			var itemStr string
			switch it.Kind {
			case ScalarNode:
				itemStr = te.formatScalar(it)
			case SequenceNode:
				nested, err := te.sequenceToInlineArray(it)
				if err != nil {
					return err
				}
				itemStr = nested
			case MappingNode:
				inline, err := te.mappingToInlineTable(it)
				if err != nil {
					return err
				}
				itemStr = inline
			case AliasNode:
				return fmt.Errorf("aliases are not supported in TOML")
			default:
				return fmt.Errorf("unsupported array item kind: %v", it.Kind)
			}

			// Always add trailing comma in multiline arrays
			itemStr += ","

			if _, err := w.Write([]byte("  " + itemStr + "\n")); err != nil {
				return err
			}

			// Add blank line between elements (except after the last one)
			if i < len(seq.Content)-1 {
				if _, err := w.Write([]byte("\n")); err != nil {
					return err
				}
			}
		}

		if _, err := w.Write([]byte("]\n")); err != nil {
			return err
		}
		return nil
	}

	// Join scalars or nested arrays recursively into TOML array syntax
	items := make([]string, 0, len(seq.Content))
	for _, it := range seq.Content {
		switch it.Kind {
		case ScalarNode:
			items = append(items, te.formatScalar(it))
		case SequenceNode:
			// Nested arrays: encode inline
			nested, err := te.sequenceToInlineArray(it)
			if err != nil {
				return err
			}
			items = append(items, nested)
		case MappingNode:
			// Inline table inside array
			inline, err := te.mappingToInlineTable(it)
			if err != nil {
				return err
			}
			items = append(items, inline)
		case AliasNode:
			return fmt.Errorf("aliases are not supported in TOML")
		default:
			return fmt.Errorf("unsupported array item kind: %v", it.Kind)
		}
	}

	line := key + " = [" + strings.Join(items, ", ") + "]"

	// Add line comment if present
	if seq.LineComment != "" {
		lineComment := strings.TrimSpace(seq.LineComment)
		if !strings.HasPrefix(lineComment, "#") {
			lineComment = "# " + lineComment
		}
		line += "  " + lineComment
	}

	_, err := w.Write([]byte(line + "\n"))
	return err
}

func (te *tomlEncoder) sequenceToInlineArray(seq *CandidateNode) (string, error) {
	items := make([]string, 0, len(seq.Content))
	for _, it := range seq.Content {
		switch it.Kind {
		case ScalarNode:
			items = append(items, te.formatScalar(it))
		case SequenceNode:
			nested, err := te.sequenceToInlineArray(it)
			if err != nil {
				return "", err
			}
			items = append(items, nested)
		case MappingNode:
			inline, err := te.mappingToInlineTable(it)
			if err != nil {
				return "", err
			}
			items = append(items, inline)
		default:
			return "", fmt.Errorf("unsupported array item kind: %v", it.Kind)
		}
	}
	return "[" + strings.Join(items, ", ") + "]", nil
}

func (te *tomlEncoder) mappingToInlineTable(m *CandidateNode) (string, error) {
	// key = { a = 1, b = "x" }
	parts := make([]string, 0, len(m.Content)/2)
	for i := 0; i < len(m.Content); i += 2 {
		k := m.Content[i].Value
		v := m.Content[i+1]
		switch v.Kind {
		case ScalarNode:
			parts = append(parts, fmt.Sprintf("%s = %s", k, te.formatScalar(v)))
		case SequenceNode:
			// inline array in inline table
			arr, err := te.sequenceToInlineArray(v)
			if err != nil {
				return "", err
			}
			parts = append(parts, fmt.Sprintf("%s = %s", k, arr))
		case MappingNode:
			// nested inline table
			inline, err := te.mappingToInlineTable(v)
			if err != nil {
				return "", err
			}
			parts = append(parts, fmt.Sprintf("%s = %s", k, inline))
		default:
			return "", fmt.Errorf("unsupported inline table value kind: %v", v.Kind)
		}
	}
	return "{ " + strings.Join(parts, ", ") + " }", nil
}

func (te *tomlEncoder) writeInlineTableAttribute(w io.Writer, key string, m *CandidateNode) error {
	inline, err := te.mappingToInlineTable(m)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(key + " = " + inline + "\n"))
	return err
}

func (te *tomlEncoder) writeTableHeader(w io.Writer, path []string, m *CandidateNode) error {
	// Add blank line before table header (or before comment if present) if we wrote root attributes
	needsBlankLine := te.wroteRootAttr
	if needsBlankLine {
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
		te.wroteRootAttr = false // Only add once
	}

	// Write head comment before the table header
	if m.HeadComment != "" {
		if err := te.writeComment(w, m.HeadComment); err != nil {
			return err
		}
	}

	// Write table header [a.b.c]
	header := "[" + strings.Join(path, ".") + "]\n"
	_, err := w.Write([]byte(header))
	return err
}

// encodeSeparateMapping handles a mapping that should be encoded as table sections.
// It emits the table header for this mapping if it has any content, then processes children.
func (te *tomlEncoder) encodeSeparateMapping(w io.Writer, path []string, m *CandidateNode) error {
	// Check if this mapping has any non-mapping, non-array-of-tables children (i.e., attributes)
	hasAttrs := false
	for i := 0; i < len(m.Content); i += 2 {
		v := m.Content[i+1]
		if v.Kind == ScalarNode {
			hasAttrs = true
			break
		}
		if v.Kind == SequenceNode {
			// Check if it's NOT an array of tables
			allMaps := true
			for _, it := range v.Content {
				if it.Kind != MappingNode {
					allMaps = false
					break
				}
			}
			if !allMaps {
				hasAttrs = true
				break
			}
		}
	}

	// If there are attributes or if the mapping is empty, emit the table header
	if hasAttrs || len(m.Content) == 0 {
		if err := te.writeTableHeader(w, path, m); err != nil {
			return err
		}
		if err := te.encodeMappingBodyWithPath(w, path, m); err != nil {
			return err
		}
		return nil
	}

	// No attributes, just nested structures - process children
	for i := 0; i < len(m.Content); i += 2 {
		k := m.Content[i].Value
		v := m.Content[i+1]
		switch v.Kind {
		case MappingNode:
			// Emit [path.k]
			newPath := append(append([]string{}, path...), k)
			if err := te.writeTableHeader(w, newPath, v); err != nil {
				return err
			}
			if err := te.encodeMappingBodyWithPath(w, newPath, v); err != nil {
				return err
			}
		case SequenceNode:
			// If sequence of maps, emit [[path.k]] per element
			allMaps := true
			for _, it := range v.Content {
				if it.Kind != MappingNode {
					allMaps = false
					break
				}
			}
			if allMaps {
				key := strings.Join(append(append([]string{}, path...), k), ".")
				for _, it := range v.Content {
					if _, err := w.Write([]byte("[[" + key + "]]\n")); err != nil {
						return err
					}
					if err := te.encodeMappingBodyWithPath(w, append(append([]string{}, path...), k), it); err != nil {
						return err
					}
				}
			} else {
				// Regular array attribute under the current table path
				if err := te.writeArrayAttribute(w, k, v); err != nil {
					return err
				}
			}
		case ScalarNode:
			// Attributes directly under the current table path
			if err := te.writeAttribute(w, k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (te *tomlEncoder) hasEncodeSeparateChild(m *CandidateNode) bool {
	for i := 0; i < len(m.Content); i += 2 {
		v := m.Content[i+1]
		if v.Kind == MappingNode && v.EncodeSeparate {
			return true
		}
	}
	return false
}

func (te *tomlEncoder) hasStructuralChildren(m *CandidateNode) bool {
	for i := 0; i < len(m.Content); i += 2 {
		v := m.Content[i+1]
		// Only consider it structural if mapping has EncodeSeparate or is non-empty
		if v.Kind == MappingNode && v.EncodeSeparate {
			return true
		}
		if v.Kind == SequenceNode {
			allMaps := true
			for _, it := range v.Content {
				if it.Kind != MappingNode {
					allMaps = false
					break
				}
			}
			if allMaps {
				return true
			}
		}
	}
	return false
}

// encodeMappingBodyWithPath encodes attributes and nested arrays of tables using full dotted path context
func (te *tomlEncoder) encodeMappingBodyWithPath(w io.Writer, path []string, m *CandidateNode) error {
	// First, attributes (scalars and non-map arrays)
	for i := 0; i < len(m.Content); i += 2 {
		k := m.Content[i].Value
		v := m.Content[i+1]
		switch v.Kind {
		case ScalarNode:
			if err := te.writeAttribute(w, k, v); err != nil {
				return err
			}
		case SequenceNode:
			allMaps := true
			for _, it := range v.Content {
				if it.Kind != MappingNode {
					allMaps = false
					break
				}
			}
			if !allMaps {
				if err := te.writeArrayAttribute(w, k, v); err != nil {
					return err
				}
			}
		}
	}

	// Then, nested arrays of tables with full path
	for i := 0; i < len(m.Content); i += 2 {
		k := m.Content[i].Value
		v := m.Content[i+1]
		if v.Kind == SequenceNode {
			allMaps := true
			for _, it := range v.Content {
				if it.Kind != MappingNode {
					allMaps = false
					break
				}
			}
			if allMaps {
				dotted := strings.Join(append(append([]string{}, path...), k), ".")
				for _, it := range v.Content {
					if _, err := w.Write([]byte("[[" + dotted + "]]\n")); err != nil {
						return err
					}
					if err := te.encodeMappingBodyWithPath(w, append(append([]string{}, path...), k), it); err != nil {
						return err
					}
				}
			}
		}
	}

	// Finally, child mappings that are not marked EncodeSeparate get inlined as attributes
	for i := 0; i < len(m.Content); i += 2 {
		k := m.Content[i].Value
		v := m.Content[i+1]
		if v.Kind == MappingNode && !v.EncodeSeparate {
			if err := te.writeInlineTableAttribute(w, k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// colorizeToml applies syntax highlighting to TOML output using fatih/color
func (te *tomlEncoder) colorizeToml(input []byte) []byte {
	toml := string(input)
	result := strings.Builder{}

	// Force color output (don't check for TTY)
	color.NoColor = false

	// Create color functions for different token types
	// Use EnableColor() to ensure colors work even when NO_COLOR env is set
	commentColorObj := color.New(color.FgHiBlack)
	commentColorObj.EnableColor()
	stringColorObj := color.New(color.FgGreen)
	stringColorObj.EnableColor()
	numberColorObj := color.New(color.FgHiMagenta)
	numberColorObj.EnableColor()
	keyColorObj := color.New(color.FgCyan)
	keyColorObj.EnableColor()
	boolColorObj := color.New(color.FgHiMagenta)
	boolColorObj.EnableColor()
	sectionColorObj := color.New(color.FgYellow, color.Bold)
	sectionColorObj.EnableColor()

	commentColor := commentColorObj.SprintFunc()
	stringColor := stringColorObj.SprintFunc()
	numberColor := numberColorObj.SprintFunc()
	keyColor := keyColorObj.SprintFunc()
	boolColor := boolColorObj.SprintFunc()
	sectionColor := sectionColorObj.SprintFunc()

	// Simple tokenization for TOML colouring
	i := 0
	for i < len(toml) {
		ch := toml[i]

		// Comments - from # to end of line
		if ch == '#' {
			end := i
			for end < len(toml) && toml[end] != '\n' {
				end++
			}
			result.WriteString(commentColor(toml[i:end]))
			i = end
			continue
		}

		// Table sections - [section] or [[array]]
		// Only treat '[' as a table section if it appears at the start of the line
		// (possibly after whitespace). This avoids mis-colouring inline arrays like
		// "ports = [8000, 8001]" as table sections.
		if ch == '[' {
			isSectionHeader := true
			if i > 0 {
				isSectionHeader = false
				j := i - 1
				for j >= 0 && toml[j] != '\n' {
					if toml[j] != ' ' && toml[j] != '\t' && toml[j] != '\r' {
						// Found a non-whitespace character before this '[' on the same line,
						// so this is not a table header.
						break
					}
					j--
				}
				if j < 0 || toml[j] == '\n' {
					// Reached the start of the string or a newline without encountering
					// any non-whitespace, so '[' is at the logical start of the line.
					isSectionHeader = true
				}
			}
			if isSectionHeader {
				end := i + 1
				// Check for [[
				if end < len(toml) && toml[end] == '[' {
					end++
				}
				// Find closing ]
				for end < len(toml) && toml[end] != ']' {
					end++
				}
				// Include closing ]
				if end < len(toml) {
					end++
					// Check for ]]
					if end < len(toml) && toml[end] == ']' {
						end++
					}
				}
				result.WriteString(sectionColor(toml[i:end]))
				i = end
				continue
			}
		}

		// Strings - quoted text (double or single quotes)
		if ch == '"' || ch == '\'' {
			quote := ch
			end := i + 1
			for end < len(toml) {
				if toml[end] == quote {
					break
				}
				if toml[end] == '\\' && end+1 < len(toml) {
					// Skip the backslash and the escaped character
					end += 2
					continue
				}
				end++
			}
			if end < len(toml) {
				end++ // include closing quote
			}
			result.WriteString(stringColor(toml[i:end]))
			i = end
			continue
		}

		// Numbers - sequences of digits, possibly with decimal point or minus
		if (ch >= '0' && ch <= '9') || (ch == '-' && i+1 < len(toml) && toml[i+1] >= '0' && toml[i+1] <= '9') {
			end := i
			if ch == '-' {
				end++
			}
			for end < len(toml) {
				c := toml[end]
				if (c >= '0' && c <= '9') || c == '.' || c == 'e' || c == 'E' {
					end++
				} else if (c == '+' || c == '-') && end > 0 && (toml[end-1] == 'e' || toml[end-1] == 'E') {
					// Only allow + or - immediately after 'e' or 'E' for scientific notation
					end++
				} else {
					break
				}
			}
			result.WriteString(numberColor(toml[i:end]))
			i = end
			continue
		}

		// Identifiers/keys - alphanumeric + underscore + dash
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' {
			end := i
			for end < len(toml) && ((toml[end] >= 'a' && toml[end] <= 'z') ||
				(toml[end] >= 'A' && toml[end] <= 'Z') ||
				(toml[end] >= '0' && toml[end] <= '9') ||
				toml[end] == '_' || toml[end] == '-') {
				end++
			}
			ident := toml[i:end]

			// Check if this is a boolean/null keyword
			switch ident {
			case "true", "false":
				result.WriteString(boolColor(ident))
			default:
				// Check if followed by = or whitespace then = (it's a key)
				j := end
				for j < len(toml) && (toml[j] == ' ' || toml[j] == '\t') {
					j++
				}
				if j < len(toml) && toml[j] == '=' {
					result.WriteString(keyColor(ident))
				} else {
					result.WriteString(ident) // plain text for other identifiers
				}
			}
			i = end
			continue
		}

		// Everything else (whitespace, operators, brackets) - no color
		result.WriteByte(ch)
		i++
	}

	return []byte(result.String())
}
