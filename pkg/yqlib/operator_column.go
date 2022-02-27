package yqlib

import (
	"container/list"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

func columnOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("columnOperator")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%v", candidate.Node.Column), Tag: "!!int"}
		result := candidate.CreateReplacement(node)
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}
