package treeops

import (
	"fmt"

	"github.com/elliotchance/orderedmap"
)

type dataTreeNavigator struct {
	leafTraverser    LeafTraverser
	operatorHandlers map[OperationType]OperatorHandler
}

type NavigationPrefs struct {
	FollowAlias bool
}

type DataTreeNavigator interface {
	GetMatchingNodes(matchingNodes []*CandidateNode, pathNode *PathTreeNode) ([]*CandidateNode, error)
}

func NewDataTreeNavigator(navigationPrefs NavigationPrefs) DataTreeNavigator {
	leafTraverser := NewLeafTraverser(navigationPrefs)
	operatorHandlers := make(map[OperationType]OperatorHandler)

	operatorHandlers[Traverse] = TraverseOperator
	operatorHandlers[Equals] = EqualsOperator
	operatorHandlers[Or] = UnionOperator
	operatorHandlers[And] = IntersectionOperator
	operatorHandlers[Assign] = AssignOperator
	operatorHandlers[DeleteChild] = DeleteChildOperator

	return &dataTreeNavigator{leafTraverser, operatorHandlers}
}

func (d *dataTreeNavigator) traverse(matchMap *orderedmap.OrderedMap, pathNode *PathElement) (*orderedmap.OrderedMap, error) {
	log.Debugf("-- Traversing")
	var matchingNodeMap = orderedmap.NewOrderedMap()
	var newNodes []*CandidateNode
	var err error

	for el := matchMap.Front(); el != nil; el = el.Next() {
		newNodes, err = d.leafTraverser.Traverse(el.Value.(*CandidateNode), pathNode)
		if err != nil {
			return nil, err
		}
		for _, n := range newNodes {
			matchingNodeMap.Set(n.getKey(), n)
		}
	}

	return matchingNodeMap, nil
}

func (d *dataTreeNavigator) GetMatchingNodes(matchingNodes []*CandidateNode, pathNode *PathTreeNode) ([]*CandidateNode, error) {
	var matchingNodeMap = orderedmap.NewOrderedMap()

	for _, n := range matchingNodes {
		matchingNodeMap.Set(n.getKey(), n)
	}

	matchedNodes, err := d.getMatchingNodes(matchingNodeMap, pathNode)
	if err != nil {
		return nil, err
	}

	values := make([]*CandidateNode, 0, matchedNodes.Len())

	for el := matchedNodes.Front(); el != nil; el = el.Next() {
		values = append(values, el.Value.(*CandidateNode))
	}
	return values, nil
}

func (d *dataTreeNavigator) getMatchingNodes(matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	log.Debugf("Processing Path: %v", pathNode.PathElement.toString())
	if pathNode.PathElement.PathElementType == SelfReference {
		return matchingNodes, nil
	} else if pathNode.PathElement.PathElementType == PathKey || pathNode.PathElement.PathElementType == ArrayIndex {
		return d.traverse(matchingNodes, pathNode.PathElement)
	} else {
		handler := d.operatorHandlers[pathNode.PathElement.OperationType]
		if handler != nil {
			return handler(d, matchingNodes, pathNode)
		}
		return nil, fmt.Errorf("Unknown operator %v", pathNode.PathElement.OperationType)
	}

}
