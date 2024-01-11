package yqlib

import (
	"container/list"
	"fmt"
)

func getVariableOperator(_ *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
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

func useWithPipe(_ *dataTreeNavigator, _ Context, _ *ExpressionNode) (Context, error) {
	return Context{}, fmt.Errorf("must use variable with a pipe, e.g. `exp as $x | ...`")
}

// variables are like loops in jq
// https://stedolan.github.io/jq/manual/#Variable
func variableLoop(d *dataTreeNavigator, context Context, originalExp *ExpressionNode) (Context, error) {
	log.Debug("variable loop!")
	results := list.New()
	var evaluateAllTogether = true
	for matchEl := context.MatchingNodes.Front(); matchEl != nil; matchEl = matchEl.Next() {
		evaluateAllTogether = evaluateAllTogether && matchEl.Value.(*CandidateNode).EvaluateTogether
		if !evaluateAllTogether {
			break
		}
	}
	if evaluateAllTogether {
		return variableLoopSingleChild(d, context, originalExp)
	}

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		result, err := variableLoopSingleChild(d, context.SingleChildContext(el.Value.(*CandidateNode)), originalExp)
		if err != nil {
			return Context{}, err
		}
		results.PushBackList(result.MatchingNodes)
	}
	return context.ChildContext(results), nil
}

func variableLoopSingleChild(d *dataTreeNavigator, context Context, originalExp *ExpressionNode) (Context, error) {

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
		log.Debug("PROCESSING VARIABLE: ", NodeToString(el.Value.(*CandidateNode)))
		var variableValue = list.New()
		if prefs.IsReference {
			variableValue.PushBack(el.Value)
		} else {
			candidateCopy := el.Value.(*CandidateNode).Copy()
			variableValue.PushBack(candidateCopy)
		}
		newContext := context.ChildContext(context.MatchingNodes)
		newContext.SetVariable(variableName, variableValue)

		rhs, err := d.GetMatchingNodes(newContext, originalExp.RHS)

		if err != nil {
			return Context{}, err
		}
		log.Debug("PROCESSING VARIABLE DONE, got back: ", rhs.MatchingNodes.Len())
		results.PushBackList(rhs.MatchingNodes)
	}

	// if there is no LHS - then I guess we just calculate originalExp.Rhs
	if lhs.MatchingNodes.Len() == 0 {
		return d.GetMatchingNodes(context, originalExp.RHS)
	}

	return context.ChildContext(results), nil

}
