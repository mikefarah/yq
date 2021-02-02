package yqlib

import (
	"container/list"

	yaml "gopkg.in/yaml.v3"
)

func assignTagOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("AssignTagOperator: %v")
	tag := ""

	if !expressionNode.Operation.UpdateAssign {
		rhs, err := d.GetMatchingNodes(context, expressionNode.Rhs)
		if err != nil {
			return Context{}, err
		}

		if rhs.MatchingNodes.Front() != nil {
			tag = rhs.MatchingNodes.Front().Value.(*CandidateNode).Node.Value
		}
	}

	lhs, err := d.GetMatchingNodes(context, expressionNode.Lhs)

	if err != nil {
		return Context{}, err
	}

	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		log.Debugf("Setting tag of : %v", candidate.GetKey())
		if expressionNode.Operation.UpdateAssign {
			rhs, err := d.GetMatchingNodes(context.SingleChildContext(candidate), expressionNode.Rhs)
			if err != nil {
				return Context{}, err
			}

			if rhs.MatchingNodes.Front() != nil {
				tag = rhs.MatchingNodes.Front().Value.(*CandidateNode).Node.Value
			}
		}
		unwrapDoc(candidate.Node).Tag = tag
	}

	return context, nil
}

func getTagOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("GetTagOperator")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := &yaml.Node{Kind: yaml.ScalarNode, Value: unwrapDoc(candidate.Node).Tag, Tag: "!!str"}
		result := candidate.CreateChild(nil, node)
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}
