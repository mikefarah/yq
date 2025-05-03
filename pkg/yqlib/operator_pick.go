package yqlib

import (
	"container/list"
	"fmt"
)

func pickMap(original *CandidateNode, indices *CandidateNode) *CandidateNode {

	filteredContent := make([]*CandidateNode, 0)
	for index := 0; index < len(indices.Content); index = index + 1 {
		keyToFind := indices.Content[index]

		indexInMap := findKeyInMap(original, keyToFind)
		if indexInMap > -1 {
			clonedKey := original.Content[indexInMap].Copy()
			clonedValue := original.Content[indexInMap+1].Copy()
			filteredContent = append(filteredContent, clonedKey, clonedValue)
		}
	}

	newNode := original.CopyWithoutContent()
	newNode.AddChildren(filteredContent)

	return newNode
}

func pickSequence(original *CandidateNode, indices *CandidateNode) (*CandidateNode, error) {

	filteredContent := make([]*CandidateNode, 0)
	for index := 0; index < len(indices.Content); index = index + 1 {
		indexInArray, err := parseInt(indices.Content[index].Value)
		if err != nil {
			return nil, fmt.Errorf("cannot index array with %v", indices.Content[index].Value)
		}

		if indexInArray > -1 && indexInArray < len(original.Content) {
			filteredContent = append(filteredContent, original.Content[indexInArray].Copy())
		}
	}

	newNode := original.CopyWithoutContent()
	newNode.AddChildren(filteredContent)

	return newNode, nil
}

func pickOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("Pick")

	contextIndicesToPick, err := d.GetMatchingNodes(context, expressionNode.RHS)

	if err != nil {
		return Context{}, err
	}
	indicesToPick := &CandidateNode{}
	if contextIndicesToPick.MatchingNodes.Len() > 0 {
		indicesToPick = contextIndicesToPick.MatchingNodes.Front().Value.(*CandidateNode)
	}

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		node := el.Value.(*CandidateNode)

		var replacement *CandidateNode
		switch node.Kind {
		case MappingNode:
			replacement = pickMap(node, indicesToPick)
		case SequenceNode:
			replacement, err = pickSequence(node, indicesToPick)
			if err != nil {
				return Context{}, err
			}

		default:
			return Context{}, fmt.Errorf("cannot pick indices from type %v (%v)", node.Tag, node.GetNicePath())
		}

		replacement.LeadingContent = node.LeadingContent
		results.PushBack(replacement)
	}

	return context.ChildContext(results), nil
}
