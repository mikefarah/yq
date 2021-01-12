package yqlib

import (
	"container/list"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

func getFilenameOperator(d *dataTreeNavigator, matchingNodes *list.List, expressionNode *ExpressionNode) (*list.List, error) {
	log.Debugf("GetFilename")

	var results = list.New()

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := &yaml.Node{Kind: yaml.ScalarNode, Value: candidate.Filename, Tag: "!!str"}
		result := candidate.CreateChild(nil, node)
		results.PushBack(result)
	}

	return results, nil
}

func getFileIndexOperator(d *dataTreeNavigator, matchingNodes *list.List, expressionNode *ExpressionNode) (*list.List, error) {
	log.Debugf("GetFileIndex")

	var results = list.New()

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%v", candidate.FileIndex), Tag: "!!int"}
		result := candidate.CreateChild(nil, node)
		results.PushBack(result)
	}

	return results, nil
}
