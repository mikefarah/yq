package yqlib

import (
	"container/list"
	"fmt"
)

func reverseOperator(_ *dataTreeNavigator, context Context, _ *ExpressionNode) (Context, error) {
	results := list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		if candidate.Kind != SequenceNode {
			return context, fmt.Errorf("node at path [%v] is not an array (it's a %v)", candidate.GetNicePath(), candidate.Tag)
		}

		reverseList := candidate.CreateReplacementWithComments(SequenceNode, "!!seq", candidate.Style)
		reverseContent := make([]*CandidateNode, len(candidate.Content))

		for i, originalNode := range candidate.Content {
			reverseContent[len(candidate.Content)-i-1] = originalNode
		}
		reverseList.AddChildren(reverseContent)
		results.PushBack(reverseList)

	}

	return context.ChildContext(results), nil

}
