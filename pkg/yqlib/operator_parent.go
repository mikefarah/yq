package yqlib

import "container/list"

type parentOpPreferences struct {
	Level int
}

func getParentOperator(_ *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("getParentOperator")

	var results = list.New()

	prefs := expressionNode.Operation.Preferences.(parentOpPreferences)

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		currentLevel := 0
		for currentLevel < prefs.Level && candidate != nil {
			log.Debugf("currentLevel: %v, desired: %v", currentLevel, prefs.Level)
			log.Debugf("candidate: %v", NodeToString(candidate))
			candidate = candidate.Parent
			currentLevel++
		}

		log.Debugf("found candidate: %v", NodeToString(candidate))
		if candidate != nil {
			results.PushBack(candidate)
		}
	}

	return context.ChildContext(results), nil

}
