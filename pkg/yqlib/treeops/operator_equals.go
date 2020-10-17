package treeops

import (
	"github.com/elliotchance/orderedmap"
)

func EqualsOperator(d *dataTreeNavigator, matchMap *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	log.Debugf("-- equalsOperation")
	var results = orderedmap.NewOrderedMap()

	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		log.Debug("equalsOperation checking %v", candidate)

		matches, errInChild := hasMatch(d, candidate, pathNode.Lhs, pathNode.Rhs)
		if errInChild != nil {
			return nil, errInChild
		}

		equalsCandidate := createBooleanCandidate(candidate, matches)
		results.Set(equalsCandidate.GetKey(), equalsCandidate)
	}

	return results, nil
}

func hasMatch(d *dataTreeNavigator, candidate *CandidateNode, lhs *PathTreeNode, rhs *PathTreeNode) (bool, error) {
	childMap := orderedmap.NewOrderedMap()
	childMap.Set(candidate.GetKey(), candidate)
	childMatches, errChild := d.getMatchingNodes(childMap, lhs)
	log.Debug("got the LHS")
	if errChild != nil {
		return false, errChild
	}

	// TODO = handle other RHS types
	return containsMatchingValue(childMatches, rhs.PathElement.StringValue), nil
}

func containsMatchingValue(matchMap *orderedmap.OrderedMap, valuePattern string) bool {
	log.Debugf("-- findMatchingValues")

	for el := matchMap.Front(); el != nil; el = el.Next() {
		node := el.Value.(*CandidateNode)
		log.Debugf("-- comparing %v to %v", node.Node.Value, valuePattern)
		if Match(node.Node.Value, valuePattern) {
			return true
		}
	}
	log.Debugf("-- done findMatchingValues")

	return false
}
