package yqlib

import (
	"container/list"
	"fmt"

	"gopkg.in/yaml.v3"
)

func getDocumentIndexOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	var results = list.New()

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%v", candidate.Document), Tag: "!!int"}
		scalar := candidate.CreateChild(nil, node)
		results.PushBack(scalar)
	}
	return results, nil
}
