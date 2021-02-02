package yqlib

import (
	"container/list"
)

func selectOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("-- selectOperation")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		rhs, err := d.GetMatchingNodes(context.SingleChildContext(candidate), expressionNode.Rhs)

		if err != nil {
			return Context{}, err
		}

		// grab the first value
		first := rhs.MatchingNodes.Front()

		if first != nil {
			result := first.Value.(*CandidateNode)
			includeResult, errDecoding := isTruthy(result)
			if errDecoding != nil {
				return Context{}, errDecoding
			}

			if includeResult {
				results.PushBack(candidate)
			}
		}
	}
	return context.ChildContext(results), nil
}
