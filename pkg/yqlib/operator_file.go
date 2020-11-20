package yqlib

import (
	"container/list"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

func GetFilenameOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("GetFilename")

	var results = list.New()

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := &yaml.Node{Kind: yaml.ScalarNode, Value: candidate.Filename, Tag: "!!str"}
		lengthCand := &CandidateNode{Node: node, Document: candidate.Document, Path: candidate.Path}
		results.PushBack(lengthCand)
	}

	return results, nil
}

func GetFileIndexOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("GetFileIndex")

	var results = list.New()

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%v", candidate.FileIndex), Tag: "!!int"}
		lengthCand := &CandidateNode{Node: node, Document: candidate.Document, Path: candidate.Path}
		results.PushBack(lengthCand)
	}

	return results, nil
}
