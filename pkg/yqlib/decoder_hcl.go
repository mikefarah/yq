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

// extractLineComment extracts any inline comment after the given position
func extractLineComment(src []byte, endPos int) string {
	// Look for # comment after the token
	for i := endPos; i < len(src); i++ {
		if src[i] == '#' {
			// Found comment, extract until end of line
			start := i
			for i < len(src) && src[i] != '\n' {
				i++
			}
			return strings.TrimSpace(string(src[start:i]))
		}
		if src[i] == '\n' {
			// Hit newline before comment
			break
		}
		// Skip whitespace and other characters
	}
	return ""
}

// extractHeadComment extracts comments before a given start position
func extractHeadComment(src []byte, startPos int) string {
	var comments []string

	// Start just before the token and skip trailing whitespace
	i := startPos - 1
	for i >= 0 && (src[i] == ' ' || src[i] == '\t' || src[i] == '\n' || src[i] == '\r') {
		i--
	}

	for i >= 0 {
		// Find line boundaries
		lineEnd := i
		for i >= 0 && src[i] != '\n' {
			i--
		}
		lineStart := i + 1

		line := strings.TrimRight(string(src[lineStart:lineEnd+1]), " \t\r")
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			break
		}

		if !strings.HasPrefix(trimmed, "#") {
			break
		}

		comments = append([]string{trimmed}, comments...)

		// Move to previous line (skip any whitespace/newlines)
		i = lineStart - 1
		for i >= 0 && (src[i] == ' ' || src[i] == '\t' || src[i] == '\n' || src[i] == '\r') {
			i--
		}
	}

	if len(comments) > 0 {
		return strings.Join(comments, "\n")
	}
	return ""
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
	firstAttr := true
	for _, attrWithName := range sortedAttributes(body.Attributes) {
		keyNode := createStringScalarNode(attrWithName.Name)
		valNode := convertHclExprToNode(attrWithName.Attr.Expr, dec.fileBytes)

		// Attach comments if any
		attrRange := attrWithName.Attr.Range()
		headComment := extractHeadComment(dec.fileBytes, attrRange.Start.Byte)
		if firstAttr && headComment != "" {
			// For the first attribute, apply its head comment to the root
			root.HeadComment = headComment
			firstAttr = false
		} else if headComment != "" {
			keyNode.HeadComment = headComment
		}
		if lineComment := extractLineComment(dec.fileBytes, attrRange.End.Byte); lineComment != "" {
			valNode.LineComment = lineComment
		}

		root.AddKeyValueChild(keyNode, valNode)
	}

	// process blocks
	// Count blocks by type at THIS level to detect multiple separate blocks
	blocksByType := make(map[string]int)
	for _, block := range body.Blocks {
		blocksByType[block.Type]++
	}

	for _, block := range body.Blocks {
		addBlockToMapping(root, block, dec.fileBytes, blocksByType[block.Type] > 1)
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

		// Attach comments if any
		attrRange := attrWithName.Attr.Range()
		if headComment := extractHeadComment(src, attrRange.Start.Byte); headComment != "" {
			key.HeadComment = headComment
		}
		if lineComment := extractLineComment(src, attrRange.End.Byte); lineComment != "" {
			val.LineComment = lineComment
		}

		node.AddKeyValueChild(key, val)
	}

	// Process nested blocks, counting blocks by type at THIS level
	// to detect which block types appear multiple times
	blocksByType := make(map[string]int)
	for _, block := range body.Blocks {
		blocksByType[block.Type]++
	}

	for _, block := range body.Blocks {
		addBlockToMapping(node, block, src, blocksByType[block.Type] > 1)
	}
	return node
}

// addBlockToMapping nests block type and labels into the parent mapping, merging children.
// isMultipleBlocksOfType indicates if there are multiple blocks of this type at THIS level
func addBlockToMapping(parent *CandidateNode, block *hclsyntax.Block, src []byte, isMultipleBlocksOfType bool) {
	bodyNode := hclBodyToNode(block.Body, src)
	current := parent

	// ensure block type mapping exists
	var typeNode *CandidateNode
	for i := 0; i < len(current.Content); i += 2 {
		if current.Content[i].Value == block.Type {
			typeNode = current.Content[i+1]
			break
		}
	}
	if typeNode == nil {
		_, typeNode = current.AddKeyValueChild(createStringScalarNode(block.Type), &CandidateNode{Kind: MappingNode})
		// Mark the type node if there are multiple blocks of this type at this level
		// This tells the encoder to emit them as separate blocks rather than consolidating them
		if isMultipleBlocksOfType {
			typeNode.EncodeSeparate = true
		}
	}
	current = typeNode

	// walk labels, creating/merging mappings
	for _, label := range block.Labels {
		var next *CandidateNode
		for i := 0; i < len(current.Content); i += 2 {
			if current.Content[i].Value == label {
				next = current.Content[i+1]
				break
			}
		}
		if next == nil {
			_, next = current.AddKeyValueChild(createStringScalarNode(label), &CandidateNode{Kind: MappingNode})
		}
		current = next
	}

	// merge body attributes/blocks into the final mapping
	for i := 0; i < len(bodyNode.Content); i += 2 {
		current.AddKeyValueChild(bodyNode.Content[i], bodyNode.Content[i+1])
	}
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
			node := createScalarNode(s, s)
			// Don't set style for regular quoted strings - let YAML handle naturally
			return node
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
				if start >= 0 && end >= start && end <= len(src) {
					keyNode := createStringScalarNode(strings.TrimSpace(string(src[start:end])))
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
		// Reconstruct template string, preserving ${} syntax for interpolations
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
				// Non-literal expression - reconstruct with ${} wrapper
				r := p.Range()
				start := r.Start.Byte
				end := r.End.Byte
				if start >= 0 && end >= start && end <= len(src) {
					exprText := string(src[start:end])
					parts = append(parts, "${"+exprText+"}")
				} else {
					parts = append(parts, fmt.Sprintf("${%v}", p))
				}
			}
		}
		combined := strings.Join(parts, "")
		node := createScalarNode(combined, combined)
		// Set DoubleQuotedStyle for all templates (which includes all quoted strings in HCL)
		// This ensures HCL roundtrips preserve quotes, and YAML properly quotes strings with ${}
		node.Style = DoubleQuotedStyle
		return node
	case *hclsyntax.ScopeTraversalExpr:
		// Simple identifier/traversal (e.g. unquoted string literal in HCL)
		r := e.Range()
		start := r.Start.Byte
		end := r.End.Byte
		if start >= 0 && end >= start && end <= len(src) {
			text := strings.TrimSpace(string(src[start:end]))
			return createStringScalarNode(text)
		}
		// Fallback to root name if source unavailable
		if len(e.Traversal) > 0 {
			if root, ok := e.Traversal[0].(hcl.TraverseRoot); ok {
				return createStringScalarNode(root.Name)
			}
		}
		return createStringScalarNode("")
	case *hclsyntax.FunctionCallExpr:
		// Preserve function calls as raw expressions for roundtrip
		r := e.Range()
		start := r.Start.Byte
		end := r.End.Byte
		if start >= 0 && end >= start && end <= len(src) {
			text := strings.TrimSpace(string(src[start:end]))
			node := createStringScalarNode(text)
			node.Style = 0
			return node
		}
		node := createStringScalarNode(e.Name)
		node.Style = 0
		return node
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
		if start >= 0 && end >= start && end <= len(src) {
			text := string(src[start:end])
			// Mark as unquoted expression so encoder emits without quoting
			node := createStringScalarNode(text)
			node.Style = 0
			return node
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
