package treeops

import (
	"fmt"

	"container/list"

	"gopkg.in/yaml.v3"
)

type TraversePreferences struct {
	DontFollowAlias bool
}

func Splat(d *dataTreeNavigator, matches *list.List) (*list.List, error) {
	preferences := &TraversePreferences{DontFollowAlias: true}
	splatOperation := &Operation{OperationType: TraversePath, Value: "[]", Preferences: preferences}
	splatTreeNode := &PathTreeNode{Operation: splatOperation}
	return TraversePathOperator(d, matches, splatTreeNode)
}

func TraversePathOperator(d *dataTreeNavigator, matchMap *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("-- Traversing")
	var matchingNodeMap = list.New()
	var newNodes []*CandidateNode
	var err error

	for el := matchMap.Front(); el != nil; el = el.Next() {
		newNodes, err = traverse(d, el.Value.(*CandidateNode), pathNode.Operation)
		if err != nil {
			return nil, err
		}
		for _, n := range newNodes {
			matchingNodeMap.PushBack(n)
		}
	}

	return matchingNodeMap, nil
}

func traverse(d *dataTreeNavigator, matchingNode *CandidateNode, operation *Operation) ([]*CandidateNode, error) {
	log.Debug("Traversing %v", NodeToString(matchingNode))
	value := matchingNode.Node

	followAlias := true

	// if operation.Preferences != nil {
	// 	followAlias = !operation.Preferences.(*TraversePreferences).DontFollowAlias
	// }

	if value.Tag == "!!null" && operation.Value != "[]" {
		log.Debugf("Guessing kind")
		// we must ahve added this automatically, lets guess what it should be now
		switch operation.Value.(type) {
		case int, int64:
			log.Debugf("probably an array")
			value.Kind = yaml.SequenceNode
		default:
			log.Debugf("probably a map")
			value.Kind = yaml.MappingNode
		}
		value.Tag = ""
	}

	switch value.Kind {
	case yaml.MappingNode:
		log.Debug("its a map with %v entries", len(value.Content)/2)
		return traverseMap(matchingNode, operation)

	case yaml.SequenceNode:
		log.Debug("its a sequence of %v things!", len(value.Content))
		return traverseArray(matchingNode, operation)

	case yaml.AliasNode:
		log.Debug("its an alias!")
		if followAlias {
			matchingNode.Node = matchingNode.Node.Alias
			return traverse(d, matchingNode, operation)
		}
		return []*CandidateNode{matchingNode}, nil
	case yaml.DocumentNode:
		log.Debug("digging into doc node")
		return traverse(d, &CandidateNode{
			Node:     matchingNode.Node.Content[0],
			Document: matchingNode.Document}, operation)
	default:
		return nil, nil
	}
}

func keyMatches(key *yaml.Node, pathNode *Operation) bool {
	return pathNode.Value == "[]" || Match(key.Value, pathNode.StringValue)
}

func traverseMap(candidate *CandidateNode, pathNode *Operation) ([]*CandidateNode, error) {
	// value.Content is a concatenated array of key, value,
	// so keys are in the even indexes, values in odd.
	// merge aliases are defined first, but we only want to traverse them
	// if we don't find a match directly on this node first.
	//TODO ALIASES, auto creation?

	var newMatches = make([]*CandidateNode, 0)

	node := candidate.Node

	var contents = node.Content
	for index := 0; index < len(contents); index = index + 2 {
		key := contents[index]
		value := contents[index+1]

		log.Debug("checking %v (%v)", key.Value, key.Tag)
		if keyMatches(key, pathNode) {
			log.Debug("MATCHED")
			newMatches = append(newMatches, &CandidateNode{
				Node:     value,
				Path:     append(candidate.Path, key.Value),
				Document: candidate.Document,
			})
		}
	}
	if len(newMatches) == 0 {
		//no matches, create one automagically
		valueNode := &yaml.Node{Tag: "!!null", Kind: yaml.ScalarNode, Value: "null"}
		node.Content = append(node.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: pathNode.StringValue}, valueNode)
		newMatches = append(newMatches, &CandidateNode{
			Node:     valueNode,
			Path:     append(candidate.Path, pathNode.StringValue),
			Document: candidate.Document,
		})

	}

	return newMatches, nil
}

func traverseArray(candidate *CandidateNode, pathNode *Operation) ([]*CandidateNode, error) {
	log.Debug("pathNode Value %v", pathNode.Value)
	if pathNode.Value == "[]" {

		var contents = candidate.Node.Content
		var newMatches = make([]*CandidateNode, len(contents))
		var index int64
		for index = 0; index < int64(len(contents)); index = index + 1 {
			newMatches[index] = &CandidateNode{
				Document: candidate.Document,
				Path:     append(candidate.Path, index),
				Node:     contents[index],
			}
		}
		return newMatches, nil

	}

	index := pathNode.Value.(int64)
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

	return []*CandidateNode{&CandidateNode{
		Node:     candidate.Node.Content[indexToUse],
		Document: candidate.Document,
		Path:     append(candidate.Path, index),
	}}, nil

}
