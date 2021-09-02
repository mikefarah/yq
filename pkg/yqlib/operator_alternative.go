package yqlib

func alternativeOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- alternative")
	return crossFunction(d, context, expressionNode, alternativeFunc, true)
}

func alternativeFunc(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	if lhs == nil {
		return rhs, nil
	}
	if rhs == nil {
		return lhs, nil
	}
	lhs.Node = unwrapDoc(lhs.Node)
	rhs.Node = unwrapDoc(rhs.Node)
	log.Debugf("Alternative LHS: %v", lhs.Node.Tag)
	log.Debugf("-          RHS: %v", rhs.Node.Tag)

	isTrue, err := isTruthy(lhs)
	if err != nil {
		return nil, err
	} else if isTrue {
		return lhs, nil
	}
	return rhs, nil
}
