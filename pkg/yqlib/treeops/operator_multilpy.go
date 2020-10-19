package treeops

import (
	"fmt"

	"github.com/elliotchance/orderedmap"
	"gopkg.in/yaml.v3"
)

func MultiplyOperator(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	lhs, err := d.getMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}

	rhs, err := d.getMatchingNodes(matchingNodes, pathNode.Rhs)

	if err != nil {
		return nil, err
	}

	var results = orderedmap.NewOrderedMap()

	for el := lhs.Front(); el != nil; el = el.Next() {
		lhsCandidate := el.Value.(*CandidateNode)

		for rightEl := rhs.Front(); rightEl != nil; rightEl = rightEl.Next() {
			rhsCandidate := rightEl.Value.(*CandidateNode)
			resultCandidate, err := multiply(d, lhsCandidate, rhsCandidate)
			if err != nil {
				return nil, err
			}
			results.Set(resultCandidate.GetKey(), resultCandidate)
		}

	}
	return matchingNodes, nil
}

func multiply(d *dataTreeNavigator, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	if lhs.Node.Kind == yaml.MappingNode && rhs.Node.Kind == yaml.MappingNode ||
		(lhs.Node.Kind == yaml.SequenceNode && rhs.Node.Kind == yaml.SequenceNode) {
		var results = orderedmap.NewOrderedMap()
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
		return &PathTreeNode{PathElement: &PathElement{PathElementType: SelfReference}}
	} else if len(path) == 1 {
		return &PathTreeNode{PathElement: &PathElement{PathElementType: PathKey, Value: path[0], StringValue: fmt.Sprintf("%v", path[0])}}
	}
	return &PathTreeNode{
		PathElement: &PathElement{PathElementType: Operation, OperationType: Pipe},
		Lhs:         createTraversalTree(path[0:1]),
		Rhs:         createTraversalTree(path[1:])}

}

func applyAssignment(d *dataTreeNavigator, pathIndexToStartFrom int, lhs *CandidateNode, rhs *CandidateNode) error {
	log.Debugf("merge - applyAssignment lhs %v, rhs: %v", NodeToString(lhs), NodeToString(rhs))

	lhsPath := rhs.Path[pathIndexToStartFrom:]

	assignmentOp := &PathElement{PathElementType: Operation, OperationType: AssignAttributes}
	if rhs.Node.Kind == yaml.ScalarNode {
		assignmentOp.OperationType = Assign
	}
	rhsOp := &PathElement{PathElementType: Value, CandidateNode: rhs}

	assignmentOpNode := &PathTreeNode{PathElement: assignmentOp, Lhs: createTraversalTree(lhsPath), Rhs: &PathTreeNode{PathElement: rhsOp}}

	_, err := d.getMatchingNodes(nodeToMap(lhs), assignmentOpNode)

	return err
}
