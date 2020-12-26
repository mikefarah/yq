package yqlib

import (
	"container/list"
	"fmt"
	"strconv"

	yaml "gopkg.in/yaml.v3"
)

func TraverseArrayOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	// lhs is an expression that will yield a bunch of arrays
	// rhs is a collect expression that will yield indexes to retreive of the arrays

	lhs, err := d.GetMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}

	rhs, err := d.GetMatchingNodes(matchingNodes, pathNode.Rhs)
	if err != nil {
		return nil, err
	}

	var indicesToTraverse = rhs.Front().Value.(*CandidateNode).Node.Content

	var matchingNodeMap = list.New()
	for el := lhs.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		if candidate.Node.Kind == yaml.SequenceNode {
			newNodes, err := traverseArrayWithIndices(candidate, indicesToTraverse)
			if err != nil {
				return nil, err
			}
			matchingNodeMap.PushBackList(newNodes)
		} else {
			log.Debugf("OperatorArray Traverse skipping %v as its a %v", candidate, candidate.Node.Tag)
		}
	}

	return matchingNodeMap, nil
}

func traverseArrayWithIndices(candidate *CandidateNode, indices []*yaml.Node) (*list.List, error) {
	log.Debug("traverseArrayWithIndices")
	var newMatches = list.New()

	if len(indices) == 0 {
		var contents = candidate.Node.Content
		var index int64
		for index = 0; index < int64(len(contents)); index = index + 1 {

			newMatches.PushBack(&CandidateNode{
				Document: candidate.Document,
				Path:     candidate.CreateChildPath(index),
				Node:     contents[index],
			})
		}
		return newMatches, nil

	}

	for _, indexNode := range indices {
		index, err := strconv.ParseInt(indexNode.Value, 10, 64)
		if err != nil {
			return nil, err
		}
		indexToUse := index
		contentLength := int64(len(candidate.Node.Content))
		for contentLength <= index {
			candidate.Node.Content = append(candidate.Node.Content, &yaml.Node{Tag: "!!null", Kind: yaml.ScalarNode, Value: "null"})
			contentLength = int64(len(candidate.Node.Content))
		}

		if indexToUse < 0 {
			indexToUse = contentLength + indexToUse
		}

		if indexToUse < 0 {
			return nil, fmt.Errorf("Index [%v] out of range, array size is %v", index, contentLength)
		}

		newMatches.PushBack(&CandidateNode{
			Node:     candidate.Node.Content[indexToUse],
			Document: candidate.Document,
			Path:     candidate.CreateChildPath(index),
		})
	}
	return newMatches, nil
}
