package yqlib

import "container/list"

func unionOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debug("unionOperator")
	log.Debug("context: %v", NodesToString(context.MatchingNodes))
	lhs, err := d.GetMatchingNodes(context, expressionNode.LHS)
	if err != nil {
		return Context{}, err
	}
	log.Debug("lhs: %v", NodesToString(lhs.MatchingNodes))
	log.Debug("rhs input: %v", NodesToString(context.MatchingNodes))
	log.Debug("rhs: %v", expressionNode.RHS.Operation.toString())
	rhs, err := d.GetMatchingNodes(context, expressionNode.RHS)

	if err != nil {
		return Context{}, err
	}
	log.Debug("lhs: %v", lhs.ToString())
	log.Debug("rhs: %v", rhs.ToString())

	results := lhs.ChildContext(list.New())
	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		node := el.Value.(*CandidateNode)
		results.MatchingNodes.PushBack(node)
	}

	// this can happen when both expressions modify the context
	// instead of creating their own.
	/// (.foo = "bar"), (.thing = "cat")
	if rhs.MatchingNodes != lhs.MatchingNodes {

		for el := rhs.MatchingNodes.Front(); el != nil; el = el.Next() {
			node := el.Value.(*CandidateNode)
			log.Debug("processing %v", NodeToString(node))

			results.MatchingNodes.PushBack(node)
		}
	}
	log.Debug("and lets print it out")
	log.Debug("all together: %v", results.ToString())
	return results, nil
}
