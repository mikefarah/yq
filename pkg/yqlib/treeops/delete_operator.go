package treeops

import (
	"github.com/elliotchance/orderedmap"
	"gopkg.in/yaml.v3"
)

func DeleteChildOperator(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	lhs, err := d.getMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	// for each lhs, splat the node,
	// the intersect it against the rhs expression
	// recreate the contents using only the intersection result.

	for el := lhs.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		elMap := orderedmap.NewOrderedMap()
		elMap.Set(candidate.getKey(), candidate)
		nodesToDelete, err := d.getMatchingNodes(elMap, pathNode.Rhs)
		log.Debug("nodesToDelete:\n%v", NodesToString(nodesToDelete))
		if err != nil {
			return nil, err
		}

		if candidate.Node.Kind == yaml.SequenceNode {
			deleteFromArray(candidate, nodesToDelete)
		} else {
			deleteFromMap(candidate, nodesToDelete)
		}
	}
	return lhs, nil
}

func deleteFromMap(candidate *CandidateNode, nodesToDelete *orderedmap.OrderedMap) {
	log.Debug("deleteFromMap")
	node := candidate.Node
	contents := node.Content
	newContents := make([]*yaml.Node, 0)

	for index := 0; index < len(contents); index = index + 2 {
		key := contents[index]
		value := contents[index+1]

		childCandidate := &CandidateNode{
			Node:     value,
			Document: candidate.Document,
			Path:     append(candidate.Path, key.Value),
		}
		_, shouldDelete := nodesToDelete.Get(childCandidate.getKey())

		log.Debugf("shouldDelete %v ? %v", childCandidate.getKey(), shouldDelete)

		if !shouldDelete {
			newContents = append(newContents, key, value)
		}
	}
	node.Content = newContents
}

func deleteFromArray(candidate *CandidateNode, nodesToDelete *orderedmap.OrderedMap) {
	log.Debug("deleteFromArray")
	node := candidate.Node
	contents := node.Content
	newContents := make([]*yaml.Node, 0)

	for index := 0; index < len(contents); index = index + 1 {
		value := contents[index]

		childCandidate := &CandidateNode{
			Node:     value,
			Document: candidate.Document,
			Path:     append(candidate.Path, index),
		}

		_, shouldDelete := nodesToDelete.Get(childCandidate.getKey())
		if !shouldDelete {
			newContents = append(newContents, value)
		}
	}
	node.Content = newContents
}
