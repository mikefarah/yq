package treeops

import (
	"fmt"

	"container/list"

	"gopkg.in/yaml.v3"
)

type CrossFunctionCalculation func(d *dataTreeNavigator, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error)

func crossFunction(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode, calculation CrossFunctionCalculation) (*list.List, error) {
	lhs, err := d.getMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}

	rhs, err := d.getMatchingNodes(matchingNodes, pathNode.Rhs)

	if err != nil {
		return nil, err
	}

	var results = list.New()

	for el := lhs.Front(); el != nil; el = el.Next() {
		lhsCandidate := el.Value.(*CandidateNode)

		for rightEl := rhs.Front(); rightEl != nil; rightEl = rightEl.Next() {
			rhsCandidate := rightEl.Value.(*CandidateNode)
			resultCandidate, err := calculation(d, lhsCandidate, rhsCandidate)
			if err != nil {
				return nil, err
			}
			results.PushBack(resultCandidate)
		}

	}
	return results, nil
}

func MultiplyOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("-- MultiplyOperator")
	return crossFunction(d, matchingNodes, pathNode, multiply)
}

func multiply(d *dataTreeNavigator, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	if lhs.Node.Kind == yaml.MappingNode && rhs.Node.Kind == yaml.MappingNode ||
		(lhs.Node.Kind == yaml.SequenceNode && rhs.Node.Kind == yaml.SequenceNode) {
		var results = list.New()
		recursiveDecent(d, results, nodeToMap(rhs))

		var pathIndexToStartFrom int = 0
		if results.Front() != nil {
			pathIndexToStartFrom = len(results.Front().Value.(*CandidateNode).Path)
		}

		for el := results.Front(); el != nil; el = el.Next() {
			err := applyAssignment(d, pathIndexToStartFrom, lhs, el.Value.(*CandidateNode))
			if err != nil {
				return nil, err
			}
		}
		return lhs, nil
	}
	return nil, fmt.Errorf("Cannot multiply %v with %v", NodeToString(lhs), NodeToString(rhs))
}

func createTraversalTree(path []interface{}) *PathTreeNode {
	if len(path) == 0 {
		return &PathTreeNode{Operation: &Operation{OperationType: SelfReference}}
	} else if len(path) == 1 {
		return &PathTreeNode{Operation: &Operation{OperationType: TraversePath, Value: path[0], StringValue: fmt.Sprintf("%v", path[0])}}
	}
	return &PathTreeNode{
		Operation: &Operation{OperationType: Pipe},
		Lhs:       createTraversalTree(path[0:1]),
		Rhs:       createTraversalTree(path[1:])}

}

func applyAssignment(d *dataTreeNavigator, pathIndexToStartFrom int, lhs *CandidateNode, rhs *CandidateNode) error {
	log.Debugf("merge - applyAssignment lhs %v, rhs: %v", NodeToString(lhs), NodeToString(rhs))

	lhsPath := rhs.Path[pathIndexToStartFrom:]

	assignmentOp := &Operation{OperationType: AssignAttributes}
	if rhs.Node.Kind == yaml.ScalarNode {
		assignmentOp.OperationType = Assign
	}
	rhsOp := &Operation{OperationType: ValueOp, CandidateNode: rhs}

	assignmentOpNode := &PathTreeNode{Operation: assignmentOp, Lhs: createTraversalTree(lhsPath), Rhs: &PathTreeNode{Operation: rhsOp}}

	_, err := d.getMatchingNodes(nodeToMap(lhs), assignmentOpNode)

	return err
}
