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
	assignmentOp.UpdateAssign = true

	assignmentOpNode := &PathTreeNode{Operation: assignmentOp, Lhs: pathNode.Lhs, Rhs: createSelfAddOp(pathNode.Rhs)}
	return d.GetMatchingNodes(matchingNodes, assignmentOpNode)
}

func toNodes(candidate *CandidateNode) []*yaml.Node {
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

	return crossFunction(d, matchingNodes, pathNode, add)
}

func add(d *dataTreeNavigator, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	lhs.Node = UnwrapDoc(lhs.Node)
	rhs.Node = UnwrapDoc(rhs.Node)

	target := &CandidateNode{
		Path:     lhs.Path,
		Document: lhs.Document,
		Filename: lhs.Filename,
		Node:     &yaml.Node{},
	}
	lhsNode := lhs.Node

	switch lhsNode.Kind {
	case yaml.MappingNode:
		return nil, fmt.Errorf("Maps not yet supported for addition")
	case yaml.SequenceNode:
		target.Node.Kind = yaml.SequenceNode
		target.Node.Style = lhsNode.Style
		target.Node.Tag = "!!seq"
		target.Node.Content = append(lhsNode.Content, toNodes(rhs)...)
	case yaml.ScalarNode:
		return nil, fmt.Errorf("Scalars not yet supported for addition")
	}

	return target, nil
}
