package yqlib

func valueOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debug("value = %v", expressionNode.Operation.CandidateNode.Node.Value)
	return context.SingleChildContext(expressionNode.Operation.CandidateNode), nil
}
