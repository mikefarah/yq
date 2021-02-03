package yqlib

import (
	"container/list"
	"fmt"
)

func getVariableOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	variableName := expressionNode.Operation.StringValue
	log.Debug("getVariableOperator %v", variableName)
	result := context.GetVariable(variableName)
	if result == nil {
		result = list.New()
	}
	return context.ChildContext(result), nil
}

func assignVariableOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	lhs, err := d.GetMatchingNodes(context, expressionNode.Lhs)
	if err != nil {
		return Context{}, nil
	}
	if expressionNode.Rhs.Operation.OperationType.Type != "GET_VARIABLE" {
		return Context{}, fmt.Errorf("RHS of 'as' operator must be a variable name e.g. $foo")
	}
	variableName := expressionNode.Rhs.Operation.StringValue
	context.SetVariable(variableName, lhs.MatchingNodes)
	return context, nil
}
