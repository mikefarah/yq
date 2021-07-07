package yqlib

func assignUpdateOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	lhs, err := d.GetMatchingNodes(context, expressionNode.Lhs)
	if err != nil {
		return Context{}, err
	}
	var rhs Context
	if !expressionNode.Operation.UpdateAssign {
		rhs, err = d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.Rhs)
		if err != nil {
			return Context{}, err
		}
	}

	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		if expressionNode.Operation.UpdateAssign {
			rhs, err = d.GetMatchingNodes(context.SingleChildContext(candidate), expressionNode.Rhs)
		}

		if err != nil {
			return Context{}, err
		}

		// grab the first value
		first := rhs.MatchingNodes.Front()

		if first != nil {
			rhsCandidate := first.Value.(*CandidateNode)
			rhsCandidate.Node = unwrapDoc(rhsCandidate.Node)
			candidate.UpdateFrom(rhsCandidate)
		}
	}

	return context, nil
}

// does not update content or values
func assignAttributesOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debug("getting lhs matching nodes for update")
	lhs, err := d.GetMatchingNodes(context, expressionNode.Lhs)
	if err != nil {
		return Context{}, err
	}
	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(candidate), expressionNode.Rhs)

		if err != nil {
			return Context{}, err
		}

		// grab the first value
		first := rhs.MatchingNodes.Front()

		if first != nil {
			candidate.UpdateAttributesFrom(first.Value.(*CandidateNode))
		}
	}
	return context, nil
}
