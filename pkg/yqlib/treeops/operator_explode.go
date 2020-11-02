package treeops

import (
	"container/list"

	"gopkg.in/yaml.v3"
)

func ExplodeOperator(d *dataTreeNavigator, matchMap *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("-- ExplodeOperation")

	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		rhs, err := d.GetMatchingNodes(nodeToMap(candidate), pathNode.Rhs)

		if err != nil {
			return nil, err
		}
		for childEl := rhs.Front(); childEl != nil; childEl = childEl.Next() {
			explodeNode(childEl.Value.(*CandidateNode).Node)
		}

	}

	return matchMap, nil
}

func explodeNode(node *yaml.Node) error {
	node.Anchor = ""
	switch node.Kind {
	case yaml.SequenceNode, yaml.DocumentNode:
		for index, contentNode := range node.Content {
			log.Debugf("exploding index %v", index)
			errorInContent := explodeNode(contentNode)
			if errorInContent != nil {
				return errorInContent
			}
		}
		return nil
	case yaml.AliasNode:
		log.Debugf("its an alias!")
		if node.Alias != nil {
			node.Kind = node.Alias.Kind
			node.Style = node.Alias.Style
			node.Tag = node.Alias.Tag
			node.Content = node.Alias.Content
			node.Value = node.Alias.Value
			node.Alias = nil
		}
		return nil
	case yaml.MappingNode:
		for index := 0; index < len(node.Content); index = index + 2 {
			keyNode := node.Content[index]
			valueNode := node.Content[index+1]
			log.Debugf("traversing %v", keyNode.Value)
			if keyNode.Value != "<<" {
				errorInContent := explodeNode(valueNode)
				if errorInContent != nil {
					return errorInContent
				}
				errorInContent = explodeNode(keyNode)
				if errorInContent != nil {
					return errorInContent
				}
			} else {
				if valueNode.Kind == yaml.SequenceNode {
					log.Debugf("an alias merge list!")
					for index := len(valueNode.Content) - 1; index >= 0; index = index - 1 {
						aliasNode := valueNode.Content[index]
						applyAlias(node, aliasNode.Alias)
					}
				} else {
					log.Debugf("an alias merge!")
					applyAlias(node, valueNode.Alias)
				}
				node.Content = append(node.Content[:index], node.Content[index+2:]...)
				//replay that index, since the array is shorter now.
				index = index - 2
			}
		}

		return nil
	default:
		return nil
	}
}

func applyAlias(node *yaml.Node, alias *yaml.Node) {
	if alias == nil {
		return
	}
	for index := 0; index < len(alias.Content); index = index + 2 {
		keyNode := alias.Content[index]
		log.Debugf("applying alias key %v", keyNode.Value)
		valueNode := alias.Content[index+1]
		setIfNotThere(node, keyNode.Value, valueNode)
	}
}

func setIfNotThere(node *yaml.Node, key string, value *yaml.Node) {
	for index := 0; index < len(node.Content); index = index + 2 {
		keyNode := node.Content[index]
		if keyNode.Value == key {
			return
		}
	}
	// need to add it to the map
	mapEntryKey := yaml.Node{Value: key, Kind: yaml.ScalarNode}
	node.Content = append(node.Content, &mapEntryKey)
	node.Content = append(node.Content, value)
}
