package yqlib

import (
	"container/list"

	yaml "gopkg.in/yaml.v3"
)

func kindToText(kind yaml.Kind) string {
	switch kind {
	case yaml.MappingNode:
		return "map"
	case yaml.SequenceNode:
		return "seq"
	case yaml.DocumentNode:
		return "doc"
	case yaml.ScalarNode:
		return "scalar"
	case yaml.AliasNode:
		return "alias"
	default:
		return "unknown"
	}
}

func getKindOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("GetKindOperator")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := &yaml.Node{Kind: yaml.ScalarNode, Value: kindToText(candidate.Node.Kind), Tag: "!!str"}
		result := candidate.CreateReplacement(node)
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}
