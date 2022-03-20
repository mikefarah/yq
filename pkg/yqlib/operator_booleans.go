package yqlib

import (
	"container/list"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

func isTruthyNode(node *yaml.Node) (bool, error) {
	value := true
	if node.Tag == "!!null" {
		return false, nil
	}
	if node.Kind == yaml.ScalarNode && node.Tag == "!!bool" {
		errDecoding := node.Decode(&value)
		if errDecoding != nil {
			return false, errDecoding
		}

	}
	return value, nil
}

func isTruthy(c *CandidateNode) (bool, error) {
	node := unwrapDoc(c.Node)
	return isTruthyNode(node)
}

func getBoolean(candidate *CandidateNode) (bool, error) {
	if candidate != nil {
		candidate.Node = unwrapDoc(candidate.Node)
		return isTruthy(candidate)
	}
	return false, nil
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

func returnRhsTruthy(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	owner := getOwner(lhs, rhs)
	rhsBool, err := getBoolean(rhs)
	if err != nil {
		return nil, err
	}

	return createBooleanCandidate(owner, rhsBool), nil
}

func returnLHSWhen(targetBool bool) func(lhs *CandidateNode) (*CandidateNode, error) {
	return func(lhs *CandidateNode) (*CandidateNode, error) {
		var err error
		var lhsBool bool

		if lhsBool, err = getBoolean(lhs); err != nil || lhsBool != targetBool {
			return nil, err
		}
		owner := &CandidateNode{}
		if lhs != nil {
			owner = lhs
		}
		return createBooleanCandidate(owner, targetBool), nil
	}
}

func findBoolean(wantBool bool, d *dataTreeNavigator, context Context, expressionNode *ExpressionNode, sequenceNode *yaml.Node) (bool, error) {
	for _, node := range sequenceNode.Content {

		if expressionNode != nil {
			//need to evaluate the expression against the node
			candidate := &CandidateNode{Node: node}
			rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(candidate), expressionNode)
			if err != nil {
				return false, err
			}
			if rhs.MatchingNodes.Len() > 0 {
				node = rhs.MatchingNodes.Front().Value.(*CandidateNode).Node
			} else {
				// no results found, ignore this entry
				continue
			}
		}

		truthy, err := isTruthyNode(node)
		if err != nil {
			return false, err
		}
		if truthy == wantBool {
			return true, nil
		}
	}
	return false, nil
}

func allOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		candidateNode := unwrapDoc(candidate.Node)
		if candidateNode.Kind != yaml.SequenceNode {
			return Context{}, fmt.Errorf("any only supports arrays, was %v", candidateNode.Tag)
		}
		booleanResult, err := findBoolean(false, d, context, expressionNode.RHS, candidateNode)
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
		candidateNode := unwrapDoc(candidate.Node)
		if candidateNode.Kind != yaml.SequenceNode {
			return Context{}, fmt.Errorf("any only supports arrays, was %v", candidateNode.Tag)
		}
		booleanResult, err := findBoolean(true, d, context, expressionNode.RHS, candidateNode)
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

func notOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- notOperation")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		log.Debug("notOperation checking %v", candidate)
		truthy, errDecoding := isTruthy(candidate)
		if errDecoding != nil {
			return Context{}, errDecoding
		}
		result := createBooleanCandidate(candidate, !truthy)
		results.PushBack(result)
	}
	return context.ChildContext(results), nil
}
