package yqlib

import (
	"container/list"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

func createPathNodeFor(pathElement interface{}) *yaml.Node {
	switch pathElement := pathElement.(type) {
	case string:
		return &yaml.Node{Kind: yaml.ScalarNode, Value: pathElement, Tag: "!!str"}
	default:
		return &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%v", pathElement), Tag: "!!int"}
	}
}

func getPathOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("GetPath")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}

		content := make([]*yaml.Node, len(candidate.Path))
		for pathIndex := 0; pathIndex < len(candidate.Path); pathIndex++ {
			path := candidate.Path[pathIndex]
			content[pathIndex] = createPathNodeFor(path)
		}
		node.Content = content
		result := candidate.CreateReplacement(node)
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}
