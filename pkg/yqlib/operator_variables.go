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

func useWithPipe(d *dataTreeNavigator, context Context, originalExp *ExpressionNode) (Context, error) {
	return Context{}, fmt.Errorf("must use variable with a pipe, e.g. `exp as $x | ...`")
}

func variableLoop(d *dataTreeNavigator, context Context, originalExp *ExpressionNode) (Context, error) {
	log.Debug("variable loop!")
	variableExp := originalExp.LHS
	lhs, err := d.GetMatchingNodes(context.ReadOnlyClone(), variableExp.LHS)
	if err != nil {
		return Context{}, err
	}
	if variableExp.RHS.Operation.OperationType.Type != "GET_VARIABLE" {
		return Context{}, fmt.Errorf("RHS of 'as' operator must be a variable name e.g. $foo")
	}
	variableName := variableExp.RHS.Operation.StringValue

	prefs := variableExp.Operation.Preferences.(assignVarPreferences)

	results := list.New()

	// now we loop over lhs, set variable to each result and calculate originalExp.Rhs
	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		var variableValue = list.New()
		if prefs.IsReference {
			variableValue.PushBack(el.Value)
		} else {
			copy, err := el.Value.(*CandidateNode).Copy()
			if err != nil {
				return Context{}, err
			}
			variableValue.PushBack(copy)
		}
		log.Debug("PROCESSING VARIABLE: ", NodeToString(el.Value.(*CandidateNode)))
		newContext := context.ChildContext(context.MatchingNodes)
		newContext.SetVariable(variableName, variableValue)

		rhs, err := d.GetMatchingNodes(newContext, originalExp.RHS)
		if err != nil {
			return Context{}, err
		}
		results.PushBackList(rhs.MatchingNodes)
	}

	// if there is no LHS - then I guess we just calculate originalExp.Rhs
	if lhs.MatchingNodes.Len() == 0 {
		return d.GetMatchingNodes(context, originalExp.RHS)
	}

	return context.ChildContext(results), nil

}
