package treeops

import (
	"fmt"

	"container/list"

	"gopkg.in/op/go-logging.v1"
)

type dataTreeNavigator struct {
	navigationPrefs NavigationPrefs
}

type NavigationPrefs struct {
	FollowAlias bool
}

type DataTreeNavigator interface {
	// given a list of CandidateEntities and a pathNode,
	// this will process the list against the given pathNode and return
	// a new list of matching candidates
	GetMatchingNodes(matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error)
}

func NewDataTreeNavigator(navigationPrefs NavigationPrefs) DataTreeNavigator {
	return &dataTreeNavigator{navigationPrefs}
}

func (d *dataTreeNavigator) GetMatchingNodes(matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	if pathNode == nil {
		log.Debugf("getMatchingNodes - nothing to do")
		return matchingNodes, nil
	}
	log.Debugf("Processing Op: %v", pathNode.Operation.toString())
	if log.IsEnabledFor(logging.DEBUG) {
		for el := matchingNodes.Front(); el != nil; el = el.Next() {
			log.Debug(NodeToString(el.Value.(*CandidateNode)))
		}
	}
	log.Debug(">>")
	handler := pathNode.Operation.OperationType.Handler
	if handler != nil {
		return handler(d, matchingNodes, pathNode)
	}
	return nil, fmt.Errorf("Unknown operator %v", pathNode.Operation.OperationType)

}
