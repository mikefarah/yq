package treeops

import (
	"github.com/elliotchance/orderedmap"
	"gopkg.in/yaml.v3"
)

func isTruthy(c *CandidateNode) (bool, error) {
	node := c.Node
	value := true
	if node.Kind == yaml.ScalarNode && node.Tag == "!!bool" {
		errDecoding := node.Decode(&value)
		if errDecoding != nil {
			return false, errDecoding
		}

	}
	return value, nil
}

type boolOp func(bool, bool) bool

func booleanOp(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode, op boolOp) (*orderedmap.OrderedMap, error) {
	var results = orderedmap.NewOrderedMap()

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

				results.Set(boolResult.GetKey(), boolResult)
			}
		}

	}
	return results, nil
}

func OrOperator(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	log.Debugf("-- orOp")
	return booleanOp(d, matchingNodes, pathNode, func(b1 bool, b2 bool) bool {
		return b1 || b2
	})
}

func AndOperator(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	log.Debugf("-- AndOp")
	return booleanOp(d, matchingNodes, pathNode, func(b1 bool, b2 bool) bool {
		return b1 && b2
	})
}
