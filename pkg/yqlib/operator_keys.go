package yqlib

import (
	"container/list"
	"fmt"

	"gopkg.in/yaml.v3"
)

func keysOperator(d *dataTreeNavigator, matchMap *list.List, expressionNode *ExpressionNode) (*list.List, error) {
	log.Debugf("-- keysOperator")

	var results = list.New()

	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := unwrapDoc(candidate.Node)
		var targetNode *yaml.Node
		if node.Kind == yaml.MappingNode {
			targetNode = getMapKeys(node)
		} else if node.Kind == yaml.SequenceNode {
			targetNode = getIndicies(node)
		} else {
			return nil, fmt.Errorf("Cannot get keys of %v, keys only works for maps and arrays", node.Tag)
		}

		result := candidate.CreateChild(nil, targetNode)
		results.PushBack(result)
	}

	return results, nil
}

func getMapKeys(node *yaml.Node) *yaml.Node {
	contents := make([]*yaml.Node, 0)
	for index := 0; index < len(node.Content); index = index + 2 {
		contents = append(contents, node.Content[index])
	}
	return &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Content: contents}
}

func getIndicies(node *yaml.Node) *yaml.Node {
	var contents = make([]*yaml.Node, len(node.Content))

	for index := range node.Content {
		contents[index] = &yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   "!!int",
			Value: fmt.Sprintf("%v", index),
		}
	}

	return &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Content: contents}
}
