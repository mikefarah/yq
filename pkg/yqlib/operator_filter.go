package yqlib

import (
	"errors"
	"container/list"
)

func filterOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- filterOperation")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		children := context.SingleChildContext(candidate)
		splatted, err := splat(children, traversePreferences{})
		if err != nil {
			return Context{}, err
		}

		if err != nil {
			return Context{}, err
		}

		for resultEl := splatted.MatchingNodes.Front(); resultEl != nil; resultEl = resultEl.Next() {
			result := resultEl.Value.(*CandidateNode)
			childCtx := context.SingleReadonlyChildContext(result)
			include, err := d.GetMatchingNodes(childCtx, expressionNode.RHS)
			if err != nil {
				return Context{}, err
			}
			var includeResult bool
			var errDecoding error
			includeEl := include.MatchingNodes.Front()
			if includeEl.Next() != nil {
				return Context{}, errors.New("Only expected one child")
			}
			includeVal := includeEl.Value.(*CandidateNode)
			includeResult, errDecoding = isTruthy(includeVal)
			if errDecoding != nil {
				return Context{}, errDecoding
			}
			log.Debug("isTruthy %v", includeResult)
			if includeResult {
				selfExpression := &ExpressionNode{Operation: &Operation{OperationType: selfReferenceOpType}}
				collected, err := collectTogether(d, childCtx, selfExpression)
				if err != nil {
					return Context{}, err
				}
				collected.Node.Style = unwrapDoc(result.Node).Style
				results.PushBack(collected)
			}
		}
	}
	return context.ChildContext(results), nil
}

