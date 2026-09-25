package yqlib

import (
	"container/list"
)

// if C then A end - there is no else, so falsy conditions return the input unchanged (like jq 1.7)
func ifThenOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("ifThenOperator")
	return ifThenElse(d, context, expressionNode.LHS, expressionNode.RHS, nil)
}

// if C then A else B end, and elif chains: the LHS is always the `then` node, the RHS is the else branch.
func ifElseOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("ifElseOperator")
	thenNode := expressionNode.LHS
	return ifThenElse(d, context, thenNode.LHS, thenNode.RHS, expressionNode.RHS)
}

func ifThenElse(d *dataTreeNavigator, context Context, conditionExp *ExpressionNode, thenExp *ExpressionNode, elseExp *ExpressionNode) (Context, error) {
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		conditions, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(candidate), conditionExp)
		if err != nil {
			return Context{}, err
		}

		// a condition with no results (e.g. a missing key) is treated as null, like `and`, `or` and `//`
		if conditions.MatchingNodes.Len() == 0 {
			conditions.MatchingNodes.PushBack((*CandidateNode)(nil))
		}

		// like jq, each result of the condition produces a result
		for conditionEl := conditions.MatchingNodes.Front(); conditionEl != nil; conditionEl = conditionEl.Next() {
			branchExp := elseExp
			if isTruthyNode(conditionEl.Value.(*CandidateNode)) {
				branchExp = thenExp
			}
			if branchExp == nil {
				results.PushBack(candidate)
				continue
			}
			branch, err := d.GetMatchingNodes(context.SingleChildContext(candidate), branchExp)
			if err != nil {
				return Context{}, err
			}
			results.PushBackList(branch.MatchingNodes)
		}
	}
	return context.ChildContext(results), nil
}
