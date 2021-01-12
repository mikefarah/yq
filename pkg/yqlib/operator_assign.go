package yqlib

import "container/list"

func assignUpdateOperator(d *dataTreeNavigator, matchingNodes *list.List, expressionNode *ExpressionNode) (*list.List, error) {
	lhs, err := d.GetMatchingNodes(matchingNodes, expressionNode.Lhs)
	if err != nil {
		return nil, err
	}
	var rhs *list.List
	if !expressionNode.Operation.UpdateAssign {
		rhs, err = d.GetMatchingNodes(matchingNodes, expressionNode.Rhs)
	}

	for el := lhs.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		if expressionNode.Operation.UpdateAssign {
			rhs, err = d.GetMatchingNodes(nodeToMap(candidate), expressionNode.Rhs)
		}

		if err != nil {
			return nil, err
		}

		// grab the first value
		first := rhs.Front()

		if first != nil {
			rhsCandidate := first.Value.(*CandidateNode)
			rhsCandidate.Node = unwrapDoc(rhsCandidate.Node)
			candidate.UpdateFrom(rhsCandidate)
		}
	}

	return matchingNodes, nil
}

// does not update content or values
func assignAttributesOperator(d *dataTreeNavigator, matchingNodes *list.List, expressionNode *ExpressionNode) (*list.List, error) {
	lhs, err := d.GetMatchingNodes(matchingNodes, expressionNode.Lhs)
	if err != nil {
		return nil, err
	}
	for el := lhs.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		rhs, err := d.GetMatchingNodes(nodeToMap(candidate), expressionNode.Rhs)

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
