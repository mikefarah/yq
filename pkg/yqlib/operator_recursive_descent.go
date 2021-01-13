package yqlib

import (
	"container/list"

	yaml "gopkg.in/yaml.v3"
)

type recursiveDescentPreferences struct {
	TraversePreferences traversePreferences
	RecurseArray        bool
}

func recursiveDescentOperator(d *dataTreeNavigator, matchMap *list.List, expressionNode *ExpressionNode) (*list.List, error) {
	var results = list.New()

	preferences := expressionNode.Operation.Preferences.(*recursiveDescentPreferences)
	err := recursiveDecent(d, results, matchMap, preferences)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func recursiveDecent(d *dataTreeNavigator, results *list.List, matchMap *list.List, preferences *recursiveDescentPreferences) error {
	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		candidate.Node = unwrapDoc(candidate.Node)

		log.Debugf("Recursive Decent, added %v", NodeToString(candidate))
		results.PushBack(candidate)

		if candidate.Node.Kind != yaml.AliasNode && len(candidate.Node.Content) > 0 &&
			(preferences.RecurseArray || candidate.Node.Kind != yaml.SequenceNode) {

			children, err := splat(d, nodeToMap(candidate), preferences.TraversePreferences)

			if err != nil {
				return err
			}
			err = recursiveDecent(d, results, children, preferences)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
