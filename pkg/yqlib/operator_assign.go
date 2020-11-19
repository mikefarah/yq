package yqlib

import "container/list"

type AssignOpPreferences struct {
	UpdateAssign bool
}

func AssignUpdateOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	lhs, err := d.GetMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	preferences := pathNode.Operation.Preferences.(*AssignOpPreferences)

	var rhs *list.List
	if !preferences.UpdateAssign {
		rhs, err = d.GetMatchingNodes(matchingNodes, pathNode.Rhs)
	}

	for el := lhs.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		if preferences.UpdateAssign {
			rhs, err = d.GetMatchingNodes(nodeToMap(candidate), pathNode.Rhs)
		}

		if err != nil {
			return nil, err
		}

		// grab the first value
		first := rhs.Front()

		if first != nil {
			candidate.UpdateFrom(first.Value.(*CandidateNode))
		}
	}
	return matchingNodes, nil
}

// does not update content or values
func AssignAttributesOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	lhs, err := d.GetMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	for el := lhs.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		rhs, err := d.GetMatchingNodes(nodeToMap(candidate), pathNode.Rhs)

		if err != nil {
			return nil, err
		}

		// grab the first value
		first := rhs.Front()

		if first != nil {
			candidate.UpdateAttributesFrom(first.Value.(*CandidateNode))
		}
	}
	return matchingNodes, nil
}
