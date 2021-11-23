package yqlib

import (
	"container/list"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

func lengthOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- lengthOperation")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		targetNode := unwrapDoc(candidate.Node)
		var length int
		switch targetNode.Kind {
		case yaml.ScalarNode:
			if targetNode.Tag == "!!null" {
				length = 0
			} else {
				length = len(targetNode.Value)
			}
		case yaml.MappingNode:
			length = len(targetNode.Content) / 2
		case yaml.SequenceNode:
			length = len(targetNode.Content)
		default:
			length = 0
		}

		node := &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%v", length), Tag: "!!int"}
		result := candidate.CreateReplacement(node)
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}
