package yqlib

import "container/list"

func unionOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debug("unionOperator--")
	log.Debugf("unionOperator: context: %v", NodesToString(context.MatchingNodes))
	lhs, err := d.GetMatchingNodes(context, expressionNode.LHS)
	if err != nil {
		return Context{}, err
	}
	log.Debugf("unionOperator: lhs: %v", NodesToString(lhs.MatchingNodes))
	log.Debugf("unionOperator: rhs input: %v", NodesToString(context.MatchingNodes))
	log.Debugf("unionOperator: rhs: %v", expressionNode.RHS.Operation.toString())
	rhs, err := d.GetMatchingNodes(context, expressionNode.RHS)

	if err != nil {
		return Context{}, err
	}
	log.Debugf("unionOperator: lhs: %v", lhs.ToString())
	log.Debugf("unionOperator: rhs: %v", rhs.ToString())

	// Both sides may modify and return the same context, rather than selecting
	// separate nodes. Keep its list identity so nested unions also return it once.
	if rhs.MatchingNodes == lhs.MatchingNodes {
		return lhs, nil
	}

	results := lhs.ChildContext(list.New())
	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		node := el.Value.(*CandidateNode)
		results.MatchingNodes.PushBack(node)
	}

	for el := rhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		node := el.Value.(*CandidateNode)
		log.Debugf("union operator rhs: processing %v", NodeToString(node))
		results.MatchingNodes.PushBack(node)
	}
	log.Debugf("union operator: all together: %v", results.ToString())
	return results, nil
}
