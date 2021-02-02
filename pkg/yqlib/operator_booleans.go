package yqlib

import (
	"container/list"

	yaml "gopkg.in/yaml.v3"
)

func isTruthy(c *CandidateNode) (bool, error) {
	node := unwrapDoc(c.Node)
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

type boolOp func(bool, bool) bool

func performBoolOp(op boolOp) func(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	return func(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
		lhs.Node = unwrapDoc(lhs.Node)
		rhs.Node = unwrapDoc(rhs.Node)

		lhsTrue, errDecoding := isTruthy(lhs)
		if errDecoding != nil {
			return nil, errDecoding
		}

		rhsTrue, errDecoding := isTruthy(rhs)
		if errDecoding != nil {
			return nil, errDecoding
		}

		return createBooleanCandidate(lhs, op(lhsTrue, rhsTrue)), nil
	}
}

func orOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- orOp")
	return crossFunction(d, context, expressionNode, performBoolOp(
		func(b1 bool, b2 bool) bool {
			return b1 || b2
		}))
}

func andOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- AndOp")
	return crossFunction(d, context, expressionNode, performBoolOp(
		func(b1 bool, b2 bool) bool {
			return b1 && b2
		}))
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
