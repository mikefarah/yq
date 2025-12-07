//go:build !yq_nohcl

package yqlib

import (
	"fmt"
	"io"
	"regexp"

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
			ctyValue, err := nodeToCtyValue(valueNode)
			if err != nil {
				return err
			}
			body.SetAttributeValue(key, ctyValue)
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

		ctyValue, err := nodeToCtyValue(valueNode)
		if err != nil {
			return err
		}
		body.SetAttributeValue(key, ctyValue)
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
