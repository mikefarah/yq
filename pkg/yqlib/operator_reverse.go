package yqlib

import (
	"container/list"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

func reverseOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	results := list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		candidateNode := unwrapDoc(candidate.Node)

		if candidateNode.Kind != yaml.SequenceNode {
			return context, fmt.Errorf("node at path [%v] is not an array (it's a %v)", candidate.GetNicePath(), candidate.GetNiceTag())
		}

		reverseList := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Style: candidateNode.Style}
		reverseList.Content = make([]*yaml.Node, len(candidateNode.Content))

		for i, originalNode := range candidateNode.Content {
			reverseList.Content[len(candidateNode.Content)-i-1] = originalNode
		}
		results.PushBack(candidate.CreateReplacement(reverseList))

	}

	return context.ChildContext(results), nil

}
