package yqlib

import (
	"container/list"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

func pickMap(original *yaml.Node, indices *yaml.Node) *yaml.Node {

	filteredContent := make([]*yaml.Node, 0)
	for index := 0; index < len(indices.Content); index = index + 1 {
		keyToFind := indices.Content[index]

		indexInMap := findKeyInMap(original, keyToFind)
		if indexInMap > -1 {
			clonedKey := deepClone(original.Content[indexInMap])
			clonedValue := deepClone(original.Content[indexInMap+1])
			filteredContent = append(filteredContent, clonedKey, clonedValue)
		}
	}

	newNode := deepCloneNoContent(original)
	newNode.Content = filteredContent

	return newNode
}

func pickSequence(original *yaml.Node, indices *yaml.Node) (*yaml.Node, error) {

	filteredContent := make([]*yaml.Node, 0)
	for index := 0; index < len(indices.Content); index = index + 1 {
		indexInArray, err := parseInt(indices.Content[index].Value)
		if err != nil {
			return nil, fmt.Errorf("cannot index array with %v", indices.Content[index].Value)
		}

		if indexInArray > -1 && indexInArray < len(original.Content) {
			filteredContent = append(filteredContent, deepClone(original.Content[indexInArray]))
		}
	}

	newNode := deepCloneNoContent(original)
	newNode.Content = filteredContent

	return newNode, nil
}

func pickOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("Pick")

	contextIndicesToPick, err := d.GetMatchingNodes(context, expressionNode.RHS)

	if err != nil {
		return Context{}, err
	}
	indicesToPick := &yaml.Node{}
	if contextIndicesToPick.MatchingNodes.Len() > 0 {
		indicesToPick = contextIndicesToPick.MatchingNodes.Front().Value.(*CandidateNode).Node
	}

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := unwrapDoc(candidate.Node)

		var replacement *yaml.Node
		if node.Kind == yaml.MappingNode {
			replacement = pickMap(node, indicesToPick)
		} else if node.Kind == yaml.SequenceNode {
			replacement, err = pickSequence(node, indicesToPick)
			if err != nil {
				return Context{}, err
			}

		} else {
			return Context{}, fmt.Errorf("cannot pick indicies from type %v (%v)", node.Tag, candidate.GetNicePath())
		}

		results.PushBack(candidate.CreateReplacementWithDocWrappers(replacement))
	}

	return context.ChildContext(results), nil
}
