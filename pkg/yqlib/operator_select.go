package yqlib

import (
	"container/list"
)

func selectOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("selectOperation")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(candidate), expressionNode.RHS)

		if err != nil {
			return Context{}, err
		}

		// find any truthy node
		includeResult := false

		for resultEl := rhs.MatchingNodes.Front(); resultEl != nil; resultEl = resultEl.Next() {
			result := resultEl.Value.(*CandidateNode)
			includeResult = isTruthyNode(result)
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
