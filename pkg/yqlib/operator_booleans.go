package yqlib

import (
	"container/list"
	"fmt"
	"strings"
)

func isTruthyNode(node *CandidateNode) bool {
	if node == nil {
		return false
	}
	if node.Tag == "!!null" {
		return false
	}
	if node.Kind == ScalarNode && node.Tag == "!!bool" {
		// yes/y/true/on
		return (strings.EqualFold(node.Value, "y") ||
			strings.EqualFold(node.Value, "yes") ||
			strings.EqualFold(node.Value, "on") ||
			strings.EqualFold(node.Value, "true"))

	}
	return true
}

func getOwner(lhs *CandidateNode, rhs *CandidateNode) *CandidateNode {
	owner := lhs

	if lhs == nil && rhs == nil {
		owner = &CandidateNode{}
	} else if lhs == nil {
		owner = rhs
	}
	return owner
}

func returnRhsTruthy(_ *dataTreeNavigator, _ Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	owner := getOwner(lhs, rhs)
	rhsBool := isTruthyNode(rhs)

	return createBooleanCandidate(owner, rhsBool), nil
}

func returnLHSWhen(targetBool bool) func(lhs *CandidateNode) (*CandidateNode, error) {
	return func(lhs *CandidateNode) (*CandidateNode, error) {
		var err error
		var lhsBool bool

		if lhsBool = isTruthyNode(lhs); lhsBool != targetBool {
			return nil, err
		}
		owner := &CandidateNode{}
		if lhs != nil {
			owner = lhs
		}
		return createBooleanCandidate(owner, targetBool), nil
	}
}

func findBoolean(wantBool bool, d *dataTreeNavigator, context Context, expressionNode *ExpressionNode, sequenceNode *CandidateNode) (bool, error) {
	for _, node := range sequenceNode.Content {

		if expressionNode != nil {
			//need to evaluate the expression against the node
			rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(node), expressionNode)
			if err != nil {
				return false, err
			}
			if rhs.MatchingNodes.Len() > 0 {
				node = rhs.MatchingNodes.Front().Value.(*CandidateNode)
			} else {
				// no results found, ignore this entry
				continue
			}
		}

		if isTruthyNode(node) == wantBool {
			return true, nil
		}
	}
	return false, nil
}

func allOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		if candidate.Kind != SequenceNode {
			return Context{}, fmt.Errorf("all only supports arrays, was %v", candidate.Tag)
		}
		booleanResult, err := findBoolean(false, d, context, expressionNode.RHS, candidate)
		if err != nil {
			return Context{}, err
		}
		result := createBooleanCandidate(candidate, !booleanResult)
		results.PushBack(result)
	}
	return context.ChildContext(results), nil
}

func anyOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		if candidate.Kind != SequenceNode {
			return Context{}, fmt.Errorf("any only supports arrays, was %v", candidate.Tag)
		}
		booleanResult, err := findBoolean(true, d, context, expressionNode.RHS, candidate)
		if err != nil {
			return Context{}, err
		}
		result := createBooleanCandidate(candidate, booleanResult)
		results.PushBack(result)
	}
	return context.ChildContext(results), nil
}

func orOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	prefs := crossFunctionPreferences{
		CalcWhenEmpty:  true,
		Calculation:    returnRhsTruthy,
		LhsResultValue: returnLHSWhen(true),
	}
	return crossFunctionWithPrefs(d, context.ReadOnlyClone(), expressionNode, prefs)
}

func andOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	prefs := crossFunctionPreferences{
		CalcWhenEmpty:  true,
		Calculation:    returnRhsTruthy,
		LhsResultValue: returnLHSWhen(false),
	}
	return crossFunctionWithPrefs(d, context.ReadOnlyClone(), expressionNode, prefs)
}

func notOperator(_ *dataTreeNavigator, context Context, _ *ExpressionNode) (Context, error) {
	log.Debugf("notOperation")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		log.Debug("notOperation checking %v", candidate)
		truthy := isTruthyNode(candidate)
		result := createBooleanCandidate(candidate, !truthy)
		results.PushBack(result)
	}
	return context.ChildContext(results), nil
}
