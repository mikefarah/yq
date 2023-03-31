package yqlib

import (
	"container/list"

	yaml "gopkg.in/yaml.v3"
)

func collectTogether(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (*CandidateNode, error) {
	collectedNode := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		collectExpResults, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(candidate), expressionNode)
		if err != nil {
			return nil, err
		}
		for result := collectExpResults.MatchingNodes.Front(); result != nil; result = result.Next() {
			resultC := result.Value.(*CandidateNode)
			log.Debugf("found this: %v", NodeToString(resultC))
			collectedNode.Content = append(collectedNode.Content, unwrapDoc(resultC.Node))
		}
	}
	return &CandidateNode{Node: collectedNode}, nil
}

func collectOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- collectOperation")

	if context.MatchingNodes.Len() == 0 {
		log.Debugf("nothing to collect")
		node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Value: "[]"}
		candidate := &CandidateNode{Node: node}
		return context.SingleChildContext(candidate), nil
	}

	var evaluateAllTogether = true
	for matchEl := context.MatchingNodes.Front(); matchEl != nil; matchEl = matchEl.Next() {
		evaluateAllTogether = evaluateAllTogether && matchEl.Value.(*CandidateNode).EvaluateTogether
		if !evaluateAllTogether {
			break
		}
	}

	if evaluateAllTogether {
		log.Debugf("collect together")
		collectedNode, err := collectTogether(d, context, expressionNode.RHS)
		if err != nil {
			return Context{}, err
		}
		return context.SingleChildContext(collectedNode), nil

	}

	var results = list.New()
	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		collectedNode := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		collectCandidate := candidate.CreateReplacement(collectedNode)

		log.Debugf("collect rhs: %v", expressionNode.RHS.Operation.toString())

		collectExpResults, err := d.GetMatchingNodes(context.SingleChildContext(candidate), expressionNode.RHS)
		if err != nil {
			return Context{}, err
		}

		for result := collectExpResults.MatchingNodes.Front(); result != nil; result = result.Next() {
			resultC := result.Value.(*CandidateNode)
			log.Debugf("found this: %v", NodeToString(resultC))
			collectedNode.Content = append(collectedNode.Content, unwrapDoc(resultC.Node))
		}
		log.Debugf("done collect rhs: %v", expressionNode.RHS.Operation.toString())

		results.PushBack(collectCandidate)
	}

	return context.ChildContext(results), nil
}
