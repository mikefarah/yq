package yqlib

import (
	"container/list"
	"strconv"
)

func omitMap(original *CandidateNode, indices *CandidateNode) *CandidateNode {
	filteredContent := make([]*CandidateNode, 0, max(0, len(original.Content)-len(indices.Content)*2))

	for index := 0; index < len(original.Content); index += 2 {
		pos := findInArray(indices, original.Content[index])
		if pos < 0 {
			clonedKey := original.Content[index].Copy()
			clonedValue := original.Content[index+1].Copy()
			filteredContent = append(filteredContent, clonedKey, clonedValue)
		}
	}
	result := original.CopyWithoutContent()
	result.AddChildren(filteredContent)
	return result
}

func omitSequence(original *CandidateNode, indices *CandidateNode) *CandidateNode {
	filteredContent := make([]*CandidateNode, 0, max(0, len(original.Content)-len(indices.Content)))

	for index := 0; index < len(original.Content); index++ {
		pos := findInArray(indices, createScalarNode(index, strconv.Itoa(index)))
		if pos < 0 {
			filteredContent = append(filteredContent, original.Content[index].Copy())
		}
	}
	result := original.CopyWithoutContent()
	result.AddChildren(filteredContent)
	return result
}

func omitOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("Omit")

	contextIndicesToOmit, err := d.GetMatchingNodes(context, expressionNode.RHS)

	if err != nil {
		return Context{}, err
	}
	indicesToOmit := &CandidateNode{}
	if contextIndicesToOmit.MatchingNodes.Len() > 0 {
		indicesToOmit = contextIndicesToOmit.MatchingNodes.Front().Value.(*CandidateNode)
	}
	if len(indicesToOmit.Content) == 0 {
		log.Debugf("No omit indices specified")
		return context, nil
	}
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		node := el.Value.(*CandidateNode)

		var replacement *CandidateNode

		switch node.Kind {
		case MappingNode:
			replacement = omitMap(node, indicesToOmit)
		case SequenceNode:
			replacement = omitSequence(node, indicesToOmit)
		default:
			log.Debugf("Omit from type %v (%v) is noop", node.Tag, node.GetNicePath())
			return context, nil
		}
		replacement.LeadingContent = node.LeadingContent
		results.PushBack(replacement)
	}
	return context.ChildContext(results), nil
}
