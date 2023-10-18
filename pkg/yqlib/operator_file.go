package yqlib

import (
	"container/list"
	"fmt"
)

func getFilenameOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("GetFilename")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		result := candidate.CreateReplacement(ScalarNode, "!!str", candidate.GetFilename())
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}

func getFileIndexOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("GetFileIndex")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		result := candidate.CreateReplacement(ScalarNode, "!!int", fmt.Sprintf("%v", candidate.GetFileIndex()))
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}
