package treeops

import (
	"github.com/elliotchance/orderedmap"
)

type dataTreeNavigator struct {
	traverser Traverser
}

type NavigationPrefs struct {
	FollowAlias bool
}

type DataTreeNavigator interface {
	GetMatchingNodes(matchingNodes []*CandidateNode, pathNode *PathTreeNode) ([]*CandidateNode, error)
}

func NewDataTreeNavigator(navigationPrefs NavigationPrefs) DataTreeNavigator {
	traverse := NewTraverser(navigationPrefs)
	return &dataTreeNavigator{traverse}
}

func (d *dataTreeNavigator) traverse(matchMap *orderedmap.OrderedMap, pathNode *PathElement) (*orderedmap.OrderedMap, error) {
	log.Debugf("-- Traversing")
	var matchingNodeMap = orderedmap.NewOrderedMap()
	var newNodes []*CandidateNode
	var err error

	for el := matchMap.Front(); el != nil; el = el.Next() {
		newNodes, err = d.traverser.Traverse(el.Value.(*CandidateNode), pathNode)
		if err != nil {
			return nil, err
		}
		for _, n := range newNodes {
			matchingNodeMap.Set(n.getKey(), n)
		}
	}

	return matchingNodeMap, nil
}

func (d *dataTreeNavigator) setFunction(op OperationType, lhs *orderedmap.OrderedMap, rhs *orderedmap.OrderedMap) *orderedmap.OrderedMap {

	if op == Or {
		for el := rhs.Front(); el != nil; el = el.Next() {
			node := el.Value.(*CandidateNode)
			lhs.Set(node.getKey(), node)
		}
		return lhs
	}
	var matchingNodeMap = orderedmap.NewOrderedMap()
	for el := lhs.Front(); el != nil; el = el.Next() {
		_, exists := rhs.Get(el.Key)
		if exists {
			matchingNodeMap.Set(el.Key, el.Value)
		}
	}
	return matchingNodeMap

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
	if pathNode.PathElement.PathElementType == PathKey || pathNode.PathElement.PathElementType == ArrayIndex {
		return d.traverse(matchingNodes, pathNode.PathElement)
	} else {
		var lhs, rhs *orderedmap.OrderedMap
		var err error
		switch pathNode.PathElement.OperationType {
		case Traverse:
			lhs, err = d.getMatchingNodes(matchingNodes, pathNode.Lhs)
			if err != nil {
				return nil, err
			}
			return d.getMatchingNodes(lhs, pathNode.Rhs)
		case Or, And:
			lhs, err = d.getMatchingNodes(matchingNodes, pathNode.Lhs)
			if err != nil {
				return nil, err
			}
			rhs, err = d.getMatchingNodes(matchingNodes, pathNode.Rhs)
			if err != nil {
				return nil, err
			}
			return d.setFunction(pathNode.PathElement.OperationType, lhs, rhs), nil
		// case Equals:
		// 	lhs, err = d.getMatchingNodes(matchingNodes, pathNode.Lhs)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	return d.findMatchingValues(lhs, pathNode.Rhs)
		// case EqualsSelf:
		// 	return d.findMatchingValues(matchingNodes, pathNode.Rhs)
		default:
			return nil, nil
		}

	}

}
