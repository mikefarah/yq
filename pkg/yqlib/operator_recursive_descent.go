package yqlib

import (
	"container/list"
)

type recursiveDescentPreferences struct {
	TraversePreferences traversePreferences
	RecurseArray        bool
}

func recursiveDescentOperator(_ *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	var results = list.New()

	preferences := expressionNode.Operation.Preferences.(recursiveDescentPreferences)
	err := recursiveDecent(results, context, preferences)
	if err != nil {
		return Context{}, err
	}

	return context.ChildContext(results), nil
}

func recursiveDecent(results *list.List, context Context, preferences recursiveDescentPreferences) error {
	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		log.Debugf("added %v", NodeToString(candidate))
		results.PushBack(candidate)

		if candidate.Kind != AliasNode && len(candidate.Content) > 0 &&
			(preferences.RecurseArray || candidate.Kind != SequenceNode) {

			children, err := splat(context.SingleChildContext(candidate), preferences.TraversePreferences)

			if err != nil {
				return err
			}
			err = recursiveDecent(results, children, preferences)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
