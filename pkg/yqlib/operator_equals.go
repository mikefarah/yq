package yqlib

func equalsOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- equalsOperation")
	return crossFunction(d, context, expressionNode, isEquals(false), true)
}

func isEquals(flip bool) func(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	return func(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
		value := false
		log.Debugf("-- isEquals cross function")
		if lhs == nil && rhs == nil {
			owner := &CandidateNode{}
			return createBooleanCandidate(owner, !flip), nil
		} else if lhs == nil {
			log.Debugf("lhs nil, but rhs is not")
			rhsNode := rhs.unwrapDocument()
			value := rhsNode.Tag == "!!null"
			if flip {
				value = !value
			}
			return createBooleanCandidate(rhs, value), nil
		} else if rhs == nil {
			log.Debugf("lhs not nil, but rhs is")
			lhsNode := lhs.unwrapDocument()
			value := lhsNode.Tag == "!!null"
			if flip {
				value = !value
			}
			return createBooleanCandidate(lhs, value), nil
		}

		lhsNode := lhs.unwrapDocument()
		rhsNode := rhs.unwrapDocument()

		if lhsNode.Tag == "!!null" {
			value = (rhsNode.Tag == "!!null")
		} else if lhsNode.Kind == ScalarNode && rhsNode.Kind == ScalarNode {
			value = matchKey(lhsNode.Value, rhsNode.Value)
		}
		log.Debugf("%v == %v ? %v", NodeToString(lhs), NodeToString(rhs), value)
		if flip {
			value = !value
		}
		return createBooleanCandidate(lhs, value), nil
	}
}

func notEqualsOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- notEqualsOperator")
	return crossFunction(d, context.ReadOnlyClone(), expressionNode, isEquals(true), true)
}
