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
		owner := lhs

		if lhs == nil && rhs == nil {
			owner = &CandidateNode{}
		} else if lhs == nil {
			owner = rhs
		}

		var errDecoding error
		lhsTrue := false
		if lhs != nil {
			lhs.Node = unwrapDoc(lhs.Node)
			lhsTrue, errDecoding = isTruthy(lhs)

			if errDecoding != nil {
				return nil, errDecoding
			}
		}
		log.Debugf("-- lhsTrue", lhsTrue)

		rhsTrue := false
		if rhs != nil {
			rhs.Node = unwrapDoc(rhs.Node)
			rhsTrue, errDecoding = isTruthy(rhs)
			if errDecoding != nil {
				return nil, errDecoding
			}
		}
		log.Debugf("-- rhsTrue", rhsTrue)

		return createBooleanCandidate(owner, op(lhsTrue, rhsTrue)), nil
	}
}

func orOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- orOp")
	return crossFunction(d, context, expressionNode, performBoolOp(
		func(b1 bool, b2 bool) bool {
			log.Debugf("-- peformingOrOp with %v and %v", b1, b2)
			return b1 || b2
		}), true)
}

func andOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- AndOp")
	return crossFunction(d, context, expressionNode, performBoolOp(
		func(b1 bool, b2 bool) bool {
			return b1 && b2
		}), true)
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
