package yqlib

func assignUpdateFunc(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	rhs.Node = unwrapDoc(rhs.Node)
	lhs.UpdateFrom(rhs)
	return lhs, nil
}

func assignUpdateOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	lhs, err := d.GetMatchingNodes(context, expressionNode.Lhs)
	if err != nil {
		return Context{}, err
	}

	if !expressionNode.Operation.UpdateAssign {
		// this works because we already ran against LHS with an editable context.
		_, err := crossFunction(d, context.ReadOnlyClone(), expressionNode, assignUpdateFunc, false)
		return context, err
	}

	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		rhs, err := d.GetMatchingNodes(context.SingleChildContext(candidate), expressionNode.Rhs)
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
