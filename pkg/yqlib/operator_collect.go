package yqlib

import (
	"container/list"

	yaml "gopkg.in/yaml.v3"
)

func collectOperator(d *dataTreeNavigator, matchMap *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("-- collectOperation")

	if matchMap.Len() == 0 {
		node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Value: "[]"}
		candidate := &CandidateNode{Node: node}
		return nodeToMap(candidate), nil
	}

	var results = list.New()

	node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	var collectC *CandidateNode
	if matchMap.Front() != nil {
		collectC = matchMap.Front().Value.(*CandidateNode).CreateChild(nil, node)
		if len(collectC.Path) > 0 {
			collectC.Path = collectC.Path[:len(collectC.Path)-1]
		}
	} else {
		collectC = &CandidateNode{Node: node}
	}

	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		log.Debugf("Collecting %v", NodeToString(candidate))
		node.Content = append(node.Content, unwrapDoc(candidate.Node))
	}

	results.PushBack(collectC)

	return results, nil
}
