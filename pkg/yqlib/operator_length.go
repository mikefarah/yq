package yqlib

import (
	"container/list"
	"fmt"
)

func lengthOperator(_ *dataTreeNavigator, context Context, _ *ExpressionNode) (Context, error) {
	log.Debugf("lengthOperation")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		var length int
		switch candidate.Kind {
		case ScalarNode:
			if candidate.Tag == "!!null" {
				length = 0
			} else {
				length = len(candidate.Value)
			}
		case MappingNode:
			length = len(candidate.Content) / 2
		case SequenceNode:
			length = len(candidate.Content)
		default:
			length = 0
		}

		result := candidate.CreateReplacement(ScalarNode, "!!int", fmt.Sprintf("%v", length))
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}
