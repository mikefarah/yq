package yqlib

import (
	"fmt"
	"io"
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

	// process attributes
	body := dec.file.Body.(*hclsyntax.Body)
	for name, attr := range body.Attributes {
		keyNode := createStringScalarNode(name)
		valNode := convertHclExprToNode(attr.Expr, dec.fileBytes)
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
	for name, attr := range body.Attributes {
		key := createStringScalarNode(name)
		val := convertHclExprToNode(attr.Expr, src)
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
			// represent numbers as float string
			bf := v.AsBigFloat()
			if bf == nil {
				// fallback to string
				return createScalarNode(nil, v.GoString())
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
			return createScalarNode(nil, s)
		}
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
		// fallback: extract source text for the expression
		r := expr.Range()
		start := r.Start.Byte
		end := r.End.Byte
		if start > 0 && end >= start && end <= len(src) {
			text := string(src[start-1 : end])
			return createScalarNode(nil, text)
		}
		return createScalarNode(nil, fmt.Sprintf("%v", expr))
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
			return createScalarNode(nil, v.GoString())
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
		return createScalarNode(nil, v.GoString())
	}

}
