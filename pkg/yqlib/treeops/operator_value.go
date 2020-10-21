package treeops

import "container/list"

func ValueOperator(d *dataTreeNavigator, matchMap *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debug("value = %v", pathNode.Operation.CandidateNode.Node.Value)
	return nodeToMap(pathNode.Operation.CandidateNode), nil
}
