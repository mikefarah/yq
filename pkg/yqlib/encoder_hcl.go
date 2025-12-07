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
			goValue, err := candidateNodeToGoValue(valueNode)
			if err != nil {
				return err
			}
			body.SetAttributeValue(key, toCtyValue(goValue))
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

		goValue, err := candidateNodeToGoValue(valueNode)
		if err != nil {
			return err
		}
		body.SetAttributeValue(key, toCtyValue(goValue))
	}
	return nil
}

// candidateNodeToGoValue converts a CandidateNode to a Go value suitable for HCL encoding
func candidateNodeToGoValue(node *CandidateNode) (interface{}, error) {
	switch node.Kind {
	case ScalarNode:
		// Parse scalar value based on its tag
		switch node.Tag {
		case "!!bool":
			return node.Value == "true", nil
		case "!!int":
			var i int64
			_, err := fmt.Sscanf(node.Value, "%d", &i)
			if err != nil {
				return nil, err
			}
			return i, nil
		case "!!float":
			var f float64
			_, err := fmt.Sscanf(node.Value, "%f", &f)
			if err != nil {
				return nil, err
			}
			return f, nil
		case "!!null":
			return nil, nil
		default:
			// Default to string
			return node.Value, nil
		}
	case MappingNode:
		m := make(map[string]interface{})
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]
			v, err := candidateNodeToGoValue(valueNode)
			if err != nil {
				return nil, err
			}
			m[keyNode.Value] = v
		}
		return m, nil
	case SequenceNode:
		arr := make([]interface{}, len(node.Content))
		for i, item := range node.Content {
			v, err := candidateNodeToGoValue(item)
			if err != nil {
				return nil, err
			}
			arr[i] = v
		}
		return arr, nil
	case AliasNode:
		return nil, fmt.Errorf("HCL encoder does not support aliases")
	default:
		return nil, fmt.Errorf("unsupported node kind: %v", node.Kind)
	}
}

// toCtyValue converts Go values to cty.Value for hclwrite
func toCtyValue(val interface{}) cty.Value {
	switch v := val.(type) {
	case string:
		return cty.StringVal(v)
	case bool:
		return cty.BoolVal(v)
	case int:
		return cty.NumberIntVal(int64(v))
	case int64:
		return cty.NumberIntVal(v)
	case float64:
		return cty.NumberFloatVal(v)
	case []interface{}:
		vals := make([]cty.Value, len(v))
		for i, item := range v {
			vals[i] = toCtyValue(item)
		}
		return cty.TupleVal(vals)
	case map[string]interface{}:
		m := make(map[string]cty.Value)
		for k, item := range v {
			m[k] = toCtyValue(item)
		}
		return cty.ObjectVal(m)
	default:
		// fallback: treat as string
		return cty.StringVal(fmt.Sprintf("%v", v))
	}
}
