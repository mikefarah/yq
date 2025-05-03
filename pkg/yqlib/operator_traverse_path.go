package yqlib

import (
	"container/list"
	"fmt"

	"github.com/elliotchance/orderedmap"
)

type traversePreferences struct {
	DontFollowAlias      bool
	IncludeMapKeys       bool
	DontAutoCreate       bool // by default, we automatically create entries on the fly.
	DontIncludeMapValues bool
	OptionalTraverse     bool // e.g. .adf?
}

func splat(context Context, prefs traversePreferences) (Context, error) {
	return traverseNodesWithArrayIndices(context, make([]*CandidateNode, 0), prefs)
}

func traversePathOperator(_ *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("traversePathOperator")
	var matches = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		newNodes, err := traverse(context, el.Value.(*CandidateNode), expressionNode.Operation)
		if err != nil {
			return Context{}, err
		}
		matches.PushBackList(newNodes)
	}

	return context.ChildContext(matches), nil
}

func traverse(context Context, matchingNode *CandidateNode, operation *Operation) (*list.List, error) {
	log.Debug("Traversing %v", NodeToString(matchingNode))

	if matchingNode.Tag == "!!null" && operation.Value != "[]" && !context.DontAutoCreate {
		log.Debugf("Guessing kind")
		// we must have added this automatically, lets guess what it should be now
		switch operation.Value.(type) {
		case int, int64:
			log.Debugf("probably an array")
			matchingNode.Kind = SequenceNode
		default:
			log.Debugf("probably a map")
			matchingNode.Kind = MappingNode
		}
		matchingNode.Tag = ""
	}

	switch matchingNode.Kind {
	case MappingNode:
		log.Debug("its a map with %v entries", len(matchingNode.Content)/2)
		return traverseMap(context, matchingNode, createStringScalarNode(operation.StringValue), operation.Preferences.(traversePreferences), false)

	case SequenceNode:
		log.Debug("its a sequence of %v things!", len(matchingNode.Content))
		return traverseArray(matchingNode, operation, operation.Preferences.(traversePreferences))

	case AliasNode:
		log.Debug("its an alias!")
		matchingNode = matchingNode.Alias
		return traverse(context, matchingNode, operation)
	default:
		return list.New(), nil
	}
}

func traverseArrayOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	//lhs may update the variable context, we should pass that into the RHS
	// BUT we still return the original context back (see jq)
	// https://stedolan.github.io/jq/manual/#Variable/SymbolicBindingOperator:...as$identifier|...

	log.Debugf("--traverseArrayOperator")

	if expressionNode.RHS != nil && expressionNode.RHS.RHS != nil && expressionNode.RHS.RHS.Operation.OperationType == createMapOpType {
		return sliceArrayOperator(d, context, expressionNode.RHS.RHS)
	}

	lhs, err := d.GetMatchingNodes(context, expressionNode.LHS)
	if err != nil {
		return Context{}, err
	}

	// rhs is a collect expression that will yield indices to retrieve of the arrays

	rhs, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.RHS)

	if err != nil {
		return Context{}, err
	}
	prefs := traversePreferences{}

	if expressionNode.Operation.Preferences != nil {
		prefs = expressionNode.Operation.Preferences.(traversePreferences)
	}
	var indicesToTraverse = rhs.MatchingNodes.Front().Value.(*CandidateNode).Content

	log.Debugf("indicesToTraverse %v", len(indicesToTraverse))

	//now we traverse the result of the lhs against the indices we found
	result, err := traverseNodesWithArrayIndices(lhs, indicesToTraverse, prefs)
	if err != nil {
		return Context{}, err
	}
	return context.ChildContext(result.MatchingNodes), nil
}

func traverseNodesWithArrayIndices(context Context, indicesToTraverse []*CandidateNode, prefs traversePreferences) (Context, error) {
	var matchingNodeMap = list.New()
	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		newNodes, err := traverseArrayIndices(context, candidate, indicesToTraverse, prefs)
		if err != nil {
			return Context{}, err
		}
		matchingNodeMap.PushBackList(newNodes)
	}

	return context.ChildContext(matchingNodeMap), nil
}

func traverseArrayIndices(context Context, matchingNode *CandidateNode, indicesToTraverse []*CandidateNode, prefs traversePreferences) (*list.List, error) { // call this if doc / alias like the other traverse
	if matchingNode.Tag == "!!null" {
		log.Debugf("OperatorArrayTraverse got a null - turning it into an empty array")
		// auto vivification
		matchingNode.Tag = ""
		matchingNode.Kind = SequenceNode
		//check that the indices are numeric, if not, then we should create an object
		if len(indicesToTraverse) != 0 && indicesToTraverse[0].Tag != "!!int" {
			matchingNode.Kind = MappingNode
		}
	}

	switch matchingNode.Kind {
	case AliasNode:
		matchingNode = matchingNode.Alias
		return traverseArrayIndices(context, matchingNode, indicesToTraverse, prefs)
	case SequenceNode:
		return traverseArrayWithIndices(matchingNode, indicesToTraverse, prefs)
	case MappingNode:
		return traverseMapWithIndices(context, matchingNode, indicesToTraverse, prefs)
	}
	log.Debugf("OperatorArrayTraverse skipping %v as its a %v", matchingNode, matchingNode.Tag)
	return list.New(), nil
}

func traverseMapWithIndices(context Context, candidate *CandidateNode, indices []*CandidateNode, prefs traversePreferences) (*list.List, error) {
	if len(indices) == 0 {
		return traverseMap(context, candidate, createStringScalarNode(""), prefs, true)
	}

	var matchingNodeMap = list.New()

	for _, indexNode := range indices {
		log.Debug("traverseMapWithIndices: %v", indexNode.Value)
		newNodes, err := traverseMap(context, candidate, indexNode, prefs, false)
		if err != nil {
			return nil, err
		}
		matchingNodeMap.PushBackList(newNodes)
	}

	return matchingNodeMap, nil
}

func traverseArrayWithIndices(node *CandidateNode, indices []*CandidateNode, prefs traversePreferences) (*list.List, error) {
	log.Debug("traverseArrayWithIndices")
	var newMatches = list.New()
	if len(indices) == 0 {
		log.Debug("splatting")
		var index int
		for index = 0; index < len(node.Content); index = index + 1 {
			newMatches.PushBack(node.Content[index])
		}
		return newMatches, nil

	}

	for _, indexNode := range indices {
		log.Debug("traverseArrayWithIndices: '%v'", indexNode.Value)
		index, err := parseInt(indexNode.Value)
		if err != nil && prefs.OptionalTraverse {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("cannot index array with '%v' (%w)", indexNode.Value, err)
		}
		indexToUse := index
		contentLength := len(node.Content)
		for contentLength <= index {
			if contentLength == 0 {
				// default to nice yaml formatting
				node.Style = 0
			}

			valueNode := createScalarNode(nil, "null")
			node.AddChild(valueNode)
			contentLength = len(node.Content)
		}

		if indexToUse < 0 {
			indexToUse = contentLength + indexToUse
		}

		if indexToUse < 0 {
			return nil, fmt.Errorf("index [%v] out of range, array size is %v", index, contentLength)
		}

		newMatches.PushBack(node.Content[indexToUse])
	}
	return newMatches, nil
}

func keyMatches(key *CandidateNode, wantedKey string) bool {
	return matchKey(key.Value, wantedKey)
}

func traverseMap(context Context, matchingNode *CandidateNode, keyNode *CandidateNode, prefs traversePreferences, splat bool) (*list.List, error) {
	var newMatches = orderedmap.NewOrderedMap()
	err := doTraverseMap(newMatches, matchingNode, keyNode.Value, prefs, splat)

	if err != nil {
		return nil, err
	}

	if !splat && !prefs.DontAutoCreate && !context.DontAutoCreate && newMatches.Len() == 0 {
		log.Debugf("no matches, creating one for %v", NodeToString(keyNode))
		//no matches, create one automagically
		valueNode := matchingNode.CreateChild()
		valueNode.Kind = ScalarNode
		valueNode.Tag = "!!null"
		valueNode.Value = "null"

		if len(matchingNode.Content) == 0 {
			matchingNode.Style = 0
		}

		keyNode, valueNode = matchingNode.AddKeyValueChild(keyNode, valueNode)

		if prefs.IncludeMapKeys {
			newMatches.Set(keyNode.GetKey(), keyNode)
		}
		if !prefs.DontIncludeMapValues {
			newMatches.Set(valueNode.GetKey(), valueNode)
		}
	}

	results := list.New()
	i := 0
	for el := newMatches.Front(); el != nil; el = el.Next() {
		results.PushBack(el.Value)
		i++
	}
	return results, nil
}

func doTraverseMap(newMatches *orderedmap.OrderedMap, node *CandidateNode, wantedKey string, prefs traversePreferences, splat bool) error {
	// value.Content is a concatenated array of key, value,
	// so keys are in the even indices, values in odd.
	// merge aliases are defined first, but we only want to traverse them
	// if we don't find a match directly on this node first.

	var contents = node.Content
	for index := 0; index+1 < len(contents); index = index + 2 {
		key := contents[index]
		value := contents[index+1]

		//skip the 'merge' tag, find a direct match first
		if key.Tag == "!!merge" && !prefs.DontFollowAlias && wantedKey != "<<" {
			log.Debug("Merge anchor")
			err := traverseMergeAnchor(newMatches, value, wantedKey, prefs, splat)
			if err != nil {
				return err
			}
		} else if splat || keyMatches(key, wantedKey) {
			log.Debug("MATCHED")
			if prefs.IncludeMapKeys {
				log.Debug("including key")
				newMatches.Set(key.GetKey(), key)
			}
			if !prefs.DontIncludeMapValues {
				log.Debug("including value")
				newMatches.Set(value.GetKey(), value)
			}
		}
	}

	return nil
}

func traverseMergeAnchor(newMatches *orderedmap.OrderedMap, value *CandidateNode, wantedKey string, prefs traversePreferences, splat bool) error {
	switch value.Kind {
	case AliasNode:
		if value.Alias.Kind != MappingNode {
			return fmt.Errorf("can only use merge anchors with maps (!!map), but got %v", value.Alias.Tag)
		}
		return doTraverseMap(newMatches, value.Alias, wantedKey, prefs, splat)
	case SequenceNode:
		for _, childValue := range value.Content {
			err := traverseMergeAnchor(newMatches, childValue, wantedKey, prefs, splat)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func traverseArray(candidate *CandidateNode, operation *Operation, prefs traversePreferences) (*list.List, error) {
	log.Debug("operation Value %v", operation.Value)
	indices := []*CandidateNode{{Value: operation.StringValue}}
	return traverseArrayWithIndices(candidate, indices, prefs)
}
