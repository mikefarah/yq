package yqlib

type expressionOpPreferences struct {
	expression string
}

func expressionOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	prefs := expressionNode.Operation.Preferences.(expressionOpPreferences)
	expNode, err := ExpressionParser.ParseExpression(prefs.expression)
	if err != nil {
		return Context{}, err
	}

	return d.GetMatchingNodes(context, expNode)
}
