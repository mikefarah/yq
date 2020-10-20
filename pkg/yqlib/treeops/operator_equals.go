package treeops

import (
	"github.com/elliotchance/orderedmap"
)

func EqualsOperator(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	log.Debugf("-- equalsOperation")
	return crossFunction(d, matchingNodes, pathNode, isEquals)
}

func isEquals(d *dataTreeNavigator, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	value := false

	if lhs.Node.Tag == "!!null" {
		value = (rhs.Node.Tag == "!!null")
	} else {
		value = Match(lhs.Node.Value, rhs.Node.Value)
	}
	log.Debugf("%v == %v ? %v", NodeToString(lhs), NodeToString(rhs), value)
	return createBooleanCandidate(lhs, value), nil
}
