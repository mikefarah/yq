package treeops

import "container/list"

func NotOperator(d *dataTreeNavigator, matchMap *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("-- notOperation")
	var results = list.New()

	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		log.Debug("notOperation checking %v", candidate)
		truthy, errDecoding := isTruthy(candidate)
		if errDecoding != nil {
			return nil, errDecoding
		}
		result := createBooleanCandidate(candidate, !truthy)
		results.PushBack(result)
	}
	return results, nil
}
