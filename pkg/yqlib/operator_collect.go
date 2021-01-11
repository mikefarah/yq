package yqlib

import (
	"container/list"

	yaml "gopkg.in/yaml.v3"
)

func CollectOperator(d *dataTreeNavigator, matchMap *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("-- collectOperation")

	if matchMap.Len() == 0 {
		node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Value: "[]"}
		candidate := &CandidateNode{Node: node}
		return nodeToMap(candidate), nil
	}

	var results = list.New()

	node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}

	var document uint = 0
	var path []interface{}

	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		log.Debugf("Collecting %v", NodeToString(candidate))
		if path == nil && candidate.Path != nil && len(candidate.Path) > 1 {
			path = candidate.Path[:len(candidate.Path)-1]
			document = candidate.Document
		}
		node.Content = append(node.Content, UnwrapDoc(candidate.Node))
	}

	collectC := &CandidateNode{Node: node, Document: document, Path: path}
	results.PushBack(collectC)

	return results, nil
}
