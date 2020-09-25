package yqlib

type dataTreeNavigator struct {
}

type DataTreeNavigator interface {
	GetMatchingNodes(matchingNodes []*NodeContext, pathNode *PathTreeNode) ([]*NodeContext, error)
}

func NewTreeNavigator() DataTreeNavigator {
	return &dataTreeNavigator{}
}

func (d *dataTreeNavigator) traverseSingle(matchingNode *NodeContext, pathNode *PathElement) ([]*NodeContext, error) {
	var value = matchingNode.Node
	// match all for splat
	// match all and recurse for deep
	// etc and so forth

}

func (d *dataTreeNavigator) traverse(matchingNodes []*NodeContext, pathNode *PathElement) ([]*NodeContext, error) {
	var newMatchingNodes = make([]*NodeContext, 0)
	var newNodes []*NodeContext
	var err error
	for _, node := range matchingNodes {

		newNodes, err = d.traverseSingle(node, pathNode)
		if err != nil {
			return nil, err
		}
		newMatchingNodes = append(newMatchingNodes, newNodes...)
	}

	return newMatchingNodes, nil
}

func (d *dataTreeNavigator) GetMatchingNodes(matchingNodes []*NodeContext, pathNode *PathTreeNode) ([]*NodeContext, error) {
	if pathNode.PathElement.PathElementType == PathKey || pathNode.PathElement.PathElementType == ArrayIndex {
		return d.traverse(matchingNodes, pathNode.PathElement)
	} else {
		var lhs, rhs []*NodeContext
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
			return d.setFunction(pathNode.PathElement, lhs, rhs), nil
		case Equals:
			lhs, err = d.GetMatchingNodes(matchingNodes, pathNode.Lhs)
			if err != nil {
				return nil, err
			}
			return d.findMatchingValues(lhs, pathNode.Rhs)
		case EqualsSelf:
			return d.findMatchingValues(matchingNodes, pathNode.Rhs)
		}

	}

}
