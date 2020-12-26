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

	rhs, err := d.GetMatchingNodes(matchingNodes, pathNode.Rhs)
	if err != nil {
		return nil, err
	}

	var indicesToTraverse = rhs.Front().Value.(*CandidateNode).Node.Content

	var matchingNodeMap = list.New()
	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := UnwrapDoc(candidate.Node)
		if node.Tag == "!!null" {
			log.Debugf("OperatorArrayTraverse got a null - turning it into an empty array")
			// auto vivification, make it into an empty array
			node.Tag = ""
			node.Kind = yaml.SequenceNode
		} else if node.Kind == yaml.AliasNode {
			candidate.Node = node.Alias
			node = node.Alias
		}

		if node.Kind == yaml.SequenceNode {
			newNodes, err := traverseArrayWithIndices(candidate, indicesToTraverse)
			if err != nil {
				return nil, err
			}
			matchingNodeMap.PushBackList(newNodes)
		} else if node.Kind == yaml.MappingNode && len(indicesToTraverse) == 0 {
			// splat the map
			newNodes, err := traverseMapWithIndices(candidate, indicesToTraverse)
			if err != nil {
				return nil, err
			}
			matchingNodeMap.PushBackList(newNodes)
		} else {
			log.Debugf("OperatorArrayTraverse skipping %v as its a %v", candidate, node.Tag)
		}
	}

	return matchingNodeMap, nil
}

func traverseMapWithIndices(candidate *CandidateNode, indices []*yaml.Node) (*list.List, error) {
	//REWRITE TO USE TRAVERSE MAP

	node := UnwrapDoc(candidate.Node)
	var contents = node.Content
	var matchingNodeMap = list.New()
	if len(indices) == 0 {
		for index := 0; index < len(contents); index = index + 2 {
			key := contents[index]
			value := contents[index+1]
			matchingNodeMap.PushBack(&CandidateNode{
				Node:     value,
				Path:     candidate.CreateChildPath(key.Value),
				Document: candidate.Document,
			})
		}
		return matchingNodeMap, nil
	}

	for index := 0; index < len(contents); index = index + 2 {
		key := contents[index]
		value := contents[index+1]
		for _, indexNode := range indices {
			if key.Value == indexNode.Value {
				matchingNodeMap.PushBack(&CandidateNode{
					Node:     value,
					Path:     candidate.CreateChildPath(key.Value),
					Document: candidate.Document,
				})
			}
		}

	}
	return matchingNodeMap, nil
}

func traverseArrayWithIndices(candidate *CandidateNode, indices []*yaml.Node) (*list.List, error) {
	log.Debug("traverseArrayWithIndices")
	var newMatches = list.New()
	node := UnwrapDoc(candidate.Node)

	if len(indices) == 0 {
		var index int64
		for index = 0; index < int64(len(node.Content)); index = index + 1 {

			newMatches.PushBack(&CandidateNode{
				Document: candidate.Document,
				Path:     candidate.CreateChildPath(index),
				Node:     node.Content[index],
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
		contentLength := int64(len(node.Content))
		for contentLength <= index {
			node.Content = append(node.Content, &yaml.Node{Tag: "!!null", Kind: yaml.ScalarNode, Value: "null"})
			contentLength = int64(len(node.Content))
		}

		if indexToUse < 0 {
			indexToUse = contentLength + indexToUse
		}

		if indexToUse < 0 {
			return nil, fmt.Errorf("Index [%v] out of range, array size is %v", index, contentLength)
		}

		newMatches.PushBack(&CandidateNode{
			Node:     node.Content[indexToUse],
			Document: candidate.Document,
			Path:     candidate.CreateChildPath(index),
		})
	}
	return newMatches, nil
}
