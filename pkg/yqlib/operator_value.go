package yqlib

func valueOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	return context.ChildContext(expressionNode.Operation.ValueNodes), nil
}
