package yqlib

import (
	"container/list"

	"gopkg.in/yaml.v3"
)

func AssignAnchorOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {

	log.Debugf("AssignAnchor operator!")

	rhs, err := d.GetMatchingNodes(matchingNodes, pathNode.Rhs)
	if err != nil {
		return nil, err
	}
	anchorName := ""
	if rhs.Front() != nil {
		anchorName = rhs.Front().Value.(*CandidateNode).Node.Value
	}

	lhs, err := d.GetMatchingNodes(matchingNodes, pathNode.Lhs)

	if err != nil {
		return nil, err
	}

	for el := lhs.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		log.Debugf("Setting anchorName of : %v", candidate.GetKey())
		candidate.Node.Anchor = anchorName
	}
	return matchingNodes, nil
}

func GetAnchorOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("GetAnchor operator!")
	var results = list.New()

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		anchor := candidate.Node.Anchor
		node := &yaml.Node{Kind: yaml.ScalarNode, Value: anchor, Tag: "!!str"}
		lengthCand := &CandidateNode{Node: node, Document: candidate.Document, Path: candidate.Path}
		results.PushBack(lengthCand)
	}
	return results, nil
}
