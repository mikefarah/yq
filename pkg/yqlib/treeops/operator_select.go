package treeops

import (
	"github.com/elliotchance/orderedmap"
)

func SelectOperator(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {

	log.Debugf("-- selectOperation")
	var results = orderedmap.NewOrderedMap()

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		rhs, err := d.getMatchingNodes(nodeToMap(candidate), pathNode.Rhs)

		if err != nil {
			return nil, err
		}

		// grab the first value
		first := rhs.Front()

		if first != nil {
			result := first.Value.(*CandidateNode)
			includeResult, errDecoding := isTruthy(result)
			if errDecoding != nil {
				return nil, errDecoding
			}

			if includeResult {
				results.Set(candidate.GetKey(), candidate)
			}
		}
	}
	return results, nil
}
