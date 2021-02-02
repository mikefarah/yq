package yqlib

import (
	"container/list"

	yaml "gopkg.in/yaml.v3"
)

func collectOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- collectOperation")

	if context.MatchingNodes.Len() == 0 {
		node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Value: "[]"}
		candidate := &CandidateNode{Node: node}
		return context.SingleChildContext(candidate), nil
	}

	var results = list.New()

	node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	var collectC *CandidateNode
	if context.MatchingNodes.Front() != nil {
		collectC = context.MatchingNodes.Front().Value.(*CandidateNode).CreateChild(nil, node)
		if len(collectC.Path) > 0 {
			collectC.Path = collectC.Path[:len(collectC.Path)-1]
		}
	} else {
		collectC = &CandidateNode{Node: node}
	}

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		log.Debugf("Collecting %v", NodeToString(candidate))
		node.Content = append(node.Content, unwrapDoc(candidate.Node))
	}

	results.PushBack(collectC)

	return context.ChildContext(results), nil
}
