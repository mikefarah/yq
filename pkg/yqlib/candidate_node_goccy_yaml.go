package yqlib

import (
	"fmt"
	"strings"

	yaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	goccyToken "github.com/goccy/go-yaml/token"
)

func (o *CandidateNode) goccyDecodeIntoChild(childNode ast.Node, cm yaml.CommentMap, anchorMap map[string]*CandidateNode) (*CandidateNode, error) {
	newChild := o.CreateChild()

	err := newChild.UnmarshalGoccyYAML(childNode, cm, anchorMap)
	return newChild, err
}

func (o *CandidateNode) UnmarshalGoccyYAML(node ast.Node, cm yaml.CommentMap, anchorMap map[string]*CandidateNode) error {
	log.Debugf("UnmarshalYAML %v", node)
	log.Debugf("UnmarshalYAML %v", node.Type().String())
	log.Debugf("UnmarshalYAML Node Value: %v", node.String())
	log.Debugf("UnmarshalYAML Node GetComment: %v", node.GetComment())

	if node.GetComment() != nil {
		commentMapComments := cm[node.GetPath()]
		for _, comment := range node.GetComment().Comments {
			// need to use the comment map to find the position :/
			log.Debugf("%v has a comment of [%v]", node.GetPath(), comment.Token.Value)
			for _, commentMapComment := range commentMapComments {
				commentMapValue := strings.Join(commentMapComment.Texts, "\n")
				if commentMapValue == comment.Token.Value {
					log.Debug("found a matching entry in comment map")
					// we found the comment in the comment map,
					// now we can process the position
					switch commentMapComment.Position {
					case yaml.CommentHeadPosition:
						o.HeadComment = comment.String()
						log.Debug("its a head comment %v", comment.String())
					case yaml.CommentLinePosition:
						o.LineComment = comment.String()
						log.Debug("its a line comment %v", comment.String())
					case yaml.CommentFootPosition:
						o.FootComment = comment.String()
						log.Debug("its a foot comment %v", comment.String())
					}
				}
			}

		}
	}

	o.Value = node.String()
	o.Line = node.GetToken().Position.Line
	o.Column = node.GetToken().Position.Column

	switch node.Type() {
	case ast.IntegerType:
		o.Kind = ScalarNode
		o.Tag = "!!int"
	case ast.FloatType:
		o.Kind = ScalarNode
		o.Tag = "!!float"
	case ast.BoolType:
		o.Kind = ScalarNode
		o.Tag = "!!bool"
	case ast.NullType:
		log.Debugf("its a null type with value %v", node.GetToken().Value)
		o.Kind = ScalarNode
		o.Tag = "!!null"
		o.Value = node.GetToken().Value
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
		log.Debugf("astLiteral.Start.Type %v", astLiteral.Start.Type)
		if astLiteral.Start.Type == goccyToken.FoldedType {
			log.Debugf("folded Type %v", astLiteral.Start.Type)
			o.Style = FoldedStyle
		}
		log.Debug("start value: %v ", node.(*ast.LiteralNode).Start.Value)
		log.Debug("start value: %v ", node.(*ast.LiteralNode).Start.Type)
		// TODO: here I could put the original value with line breaks
		// to solve the multiline > problem
		o.Value = astLiteral.Value.Value
	case ast.TagType:
		if err := o.UnmarshalGoccyYAML(node.(*ast.TagNode).Value, cm, anchorMap); err != nil {
			return err
		}
		o.Tag = node.(*ast.TagNode).Start.Value
	case ast.MappingType:
		log.Debugf("UnmarshalYAML -  a mapping node")
		o.Kind = MappingNode
		o.Tag = "!!map"

		mappingNode := node.(*ast.MappingNode)
		if mappingNode.IsFlowStyle {
			o.Style = FlowStyle
		}
		for _, mappingValueNode := range mappingNode.Values {
			err := o.goccyProcessMappingValueNode(mappingValueNode, cm, anchorMap)
			if err != nil {
				return ast.ErrInvalidAnchorName
			}
		}
		if mappingNode.FootComment != nil {
			log.Debugf("mapping node has a foot comment of: %v", mappingNode.FootComment)
			o.FootComment = mappingNode.FootComment.String()
		}
	case ast.MappingValueType:
		log.Debugf("UnmarshalYAML -  a mapping node")
		o.Kind = MappingNode
		o.Tag = "!!map"
		mappingValueNode := node.(*ast.MappingValueNode)
		err := o.goccyProcessMappingValueNode(mappingValueNode, cm, anchorMap)
		if err != nil {
			return ast.ErrInvalidAnchorName
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

			valueNode, err := o.goccyDecodeIntoChild(astSeq[i], cm, anchorMap)
			if err != nil {
				return err
			}

			valueNode.Key = keyNode
			o.Content[i] = valueNode
		}
	case ast.AnchorType:
		log.Debugf("UnmarshalYAML -  an anchor node")
		anchorNode := node.(*ast.AnchorNode)
		err := o.UnmarshalGoccyYAML(anchorNode.Value, cm, anchorMap)
		if err != nil {
			return err
		}
		o.Anchor = anchorNode.Name.String()
		anchorMap[o.Anchor] = o

	case ast.AliasType:
		log.Debugf("UnmarshalYAML -  an alias node")
		aliasNode := node.(*ast.AliasNode)
		o.Kind = AliasNode
		o.Value = aliasNode.Value.String()
		o.Alias = anchorMap[o.Value]

	case ast.MergeKeyType:
		log.Debugf("UnmarshalYAML -  a merge key")
		o.Kind = ScalarNode
		o.Tag = "!!merge" // note - I should be able to get rid of this.
		o.Value = "<<"

	default:
		log.Debugf("UnmarshalYAML -  no idea of the type!!\n%v: %v", node.Type(), node.String())
	}
	log.Debugf("KIND: %v", o.Kind)
	return nil
}

func (o *CandidateNode) goccyProcessMappingValueNode(mappingEntry *ast.MappingValueNode, cm yaml.CommentMap, anchorMap map[string]*CandidateNode) error {
	log.Debug("UnmarshalYAML MAP KEY entry %v", mappingEntry.Key)

	// AddKeyValueFirst because it clones the nodes, and we want to have the real refs when Unmarshalling
	// particularly for the anchorMap
	keyNode, valueNode := o.AddKeyValueChild(&CandidateNode{}, &CandidateNode{})

	if err := keyNode.UnmarshalGoccyYAML(mappingEntry.Key, cm, anchorMap); err != nil {
		return err
	}

	log.Debug("UnmarshalYAML MAP VALUE entry %v", mappingEntry.Value)
	if err := valueNode.UnmarshalGoccyYAML(mappingEntry.Value, cm, anchorMap); err != nil {
		return err
	}

	if mappingEntry.FootComment != nil {
		valueNode.FootComment = mappingEntry.FootComment.String()
	}

	return nil
}
