package yqlib

import (
	"errors"
)

func errorOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("errorOperation")

	rhs, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.RHS)
	if err != nil {
		return Context{}, err
	}
	errorMessage := "aborted"
	if rhs.MatchingNodes.Len() > 0 {
		errorMessage = rhs.MatchingNodes.Front().Value.(*CandidateNode).Value
	}
	return Context{}, errors.New(errorMessage)
}
