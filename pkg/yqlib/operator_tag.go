package yqlib

import (
	"container/list"

	yaml "gopkg.in/yaml.v3"
)

func assignTagOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {

	log.Debugf("AssignTagOperator: %v")
	tag := ""

	if !pathNode.Operation.UpdateAssign {
		rhs, err := d.GetMatchingNodes(matchingNodes, pathNode.Rhs)
		if err != nil {
			return nil, err
		}

		if rhs.Front() != nil {
			tag = rhs.Front().Value.(*CandidateNode).Node.Value
		}
	}

	lhs, err := d.GetMatchingNodes(matchingNodes, pathNode.Lhs)

	if err != nil {
		return nil, err
	}

	for el := lhs.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		log.Debugf("Setting tag of : %v", candidate.GetKey())
		if pathNode.Operation.UpdateAssign {
			rhs, err := d.GetMatchingNodes(nodeToMap(candidate), pathNode.Rhs)
			if err != nil {
				return nil, err
			}

			if rhs.Front() != nil {
				tag = rhs.Front().Value.(*CandidateNode).Node.Value
			}
		}
		unwrapDoc(candidate.Node).Tag = tag
	}

	return matchingNodes, nil
}

func getTagOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("GetTagOperator")

	var results = list.New()

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := &yaml.Node{Kind: yaml.ScalarNode, Value: unwrapDoc(candidate.Node).Tag, Tag: "!!str"}
		result := candidate.CreateChild(nil, node)
		results.PushBack(result)
	}

	return results, nil
}
