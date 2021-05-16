package yqlib

import (
	"container/list"
)

func selectOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("-- selectOperation")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		childContext := context.SingleChildContext(candidate)
		childContext.DontAutoCreate = true
		rhs, err := d.GetMatchingNodes(childContext, expressionNode.Rhs)

		if err != nil {
			return Context{}, err
		}

		// grab the first value
		first := rhs.MatchingNodes.Front()

		if first != nil {
			result := first.Value.(*CandidateNode)
			log.Debugf("result %v", NodeToString(result))
			includeResult, errDecoding := isTruthy(result)
			log.Debugf("isTruthy %v", includeResult)
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
