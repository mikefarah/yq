package yqlib

import (
	"fmt"

	"container/list"

	"github.com/elliotchance/orderedmap"
	yaml "gopkg.in/yaml.v3"
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

	for el := matchMap.Front(); el != nil; el = el.Next() {
		newNodes, err := traverse(d, el.Value.(*CandidateNode), pathNode.Operation)
		if err != nil {
			return nil, err
		}
		matchingNodeMap.PushBackList(newNodes)
	}

	return matchingNodeMap, nil
}

func traverse(d *dataTreeNavigator, matchingNode *CandidateNode, operation *Operation) (*list.List, error) {
	log.Debug("Traversing %v", NodeToString(matchingNode))
	value := matchingNode.Node

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
		matchingNode.Node = matchingNode.Node.Alias
		return traverse(d, matchingNode, operation)
	case yaml.DocumentNode:
		log.Debug("digging into doc node")
		return traverse(d, &CandidateNode{
			Node:     matchingNode.Node.Content[0],
			Document: matchingNode.Document}, operation)
	default:
		return list.New(), nil
	}
}

func keyMatches(key *yaml.Node, pathNode *Operation) bool {
	return pathNode.Value == "[]" || Match(key.Value, pathNode.StringValue)
}

func traverseMap(matchingNode *CandidateNode, operation *Operation) (*list.List, error) {
	var newMatches = orderedmap.NewOrderedMap()
	err := doTraverseMap(newMatches, matchingNode, operation)

	if err != nil {
		return nil, err
	}

	if newMatches.Len() == 0 {
		//no matches, create one automagically
		valueNode := &yaml.Node{Tag: "!!null", Kind: yaml.ScalarNode, Value: "null"}
		node := matchingNode.Node
		node.Content = append(node.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: operation.StringValue}, valueNode)
		candidateNode := &CandidateNode{
			Node:     valueNode,
			Path:     append(matchingNode.Path, operation.StringValue),
			Document: matchingNode.Document,
		}
		newMatches.Set(candidateNode.GetKey(), candidateNode)

	}

	results := list.New()
	i := 0
	for el := newMatches.Front(); el != nil; el = el.Next() {
		results.PushBack(el.Value)
		i++
	}
	return results, nil
}

func doTraverseMap(newMatches *orderedmap.OrderedMap, candidate *CandidateNode, operation *Operation) error {
	// value.Content is a concatenated array of key, value,
	// so keys are in the even indexes, values in odd.
	// merge aliases are defined first, but we only want to traverse them
	// if we don't find a match directly on this node first.
	//TODO ALIASES, auto creation?

	node := candidate.Node

	followAlias := true

	if operation.Preferences != nil {
		followAlias = !operation.Preferences.(*TraversePreferences).DontFollowAlias
	}

	var contents = node.Content
	for index := 0; index < len(contents); index = index + 2 {
		key := contents[index]
		value := contents[index+1]

		log.Debug("checking %v (%v)", key.Value, key.Tag)
		//skip the 'merge' tag, find a direct match first
		if key.Tag == "!!merge" && followAlias {
			log.Debug("Merge anchor")
			err := traverseMergeAnchor(newMatches, candidate, value, operation)
			if err != nil {
				return err
			}
		} else if keyMatches(key, operation) {
			log.Debug("MATCHED")
			candidateNode := &CandidateNode{
				Node:     value,
				Path:     candidate.CreateChildPath(key.Value),
				Document: candidate.Document,
			}
			newMatches.Set(candidateNode.GetKey(), candidateNode)
		}
	}

	return nil
}

func traverseMergeAnchor(newMatches *orderedmap.OrderedMap, originalCandidate *CandidateNode, value *yaml.Node, operation *Operation) error {
	switch value.Kind {
	case yaml.AliasNode:
		candidateNode := &CandidateNode{
			Node:     value.Alias,
			Path:     originalCandidate.Path,
			Document: originalCandidate.Document,
		}
		return doTraverseMap(newMatches, candidateNode, operation)
	case yaml.SequenceNode:
		for _, childValue := range value.Content {
			err := traverseMergeAnchor(newMatches, originalCandidate, childValue, operation)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func traverseArray(candidate *CandidateNode, operation *Operation) (*list.List, error) {
	log.Debug("operation Value %v", operation.Value)
	if operation.Value == "[]" {

		var contents = candidate.Node.Content
		var newMatches = list.New()
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

	switch operation.Value.(type) {
	case int64:
		index := operation.Value.(int64)
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

		return nodeToMap(&CandidateNode{
			Node:     candidate.Node.Content[indexToUse],
			Document: candidate.Document,
			Path:     candidate.CreateChildPath(index),
		}), nil
	default:
		log.Debug("argument not an int (%v), no array matches", operation.Value)
		return list.New(), nil
	}

}
