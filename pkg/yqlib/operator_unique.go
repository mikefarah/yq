package yqlib

import (
	"container/list"
	"fmt"

	"github.com/elliotchance/orderedmap"
)

func unique(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	selfExpression := &ExpressionNode{Operation: &Operation{OperationType: selfReferenceOpType}}
	uniqueByExpression := &ExpressionNode{Operation: &Operation{OperationType: uniqueByOpType}, RHS: selfExpression}
	return uniqueBy(d, context, uniqueByExpression)

}

func uniqueBy(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("-- uniqueBy Operator")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		candidateNode := candidate.unwrapDocument()

		if candidateNode.Kind != SequenceNode {
			return Context{}, fmt.Errorf("Only arrays are supported for unique")
		}

		var newMatches = orderedmap.NewOrderedMap()
		for _, child := range candidateNode.Content {
			rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(child), expressionNode.RHS)

			if err != nil {
				return Context{}, err
			}

			keyValue := "null"

			if rhs.MatchingNodes.Len() > 0 {
				first := rhs.MatchingNodes.Front()
				keyCandidate := first.Value.(*CandidateNode)
				keyValue = keyCandidate.Value
			}

			_, exists := newMatches.Get(keyValue)

			if !exists {
				newMatches.Set(keyValue, child)
			}
		}
		resultNode := candidate.CreateReplacementWithDocWrappers(SequenceNode, "!!seq", "")
		for el := newMatches.Front(); el != nil; el = el.Next() {
			resultNode.Content = append(resultNode.Content, el.Value.(*CandidateNode))
		}

		results.PushBack(resultNode)
	}

	return context.ChildContext(results), nil

}
