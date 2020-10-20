package treeops

import "github.com/elliotchance/orderedmap"

func ValueOperator(d *dataTreeNavigator, matchMap *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	log.Debug("value = %v", pathNode.Operation.CandidateNode.Node.Value)
	return nodeToMap(pathNode.Operation.CandidateNode), nil
}
