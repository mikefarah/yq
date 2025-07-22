package yqlib

import (
	"container/list"
	"fmt"
)

var showMergeAnchorToSpecWarning = true

func assignAliasOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("AssignAlias operator!")

	aliasName := ""
	if !expressionNode.Operation.UpdateAssign {
		rhs, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.RHS)
		if err != nil {
			return Context{}, err
		}
		if rhs.MatchingNodes.Front() != nil {
			aliasName = rhs.MatchingNodes.Front().Value.(*CandidateNode).Value
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
				aliasName = rhs.MatchingNodes.Front().Value.(*CandidateNode).Value
			}
		}

		if aliasName != "" {
			candidate.Kind = AliasNode
			candidate.Value = aliasName
		}
	}
	return context, nil
}

func getAliasOperator(_ *dataTreeNavigator, context Context, _ *ExpressionNode) (Context, error) {
	log.Debugf("GetAlias operator!")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		result := candidate.CreateReplacement(ScalarNode, "!!str", candidate.Value)
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
			anchorName = rhs.MatchingNodes.Front().Value.(*CandidateNode).Value
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
				anchorName = rhs.MatchingNodes.Front().Value.(*CandidateNode).Value
			}
		}

		candidate.Anchor = anchorName
	}
	return context, nil
}

func getAnchorOperator(_ *dataTreeNavigator, context Context, _ *ExpressionNode) (Context, error) {
	log.Debugf("GetAnchor operator!")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		anchor := candidate.Anchor
		result := candidate.CreateReplacement(ScalarNode, "!!str", anchor)
		results.PushBack(result)
	}
	return context.ChildContext(results), nil
}

func explodeOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("ExplodeOperation")

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		rhs, err := d.GetMatchingNodes(context.SingleChildContext(candidate), expressionNode.RHS)

		if err != nil {
			return Context{}, err
		}
		for childEl := rhs.MatchingNodes.Front(); childEl != nil; childEl = childEl.Next() {
			err = explodeNode(childEl.Value.(*CandidateNode), context)
			if err != nil {
				return Context{}, err
			}
		}

	}

	return context, nil
}

func fixedReconstructAliasedMap(node *CandidateNode) error {
	var newContent = []*CandidateNode{}

	for index := 0; index < len(node.Content); index = index + 2 {
		keyNode := node.Content[index]
		valueNode := node.Content[index+1]
		if keyNode.Tag != "!!merge" {
			// always add in plain nodes
			// explode both the key and value nodes
			if err := explodeNode(keyNode, Context{}); err != nil {
				return err
			}
			if err := explodeNode(valueNode, Context{}); err != nil {
				return err
			}
			newContent = append(newContent, keyNode, valueNode)
		} else {
			sequence := valueNode
			if sequence.Kind == AliasNode {
				sequence = sequence.Alias
			}
			if sequence.Kind != SequenceNode {
				sequence = &CandidateNode{Content: []*CandidateNode{sequence}}
			}
			for index := 0; index < len(sequence.Content); index = index + 1 {
				// for merge anchors, we only set them if the key is not already in node or the newContent
				mergeNodeSeq := sequence.Content[index]
				if mergeNodeSeq.Kind == AliasNode {
					mergeNodeSeq = mergeNodeSeq.Alias
				}
				if mergeNodeSeq.Kind != MappingNode {
					return fmt.Errorf("can only use merge anchors with maps (!!map) or sequences (!!seq) of maps, but got sequence containing %v", mergeNodeSeq.Tag)
				}
				itemsToAdd := mergeNodeSeq.FilterMapContentByKey(func(keyNode *CandidateNode) bool {
					return getContentValueByKey(node.Content, keyNode.Value) == nil &&
						getContentValueByKey(newContent, keyNode.Value) == nil
				})

				for _, item := range itemsToAdd {
					// copy to ensure exploding doesn't modify the original node
					itemCopy := item.Copy()
					if err := explodeNode(itemCopy, Context{}); err != nil {
						return err
					}
					newContent = append(newContent, itemCopy)
				}
			}
		}
	}
	node.Content = newContent
	return nil
}

func reconstructAliasedMap(node *CandidateNode, context Context) error {
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
			if valueNode.Kind == SequenceNode {
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
	node.Content = make([]*CandidateNode, 0)
	for newEl := newContent.Front(); newEl != nil; newEl = newEl.Next() {
		node.AddChild(newEl.Value.(*CandidateNode))
	}
	return nil
}

func explodeNode(node *CandidateNode, context Context) error {
	log.Debugf("explodeNode -  %v", NodeToString(node))
	node.Anchor = ""
	switch node.Kind {
	case SequenceNode:
		for index, contentNode := range node.Content {
			log.Debugf("explodeNode -  index %v", index)
			errorInContent := explodeNode(contentNode, context)
			if errorInContent != nil {
				return errorInContent
			}
		}
		return nil
	case AliasNode:
		log.Debugf("explodeNode - an alias to %v", NodeToString(node.Alias))
		if node.Alias != nil {
			node.Kind = node.Alias.Kind
			node.Style = node.Alias.Style
			node.Tag = node.Alias.Tag
			node.AddChildren(node.Alias.Content)
			node.Value = node.Alias.Value
			node.Alias = nil
		}
		log.Debug("now I'm %v", NodeToString(node))
		return nil
	case MappingNode:
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
			if ConfiguredYamlPreferences.FixMergeAnchorToSpec {
				return fixedReconstructAliasedMap(node)
			}
			if showMergeAnchorToSpecWarning {
				log.Warning("--yaml-fix-merge-anchor-to-spec is false; causing merge anchors to override the existing values which isn't to the yaml spec. This flag will default to true in late 2025.")
				showMergeAnchorToSpecWarning = false
			}
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

func applyAlias(node *CandidateNode, alias *CandidateNode, aliasIndex int, newContent Context) error {
	log.Debug("alias is nil ?")
	if alias == nil {
		return nil
	}
	log.Debug("alias: %v", NodeToString(alias))
	if alias.Kind != MappingNode {
		return fmt.Errorf("can only use merge anchors with maps (!!map) or sequences (!!seq) of maps, but got sequence containing %v", alias.Tag)
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

func overrideEntry(node *CandidateNode, key *CandidateNode, value *CandidateNode, startIndex int, newContent Context) error {

	err := explodeNode(value, newContent)

	if err != nil {
		return err
	}

	for newEl := newContent.MatchingNodes.Front(); newEl != nil; newEl = newEl.Next() {
		valueEl := newEl.Next() // move forward twice
		keyNode := newEl.Value.(*CandidateNode)
		log.Debugf("checking new content %v:%v", keyNode.Value, valueEl.Value.(*CandidateNode).Value)
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
