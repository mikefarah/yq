package yqlib

import (
	"container/list"
	"fmt"

	"gopkg.in/yaml.v3"
)

func getDocumentIndexOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%v", candidate.Document), Tag: "!!int"}
		scalar := candidate.CreateReplacement(node)
		results.PushBack(scalar)
	}
	return context.ChildContext(results), nil
}
