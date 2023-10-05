package yqlib

import "github.com/goccy/go-yaml/ast"

func (o *CandidateNode) goccyDecodeIntoChild(childNode ast.Node, anchorMap map[string]*CandidateNode) (*CandidateNode, error) {
	newChild := o.CreateChild()

	err := newChild.UnmarshalGoccyYAML(childNode, anchorMap)
	return newChild, err
}

func (o *CandidateNode) UnmarshalGoccyYAML(node ast.Node, anchorMap map[string]*CandidateNode) error {
	log.Debugf("UnmarshalYAML %v", node)
	log.Debugf("UnmarshalYAML %v", node.Type().String())

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
	case ast.TagType:
		o.UnmarshalGoccyYAML(node.(*ast.TagNode).Value, anchorMap)
		o.Tag = node.(*ast.TagNode).Start.Value
	case ast.MappingValueType:
		log.Debugf("UnmarshalYAML -  a mapping node")
		o.Kind = MappingNode
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
	default:
		log.Debugf("UnmarshalYAML -  node idea of the type!!")
	}
	log.Debugf("KIND: %v", o.Kind)
	return nil
}
