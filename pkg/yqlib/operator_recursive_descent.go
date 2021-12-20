package yqlib

import (
	"container/list"

	yaml "gopkg.in/yaml.v3"
)

type recursiveDescentPreferences struct {
	TraversePreferences traversePreferences
	RecurseArray        bool
}

func recursiveDescentOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
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

		candidate.Node = unwrapDoc(candidate.Node)

		log.Debugf("Recursive Decent, added %v", NodeToString(candidate))
		results.PushBack(candidate)

		if candidate.Node.Kind != yaml.AliasNode && len(candidate.Node.Content) > 0 &&
			(preferences.RecurseArray || candidate.Node.Kind != yaml.SequenceNode) {

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
