package yqlib

func equalsOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("equalsOperation")
	return crossFunction(d, context, expressionNode, isEquals(false), true)
}

func isEquals(flip bool) func(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	return func(_ *dataTreeNavigator, _ Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
		value := false
		log.Debugf("isEquals cross function")
		if lhs == nil && rhs == nil {
			log.Debugf("both are nil")
			owner := &CandidateNode{}
			return createBooleanCandidate(owner, !flip), nil
		} else if lhs == nil {
			log.Debugf("lhs nil, but rhs is not")
			value := rhs.Tag == "!!null"
			if flip {
				value = !value
			}
			return createBooleanCandidate(rhs, value), nil
		} else if rhs == nil {
			log.Debugf("lhs not nil, but rhs is")
			value := lhs.Tag == "!!null"
			if flip {
				value = !value
			}
			return createBooleanCandidate(lhs, value), nil
		}

		if lhs.Tag == "!!null" {
			value = (rhs.Tag == "!!null")
		} else if lhs.Kind == ScalarNode && rhs.Kind == ScalarNode {
			value = matchKey(lhs.Value, rhs.Value)
		}
		log.Debugf("%v == %v ? %v", NodeToString(lhs), NodeToString(rhs), value)
		if flip {
			value = !value
		}
		return createBooleanCandidate(lhs, value), nil
	}
}

func notEqualsOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("notEqualsOperator")
	return crossFunction(d, context.ReadOnlyClone(), expressionNode, isEquals(true), true)
}
