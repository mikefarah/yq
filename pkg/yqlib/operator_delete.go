package yqlib

import (
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

func deleteChildOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	nodesToDelete, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.RHS)

	if err != nil {
		return Context{}, err
	}
	//need to iterate backwards to ensure correct indices when deleting multiple
	for el := nodesToDelete.MatchingNodes.Back(); el != nil; el = el.Prev() {
		candidate := el.Value.(*CandidateNode)

		//problem: context may already be '.a' and then I pass in '.a.a2'.
		// should pass in .a2.
		if candidate.Parent == nil {
			log.Info("Could not find parent of %v", candidate.GetKey())
			return context, nil
		}

		parentNode := candidate.Parent.Node
		childPath := candidate.Path[len(candidate.Path)-1]

		if parentNode.Kind == yaml.MappingNode {
			deleteFromMap(candidate.Parent, childPath)
		} else if parentNode.Kind == yaml.SequenceNode {
			deleteFromArray(candidate.Parent, childPath)
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

		childCandidate := candidate.CreateChildInMap(key, value)

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
