package treeops

import "github.com/elliotchance/orderedmap"

func ValueOperator(d *dataTreeNavigator, matchMap *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	return nodeToMap(pathNode.Operation.CandidateNode), nil
}
