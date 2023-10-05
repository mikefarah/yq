package yqlib

import (
	"fmt"

	"github.com/goccy/go-yaml/ast"
	goccyToken "github.com/goccy/go-yaml/token"
)

func (o *CandidateNode) goccyDecodeIntoChild(childNode ast.Node, anchorMap map[string]*CandidateNode) (*CandidateNode, error) {
	newChild := o.CreateChild()

	err := newChild.UnmarshalGoccyYAML(childNode, anchorMap)
	return newChild, err
}

func (o *CandidateNode) UnmarshalGoccyYAML(node ast.Node, anchorMap map[string]*CandidateNode) error {
	log.Debugf("UnmarshalYAML %v", node)
	log.Debugf("UnmarshalYAML %v", node.Type().String())
	log.Debugf("UnmarshalYAML Value: %v", node.String())

	o.Value = node.String()
	switch node.Type() {
	case ast.IntegerType:
		o.Kind = ScalarNode
		o.Tag = "!!int"
	case ast.FloatType:
		o.Kind = ScalarNode
		o.Tag = "!!float"
	case ast.StringType:
		o.Kind = ScalarNode
		o.Tag = "!!str"
		switch node.GetToken().Type {
		case goccyToken.SingleQuoteType:
			o.Style = SingleQuotedStyle
		case goccyToken.DoubleQuoteType:
			o.Style = DoubleQuotedStyle
		}
		o.Value = node.(*ast.StringNode).Value
		log.Debugf("string value %v", node.(*ast.StringNode).Value)
	case ast.LiteralType:
		o.Kind = ScalarNode
		o.Tag = "!!str"
		o.Style = LiteralStyle
		astLiteral := node.(*ast.LiteralNode)
		if astLiteral.Start.Type == goccyToken.FoldedType {
			o.Style = FoldedStyle
		}
		log.Debug("startvalue: %v ", node.(*ast.LiteralNode).Start.Value)
		log.Debug("startvalue: %v ", node.(*ast.LiteralNode).Start.Type)
		o.Value = astLiteral.Value.Value
	case ast.TagType:
		o.UnmarshalGoccyYAML(node.(*ast.TagNode).Value, anchorMap)
		o.Tag = node.(*ast.TagNode).Start.Value
	case ast.MappingValueType, ast.MappingType:
		log.Debugf("UnmarshalYAML -  a mapping node")
		o.Kind = MappingNode
		o.Tag = "!!map"

		if node.Type() == ast.MappingType {
			o.Style = FlowStyle
		}

		astMapIter := node.(ast.MapNode).MapRange()
		for astMapIter.Next() {
			log.Debug("UnmarshalYAML map entry %v", astMapIter.Key().String())
			keyNode, err := o.goccyDecodeIntoChild(astMapIter.Key(), anchorMap)
			if err != nil {
				return err
			}

			keyNode.IsMapKey = true
			log.Debug("UnmarshalYAML map value %v", astMapIter.Value().String())
			valueNode, err := o.goccyDecodeIntoChild(astMapIter.Value(), anchorMap)
			if err != nil {
				return err
			}

			o.Content = append(o.Content, keyNode, valueNode)
		}
	case ast.SequenceType:
		log.Debugf("UnmarshalYAML -  a sequence node")
		o.Kind = SequenceNode
		o.Tag = "!!seq"
		sequenceNode := node.(*ast.SequenceNode)
		if sequenceNode.IsFlowStyle {
			o.Style = FlowStyle
		}
		astSeq := sequenceNode.Values
		o.Content = make([]*CandidateNode, len(astSeq))
		for i := 0; i < len(astSeq); i++ {
			keyNode := o.CreateChild()
			keyNode.IsMapKey = true
			keyNode.Tag = "!!int"
			keyNode.Kind = ScalarNode
			keyNode.Value = fmt.Sprintf("%v", i)

			valueNode, err := o.goccyDecodeIntoChild(astSeq[i], anchorMap)
			if err != nil {
				return err
			}

			valueNode.Key = keyNode
			o.Content[i] = valueNode
		}

	default:
		log.Debugf("UnmarshalYAML -  node idea of the type!!")
	}
	log.Debugf("KIND: %v", o.Kind)
	return nil
}
