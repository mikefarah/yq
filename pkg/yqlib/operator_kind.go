package yqlib

import (
	"container/list"
)

func kindToText(kind Kind) string {
	switch kind {
	case MappingNode:
		return "map"
	case SequenceNode:
		return "seq"
	case ScalarNode:
		return "scalar"
	case AliasNode:
		return "alias"
	default:
		return "unknown"
	}
}

func getKindOperator(_ *dataTreeNavigator, context Context, _ *ExpressionNode) (Context, error) {
	log.Debugf("GetKindOperator")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		result := candidate.CreateReplacement(ScalarNode, "!!str", kindToText(candidate.Kind))
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}
