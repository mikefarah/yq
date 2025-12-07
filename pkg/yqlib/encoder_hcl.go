//go:build !yq_nohcl

package yqlib

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	hclwrite "github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type hclEncoder struct {
}

// NewHclEncoder creates a new HCL encoder
func NewHclEncoder() Encoder {
	return &hclEncoder{}
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

	f := hclwrite.NewEmptyFile()
	body := f.Body()
	if err := he.encodeNode(body, node); err != nil {
		return fmt.Errorf("failed to encode HCL: %w", err)
	}

	// Get the formatted output and remove extra spacing before '='
	output := f.Bytes()
	compactOutput := he.compactSpacing(output)

	_, err := writer.Write(compactOutput)
	return err
}

// compactSpacing removes extra whitespace before '=' in attribute assignments
func (he *hclEncoder) compactSpacing(input []byte) []byte {
	// Use regex to replace multiple spaces before = with single space
	re := regexp.MustCompile(`(\S)\s{2,}=`)
	return re.ReplaceAll(input, []byte("$1 ="))
}

// Helper runes for unquoted identifiers
func isHCLIdentifierStart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func isHCLIdentifierPart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-'
}

// isValidHCLIdentifier checks if a string is a valid HCL identifier (unquoted)
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

// encodeNode encodes a CandidateNode directly to HCL, preserving style information
func (he *hclEncoder) encodeNode(body *hclwrite.Body, node *CandidateNode) error {
	if node.Kind != MappingNode {
		return fmt.Errorf("HCL encoder expects a mapping at the root level")
	}

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]
		key := keyNode.Value

		// Check if value is a mapping without FlowStyle -> render as block
		if valueNode.Kind == MappingNode && valueNode.Style != FlowStyle {
			// Render as block: key { ... }
			block := body.AppendNewBlock(key, nil)
			if err := he.encodeNodeAttributes(block.Body(), valueNode); err != nil {
				return err
			}
		} else {
			// Render as attribute: key = value
			// Check the style to determine how to encode strings
			if valueNode.Kind == ScalarNode && valueNode.Tag == "!!str" {
				// Check style: DoubleQuotedStyle means template, no style could be unquoted or regular
				// To distinguish unquoted from regular, we check if the value is a valid identifier
				if valueNode.Style&DoubleQuotedStyle != 0 && strings.Contains(valueNode.Value, "${") {
					// Template string - use raw tokens to preserve ${} syntax
					tokens := hclwrite.Tokens{
						{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
					}
					// Parse the string and add tokens
					for i := 0; i < len(valueNode.Value); i++ {
						if i < len(valueNode.Value)-1 && valueNode.Value[i] == '$' && valueNode.Value[i+1] == '{' {
							// Start of template interpolation
							tokens = append(tokens, &hclwrite.Token{
								Type:  hclsyntax.TokenTemplateInterp,
								Bytes: []byte("${"),
							})
							i++ // skip the '{'
							// Find the matching '}'
							start := i + 1
							depth := 1
							for i++; i < len(valueNode.Value) && depth > 0; i++ {
								switch valueNode.Value[i] {
								case '{':
									depth++
								case '}':
									depth--
								}
							}
							i-- // back up to the '}'
							interpExpr := valueNode.Value[start:i]
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
								Bytes: []byte{valueNode.Value[i]},
							})
						}
					}
					tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}})
					body.SetAttributeRaw(key, tokens)
				} else if isValidHCLIdentifier(valueNode.Value) && valueNode.Style == 0 {
					// Could be unquoted identifier - but only if it came from HCL originally
					// For safety, only use traversal if style is explicitly 0 (not set)
					// This avoids treating strings from YAML as unquoted
					traversal := hcl.Traversal{
						hcl.TraverseRoot{Name: valueNode.Value},
					}
					body.SetAttributeTraversal(key, traversal)
				} else {
					// Regular quoted string - use cty.Value
					ctyValue, err := nodeToCtyValue(valueNode)
					if err != nil {
						return err
					}
					body.SetAttributeValue(key, ctyValue)
				}
			} else {
				// Non-string value - use cty.Value
				ctyValue, err := nodeToCtyValue(valueNode)
				if err != nil {
					return err
				}
				body.SetAttributeValue(key, ctyValue)
			}
		}
	}
	return nil
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

		// Check if this is an unquoted identifier (no DoubleQuotedStyle)
		if valueNode.Kind == ScalarNode && valueNode.Tag == "!!str" && valueNode.Style&DoubleQuotedStyle == 0 {
			// Unquoted identifier - use traversal
			traversal := hcl.Traversal{
				hcl.TraverseRoot{Name: valueNode.Value},
			}
			body.SetAttributeTraversal(key, traversal)
		} else {
			// Quoted value or non-string - use cty.Value
			ctyValue, err := nodeToCtyValue(valueNode)
			if err != nil {
				return err
			}
			body.SetAttributeValue(key, ctyValue)
		}
	}
	return nil
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
