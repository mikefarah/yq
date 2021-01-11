package yqlib

import (
	"container/list"
)

func equalsOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("-- equalsOperation")
	return crossFunction(d, matchingNodes, pathNode, isEquals)
}

func isEquals(d *dataTreeNavigator, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	value := false

	lhsNode := UnwrapDoc(lhs.Node)
	rhsNode := UnwrapDoc(rhs.Node)

	if lhsNode.Tag == "!!null" {
		value = (rhsNode.Tag == "!!null")
	} else {
		value = matchKey(lhsNode.Value, rhsNode.Value)
	}
	log.Debugf("%v == %v ? %v", NodeToString(lhs), NodeToString(rhs), value)
	return createBooleanCandidate(lhs, value), nil
}
