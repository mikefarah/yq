package yqlib

import (
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

func deleteChildOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	contextToUse := context.Clone()
	contextToUse.DontAutoCreate = true
	nodesToDelete, err := d.GetMatchingNodes(contextToUse, expressionNode.Rhs)

	if err != nil {
		return Context{}, err
	}
	//need to iterate backwards to ensure correct indices when deleting multiple
	for el := nodesToDelete.MatchingNodes.Back(); el != nil; el = el.Prev() {
		candidate := el.Value.(*CandidateNode)

		if len(candidate.Path) > 0 {
			deleteImmediateChildOp := &Operation{
				OperationType: deleteImmediateChildOpType,
				Value:         candidate.Path[len(candidate.Path)-1],
			}

			deleteImmediateChildOpNode := &ExpressionNode{
				Operation: deleteImmediateChildOp,
				Rhs:       createTraversalTree(candidate.Path[0:len(candidate.Path)-1], traversePreferences{}, false),
			}

			_, err := d.GetMatchingNodes(contextToUse, deleteImmediateChildOpNode)
			if err != nil {
				return Context{}, err
			}
		}
	}
	return context, nil
}

func deleteImmediateChildOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	parents, err := d.GetMatchingNodes(context, expressionNode.Rhs)

	if err != nil {
		return Context{}, err
	}

	childPath := expressionNode.Operation.Value

	log.Debug("childPath to remove %v", childPath)

	for el := parents.MatchingNodes.Front(); el != nil; el = el.Next() {
		parent := el.Value.(*CandidateNode)
		parentNode := unwrapDoc(parent.Node)
		if parentNode.Kind == yaml.MappingNode {
			deleteFromMap(parent, childPath)
		} else if parentNode.Kind == yaml.SequenceNode {
			deleteFromArray(parent, childPath)
		} else {
			return Context{}, fmt.Errorf("Cannot delete nodes from parent of tag %v", parentNode.Tag)
		}

	}
	return context, nil
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
