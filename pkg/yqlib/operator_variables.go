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

type assignVarPreferences struct {
	IsReference bool
}

func assignVariableOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	lhs, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.LHS)
	if err != nil {
		return Context{}, err
	}
	if expressionNode.RHS.Operation.OperationType.Type != "GET_VARIABLE" {
		return Context{}, fmt.Errorf("RHS of 'as' operator must be a variable name e.g. $foo")
	}
	variableName := expressionNode.RHS.Operation.StringValue

	prefs := expressionNode.Operation.Preferences.(assignVarPreferences)

	var variableValue *list.List
	if prefs.IsReference {
		variableValue = lhs.MatchingNodes
	} else {
		variableValue = lhs.DeepClone().MatchingNodes
	}
	context.SetVariable(variableName, variableValue)
	return context, nil
}
