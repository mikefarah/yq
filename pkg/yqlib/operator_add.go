package yqlib

import (
	"fmt"

	"container/list"

	yaml "gopkg.in/yaml.v3"
)

func createSelfAddOp(rhs *PathTreeNode) *PathTreeNode {
	return &PathTreeNode{Operation: &Operation{OperationType: Add},
		Lhs: &PathTreeNode{Operation: &Operation{OperationType: SelfReference}},
		Rhs: rhs}
}

func AddAssignOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	assignmentOp := &Operation{OperationType: Assign}
	assignmentOp.Preferences = &AssignOpPreferences{true}

	assignmentOpNode := &PathTreeNode{Operation: assignmentOp, Lhs: pathNode.Lhs, Rhs: createSelfAddOp(pathNode.Rhs)}
	return d.GetMatchingNodes(matchingNodes, assignmentOpNode)
}

func toNodes(candidates *list.List) []*yaml.Node {

	if candidates.Len() == 0 {
		return []*yaml.Node{}
	}
	candidate := candidates.Front().Value.(*CandidateNode)

	if candidate.Node.Tag == "!!null" {
		return []*yaml.Node{}
	}

	switch candidate.Node.Kind {
	case yaml.SequenceNode:
		return candidate.Node.Content
	default:
		return []*yaml.Node{candidate.Node}
	}

}

func AddOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("Add operator")
	var results = list.New()
	lhs, err := d.GetMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	rhs, err := d.GetMatchingNodes(matchingNodes, pathNode.Rhs)

	if err != nil {
		return nil, err
	}

	for el := lhs.Front(); el != nil; el = el.Next() {
		lhsCandidate := el.Value.(*CandidateNode)
		lhsNode := UnwrapDoc(lhsCandidate.Node)

		target := &CandidateNode{
			Path:     lhsCandidate.Path,
			Document: lhsCandidate.Document,
			Filename: lhsCandidate.Filename,
			Node:     &yaml.Node{},
		}

		switch lhsNode.Kind {
		case yaml.MappingNode:
			return nil, fmt.Errorf("Maps not yet supported for addition")
		case yaml.SequenceNode:
			target.Node.Kind = yaml.SequenceNode
			target.Node.Style = lhsNode.Style
			target.Node.Tag = "!!seq"
			target.Node.Content = append(lhsNode.Content, toNodes(rhs)...)
			results.PushBack(target)
		case yaml.ScalarNode:
			return nil, fmt.Errorf("Scalars not yet supported for addition")
		}
	}
	return results, nil
}
