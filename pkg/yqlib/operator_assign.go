package yqlib

type assignPreferences struct {
	DontOverWriteAnchor bool
	OnlyWriteNull       bool
	ClobberCustomTags   bool
}

func assignUpdateFunc(prefs assignPreferences) crossFunctionCalculation {
	return func(_ *dataTreeNavigator, _ Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
		if !prefs.OnlyWriteNull || lhs.Tag == "!!null" {
			lhs.UpdateFrom(rhs, prefs)
		}
		return lhs, nil
	}
}

// they way *= (multipleAssign) is handled, we set the multiplePrefs
// on the node, not assignPrefs. Long story.
func getAssignPreferences(preferences interface{}) assignPreferences {
	prefs := assignPreferences{}

	switch typedPref := preferences.(type) {
	case assignPreferences:
		prefs = typedPref
	case multiplyPreferences:
		prefs = typedPref.AssignPrefs
	}
	return prefs
}

func assignUpdateOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	lhs, err := d.GetMatchingNodes(context, expressionNode.LHS)
	if err != nil {
		return Context{}, err
	}

	prefs := getAssignPreferences(expressionNode.Operation.Preferences)

	log.Debug("assignUpdateOperator prefs: %v", prefs)

	if !expressionNode.Operation.UpdateAssign {
		// this works because we already ran against LHS with an editable context.
		_, err := crossFunction(d, context.ReadOnlyClone(), expressionNode, assignUpdateFunc(prefs), false)
		return context, err
	}

	//traverse backwards through the context -
	// like delete, we need to run against the children first.
	// (e.g. consider when running with expression '.. |= [.]' - we need
	// to wrap the children first
	for el := lhs.MatchingNodes.Back(); el != nil; el = el.Prev() {
		candidate := el.Value.(*CandidateNode)

		rhs, err := d.GetMatchingNodes(context.SingleChildContext(candidate), expressionNode.RHS)

		if err != nil {
			return Context{}, err
		}

		// grab the first value
		first := rhs.MatchingNodes.Front()

		if first != nil {
			rhsCandidate := first.Value.(*CandidateNode)
			candidate.UpdateFrom(rhsCandidate, prefs)
		}
	}

	return context, nil
}

// does not update content or values
func assignAttributesOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debug("getting lhs matching nodes for update")
	lhs, err := d.GetMatchingNodes(context, expressionNode.LHS)
	if err != nil {
		return Context{}, err
	}
	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(candidate), expressionNode.RHS)

		if err != nil {
			return Context{}, err
		}

		// grab the first value
		first := rhs.MatchingNodes.Front()

		if first != nil {
			prefs := assignPreferences{}
			if expressionNode.Operation.Preferences != nil {
				prefs = expressionNode.Operation.Preferences.(assignPreferences)
			}
			if !prefs.OnlyWriteNull || candidate.Tag == "!!null" {
				candidate.UpdateAttributesFrom(first.Value.(*CandidateNode), prefs)
			}
		}
	}
	return context, nil
}
