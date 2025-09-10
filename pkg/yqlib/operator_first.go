package yqlib

import "container/list"

func firstOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	results := list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		// If no RHS expression is provided, simply return the first entry in candidate.Content
		if expressionNode == nil || expressionNode.RHS == nil {
			if len(candidate.Content) > 0 {
				results.PushBack(candidate.Content[0])
			}
			continue
		}

		splatted, err := splat(context.SingleChildContext(candidate), traversePreferences{})
		if err != nil {
			return Context{}, err
		}

		for splatEl := splatted.MatchingNodes.Front(); splatEl != nil; splatEl = splatEl.Next() {
			splatCandidate := splatEl.Value.(*CandidateNode)
			// Create a new context for this splatted candidate
			splatContext := context.SingleChildContext(splatCandidate)
			// Evaluate the RHS expression against this splatted candidate
			rhs, err := d.GetMatchingNodes(splatContext, expressionNode.RHS)
			if err != nil {
				return Context{}, err
			}

			includeResult := false

			for resultEl := rhs.MatchingNodes.Front(); resultEl != nil; resultEl = resultEl.Next() {
				result := resultEl.Value.(*CandidateNode)
				includeResult = isTruthyNode(result)
				if includeResult {
					break
				}
			}
			if includeResult {
				results.PushBack(splatCandidate)
				break
			}
		}

	}
	return context.ChildContext(results), nil
}
