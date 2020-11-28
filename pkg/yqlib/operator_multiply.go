package yqlib

import (
	"fmt"

	"container/list"

	yaml "gopkg.in/yaml.v3"
)

type CrossFunctionCalculation func(d *dataTreeNavigator, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error)

func crossFunction(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode, calculation CrossFunctionCalculation) (*list.List, error) {
	lhs, err := d.GetMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	log.Debugf("crossFunction LHS len: %v", lhs.Len())

	rhs, err := d.GetMatchingNodes(matchingNodes, pathNode.Rhs)

	if err != nil {
		return nil, err
	}
	log.Debugf("crossFunction RHS len: %v", rhs.Len())

	var results = list.New()

	for el := lhs.Front(); el != nil; el = el.Next() {
		lhsCandidate := el.Value.(*CandidateNode)

		for rightEl := rhs.Front(); rightEl != nil; rightEl = rightEl.Next() {
			log.Debugf("Applying calc")
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

type MultiplyPreferences struct {
	AppendArrays bool
}

func MultiplyOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("-- MultiplyOperator")
	return crossFunction(d, matchingNodes, pathNode, multiply(pathNode.Operation.Preferences.(*MultiplyPreferences)))
}

func multiply(preferences *MultiplyPreferences) func(d *dataTreeNavigator, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	return func(d *dataTreeNavigator, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
		lhs.Node = UnwrapDoc(lhs.Node)
		rhs.Node = UnwrapDoc(rhs.Node)
		log.Debugf("Multipling LHS: %v", lhs.Node.Tag)
		log.Debugf("-          RHS: %v", rhs.Node.Tag)

		shouldAppendArrays := preferences.AppendArrays

		if lhs.Node.Kind == yaml.MappingNode && rhs.Node.Kind == yaml.MappingNode ||
			(lhs.Node.Kind == yaml.SequenceNode && rhs.Node.Kind == yaml.SequenceNode) {

			var newBlank = &CandidateNode{
				Path:     lhs.Path,
				Document: lhs.Document,
				Filename: lhs.Filename,
				Node:     &yaml.Node{},
			}
			var newThing, err = mergeObjects(d, newBlank, lhs, false)
			if err != nil {
				return nil, err
			}
			return mergeObjects(d, newThing, rhs, shouldAppendArrays)

		}
		return nil, fmt.Errorf("Cannot multiply %v with %v", lhs.Node.Tag, rhs.Node.Tag)
	}
}

func mergeObjects(d *dataTreeNavigator, lhs *CandidateNode, rhs *CandidateNode, shouldAppendArrays bool) (*CandidateNode, error) {
	var results = list.New()

	// shouldn't recurse arrays if appending
	err := recursiveDecent(d, results, nodeToMap(rhs), !shouldAppendArrays)
	if err != nil {
		return nil, err
	}

	var pathIndexToStartFrom int = 0
	if results.Front() != nil {
		pathIndexToStartFrom = len(results.Front().Value.(*CandidateNode).Path)
	}

	for el := results.Front(); el != nil; el = el.Next() {
		err := applyAssignment(d, pathIndexToStartFrom, lhs, el.Value.(*CandidateNode), shouldAppendArrays)
		if err != nil {
			return nil, err
		}
	}
	return lhs, nil
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

func applyAssignment(d *dataTreeNavigator, pathIndexToStartFrom int, lhs *CandidateNode, rhs *CandidateNode, shouldAppendArrays bool) error {

	log.Debugf("merge - applyAssignment lhs %v, rhs: %v", NodeToString(lhs), NodeToString(rhs))

	lhsPath := rhs.Path[pathIndexToStartFrom:]

	assignmentOp := &Operation{OperationType: AssignAttributes}
	if rhs.Node.Kind == yaml.ScalarNode || rhs.Node.Kind == yaml.AliasNode {
		assignmentOp.OperationType = Assign
		assignmentOp.Preferences = &AssignOpPreferences{false}
	} else if shouldAppendArrays && rhs.Node.Kind == yaml.SequenceNode {
		assignmentOp.OperationType = AddAssign
	}
	rhsOp := &Operation{OperationType: ValueOp, CandidateNode: rhs}

	assignmentOpNode := &PathTreeNode{Operation: assignmentOp, Lhs: createTraversalTree(lhsPath), Rhs: &PathTreeNode{Operation: rhsOp}}

	_, err := d.GetMatchingNodes(nodeToMap(lhs), assignmentOpNode)

	return err
}
