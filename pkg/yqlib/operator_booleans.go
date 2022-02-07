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
	log.Debugf("-- orOp")
	return crossFunction(d, context.ReadOnlyClone(), expressionNode, performBoolOp(
		func(b1 bool, b2 bool) bool {
			log.Debugf("-- peformingOrOp with %v and %v", b1, b2)
			return b1 || b2
		}), true)
}

func andOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- AndOp")
	return crossFunction(d, context.ReadOnlyClone(), expressionNode, performBoolOp(
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
