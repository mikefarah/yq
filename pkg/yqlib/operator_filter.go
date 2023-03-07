package yqlib

import (
	"container/list"
)

func filterOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Warningf("-- filterOperation")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		log.Warningf("candidate %#v", candidate)
		children := context.SingleChildContext(candidate)
		splatted, err := splat(children, traversePreferences{})
		if err != nil {
			return Context{}, err
		}
		filtered, err := selectOperator(d, splatted, expressionNode)
		if err != nil {
			return Context{}, err
		}
		for resultEl := filtered.MatchingNodes.Front(); resultEl != nil; resultEl = resultEl.Next() {
			result := resultEl.Value.(*CandidateNode)
			log.Warningf("filtered %#v", result)
			selfExpression := &ExpressionNode{Operation: &Operation{OperationType: selfReferenceOpType}}
			childCtx := context.SingleReadonlyChildContext(result)
			collected, err := collectTogether(d, childCtx, selfExpression)
			if err != nil {
				return Context{}, err
			}
			collected.Node.Style = unwrapDoc(result.Node).Style
			results.PushBack(collected)
		}
	}
	return context.ChildContext(results), nil
}

