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

// bodyItem represents either an attribute or a block at a given byte position in the source,
// allowing attributes and blocks to be processed together in source order.
type bodyItem struct {
	startByte int
	attr      *attributeWithName // non-nil for attributes
	block     *hclsyntax.Block   // non-nil for blocks
}

// sortedBodyItems returns attributes and blocks interleaved in source declaration order.
func sortedBodyItems(attrs hclsyntax.Attributes, blocks hclsyntax.Blocks) []bodyItem {
	var items []bodyItem
	for name, attr := range attrs {
		items = append(items, bodyItem{
			startByte: attr.Range().Start.Byte,
			attr:      &attributeWithName{Name: name, Attr: attr},
		})
	}
	for _, block := range blocks {
		b := block
		items = append(items, bodyItem{
			startByte: b.TypeRange.Start.Byte,
			block:     b,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].startByte < items[j].startByte
	})
	return items
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

// hasPrecedingBlankLine reports whether there is a blank line immediately before startPos,
// skipping over any immediately preceding comment lines and whitespace.
func hasPrecedingBlankLine(src []byte, startPos int) bool {
	i := startPos - 1

	// Skip trailing spaces/tabs on the current token's preceding content
	for i >= 0 && (src[i] == ' ' || src[i] == '\t') {
		i--
	}

	// We expect to be sitting just before a newline that ends the previous line.
	// Walk backwards skipping comment lines until we find a blank line or a non-comment line.
	for i >= 0 {
		// We should be pointing at '\n' (end of previous line) or start of file.
		if src[i] != '\n' {
			return false
		}
		i-- // step past the '\n'

		// Skip '\r' for Windows line endings
		if i >= 0 && src[i] == '\r' {
			i--
		}

		// If immediately another '\n', this is a blank line.
		if i < 0 || src[i] == '\n' {
			return true
		}

		// Read the previous line to see if it's a comment or blank.
		lineEnd := i
		for i >= 0 && src[i] != '\n' {
			i--
		}
		lineStart := i + 1
		line := strings.TrimSpace(string(src[lineStart : lineEnd+1]))

		if line == "" {
			return true
		}

		if strings.HasPrefix(line, "#") {
			// This line is a comment belonging to the current element; keep scanning upward.
			continue
		}

		// A non-blank, non-comment line: no blank line precedes this element.
		return false
	}

	return false
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

	body := dec.file.Body.(*hclsyntax.Body)

	// Count blocks by type at THIS level to detect multiple separate blocks of the same type.
	blocksByType := make(map[string]int)
	for _, block := range body.Blocks {
		blocksByType[block.Type]++
	}

	// Process attributes and blocks together in source declaration order.
	isFirst := true
	for _, item := range sortedBodyItems(body.Attributes, body.Blocks) {
		if item.attr != nil {
			aw := item.attr
			keyNode := createStringScalarNode(aw.Name)
			valNode := convertHclExprToNode(aw.Attr.Expr, dec.fileBytes)

			attrRange := aw.Attr.Range()
			headComment := extractHeadComment(dec.fileBytes, attrRange.Start.Byte)
			if isFirst && headComment != "" {
				// For the first element, apply its head comment to the root node
				root.HeadComment = headComment
			} else if headComment != "" {
				keyNode.HeadComment = headComment
			}
			if lineComment := extractLineComment(dec.fileBytes, attrRange.End.Byte); lineComment != "" {
				valNode.LineComment = lineComment
			}
			if !isFirst && hasPrecedingBlankLine(dec.fileBytes, attrRange.Start.Byte) {
				keyNode.BlankLineBefore = true
			}

			root.AddKeyValueChild(keyNode, valNode)
		} else {
			block := item.block
			headComment := extractHeadComment(dec.fileBytes, block.TypeRange.Start.Byte)
			if isFirst && headComment != "" {
				root.HeadComment = headComment
			}
			addBlockToMappingOrdered(root, block, dec.fileBytes, blocksByType[block.Type] > 1, isFirst, headComment)
		}
		isFirst = false
	}

	dec.documentIndex++
	root.document = dec.documentIndex - 1
	return root, nil
}

func hclBodyToNode(body *hclsyntax.Body, src []byte) *CandidateNode {
	node := &CandidateNode{Kind: MappingNode}

	blocksByType := make(map[string]int)
	for _, block := range body.Blocks {
		blocksByType[block.Type]++
	}

	isFirst := true
	for _, item := range sortedBodyItems(body.Attributes, body.Blocks) {
		if item.attr != nil {
			aw := item.attr
			key := createStringScalarNode(aw.Name)
			val := convertHclExprToNode(aw.Attr.Expr, src)

			attrRange := aw.Attr.Range()
			if headComment := extractHeadComment(src, attrRange.Start.Byte); headComment != "" {
				key.HeadComment = headComment
			}
			if lineComment := extractLineComment(src, attrRange.End.Byte); lineComment != "" {
				val.LineComment = lineComment
			}
			if !isFirst && hasPrecedingBlankLine(src, attrRange.Start.Byte) {
				key.BlankLineBefore = true
			}

			node.AddKeyValueChild(key, val)
		} else {
			block := item.block
			headComment := extractHeadComment(src, block.TypeRange.Start.Byte)
			addBlockToMappingOrdered(node, block, src, blocksByType[block.Type] > 1, isFirst, headComment)
		}
		isFirst = false
	}
	return node
}

// addBlockToMappingOrdered nests a block's type and labels into the parent mapping, merging children.
// isMultipleBlocksOfType: there are multiple blocks of this type at this level.
// isFirstInParent: this block is the first element in the parent (no preceding sibling).
// headComment: any comment extracted before this block's type keyword.
func addBlockToMappingOrdered(parent *CandidateNode, block *hclsyntax.Block, src []byte, isMultipleBlocksOfType bool, isFirstInParent bool, headComment string) {
	bodyNode := hclBodyToNode(block.Body, src)
	current := parent

	// ensure block type mapping exists
	var typeNode *CandidateNode
	var typeKeyNode *CandidateNode
	for i := 0; i < len(current.Content); i += 2 {
		if current.Content[i].Value == block.Type {
			typeKeyNode = current.Content[i]
			typeNode = current.Content[i+1]
			break
		}
	}
	if typeNode == nil {
		var newTypeKey *CandidateNode
		newTypeKey, typeNode = current.AddKeyValueChild(createStringScalarNode(block.Type), &CandidateNode{Kind: MappingNode})
		typeKeyNode = newTypeKey
		// Mark the type node if there are multiple blocks of this type at this level.
		// This tells the encoder to emit them as separate blocks rather than consolidating them.
		if isMultipleBlocksOfType {
			typeNode.EncodeSeparate = true
		}
		// Store the head comment on the type key (non-first elements only; first element's
		// comment is handled by the caller and applied to the root node).
		if !isFirstInParent && headComment != "" {
			typeKeyNode.HeadComment = headComment
		}
		// Detect blank line before this block in the source.
		// Only set it when this is not the first element (i.e. something already precedes it).
		if !isFirstInParent && hasPrecedingBlankLine(src, block.TypeRange.Start.Byte) {
			typeKeyNode.BlankLineBefore = true
		}
	}
	current = typeNode

	// walk labels, creating/merging mappings
	for labelIdx, label := range block.Labels {
		var next *CandidateNode
		var labelKey *CandidateNode
		for i := 0; i < len(current.Content); i += 2 {
			if current.Content[i].Value == label {
				labelKey = current.Content[i]
				next = current.Content[i+1]
				break
			}
		}
		if next == nil {
			var newLabelKey *CandidateNode
			newLabelKey, next = current.AddKeyValueChild(createStringScalarNode(label), &CandidateNode{Kind: MappingNode})
			labelKey = newLabelKey
			// For same-type blocks: mark the first label key with BlankLineBefore when
			// there is a blank line before this block in the source.
			if labelIdx == 0 && len(current.Content) > 2 {
				if hasPrecedingBlankLine(src, block.TypeRange.Start.Byte) {
					labelKey.BlankLineBefore = true
				}
			}
		}
		_ = labelKey
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
