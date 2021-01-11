package yqlib

import (
	"fmt"

	"container/list"
	"strconv"

	yaml "gopkg.in/yaml.v3"
)

func createAddOp(lhs *PathTreeNode, rhs *PathTreeNode) *PathTreeNode {
	return &PathTreeNode{Operation: &Operation{OperationType: addOpType},
		Lhs: lhs,
		Rhs: rhs}
}

func addAssignOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	assignmentOp := &Operation{OperationType: assignOpType}
	assignmentOp.UpdateAssign = false

	assignmentOpNode := &PathTreeNode{Operation: assignmentOp, Lhs: pathNode.Lhs, Rhs: createAddOp(pathNode.Lhs, pathNode.Rhs)}
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

func addOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
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
		if rhs.Node.Kind != yaml.ScalarNode {
			return nil, fmt.Errorf("%v (%v) cannot be added to a %v", rhs.Node.Tag, rhs.Path, lhsNode.Tag)
		}
		target.Node.Kind = yaml.ScalarNode
		target.Node.Style = lhsNode.Style
		return addScalars(target, lhsNode, rhs.Node)
	}

	return target, nil
}

func addScalars(target *CandidateNode, lhs *yaml.Node, rhs *yaml.Node) (*CandidateNode, error) {

	if lhs.Tag == "!!str" {
		target.Node.Tag = "!!str"
		target.Node.Value = lhs.Value + rhs.Value
	} else if lhs.Tag == "!!int" && rhs.Tag == "!!int" {
		lhsNum, err := strconv.Atoi(lhs.Value)
		if err != nil {
			return nil, err
		}
		rhsNum, err := strconv.Atoi(rhs.Value)
		if err != nil {
			return nil, err
		}
		sum := lhsNum + rhsNum
		target.Node.Tag = "!!int"
		target.Node.Value = fmt.Sprintf("%v", sum)
	} else if (lhs.Tag == "!!int" || lhs.Tag == "!!float") && (rhs.Tag == "!!int" || rhs.Tag == "!!float") {
		lhsNum, err := strconv.ParseFloat(lhs.Value, 64)
		if err != nil {
			return nil, err
		}
		rhsNum, err := strconv.ParseFloat(rhs.Value, 64)
		if err != nil {
			return nil, err
		}
		sum := lhsNum + rhsNum
		target.Node.Tag = "!!float"
		target.Node.Value = fmt.Sprintf("%v", sum)
	} else {
		return nil, fmt.Errorf("%v cannot be added to %v", lhs.Tag, rhs.Tag)
	}

	return target, nil
}
