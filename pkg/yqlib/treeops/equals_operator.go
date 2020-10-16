package treeops

import (
	"github.com/elliotchance/orderedmap"
	"gopkg.in/yaml.v3"
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

		matchString := "true"
		if !matches {
			matchString = "false"
		}

		node := &yaml.Node{Kind: yaml.ScalarNode, Value: matchString, Tag: "!!bool"}
		lengthCand := &CandidateNode{Node: node, Document: candidate.Document, Path: candidate.Path}
		results.Set(candidate.GetKey(), lengthCand)

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
