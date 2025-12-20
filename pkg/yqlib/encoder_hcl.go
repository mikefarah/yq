//go:build !yq_nohcl

package yqlib

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	hclwrite "github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type hclEncoder struct {
	prefs HclPreferences
}

// commentPathSep is used to join path segments when collecting comments.
// It uses a rarely used ASCII control character to avoid collisions with
// normal key names (including dots).
const commentPathSep = "\x1e"

// NewHclEncoder creates a new HCL encoder
func NewHclEncoder(prefs HclPreferences) Encoder {
	return &hclEncoder{prefs: prefs}
}

func (he *hclEncoder) CanHandleAliases() bool {
	return false
}

func (he *hclEncoder) PrintDocumentSeparator(_ io.Writer) error {
	return nil
}

func (he *hclEncoder) PrintLeadingContent(_ io.Writer, _ string) error {
	return nil
}

func (he *hclEncoder) Encode(writer io.Writer, node *CandidateNode) error {
	log.Debugf("I need to encode %v", NodeToString(node))
	if node.Kind == ScalarNode {
		return writeString(writer, node.Value+"\n")
	}

	f := hclwrite.NewEmptyFile()
	body := f.Body()

	// Collect comments as we encode
	commentMap := make(map[string]string)
	he.collectComments(node, "", commentMap)

	if err := he.encodeNode(body, node); err != nil {
		return fmt.Errorf("failed to encode HCL: %w", err)
	}

	// Get the formatted output and remove extra spacing before '='
	output := f.Bytes()
	compactOutput := he.compactSpacing(output)

	// Inject comments back into the output
	finalOutput := he.injectComments(compactOutput, commentMap)

	if he.prefs.ColorsEnabled {
		colourized := he.colorizeHcl(finalOutput)
		_, err := writer.Write(colourized)
		return err
	}

	_, err := writer.Write(finalOutput)
	return err
}

// compactSpacing removes extra whitespace before '=' in attribute assignments
func (he *hclEncoder) compactSpacing(input []byte) []byte {
	// Use regex to replace multiple spaces before = with single space
	re := regexp.MustCompile(`(\S)\s{2,}=`)
	return re.ReplaceAll(input, []byte("$1 ="))
}

// collectComments recursively collects comments from nodes for later injection
func (he *hclEncoder) collectComments(node *CandidateNode, prefix string, commentMap map[string]string) {
	if node == nil {
		return
	}

	// For mapping nodes, collect comments from keys and values
	if node.Kind == MappingNode {
		// Collect root-level head comment if at root (prefix is empty)
		if prefix == "" && node.HeadComment != "" {
			commentMap[joinCommentPath("__root__", "head")] = node.HeadComment
		}

		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]
			key := keyNode.Value

			// Create a path for this key
			path := joinCommentPath(prefix, key)

			// Store comments from the key (head comments appear before the attribute)
			if keyNode.HeadComment != "" {
				commentMap[joinCommentPath(path, "head")] = keyNode.HeadComment
			}
			// Store comments from the value (line comments appear after the value)
			if valueNode.LineComment != "" {
				commentMap[joinCommentPath(path, "line")] = valueNode.LineComment
			}
			if valueNode.FootComment != "" {
				commentMap[joinCommentPath(path, "foot")] = valueNode.FootComment
			}

			// Recurse into nested mappings
			if valueNode.Kind == MappingNode {
				he.collectComments(valueNode, path, commentMap)
			}
		}
	}
}

// joinCommentPath concatenates path segments using commentPathSep, safely handling empty prefixes.
func joinCommentPath(prefix, segment string) string {
	if prefix == "" {
		return segment
	}
	return prefix + commentPathSep + segment
}

// injectComments adds collected comments back into the HCL output
func (he *hclEncoder) injectComments(output []byte, commentMap map[string]string) []byte {
	// Convert output to string for easier manipulation
	result := string(output)

	// Root-level head comment (stored on the synthetic __root__/head path)
	for path, comment := range commentMap {
		if path == joinCommentPath("__root__", "head") {
			trimmed := strings.TrimSpace(comment)
			if trimmed != "" && !strings.HasPrefix(result, trimmed) {
				result = trimmed + "\n" + result
			}
		}
	}

	// Attribute head comments: insert above matching assignment
	for path, comment := range commentMap {
		parts := strings.Split(path, commentPathSep)
		if len(parts) < 2 {
			continue
		}

		commentType := parts[len(parts)-1]
		key := parts[len(parts)-2]
		if commentType != "head" || key == "" {
			continue
		}

		trimmed := strings.TrimSpace(comment)
		if trimmed == "" {
			continue
		}

		re := regexp.MustCompile(`(?m)^(\s*)` + regexp.QuoteMeta(key) + `\s*=`)
		if re.MatchString(result) {
			result = re.ReplaceAllString(result, "$1"+trimmed+"\n$0")
		}
	}

	return []byte(result)
}

func (he *hclEncoder) colorizeHcl(input []byte) []byte {
	hcl := string(input)
	result := strings.Builder{}

	// Create colour functions for different token types
	commentColor := color.New(color.FgHiBlack).SprintFunc()
	stringColor := color.New(color.FgGreen).SprintFunc()
	numberColor := color.New(color.FgHiMagenta).SprintFunc()
	keyColor := color.New(color.FgCyan).SprintFunc()
	boolColor := color.New(color.FgHiMagenta).SprintFunc()

	// Simple tokenization for HCL colouring
	i := 0
	for i < len(hcl) {
		ch := hcl[i]

		// Comments - from # to end of line
		if ch == '#' {
			end := i
			for end < len(hcl) && hcl[end] != '\n' {
				end++
			}
			result.WriteString(commentColor(hcl[i:end]))
			i = end
			continue
		}

		// Strings - quoted text
		if ch == '"' || ch == '\'' {
			quote := ch
			end := i + 1
			for end < len(hcl) && hcl[end] != quote {
				if hcl[end] == '\\' {
					end++ // skip escaped char
				}
				end++
			}
			if end < len(hcl) {
				end++ // include closing quote
			}
			result.WriteString(stringColor(hcl[i:end]))
			i = end
			continue
		}

		// Numbers - sequences of digits, possibly with decimal point or minus
		if (ch >= '0' && ch <= '9') || (ch == '-' && i+1 < len(hcl) && hcl[i+1] >= '0' && hcl[i+1] <= '9') {
			end := i
			if ch == '-' {
				end++
			}
			for end < len(hcl) && ((hcl[end] >= '0' && hcl[end] <= '9') || hcl[end] == '.') {
				end++
			}
			result.WriteString(numberColor(hcl[i:end]))
			i = end
			continue
		}

		// Identifiers/keys - alphanumeric + underscore
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' {
			end := i
			for end < len(hcl) && ((hcl[end] >= 'a' && hcl[end] <= 'z') ||
				(hcl[end] >= 'A' && hcl[end] <= 'Z') ||
				(hcl[end] >= '0' && hcl[end] <= '9') ||
				hcl[end] == '_' || hcl[end] == '-') {
				end++
			}
			ident := hcl[i:end]

			// Check if this is a keyword/reserved word
			switch ident {
			case "true", "false", "null":
				result.WriteString(boolColor(ident))
			default:
				// Check if followed by = (it's a key)
				j := end
				for j < len(hcl) && (hcl[j] == ' ' || hcl[j] == '\t') {
					j++
				}
				if j < len(hcl) && hcl[j] == '=' {
					result.WriteString(keyColor(ident))
				} else if j < len(hcl) && hcl[j] == '{' {
					// Block type
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

// Helper runes for unquoted identifiers
func isHCLIdentifierStart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func isHCLIdentifierPart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-'
}

func isValidHCLIdentifier(s string) bool {
	if s == "" {
		return false
	}
	// HCL identifiers must start with a letter or underscore
	// and contain only letters, digits, underscores, and hyphens
	for i, r := range s {
		if i == 0 {
			if !isHCLIdentifierStart(r) {
				return false
			}
			continue
		}
		if !isHCLIdentifierPart(r) {
			return false
		}
	}
	return true
}

// tokensForRawHCLExpr produces a minimal token stream for a simple HCL expression so we can
// write it without introducing quotes (e.g. function calls like upper(message)).
func tokensForRawHCLExpr(expr string) (hclwrite.Tokens, error) {
	var tokens hclwrite.Tokens
	for i := 0; i < len(expr); {
		ch := expr[i]
		switch {
		case ch == ' ' || ch == '\t':
			i++
			continue
		case isHCLIdentifierStart(rune(ch)):
			start := i
			i++
			for i < len(expr) && isHCLIdentifierPart(rune(expr[i])) {
				i++
			}
			tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenIdent, Bytes: []byte(expr[start:i])})
			continue
		case ch >= '0' && ch <= '9':
			start := i
			i++
			for i < len(expr) && ((expr[i] >= '0' && expr[i] <= '9') || expr[i] == '.') {
				i++
			}
			tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenNumberLit, Bytes: []byte(expr[start:i])})
			continue
		case ch == '(':
			tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenOParen, Bytes: []byte{'('}})
		case ch == ')':
			tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenCParen, Bytes: []byte{')'}})
		case ch == ',':
			tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenComma, Bytes: []byte{','}})
		case ch == '.':
			tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenDot, Bytes: []byte{'.'}})
		case ch == '+':
			tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenPlus, Bytes: []byte{'+'}})
		case ch == '-':
			tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenMinus, Bytes: []byte{'-'}})
		case ch == '*':
			tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenStar, Bytes: []byte{'*'}})
		case ch == '/':
			tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenSlash, Bytes: []byte{'/'}})
		default:
			return nil, fmt.Errorf("unsupported character %q in raw HCL expression", ch)
		}
		i++
	}
	return tokens, nil
}

// encodeAttribute encodes a value as an HCL attribute
func (he *hclEncoder) encodeAttribute(body *hclwrite.Body, key string, valueNode *CandidateNode) error {
	if valueNode.Kind == ScalarNode && valueNode.Tag == "!!str" {
		// Handle unquoted expressions (as-is, without quotes)
		if valueNode.Style == 0 {
			tokens, err := tokensForRawHCLExpr(valueNode.Value)
			if err != nil {
				return err
			}
			body.SetAttributeRaw(key, tokens)
			return nil
		}
		if valueNode.Style&LiteralStyle != 0 {
			tokens, err := tokensForRawHCLExpr(valueNode.Value)
			if err != nil {
				return err
			}
			body.SetAttributeRaw(key, tokens)
			return nil
		}
		// Check if template with interpolation
		if valueNode.Style&DoubleQuotedStyle != 0 && strings.Contains(valueNode.Value, "${") {
			return he.encodeTemplateAttribute(body, key, valueNode.Value)
		}
		// Check if unquoted identifier
		if isValidHCLIdentifier(valueNode.Value) && valueNode.Style == 0 {
			traversal := hcl.Traversal{
				hcl.TraverseRoot{Name: valueNode.Value},
			}
			body.SetAttributeTraversal(key, traversal)
			return nil
		}
	}
	// Default: use cty.Value for quoted strings and all other types
	ctyValue, err := nodeToCtyValue(valueNode)
	if err != nil {
		return err
	}
	body.SetAttributeValue(key, ctyValue)
	return nil
}

// encodeTemplateAttribute encodes a template string with ${} interpolations
func (he *hclEncoder) encodeTemplateAttribute(body *hclwrite.Body, key string, templateStr string) error {
	tokens := hclwrite.Tokens{
		{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
	}

	for i := 0; i < len(templateStr); i++ {
		if i < len(templateStr)-1 && templateStr[i] == '$' && templateStr[i+1] == '{' {
			// Start of template interpolation
			tokens = append(tokens, &hclwrite.Token{
				Type:  hclsyntax.TokenTemplateInterp,
				Bytes: []byte("${"),
			})
			i++ // skip the '{'
			// Find the matching '}'
			start := i + 1
			depth := 1
			for i++; i < len(templateStr) && depth > 0; i++ {
				switch templateStr[i] {
				case '{':
					depth++
				case '}':
					depth--
				}
			}
			i-- // back up to the '}'
			interpExpr := templateStr[start:i]
			tokens = append(tokens, &hclwrite.Token{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte(interpExpr),
			})
			tokens = append(tokens, &hclwrite.Token{
				Type:  hclsyntax.TokenTemplateSeqEnd,
				Bytes: []byte("}"),
			})
		} else {
			// Regular character
			tokens = append(tokens, &hclwrite.Token{
				Type:  hclsyntax.TokenQuotedLit,
				Bytes: []byte{templateStr[i]},
			})
		}
	}
	tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}})
	body.SetAttributeRaw(key, tokens)
	return nil
}

// encodeBlockIfMapping attempts to encode a value as a block. Returns true if it was encoded as a block.
func (he *hclEncoder) encodeBlockIfMapping(body *hclwrite.Body, key string, valueNode *CandidateNode) bool {
	if valueNode.Kind != MappingNode || valueNode.Style == FlowStyle {
		return false
	}

	// If EncodeSeparate is set, emit children as separate blocks regardless of label extraction
	if valueNode.EncodeSeparate {
		if handled, _ := he.encodeMappingChildrenAsBlocks(body, key, valueNode); handled {
			return true
		}
	}

	// Try to extract block labels from a single-entry mapping chain
	if labels, bodyNode, ok := extractBlockLabels(valueNode); ok {
		if len(labels) > 1 && mappingChildrenAllMappings(bodyNode) {
			primaryLabels := labels[:len(labels)-1]
			nestedType := labels[len(labels)-1]
			block := body.AppendNewBlock(key, primaryLabels)
			if handled, err := he.encodeMappingChildrenAsBlocks(block.Body(), nestedType, bodyNode); err == nil && handled {
				return true
			}
			if err := he.encodeNodeAttributes(block.Body(), bodyNode); err == nil {
				return true
			}
		}
		block := body.AppendNewBlock(key, labels)
		if err := he.encodeNodeAttributes(block.Body(), bodyNode); err == nil {
			return true
		}
	}

	// If all child values are mappings, treat each child key as a labelled instance of this block type
	if handled, _ := he.encodeMappingChildrenAsBlocks(body, key, valueNode); handled {
		return true
	}

	// No labels detected, render as unlabelled block
	block := body.AppendNewBlock(key, nil)
	if err := he.encodeNodeAttributes(block.Body(), valueNode); err == nil {
		return true
	}

	return false
}

// encodeNode encodes a CandidateNode directly to HCL, preserving style information
func (he *hclEncoder) encodeNode(body *hclwrite.Body, node *CandidateNode) error {
	if node.Kind != MappingNode {
		return fmt.Errorf("HCL encoder expects a mapping at the root level, got %v", kindToString(node.Kind))
	}

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]
		key := keyNode.Value

		// Render as block or attribute depending on value type
		if he.encodeBlockIfMapping(body, key, valueNode) {
			continue
		}

		// Render as attribute: key = value
		if err := he.encodeAttribute(body, key, valueNode); err != nil {
			return err
		}
	}
	return nil
}

// mappingChildrenAllMappings reports whether all values in a mapping node are non-flow mappings.
func mappingChildrenAllMappings(node *CandidateNode) bool {
	if node == nil || node.Kind != MappingNode || node.Style == FlowStyle {
		return false
	}
	if len(node.Content) == 0 {
		return false
	}
	for i := 0; i < len(node.Content); i += 2 {
		childVal := node.Content[i+1]
		if childVal.Kind != MappingNode || childVal.Style == FlowStyle {
			return false
		}
	}
	return true
}

// encodeMappingChildrenAsBlocks emits a block for each mapping child, treating the child key as a label.
// Returns handled=true when it emitted blocks.
func (he *hclEncoder) encodeMappingChildrenAsBlocks(body *hclwrite.Body, blockType string, valueNode *CandidateNode) (bool, error) {
	if !mappingChildrenAllMappings(valueNode) {
		return false, nil
	}

	// Only emit as separate blocks if EncodeSeparate is true
	// This allows the encoder to respect the original block structure preserved by the decoder
	if !valueNode.EncodeSeparate {
		return false, nil
	}

	for i := 0; i < len(valueNode.Content); i += 2 {
		childKey := valueNode.Content[i].Value
		childVal := valueNode.Content[i+1]

		// Check if this child also represents multiple blocks (all children are mappings)
		if mappingChildrenAllMappings(childVal) {
			// Recursively emit each grandchild as a separate block with extended labels
			for j := 0; j < len(childVal.Content); j += 2 {
				grandchildKey := childVal.Content[j].Value
				grandchildVal := childVal.Content[j+1]
				labels := []string{childKey, grandchildKey}

				// Try to extract additional labels if this is a single-entry chain
				if extraLabels, bodyNode, ok := extractBlockLabels(grandchildVal); ok {
					labels = append(labels, extraLabels...)
					grandchildVal = bodyNode
				}

				block := body.AppendNewBlock(blockType, labels)
				if err := he.encodeNodeAttributes(block.Body(), grandchildVal); err != nil {
					return true, err
				}
			}
		} else {
			// Single block with this child as label(s)
			labels := []string{childKey}
			if extraLabels, bodyNode, ok := extractBlockLabels(childVal); ok {
				labels = append(labels, extraLabels...)
				childVal = bodyNode
			}
			block := body.AppendNewBlock(blockType, labels)
			if err := he.encodeNodeAttributes(block.Body(), childVal); err != nil {
				return true, err
			}
		}
	}

	return true, nil
}

// encodeNodeAttributes encodes the attributes of a mapping node (used for blocks)
func (he *hclEncoder) encodeNodeAttributes(body *hclwrite.Body, node *CandidateNode) error {
	if node.Kind != MappingNode {
		return fmt.Errorf("expected mapping node for block body")
	}

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]
		key := keyNode.Value

		// Render as block or attribute depending on value type
		if he.encodeBlockIfMapping(body, key, valueNode) {
			continue
		}

		// Render attribute for non-block value
		if err := he.encodeAttribute(body, key, valueNode); err != nil {
			return err
		}
	}
	return nil
}

// extractBlockLabels detects a chain of single-entry mappings that encode block labels.
// It returns the collected labels and the final mapping to be used as the block body.
// Pattern: {label1: {label2: { ... {bodyMap} }}}
func extractBlockLabels(node *CandidateNode) ([]string, *CandidateNode, bool) {
	var labels []string
	current := node
	for current != nil && current.Kind == MappingNode && len(current.Content) == 2 {
		keyNode := current.Content[0]
		valNode := current.Content[1]
		if valNode.Kind != MappingNode {
			break
		}
		labels = append(labels, keyNode.Value)
		// If the child is itself a single mapping entry with a mapping value, keep descending.
		if len(valNode.Content) == 2 && valNode.Content[1].Kind == MappingNode {
			current = valNode
			continue
		}
		// Otherwise, we have reached the body mapping.
		return labels, valNode, true
	}
	return nil, nil, false
}

// nodeToCtyValue converts a CandidateNode directly to cty.Value, preserving order
func nodeToCtyValue(node *CandidateNode) (cty.Value, error) {
	switch node.Kind {
	case ScalarNode:
		// Parse scalar value based on its tag
		switch node.Tag {
		case "!!bool":
			return cty.BoolVal(node.Value == "true"), nil
		case "!!int":
			var i int64
			_, err := fmt.Sscanf(node.Value, "%d", &i)
			if err != nil {
				return cty.NilVal, err
			}
			return cty.NumberIntVal(i), nil
		case "!!float":
			var f float64
			_, err := fmt.Sscanf(node.Value, "%f", &f)
			if err != nil {
				return cty.NilVal, err
			}
			return cty.NumberFloatVal(f), nil
		case "!!null":
			return cty.NullVal(cty.DynamicPseudoType), nil
		default:
			// Default to string
			return cty.StringVal(node.Value), nil
		}
	case MappingNode:
		// Preserve order by iterating Content directly
		m := make(map[string]cty.Value)
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]
			v, err := nodeToCtyValue(valueNode)
			if err != nil {
				return cty.NilVal, err
			}
			m[keyNode.Value] = v
		}
		return cty.ObjectVal(m), nil
	case SequenceNode:
		vals := make([]cty.Value, len(node.Content))
		for i, item := range node.Content {
			v, err := nodeToCtyValue(item)
			if err != nil {
				return cty.NilVal, err
			}
			vals[i] = v
		}
		return cty.TupleVal(vals), nil
	case AliasNode:
		return cty.NilVal, fmt.Errorf("HCL encoder does not support aliases")
	default:
		return cty.NilVal, fmt.Errorf("unsupported node kind: %v", node.Kind)
	}
}
