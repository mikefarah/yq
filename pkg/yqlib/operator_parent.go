package yqlib

import "container/list"

type parentOpPreferences struct {
	Level int
}

func getParentsOperator(_ *dataTreeNavigator, context Context, _ *ExpressionNode) (Context, error) {
	log.Debugf("getParentsOperator")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		parentsList := &CandidateNode{Kind: SequenceNode, Tag: "!!seq"}
		parent := candidate.Parent
		for parent != nil {
			parentsList.AddChild(parent)
			parent = parent.Parent
		}
		results.PushBack(parentsList)
	}

	return context.ChildContext(results), nil

}

func getParentOperator(_ *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("getParentOperator")

	var results = list.New()

	prefs := expressionNode.Operation.Preferences.(parentOpPreferences)

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		// Handle negative levels: count total parents first
		levelsToGoUp := prefs.Level
		if prefs.Level < 0 {
			// Count all parents
			totalParents := 0
			temp := candidate.Parent
			for temp != nil {
				totalParents++
				temp = temp.Parent
			}
			// Convert negative index to positive
			// -1 means last parent (root), -2 means second to last, etc.
			levelsToGoUp = totalParents + prefs.Level + 1
			if levelsToGoUp < 0 {
				levelsToGoUp = 0
			}
		}

		currentLevel := 0
		for currentLevel < levelsToGoUp && candidate != nil {
			log.Debugf("currentLevel: %v, desired: %v", currentLevel, levelsToGoUp)
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
