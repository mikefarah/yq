package yqlib

import "fmt"

func withOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- withOperator")
	// with(path, exp)

	if expressionNode.Rhs.Operation.OperationType != blockOpType {
		return Context{}, fmt.Errorf("with must be given a block, got %v instead", expressionNode.Rhs.Operation.OperationType.Type)
	}

	pathExp := expressionNode.Rhs.Lhs

	updateContext, err := d.GetMatchingNodes(context, pathExp)
	if err != nil {
		return Context{}, err
	}

	updateExp := expressionNode.Rhs.Rhs

	_, err = d.GetMatchingNodes(updateContext, updateExp)
	if err != nil {
		return Context{}, err
	}

	return context, nil
}
