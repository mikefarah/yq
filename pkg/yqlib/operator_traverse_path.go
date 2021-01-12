package yqlib

import (
	"container/list"
	"fmt"
	"strconv"

	"github.com/elliotchance/orderedmap"
	yaml "gopkg.in/yaml.v3"
)

type traversePreferences struct {
	FollowAlias    bool
	IncludeMapKeys bool
}

func splat(d *dataTreeNavigator, matches *list.List, prefs *traversePreferences) (*list.List, error) {
	return traverseNodesWithArrayIndices(matches, make([]*yaml.Node, 0), prefs)
}

func traversePathOperator(d *dataTreeNavigator, matchMap *list.List, expressionNode *ExpressionNode) (*list.List, error) {
	log.Debugf("-- Traversing")
	var matchingNodeMap = list.New()

	for el := matchMap.Front(); el != nil; el = el.Next() {
		newNodes, err := traverse(d, el.Value.(*CandidateNode), expressionNode.Operation)
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
		prefs := &traversePreferences{FollowAlias: true}
		return traverseMap(matchingNode, operation.StringValue, prefs, false)

	case yaml.SequenceNode:
		log.Debug("its a sequence of %v things!", len(value.Content))
		return traverseArray(matchingNode, operation)

	case yaml.AliasNode:
		log.Debug("its an alias!")
		matchingNode.Node = matchingNode.Node.Alias
		return traverse(d, matchingNode, operation)
	case yaml.DocumentNode:
		log.Debug("digging into doc node")

		return traverse(d, matchingNode.CreateChild(nil, matchingNode.Node.Content[0]), operation)
	default:
		return list.New(), nil
	}
}

func traverseArrayOperator(d *dataTreeNavigator, matchingNodes *list.List, expressionNode *ExpressionNode) (*list.List, error) {
	// rhs is a collect expression that will yield indexes to retreive of the arrays

	rhs, err := d.GetMatchingNodes(matchingNodes, expressionNode.Rhs)
	if err != nil {
		return nil, err
	}

	var indicesToTraverse = rhs.Front().Value.(*CandidateNode).Node.Content
	prefs := &traversePreferences{FollowAlias: true}
	return traverseNodesWithArrayIndices(matchingNodes, indicesToTraverse, prefs)
}

func traverseNodesWithArrayIndices(matchingNodes *list.List, indicesToTraverse []*yaml.Node, prefs *traversePreferences) (*list.List, error) {
	var matchingNodeMap = list.New()
	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		newNodes, err := traverseArrayIndices(candidate, indicesToTraverse, prefs)
		if err != nil {
			return nil, err
		}
		matchingNodeMap.PushBackList(newNodes)
	}

	return matchingNodeMap, nil
}

func traverseArrayIndices(matchingNode *CandidateNode, indicesToTraverse []*yaml.Node, prefs *traversePreferences) (*list.List, error) { // call this if doc / alias like the other traverse
	node := matchingNode.Node
	if node.Tag == "!!null" {
		log.Debugf("OperatorArrayTraverse got a null - turning it into an empty array")
		// auto vivification, make it into an empty array
		node.Tag = ""
		node.Kind = yaml.SequenceNode
	}

	if node.Kind == yaml.AliasNode {
		matchingNode.Node = node.Alias
		return traverseArrayIndices(matchingNode, indicesToTraverse, prefs)
	} else if node.Kind == yaml.SequenceNode {
		return traverseArrayWithIndices(matchingNode, indicesToTraverse)
	} else if node.Kind == yaml.MappingNode {
		return traverseMapWithIndices(matchingNode, indicesToTraverse, prefs)
	} else if node.Kind == yaml.DocumentNode {
		return traverseArrayIndices(matchingNode.CreateChild(nil, matchingNode.Node.Content[0]), indicesToTraverse, prefs)
	}
	log.Debugf("OperatorArrayTraverse skipping %v as its a %v", matchingNode, node.Tag)
	return list.New(), nil
}

func traverseMapWithIndices(candidate *CandidateNode, indices []*yaml.Node, prefs *traversePreferences) (*list.List, error) {
	if len(indices) == 0 {
		return traverseMap(candidate, "", prefs, true)
	}

	var matchingNodeMap = list.New()

	for _, indexNode := range indices {
		log.Debug("traverseMapWithIndices: %v", indexNode.Value)
		newNodes, err := traverseMap(candidate, indexNode.Value, prefs, false)
		if err != nil {
			return nil, err
		}
		matchingNodeMap.PushBackList(newNodes)
	}

	return matchingNodeMap, nil
}

func traverseArrayWithIndices(candidate *CandidateNode, indices []*yaml.Node) (*list.List, error) {
	log.Debug("traverseArrayWithIndices")
	var newMatches = list.New()
	node := unwrapDoc(candidate.Node)
	if len(indices) == 0 {
		log.Debug("splatting")
		var index int64
		for index = 0; index < int64(len(node.Content)); index = index + 1 {

			newMatches.PushBack(candidate.CreateChild(index, node.Content[index]))
		}
		return newMatches, nil

	}

	for _, indexNode := range indices {
		log.Debug("traverseArrayWithIndices: '%v'", indexNode.Value)
		index, err := strconv.ParseInt(indexNode.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Cannot index array with '%v' (%v)", indexNode.Value, err)
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

		newMatches.PushBack(candidate.CreateChild(index, node.Content[indexToUse]))
	}
	return newMatches, nil
}

func keyMatches(key *yaml.Node, wantedKey string) bool {
	return matchKey(key.Value, wantedKey)
}

func traverseMap(matchingNode *CandidateNode, key string, prefs *traversePreferences, splat bool) (*list.List, error) {
	var newMatches = orderedmap.NewOrderedMap()
	err := doTraverseMap(newMatches, matchingNode, key, prefs, splat)

	if err != nil {
		return nil, err
	}

	if newMatches.Len() == 0 {
		//no matches, create one automagically
		valueNode := &yaml.Node{Tag: "!!null", Kind: yaml.ScalarNode, Value: "null"}
		node := matchingNode.Node
		node.Content = append(node.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: key}, valueNode)
		candidateNode := matchingNode.CreateChild(key, valueNode)
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

func doTraverseMap(newMatches *orderedmap.OrderedMap, candidate *CandidateNode, wantedKey string, prefs *traversePreferences, splat bool) error {
	// value.Content is a concatenated array of key, value,
	// so keys are in the even indexes, values in odd.
	// merge aliases are defined first, but we only want to traverse them
	// if we don't find a match directly on this node first.

	node := candidate.Node

	var contents = node.Content
	for index := 0; index < len(contents); index = index + 2 {
		key := contents[index]
		value := contents[index+1]

		log.Debug("checking %v (%v)", key.Value, key.Tag)
		//skip the 'merge' tag, find a direct match first
		if key.Tag == "!!merge" && prefs.FollowAlias {
			log.Debug("Merge anchor")
			err := traverseMergeAnchor(newMatches, candidate, value, wantedKey, prefs, splat)
			if err != nil {
				return err
			}
		} else if splat || keyMatches(key, wantedKey) {
			log.Debug("MATCHED")
			if prefs.IncludeMapKeys {
				candidateNode := candidate.CreateChild(key.Value, key)
				newMatches.Set(fmt.Sprintf("keyOf-%v", candidateNode.GetKey()), candidateNode)
			}
			candidateNode := candidate.CreateChild(key.Value, value)
			newMatches.Set(candidateNode.GetKey(), candidateNode)
		}
	}

	return nil
}

func traverseMergeAnchor(newMatches *orderedmap.OrderedMap, originalCandidate *CandidateNode, value *yaml.Node, wantedKey string, prefs *traversePreferences, splat bool) error {
	switch value.Kind {
	case yaml.AliasNode:
		candidateNode := originalCandidate.CreateChild(nil, value.Alias)
		return doTraverseMap(newMatches, candidateNode, wantedKey, prefs, splat)
	case yaml.SequenceNode:
		for _, childValue := range value.Content {
			err := traverseMergeAnchor(newMatches, originalCandidate, childValue, wantedKey, prefs, splat)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func traverseArray(candidate *CandidateNode, operation *Operation) (*list.List, error) {
	log.Debug("operation Value %v", operation.Value)
	indices := []*yaml.Node{&yaml.Node{Value: operation.StringValue}}
	return traverseArrayWithIndices(candidate, indices)
}
