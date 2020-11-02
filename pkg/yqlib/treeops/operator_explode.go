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
		var newContent = list.New()
		for index := 0; index < len(node.Content); index = index + 2 {
			keyNode := node.Content[index]
			valueNode := node.Content[index+1]
			log.Debugf("traversing %v", keyNode.Value)
			if keyNode.Value != "<<" {
				err := overrideEntry(node, keyNode, valueNode, index, newContent)
				if err != nil {
					return err
				}
			} else {
				if valueNode.Kind == yaml.SequenceNode {
					log.Debugf("an alias merge list!")
					for index := 0; index < len(valueNode.Content); index = index + 1 {
						aliasNode := valueNode.Content[index]
						applyAlias(node, aliasNode.Alias, index, newContent)
					}
				} else {
					log.Debugf("an alias merge!")
					applyAlias(node, valueNode.Alias, index, newContent)
				}
			}
		}
		node.Content = make([]*yaml.Node, newContent.Len())
		index := 0
		for newEl := newContent.Front(); newEl != nil; newEl = newEl.Next() {
			node.Content[index] = newEl.Value.(*yaml.Node)
			index++
		}

		return nil
	default:
		return nil
	}
}

func applyAlias(node *yaml.Node, alias *yaml.Node, aliasIndex int, newContent *list.List) error {
	if alias == nil {
		return nil
	}
	for index := 0; index < len(alias.Content); index = index + 2 {
		keyNode := alias.Content[index]
		log.Debugf("applying alias key %v", keyNode.Value)
		valueNode := alias.Content[index+1]
		err := overrideEntry(node, keyNode, valueNode, aliasIndex, newContent)
		if err != nil {
			return err
		}
	}
	return nil
}

func overrideEntry(node *yaml.Node, key *yaml.Node, value *yaml.Node, startIndex int, newContent *list.List) error {

	err := explodeNode(value)

	if err != nil {
		return err
	}

	for newEl := newContent.Front(); newEl != nil; newEl = newEl.Next() {
		valueEl := newEl.Next() // move forward twice
		keyNode := newEl.Value.(*yaml.Node)
		log.Debugf("checking new content %v:%v", keyNode.Value, valueEl.Value.(*yaml.Node).Value)
		if keyNode.Value == key.Value && keyNode.Alias == nil && key.Alias == nil {
			log.Debugf("overridign new content")
			valueEl.Value = value
			return nil
		}
		newEl = valueEl // move forward twice
	}

	for index := startIndex + 2; index < len(node.Content); index = index + 2 {
		keyNode := node.Content[index]

		if keyNode.Value == key.Value && keyNode.Alias == nil {
			log.Debugf("content will be overridden at index %v", index)
			return nil
		}
	}

	err = explodeNode(key)
	if err != nil {
		return err
	}
	log.Debugf("adding %v:%v", key.Value, value.Value)
	newContent.PushBack(key)
	newContent.PushBack(value)
	return nil
}
