package yqlib

import (
	"container/list"

	yaml "gopkg.in/yaml.v3"
)

func assignAliasOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("AssignAlias operator!")

	aliasName := ""
	if !expressionNode.Operation.UpdateAssign {
		rhs, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.RHS)
		if err != nil {
			return Context{}, err
		}
		if rhs.MatchingNodes.Front() != nil {
			aliasName = rhs.MatchingNodes.Front().Value.(*CandidateNode).Node.Value
		}
	}

	lhs, err := d.GetMatchingNodes(context, expressionNode.LHS)

	if err != nil {
		return Context{}, err
	}

	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		log.Debugf("Setting aliasName : %v", candidate.GetKey())

		if expressionNode.Operation.UpdateAssign {
			rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(candidate), expressionNode.RHS)
			if err != nil {
				return Context{}, err
			}
			if rhs.MatchingNodes.Front() != nil {
				aliasName = rhs.MatchingNodes.Front().Value.(*CandidateNode).Node.Value
			}
		}

		if aliasName != "" {
			candidate.Node.Kind = yaml.AliasNode
			candidate.Node.Value = aliasName
		}
	}
	return context, nil
}

func getAliasOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("GetAlias operator!")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := &yaml.Node{Kind: yaml.ScalarNode, Value: candidate.Node.Value, Tag: "!!str"}
		result := candidate.CreateReplacement(node)
		results.PushBack(result)
	}
	return context.ChildContext(results), nil
}

func assignAnchorOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("AssignAnchor operator!")

	anchorName := ""
	if !expressionNode.Operation.UpdateAssign {
		rhs, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.RHS)
		if err != nil {
			return Context{}, err
		}

		if rhs.MatchingNodes.Front() != nil {
			anchorName = rhs.MatchingNodes.Front().Value.(*CandidateNode).Node.Value
		}
	}

	lhs, err := d.GetMatchingNodes(context, expressionNode.LHS)

	if err != nil {
		return Context{}, err
	}

	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		log.Debugf("Setting anchorName of : %v", candidate.GetKey())

		if expressionNode.Operation.UpdateAssign {
			rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(candidate), expressionNode.RHS)
			if err != nil {
				return Context{}, err
			}

			if rhs.MatchingNodes.Front() != nil {
				anchorName = rhs.MatchingNodes.Front().Value.(*CandidateNode).Node.Value
			}
		}

		candidate.Node.Anchor = anchorName
	}
	return context, nil
}

func getAnchorOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("GetAnchor operator!")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		anchor := candidate.Node.Anchor
		node := &yaml.Node{Kind: yaml.ScalarNode, Value: anchor, Tag: "!!str"}
		result := candidate.CreateReplacement(node)
		results.PushBack(result)
	}
	return context.ChildContext(results), nil
}

func explodeOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- ExplodeOperation")

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		rhs, err := d.GetMatchingNodes(context.SingleChildContext(candidate), expressionNode.RHS)

		if err != nil {
			return Context{}, err
		}
		for childEl := rhs.MatchingNodes.Front(); childEl != nil; childEl = childEl.Next() {
			err = explodeNode(childEl.Value.(*CandidateNode).Node, context)
			if err != nil {
				return Context{}, err
			}
		}

	}

	return context, nil
}

func reconstructAliasedMap(node *yaml.Node, context Context) error {
	var newContent = list.New()
	// can I short cut here by prechecking if there's an anchor in the map?
	// no it needs to recurse in overrideEntry.

	for index := 0; index < len(node.Content); index = index + 2 {
		keyNode := node.Content[index]
		valueNode := node.Content[index+1]
		log.Debugf("traversing %v", keyNode.Value)
		if keyNode.Value != "<<" {
			err := overrideEntry(node, keyNode, valueNode, index, context.ChildContext(newContent))
			if err != nil {
				return err
			}
		} else {
			if valueNode.Kind == yaml.SequenceNode {
				log.Debugf("an alias merge list!")
				for index := len(valueNode.Content) - 1; index >= 0; index = index - 1 {
					aliasNode := valueNode.Content[index]
					err := applyAlias(node, aliasNode.Alias, index, context.ChildContext(newContent))
					if err != nil {
						return err
					}
				}
			} else {
				log.Debugf("an alias merge!")
				err := applyAlias(node, valueNode.Alias, index, context.ChildContext(newContent))
				if err != nil {
					return err
				}
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
}

func explodeNode(node *yaml.Node, context Context) error {
	node.Anchor = ""
	switch node.Kind {
	case yaml.SequenceNode, yaml.DocumentNode:
		for index, contentNode := range node.Content {
			log.Debugf("exploding index %v", index)
			errorInContent := explodeNode(contentNode, context)
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
			node.Content = deepCloneContent(node.Alias.Content)
			node.Value = node.Alias.Value
			node.Alias = nil
		}
		return nil
	case yaml.MappingNode:
		// //check the map has an alias in it
		hasAlias := false
		for index := 0; index < len(node.Content); index = index + 2 {
			keyNode := node.Content[index]
			if keyNode.Value == "<<" {
				hasAlias = true
				break
			}
		}

		if hasAlias {
			// this is a slow op, which is why we want to check before running it.
			return reconstructAliasedMap(node, context)
		}
		// this map has no aliases, but it's kids might
		for index := 0; index < len(node.Content); index = index + 2 {
			keyNode := node.Content[index]
			valueNode := node.Content[index+1]
			err := explodeNode(keyNode, context)
			if err != nil {
				return err
			}
			err = explodeNode(valueNode, context)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return nil
	}
}

func applyAlias(node *yaml.Node, alias *yaml.Node, aliasIndex int, newContent Context) error {
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

func overrideEntry(node *yaml.Node, key *yaml.Node, value *yaml.Node, startIndex int, newContent Context) error {

	err := explodeNode(value, newContent)

	if err != nil {
		return err
	}

	for newEl := newContent.MatchingNodes.Front(); newEl != nil; newEl = newEl.Next() {
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

	err = explodeNode(key, newContent)
	if err != nil {
		return err
	}
	log.Debugf("adding %v:%v", key.Value, value.Value)
	newContent.MatchingNodes.PushBack(key)
	newContent.MatchingNodes.PushBack(value)
	return nil
}
