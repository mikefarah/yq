package yqlib

import "fmt"

func withOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("withOperator")
	// with(path, exp)

	if expressionNode.RHS.Operation.OperationType != blockOpType {
		return Context{}, fmt.Errorf("with must be given a block (;), got %v instead", expressionNode.RHS.Operation.OperationType.Type)
	}

	pathExp := expressionNode.RHS.LHS

	updateContext, err := d.GetMatchingNodes(context, pathExp)

	if err != nil {
		return Context{}, err
	}

	updateExp := expressionNode.RHS.RHS

	for el := updateContext.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		_, err = d.GetMatchingNodes(updateContext.SingleChildContext(candidate), updateExp)
		if err != nil {
			return Context{}, err
		}

	}

	return context, nil

}
