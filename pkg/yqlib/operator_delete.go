package yqlib

import (
	"container/list"
	"fmt"
)

func deleteChildOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	nodesToDelete, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.RHS)

	if err != nil {
		return Context{}, err
	}
	//need to iterate backwards to ensure correct indices when deleting multiple
	for el := nodesToDelete.MatchingNodes.Back(); el != nil; el = el.Prev() {
		candidate := el.Value.(*CandidateNode)

		if candidate.Parent == nil {
			// must be a top level thing, delete it
			return removeFromContext(context, candidate)
		}
		log.Debugf("processing deletion of candidate %v", NodeToString(candidate))

		parentNode := candidate.Parent

		candidatePath := candidate.GetPath()
		childPath := candidatePath[len(candidatePath)-1]

		switch parentNode.Kind {
		case MappingNode:
			deleteFromMap(candidate.Parent, childPath)
		case SequenceNode:
			deleteFromArray(candidate.Parent, childPath)
		default:
			return Context{}, fmt.Errorf("cannot delete nodes from parent of tag %v", parentNode.Tag)
		}
	}
	return context, nil
}

func removeFromContext(context Context, candidate *CandidateNode) (Context, error) {
	newResults := list.New()
	for item := context.MatchingNodes.Front(); item != nil; item = item.Next() {
		nodeInContext := item.Value.(*CandidateNode)
		if nodeInContext != candidate {
			newResults.PushBack(nodeInContext)
		} else {
			log.Info("Need to delete this %v", NodeToString(nodeInContext))
		}
	}
	return context.ChildContext(newResults), nil
}

func deleteFromMap(node *CandidateNode, childPath interface{}) {
	log.Debug("deleteFromMap")
	contents := node.Content
	newContents := make([]*CandidateNode, 0)

	for index := 0; index < len(contents); index = index + 2 {
		key := contents[index]
		value := contents[index+1]

		shouldDelete := key.Value == childPath

		log.Debugf("shouldDelete %v? %v == %v = %v", NodeToString(value), key.Value, childPath, shouldDelete)

		if !shouldDelete {
			newContents = append(newContents, key, value)
		}
	}
	node.Content = newContents
}

func deleteFromArray(node *CandidateNode, childPath interface{}) {
	log.Debug("deleteFromArray")
	contents := node.Content
	newContents := make([]*CandidateNode, 0)

	for index := 0; index < len(contents); index = index + 1 {
		value := contents[index]

		shouldDelete := fmt.Sprintf("%v", index) == fmt.Sprintf("%v", childPath)

		if !shouldDelete {
			value.Key.Value = fmt.Sprintf("%v", len(newContents))
			newContents = append(newContents, value)
		}
	}
	node.Content = newContents
}
