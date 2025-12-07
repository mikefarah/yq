//go:build !yq_nohcl

package yqlib

import (
	"fmt"
	"io"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

type hclDecoder struct {
	file          *hcl.File
	fileBytes     []byte
	readAnything  bool
	documentIndex uint
}

func NewHclDecoder() Decoder {
	return &hclDecoder{}
}

// sortedAttributes returns attributes in declaration order by source position
func sortedAttributes(attrs hclsyntax.Attributes) []*attributeWithName {
	var sorted []*attributeWithName
	for name, attr := range attrs {
		sorted = append(sorted, &attributeWithName{Name: name, Attr: attr})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Attr.Range().Start.Byte < sorted[j].Attr.Range().Start.Byte
	})
	return sorted
}

type attributeWithName struct {
	Name string
	Attr *hclsyntax.Attribute
}

func (dec *hclDecoder) Init(reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	file, diags := hclsyntax.ParseConfig(data, "input.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags != nil && diags.HasErrors() {
		return fmt.Errorf("hcl parse error: %w", diags)
	}
	dec.file = file
	dec.fileBytes = data
	dec.readAnything = false
	dec.documentIndex = 0
	return nil
}

func (dec *hclDecoder) Decode() (*CandidateNode, error) {
	if dec.readAnything {
		return nil, io.EOF
	}
	dec.readAnything = true

	if dec.file == nil {
		return nil, fmt.Errorf("no hcl file parsed")
	}

	root := &CandidateNode{Kind: MappingNode}

	// process attributes in declaration order
	body := dec.file.Body.(*hclsyntax.Body)
	for _, attrWithName := range sortedAttributes(body.Attributes) {
		keyNode := createStringScalarNode(attrWithName.Name)
		valNode := convertHclExprToNode(attrWithName.Attr.Expr, dec.fileBytes)
		root.AddKeyValueChild(keyNode, valNode)
	}

	// process blocks
	for _, block := range body.Blocks {
		// build a key from type and labels to preserve identity
		keyName := block.Type
		if len(block.Labels) > 0 {
			keyName = keyName + " " + strings.Join(block.Labels, " ")
		}
		keyNode := createStringScalarNode(keyName)
		valueNode := hclBodyToNode(block.Body, dec.fileBytes)
		root.AddKeyValueChild(keyNode, valueNode)
	}

	dec.documentIndex++
	root.document = dec.documentIndex - 1
	return root, nil
}

func hclBodyToNode(body *hclsyntax.Body, src []byte) *CandidateNode {
	node := &CandidateNode{Kind: MappingNode}
	for _, attrWithName := range sortedAttributes(body.Attributes) {
		key := createStringScalarNode(attrWithName.Name)
		val := convertHclExprToNode(attrWithName.Attr.Expr, src)
		node.AddKeyValueChild(key, val)
	}
	for _, block := range body.Blocks {
		keyName := block.Type
		if len(block.Labels) > 0 {
			keyName = keyName + " " + strings.Join(block.Labels, " ")
		}
		key := createStringScalarNode(keyName)
		val := hclBodyToNode(block.Body, src)
		node.AddKeyValueChild(key, val)
	}
	return node
}

func convertHclExprToNode(expr hclsyntax.Expression, src []byte) *CandidateNode {
	// handle literal values directly
	switch e := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		v := e.Val
		if v.IsNull() {
			return createScalarNode(nil, "")
		}
		switch {
		case v.Type().Equals(cty.String):
			// prefer to extract exact source (to avoid extra quoting) when available
			// Prefer the actual cty string value
			s := v.AsString()
			return createScalarNode(s, s)
		case v.Type().Equals(cty.Bool):
			b := v.True()
			return createScalarNode(b, strconv.FormatBool(b))
		case v.Type() == cty.Number:
			// prefer integers when the numeric value is integral
			bf := v.AsBigFloat()
			if bf == nil {
				// fallback to string
				return createStringScalarNode(v.GoString())
			}
			// check if bf represents an exact integer
			if intVal, acc := bf.Int(nil); acc == big.Exact {
				s := intVal.String()
				return createScalarNode(intVal.Int64(), s)
			}
			s := bf.Text('g', -1)
			return createScalarNode(0.0, s)
		case v.Type().IsTupleType() || v.Type().IsListType() || v.Type().IsSetType():
			seq := &CandidateNode{Kind: SequenceNode}
			it := v.ElementIterator()
			for it.Next() {
				_, val := it.Element()
				// convert cty.Value to a node by wrapping in literal expr via string representation
				child := convertCtyValueToNode(val)
				seq.AddChild(child)
			}
			return seq
		case v.Type().IsMapType() || v.Type().IsObjectType():
			m := &CandidateNode{Kind: MappingNode}
			it := v.ElementIterator()
			for it.Next() {
				key, val := it.Element()
				keyStr := key.AsString()
				keyNode := createStringScalarNode(keyStr)
				valNode := convertCtyValueToNode(val)
				m.AddKeyValueChild(keyNode, valNode)
			}
			return m
		default:
			// fallback to string
			s := v.GoString()
			return createStringScalarNode(s)
		}
	case *hclsyntax.TupleConsExpr:
		// parse tuple/list into YAML sequence
		seq := &CandidateNode{Kind: SequenceNode}
		for _, exprVal := range e.Exprs {
			child := convertHclExprToNode(exprVal, src)
			seq.AddChild(child)
		}
		return seq
	case *hclsyntax.ObjectConsExpr:
		// parse object into YAML mapping
		m := &CandidateNode{Kind: MappingNode}
		m.Style = FlowStyle // Mark as inline object (flow style) for encoder
		for _, item := range e.Items {
			// evaluate key expression to get the key string
			keyVal, keyDiags := item.KeyExpr.Value(nil)
			if keyDiags != nil && keyDiags.HasErrors() {
				// fallback: try to extract key from source
				r := item.KeyExpr.Range()
				start := r.Start.Byte
				end := r.End.Byte
				if start > 0 && end >= start && end <= len(src) {
					keyNode := createStringScalarNode(strings.TrimSpace(string(src[start-1 : end])))
					valNode := convertHclExprToNode(item.ValueExpr, src)
					m.AddKeyValueChild(keyNode, valNode)
				}
				continue
			}
			keyStr := keyVal.AsString()
			keyNode := createStringScalarNode(keyStr)
			valNode := convertHclExprToNode(item.ValueExpr, src)
			m.AddKeyValueChild(keyNode, valNode)
		}
		return m
	case *hclsyntax.TemplateExpr:
		// join parts; if single literal, return that string
		var parts []string
		for _, p := range e.Parts {
			switch lp := p.(type) {
			case *hclsyntax.LiteralValueExpr:
				if lp.Val.Type().Equals(cty.String) {
					parts = append(parts, lp.Val.AsString())
				} else {
					parts = append(parts, lp.Val.GoString())
				}
			default:
				r := p.Range()
				start := r.Start.Byte
				end := r.End.Byte
				if start > 0 && end >= start && end <= len(src) {
					parts = append(parts, strings.TrimSpace(string(src[start-1:end])))
				} else {
					parts = append(parts, fmt.Sprintf("%v", p))
				}
			}
		}
		combined := strings.Join(parts, "")
		return createScalarNode(combined, combined)
	default:
		// try to evaluate the expression (handles unary, binary ops, etc.)
		val, diags := expr.Value(nil)
		if diags == nil || !diags.HasErrors() {
			// successfully evaluated, convert cty.Value to node
			return convertCtyValueToNode(val)
		}
		// fallback: extract source text for the expression
		r := expr.Range()
		start := r.Start.Byte
		end := r.End.Byte
		if start > 0 && end >= start && end <= len(src) {
			text := string(src[start-1 : end])
			return createStringScalarNode(text)
		}
		return createStringScalarNode(fmt.Sprintf("%v", expr))
	}
}

func convertCtyValueToNode(v cty.Value) *CandidateNode {
	if v.IsNull() {
		return createScalarNode(nil, "")
	}
	switch {
	case v.Type().Equals(cty.String):
		return createScalarNode("", v.AsString())
	case v.Type().Equals(cty.Bool):
		b := v.True()
		return createScalarNode(b, strconv.FormatBool(b))
	case v.Type() == cty.Number:
		bf := v.AsBigFloat()
		if bf == nil {
			return createStringScalarNode(v.GoString())
		}
		if intVal, acc := bf.Int(nil); acc == big.Exact {
			s := intVal.String()
			return createScalarNode(intVal.Int64(), s)
		}
		s := bf.Text('g', -1)
		return createScalarNode(0.0, s)
	case v.Type().IsTupleType() || v.Type().IsListType() || v.Type().IsSetType():
		seq := &CandidateNode{Kind: SequenceNode}
		it := v.ElementIterator()
		for it.Next() {
			_, val := it.Element()
			seq.AddChild(convertCtyValueToNode(val))
		}
		return seq
	case v.Type().IsMapType() || v.Type().IsObjectType():
		m := &CandidateNode{Kind: MappingNode}
		it := v.ElementIterator()
		for it.Next() {
			key, val := it.Element()
			keyNode := createStringScalarNode(key.AsString())
			valNode := convertCtyValueToNode(val)
			m.AddKeyValueChild(keyNode, valNode)
		}
		return m
	default:
		return createStringScalarNode(v.GoString())
	}

}
