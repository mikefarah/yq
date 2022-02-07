package yqlib

func pipeOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	//lhs may update the variable context, we should pass that into the RHS
	// BUT we still return the original context back (see jq)
	// https://stedolan.github.io/jq/manual/#Variable/SymbolicBindingOperator:...as$identifier|...

	lhs, err := d.GetMatchingNodes(context, expressionNode.LHS)
	if err != nil {
		return Context{}, err
	}
	rhs, err := d.GetMatchingNodes(lhs, expressionNode.RHS)
	if err != nil {
		return Context{}, err
	}
	return context.ChildContext(rhs.MatchingNodes), nil
}
