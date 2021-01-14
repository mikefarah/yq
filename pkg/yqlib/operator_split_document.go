package yqlib

import (
	"container/list"
)

func splitDocumentOperator(d *dataTreeNavigator, matchMap *list.List, expressionNode *ExpressionNode) (*list.List, error) {
	log.Debugf("-- splitDocumentOperator")

	var index uint = 0
	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		candidate.Document = index
		index = index + 1
	}

	return matchMap, nil
}
