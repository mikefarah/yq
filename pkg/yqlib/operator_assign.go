package yqlib

import "container/list"

func AssignUpdateOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	lhs, err := d.GetMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	var rhs *list.List
	if !pathNode.Operation.UpdateAssign {
		rhs, err = d.GetMatchingNodes(matchingNodes, pathNode.Rhs)
	}

	for el := lhs.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		if pathNode.Operation.UpdateAssign {
			rhs, err = d.GetMatchingNodes(nodeToMap(candidate), pathNode.Rhs)
		}

		if err != nil {
			return nil, err
		}

		// grab the first value
		first := rhs.Front()

		if first != nil {
			rhsCandidate := first.Value.(*CandidateNode)
			rhsCandidate.Node = UnwrapDoc(rhsCandidate.Node)
			candidate.UpdateFrom(rhsCandidate)
		}
	}
	// // if there was nothing given, perhaps we are creating a new yaml doc
	// if matchingNodes.Len() == 0 {
	// 	log.Debug("started with nothing, returning LHS, %v", lhs.Len())
	// 	return lhs, nil
	// }
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
