package yqlib

func alternativeOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("alternative")
	prefs := crossFunctionPreferences{
		CalcWhenEmpty: true,
		Calculation:   alternativeFunc,
		LhsResultValue: func(lhs *CandidateNode) (*CandidateNode, error) {
			if lhs == nil {
				return nil, nil
			}
			truthy := isTruthyNode(lhs)
			if truthy {
				return lhs, nil
			}
			return nil, nil
		},
	}
	return crossFunctionWithPrefs(d, context, expressionNode, prefs)
}

func alternativeFunc(_ *dataTreeNavigator, _ Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	if lhs == nil {
		return rhs, nil
	}
	if rhs == nil {
		return lhs, nil
	}

	isTrue := isTruthyNode(lhs)
	if isTrue {
		return lhs, nil
	}
	return rhs, nil
}
