package yqlib

import (
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

func (o *orderedMap) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.DocumentNode:
		if len(node.Content) == 0 {
			return nil
		}
		return o.UnmarshalYAML(node.Content[0])
	case yaml.AliasNode:
		return o.UnmarshalYAML(node.Alias)
	case yaml.ScalarNode:
		return node.Decode(&o.altVal)
	case yaml.MappingNode:
		// set kv to non-nil
		o.kv = []orderedMapKV{}
		for i := 0; i < len(node.Content); i += 2 {
			var key string
			var val orderedMap
			if err := node.Content[i].Decode(&key); err != nil {
				return err
			}
			if err := node.Content[i+1].Decode(&val); err != nil {
				return err
			}
			o.kv = append(o.kv, orderedMapKV{
				K: key,
				V: val,
			})
		}
		return nil
	case yaml.SequenceNode:
		// note that this has to be a pointer, so that nulls can be represented.
		var res []*orderedMap
		if err := node.Decode(&res); err != nil {
			return err
		}
		o.altVal = res
		o.kv = nil
		return nil
	case 0:
		// null
		o.kv = nil
		o.altVal = nil
		return nil
	default:
		return fmt.Errorf("orderedMap: invalid yaml node")
	}
}

func (o *orderedMap) MarshalYAML() (interface{}, error) {
	// fast path: kv is nil, use altVal
	if o.kv == nil {
		return o.altVal, nil
	}
	content := make([]*yaml.Node, 0, len(o.kv)*2)
	for _, val := range o.kv {
		n := new(yaml.Node)
		if err := n.Encode(val.V); err != nil {
			return nil, err
		}
		content = append(content, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   "!!str",
			Value: val.K,
		}, n)
	}
	return &yaml.Node{
		Kind:    yaml.MappingNode,
		Tag:     "!!map",
		Content: content,
	}, nil
}
