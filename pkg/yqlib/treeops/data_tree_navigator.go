package treeops

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

func (d *dataTreeNavigator) traverse(matchingNodes []*CandidateNode, pathNode *PathElement) ([]*CandidateNode, error) {
	log.Debugf("-- Traversing")
	var newMatchingNodes = make([]*CandidateNode, 0)
	var newNodes []*CandidateNode
	var err error
	for _, node := range matchingNodes {

		newNodes, err = d.traverser.Traverse(node, pathNode)
		if err != nil {
			return nil, err
		}
		newMatchingNodes = append(newMatchingNodes, newNodes...)
	}

	return newMatchingNodes, nil
}

func (d *dataTreeNavigator) setFunction(op OperationType, lhs []*CandidateNode, rhs []*CandidateNode) []*CandidateNode {

	return append(lhs, rhs...)
}

func (d *dataTreeNavigator) GetMatchingNodes(matchingNodes []*CandidateNode, pathNode *PathTreeNode) ([]*CandidateNode, error) {
	log.Debugf("Processing Path: %v", pathNode.PathElement.toString())
	if pathNode.PathElement.PathElementType == PathKey || pathNode.PathElement.PathElementType == ArrayIndex {
		return d.traverse(matchingNodes, pathNode.PathElement)
	} else {
		var lhs, rhs []*CandidateNode
		var err error
		switch pathNode.PathElement.OperationType {
		case Traverse:
			lhs, err = d.GetMatchingNodes(matchingNodes, pathNode.Lhs)
			if err != nil {
				return nil, err
			}
			return d.GetMatchingNodes(lhs, pathNode.Rhs)
		case Or, And:
			lhs, err = d.GetMatchingNodes(matchingNodes, pathNode.Lhs)
			if err != nil {
				return nil, err
			}
			rhs, err = d.GetMatchingNodes(matchingNodes, pathNode.Rhs)
			if err != nil {
				return nil, err
			}
			return d.setFunction(pathNode.PathElement.OperationType, lhs, rhs), nil
		// case Equals:
		// 	lhs, err = d.GetMatchingNodes(matchingNodes, pathNode.Lhs)
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
