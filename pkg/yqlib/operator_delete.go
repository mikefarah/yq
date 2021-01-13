package yqlib

import (
	"container/list"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

func deleteChildOperator(d *dataTreeNavigator, matchingNodes *list.List, expressionNode *ExpressionNode) (*list.List, error) {

	nodesToDelete, err := d.GetMatchingNodes(matchingNodes, expressionNode.Rhs)

	if err != nil {
		return nil, err
	}

	for el := nodesToDelete.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		deleteImmediateChildOp := &Operation{
			OperationType: deleteImmediateChildOpType,
			Value:         candidate.Path[len(candidate.Path)-1],
		}

		deleteImmediateChildOpNode := &ExpressionNode{
			Operation: deleteImmediateChildOp,
			Rhs:       createTraversalTree(candidate.Path[0:len(candidate.Path)-1], traversePreferences{}),
		}

		_, err := d.GetMatchingNodes(matchingNodes, deleteImmediateChildOpNode)
		if err != nil {
			return nil, err
		}
	}
	return matchingNodes, nil
}

func deleteImmediateChildOperator(d *dataTreeNavigator, matchingNodes *list.List, expressionNode *ExpressionNode) (*list.List, error) {
	parents, err := d.GetMatchingNodes(matchingNodes, expressionNode.Rhs)

	if err != nil {
		return nil, err
	}

	childPath := expressionNode.Operation.Value

	log.Debug("childPath to remove %v", childPath)

	for el := parents.Front(); el != nil; el = el.Next() {
		parent := el.Value.(*CandidateNode)
		parentNode := unwrapDoc(parent.Node)
		if parentNode.Kind == yaml.MappingNode {
			deleteFromMap(parent, childPath)
		} else if parentNode.Kind == yaml.SequenceNode {
			deleteFromArray(parent, childPath)
		} else {
			return nil, fmt.Errorf("Cannot delete nodes from parent of tag %v", parentNode.Tag)
		}

	}
	return matchingNodes, nil
}

func deleteFromMap(candidate *CandidateNode, childPath interface{}) {
	log.Debug("deleteFromMap")
	node := unwrapDoc(candidate.Node)
	contents := node.Content
	newContents := make([]*yaml.Node, 0)

	for index := 0; index < len(contents); index = index + 2 {
		key := contents[index]
		value := contents[index+1]

		childCandidate := candidate.CreateChild(key.Value, value)

		shouldDelete := key.Value == childPath

		log.Debugf("shouldDelete %v ? %v", childCandidate.GetKey(), shouldDelete)

		if !shouldDelete {
			newContents = append(newContents, key, value)
		}
	}
	node.Content = newContents
}

func deleteFromArray(candidate *CandidateNode, childPath interface{}) {
	log.Debug("deleteFromArray")
	node := unwrapDoc(candidate.Node)
	contents := node.Content
	newContents := make([]*yaml.Node, 0)

	for index := 0; index < len(contents); index = index + 1 {
		value := contents[index]

		shouldDelete := fmt.Sprintf("%v", index) == fmt.Sprintf("%v", childPath)

		if !shouldDelete {
			newContents = append(newContents, value)
		}
	}
	node.Content = newContents
}
