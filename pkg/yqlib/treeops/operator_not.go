package treeops

import "github.com/elliotchance/orderedmap"

func NotOperator(d *dataTreeNavigator, matchMap *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	log.Debugf("-- notOperation")
	var results = orderedmap.NewOrderedMap()

	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		log.Debug("notOperation checking %v", candidate)
		truthy, errDecoding := isTruthy(candidate)
		if errDecoding != nil {
			return nil, errDecoding
		}
		result := createBooleanCandidate(candidate, !truthy)
		results.Set(result.GetKey(), result)
	}
	return results, nil
}
