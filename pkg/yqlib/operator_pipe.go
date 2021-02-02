package yqlib

func pipeOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	//lhs may update the variable context, we should pass that into the RHS
	// BUT we still return the original context back (see jq)
	// https://stedolan.github.io/jq/manual/#Variable/SymbolicBindingOperator:...as$identifier|...

	lhs, err := d.GetMatchingNodes(context, expressionNode.Lhs)
	if err != nil {
		return Context{}, err
	}
	rhs, err := d.GetMatchingNodes(lhs, expressionNode.Rhs)
	if err != nil {
		return Context{}, err
	}
	return context.ChildContext(rhs.MatchingNodes), nil
}
