package yqlib

import (
	"container/list"
	"fmt"
)

func lengthOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- lengthOperation")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		targetNode := candidate.unwrapDocument()
		var length int
		switch targetNode.Kind {
		case ScalarNode:
			if targetNode.Tag == "!!null" {
				length = 0
			} else {
				length = len(targetNode.Value)
			}
		case MappingNode:
			length = len(targetNode.Content) / 2
		case SequenceNode:
			length = len(targetNode.Content)
		default:
			length = 0
		}

		result := candidate.CreateReplacement()
		result.Kind = ScalarNode
		result.Value = fmt.Sprintf("%v", length)
		result.Tag = "!!int"
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}
