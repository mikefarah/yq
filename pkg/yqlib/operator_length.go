package yqlib

import (
	"container/list"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

func LengthOperator(d *dataTreeNavigator, matchMap *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("-- lengthOperation")
	var results = list.New()

	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		targetNode := UnwrapDoc(candidate.Node)
		var length int
		switch targetNode.Kind {
		case yaml.ScalarNode:
			length = len(targetNode.Value)
		case yaml.MappingNode:
			length = len(targetNode.Content) / 2
		case yaml.SequenceNode:
			length = len(targetNode.Content)
		default:
			length = 0
		}

		node := &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%v", length), Tag: "!!int"}
		lengthCand := &CandidateNode{Node: node, Document: candidate.Document, Path: candidate.Path}
		results.PushBack(lengthCand)
	}

	return results, nil
}
