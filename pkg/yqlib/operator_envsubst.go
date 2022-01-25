package yqlib

import (
	"container/list"
	"fmt"

	envsubst "github.com/a8m/envsubst"
	yaml "gopkg.in/yaml.v3"
)

func envsubstOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := unwrapDoc(candidate.Node)
		if node.Tag != "!!str" {
			log.Warning("EnvSubstOperator, env name:", node.Tag, node.Value)
			return Context{}, fmt.Errorf("cannot substitute with %v, can only substitute strings. Hint: Most often you'll want to use '|=' over '=' for this operation", node.Tag)
		}

		value, err := envsubst.String(node.Value)
		if err != nil {
			return Context{}, err
		}
		targetNode := &yaml.Node{Kind: yaml.ScalarNode, Value: value, Tag: "!!str"}
		result := candidate.CreateReplacement(targetNode)
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}
