package yqlib

import (
	"fmt"
	"io"
	"strconv"
)

type hclEncoder struct {
	indentString string
}

// NewHclEncoder creates a new HCL encoder
func NewHclEncoder() Encoder {
	return &hclEncoder{
		indentString: "  ", // 2 spaces for HCL indentation
	}
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

	return he.encodeNodeInContext(writer, node, "", false)
}

func (he *hclEncoder) encodeNodeInContext(writer io.Writer, node *CandidateNode, indent string, isInAttribute bool) error {
	switch node.Kind {
	case ScalarNode:
		return writeString(writer, he.formatScalarValue(node.Value))
	case MappingNode:
		return he.encodeMappingInContext(writer, node, indent, isInAttribute)
	case SequenceNode:
		return he.encodeSequence(writer, node, indent)
	case AliasNode:
		return fmt.Errorf("HCL encoder does not support aliases")
	default:
		return fmt.Errorf("unsupported node kind: %v", node.Kind)
	}
}

func (he *hclEncoder) encodeMappingInContext(writer io.Writer, node *CandidateNode, indent string, isInAttribute bool) error {
	if len(node.Content) == 0 {
		return writeString(writer, "{}")
	}

	// If this mapping is an attribute value or flow-styled, render as inline object: { a = 1, b = "two" }
	if isInAttribute || node.Style == FlowStyle {
		return he.encodeInlineMapping(writer, node, indent)
	} // If we're at the top level (indent == "") AND all values are scalars OR mappings,
	// render as attributes (key = value) or blocks (key { ... })
	if indent == "" {
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]
			key := keyNode.Value

			// Block-style for nested mappings (unless they're inline objects), attribute-style for scalars/sequences
			if valueNode.Kind == MappingNode && valueNode.Style != FlowStyle {
				// Block: key { ... }
				if err := writeString(writer, key); err != nil {
					return err
				}
				if err := writeString(writer, " {\n"); err != nil {
					return err
				}

				nextIndent := he.indentString
				for j := 0; j < len(valueNode.Content); j += 2 {
					nestedKeyNode := valueNode.Content[j]
					nestedValueNode := valueNode.Content[j+1]
					nestedKey := nestedKeyNode.Value

					if err := writeString(writer, nextIndent); err != nil {
						return err
					}
					if err := writeString(writer, nestedKey); err != nil {
						return err
					}
					if err := writeString(writer, " = "); err != nil {
						return err
					}
					if err := he.encodeNodeInContext(writer, nestedValueNode, nextIndent, true); err != nil {
						return err
					}
					if err := writeString(writer, "\n"); err != nil {
						return err
					}
				}

				if err := writeString(writer, "}\n"); err != nil {
					return err
				}
			} else {
				// Attribute: key = value
				if err := writeString(writer, key); err != nil {
					return err
				}
				if err := writeString(writer, " = "); err != nil {
					return err
				}
				if err := he.encodeNodeInContext(writer, valueNode, "", true); err != nil {
					return err
				}
				if err := writeString(writer, "\n"); err != nil {
					return err
				}
			}
		}
		return nil
	}

	// Otherwise, this shouldn't happen at nested levels in top-level syntax
	return writeString(writer, "{}")
}

func (he *hclEncoder) encodeInlineMapping(writer io.Writer, node *CandidateNode, indent string) error {
	if len(node.Content) == 0 {
		return writeString(writer, "{}")
	}

	if err := writeString(writer, "{\n"); err != nil {
		return err
	}

	nextIndent := indent + he.indentString
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]
		key := keyNode.Value

		if err := writeString(writer, nextIndent); err != nil {
			return err
		}
		if err := writeString(writer, key); err != nil {
			return err
		}
		if err := writeString(writer, " = "); err != nil {
			return err
		}
		if err := he.encodeNodeInContext(writer, valueNode, nextIndent, true); err != nil {
			return err
		}
		if err := writeString(writer, "\n"); err != nil {
			return err
		}
	}

	if err := writeString(writer, indent+"}"); err != nil {
		return err
	}

	return nil
}

func (he *hclEncoder) encodeSequence(writer io.Writer, node *CandidateNode, indent string) error {
	if len(node.Content) == 0 {
		return writeString(writer, "[]")
	}

	// Check if we should use inline format (simple values only)
	useInline := true
	for _, item := range node.Content {
		if item.Kind != ScalarNode {
			useInline = false
			break
		}
	}

	if useInline {
		// Inline format: ["a", "b", "c"]
		if err := writeString(writer, "["); err != nil {
			return err
		}
		for i, item := range node.Content {
			if i > 0 {
				if err := writeString(writer, ", "); err != nil {
					return err
				}
			}
			if err := he.encodeNodeInContext(writer, item, indent, true); err != nil {
				return err
			}
		}
		if err := writeString(writer, "]"); err != nil {
			return err
		}
		return nil
	}

	// Multi-line format for complex items
	if err := writeString(writer, "[\n"); err != nil {
		return err
	}

	nextIndent := indent + he.indentString
	for _, item := range node.Content {
		if err := writeString(writer, nextIndent); err != nil {
			return err
		}
		if err := he.encodeNodeInContext(writer, item, nextIndent, true); err != nil {
			return err
		}
		if err := writeString(writer, ",\n"); err != nil {
			return err
		}
	}

	if err := writeString(writer, indent+"]"); err != nil {
		return err
	}

	return nil
}

func (he *hclEncoder) formatScalarValue(value string) string {
	// Check if value is a boolean
	if value == "true" || value == "false" {
		return value
	}

	// Check if value is a number
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return value
	}

	// Treat as string, quote it
	return strconv.Quote(value)
}
