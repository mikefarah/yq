package yqlib

import (
	"container/list"
)

// corssFunction no matches
// can boolean use crossfunction

func AlternativeOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("-- alternative")
	return crossFunction(d, matchingNodes, pathNode, alternativeFunc)
}

func alternativeFunc(d *dataTreeNavigator, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	lhs.Node = UnwrapDoc(lhs.Node)
	rhs.Node = UnwrapDoc(rhs.Node)
	log.Debugf("Alternative LHS: %v", lhs.Node.Tag)
	log.Debugf("-          RHS: %v", rhs.Node.Tag)

	isTrue, err := isTruthy(lhs)
	if err != nil {
		return nil, err
	} else if isTrue {
		return lhs, nil
	}
	return rhs, nil
}
