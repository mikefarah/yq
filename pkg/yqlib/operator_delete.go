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

		if candidate.Kind == DocumentNode {
			//need to delete this node from context.
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
		} else if candidate.Parent == nil {
			//problem: context may already be '.a' and then I pass in '.a.a2'.
			// should pass in .a2.
			log.Info("Could not find parent of %v", NodeToString(candidate))
			return context, nil
		}

		parentNode := candidate.Parent
		childPath := candidate.Path[len(candidate.Path)-1]

		if parentNode.Kind == MappingNode {
			deleteFromMap(candidate.Parent, childPath)
		} else if parentNode.Kind == SequenceNode {
			deleteFromArray(candidate.Parent, childPath)
		} else {
			return Context{}, fmt.Errorf("Cannot delete nodes from parent of tag %v", parentNode.Tag)
		}
	}
	return context, nil
}

func deleteFromMap(candidate *CandidateNode, childPath interface{}) {
	log.Debug("deleteFromMap")
	node := candidate.unwrapDocument()
	contents := node.Content
	newContents := make([]*CandidateNode, 0)

	for index := 0; index < len(contents); index = index + 2 {
		key := contents[index]
		value := contents[index+1]

		shouldDelete := key.Value == childPath

		log.Debugf("shouldDelete %v ? %v", value.GetKey(), shouldDelete)

		if !shouldDelete {
			newContents = append(newContents, key, value)
		}
	}
	node.Content = newContents
}

func deleteFromArray(candidate *CandidateNode, childPath interface{}) {
	log.Debug("deleteFromArray")
	node := candidate.unwrapDocument()
	contents := node.Content
	newContents := make([]*CandidateNode, 0)

	for index := 0; index < len(contents); index = index + 1 {
		value := contents[index]

		shouldDelete := fmt.Sprintf("%v", index) == fmt.Sprintf("%v", childPath)

		if !shouldDelete {
			newContents = append(newContents, value)
		}
	}
	node.Content = newContents
}
