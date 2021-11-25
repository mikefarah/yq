package yqlib

import (
	"container/list"
)

func selectOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- selectOperation")
	results := list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(candidate), expressionNode.Rhs)
		if err != nil {
			return Context{}, err
		}

		// find any truthy node
		var errDecoding error
		includeResult := false

		for resultEl := rhs.MatchingNodes.Front(); resultEl != nil; resultEl = resultEl.Next() {
			result := resultEl.Value.(*CandidateNode)
			includeResult, errDecoding = isTruthy(result)
			log.Debugf("isTruthy %v", includeResult)
			if errDecoding != nil {
				return Context{}, errDecoding
			}
			if includeResult {
				break
			}
		}

		if includeResult {
			results.PushBack(candidate)
		}
	}
	return context.ChildContext(results), nil
}
