package treeops

import (
	"container/list"

	"gopkg.in/yaml.v3"
)

func isTruthy(c *CandidateNode) (bool, error) {
	node := c.Node
	value := true
	if node.Tag == "!!null" {
		return false, nil
	}
	if node.Kind == yaml.ScalarNode && node.Tag == "!!bool" {
		errDecoding := node.Decode(&value)
		if errDecoding != nil {
			return false, errDecoding
		}

	}
	return value, nil
}

type boolOp func(bool, bool) bool

func booleanOp(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode, op boolOp) (*list.List, error) {
	var results = list.New()

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		lhs, err := d.getMatchingNodes(nodeToMap(candidate), pathNode.Lhs)
		if err != nil {
			return nil, err
		}
		rhs, err := d.getMatchingNodes(nodeToMap(candidate), pathNode.Rhs)
		if err != nil {
			return nil, err
		}

		for lhsChild := lhs.Front(); lhsChild != nil; lhsChild = lhsChild.Next() {
			lhsCandidate := lhsChild.Value.(*CandidateNode)
			lhsTrue, errDecoding := isTruthy(lhsCandidate)
			if errDecoding != nil {
				return nil, errDecoding
			}

			for rhsChild := rhs.Front(); rhsChild != nil; rhsChild = rhsChild.Next() {
				rhsCandidate := rhsChild.Value.(*CandidateNode)
				rhsTrue, errDecoding := isTruthy(rhsCandidate)
				if errDecoding != nil {
					return nil, errDecoding
				}
				boolResult := createBooleanCandidate(lhsCandidate, op(lhsTrue, rhsTrue))

				results.PushBack(boolResult)
			}
		}

	}
	return results, nil
}

func OrOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("-- orOp")
	return booleanOp(d, matchingNodes, pathNode, func(b1 bool, b2 bool) bool {
		return b1 || b2
	})
}

func AndOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("-- AndOp")
	return booleanOp(d, matchingNodes, pathNode, func(b1 bool, b2 bool) bool {
		return b1 && b2
	})
}
