package yqlib

import (
	"container/list"
	"fmt"

	"github.com/elliotchance/orderedmap"
	yaml "gopkg.in/yaml.v3"
)

type traversePreferences struct {
	DontFollowAlias      bool
	IncludeMapKeys       bool
	DontAutoCreate       bool // by default, we automatically create entries on the fly.
	DontIncludeMapValues bool
	OptionalTraverse     bool // e.g. .adf?
}

func splat(context Context, prefs traversePreferences) (Context, error) {
	return traverseNodesWithArrayIndices(context, make([]*yaml.Node, 0), prefs)
}

func traversePathOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- traversePathOperator")
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
	value := matchingNode.Node

	if value.Tag == "!!null" && operation.Value != "[]" && !context.DontAutoCreate {
		log.Debugf("Guessing kind")
		// we must have added this automatically, lets guess what it should be now
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
		return traverseMap(context, matchingNode, createStringScalarNode(operation.StringValue), operation.Preferences.(traversePreferences), false)

	case yaml.SequenceNode:
		log.Debug("its a sequence of %v things!", len(value.Content))
		return traverseArray(matchingNode, operation, operation.Preferences.(traversePreferences))

	case yaml.AliasNode:
		log.Debug("its an alias!")
		matchingNode.Node = matchingNode.Node.Alias
		return traverse(context, matchingNode, operation)
	case yaml.DocumentNode:
		log.Debug("digging into doc node")

		return traverse(context, matchingNode.CreateChildInMap(nil, matchingNode.Node.Content[0]), operation)
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
	var indicesToTraverse = rhs.MatchingNodes.Front().Value.(*CandidateNode).Node.Content

	log.Debugf("indicesToTraverse %v", len(indicesToTraverse))

	//now we traverse the result of the lhs against the indices we found
	result, err := traverseNodesWithArrayIndices(lhs, indicesToTraverse, prefs)
	if err != nil {
		return Context{}, err
	}
	return context.ChildContext(result.MatchingNodes), nil
}

func traverseNodesWithArrayIndices(context Context, indicesToTraverse []*yaml.Node, prefs traversePreferences) (Context, error) {
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

func traverseArrayIndices(context Context, matchingNode *CandidateNode, indicesToTraverse []*yaml.Node, prefs traversePreferences) (*list.List, error) { // call this if doc / alias like the other traverse
	node := matchingNode.Node
	if node.Tag == "!!null" {
		log.Debugf("OperatorArrayTraverse got a null - turning it into an empty array")
		// auto vivification
		node.Tag = ""
		node.Kind = yaml.SequenceNode
		//check that the indices are numeric, if not, then we should create an object
		if len(indicesToTraverse) != 0 && indicesToTraverse[0].Tag != "!!int" {
			node.Kind = yaml.MappingNode
		}
	}

	if node.Kind == yaml.AliasNode {
		matchingNode.Node = node.Alias
		return traverseArrayIndices(context, matchingNode, indicesToTraverse, prefs)
	} else if node.Kind == yaml.SequenceNode {
		return traverseArrayWithIndices(matchingNode, indicesToTraverse, prefs)
	} else if node.Kind == yaml.MappingNode {
		return traverseMapWithIndices(context, matchingNode, indicesToTraverse, prefs)
	} else if node.Kind == yaml.DocumentNode {
		return traverseArrayIndices(context, matchingNode.CreateChildInMap(nil, matchingNode.Node.Content[0]), indicesToTraverse, prefs)
	}
	log.Debugf("OperatorArrayTraverse skipping %v as its a %v", matchingNode, node.Tag)
	return list.New(), nil
}

func traverseMapWithIndices(context Context, candidate *CandidateNode, indices []*yaml.Node, prefs traversePreferences) (*list.List, error) {
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

func traverseArrayWithIndices(candidate *CandidateNode, indices []*yaml.Node, prefs traversePreferences) (*list.List, error) {
	log.Debug("traverseArrayWithIndices")
	var newMatches = list.New()
	node := unwrapDoc(candidate.Node)
	if len(indices) == 0 {
		log.Debug("splatting")
		var index int
		for index = 0; index < len(node.Content); index = index + 1 {
			newMatches.PushBack(candidate.CreateChildInArray(index, node.Content[index]))
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
				// default to nice yaml formating
				node.Style = 0
			}

			node.Content = append(node.Content, &yaml.Node{Tag: "!!null", Kind: yaml.ScalarNode, Value: "null"})
			contentLength = len(node.Content)
		}

		if indexToUse < 0 {
			indexToUse = contentLength + indexToUse
		}

		if indexToUse < 0 {
			return nil, fmt.Errorf("index [%v] out of range, array size is %v", index, contentLength)
		}

		newMatches.PushBack(candidate.CreateChildInArray(index, node.Content[indexToUse]))
	}
	return newMatches, nil
}

func keyMatches(key *yaml.Node, wantedKey string) bool {
	return matchKey(key.Value, wantedKey)
}

func traverseMap(context Context, matchingNode *CandidateNode, keyNode *yaml.Node, prefs traversePreferences, splat bool) (*list.List, error) {
	var newMatches = orderedmap.NewOrderedMap()
	err := doTraverseMap(newMatches, matchingNode, keyNode.Value, prefs, splat)

	if err != nil {
		return nil, err
	}

	if !splat && !prefs.DontAutoCreate && !context.DontAutoCreate && newMatches.Len() == 0 {
		log.Debugf("no matches, creating one")
		//no matches, create one automagically
		valueNode := &yaml.Node{Tag: "!!null", Kind: yaml.ScalarNode, Value: "null"}

		node := matchingNode.Node

		if len(node.Content) == 0 {
			node.Style = 0
		}

		node.Content = append(node.Content, keyNode, valueNode)

		if prefs.IncludeMapKeys {
			log.Debug("including key")
			candidateNode := matchingNode.CreateChildInMap(keyNode, keyNode)
			candidateNode.IsMapKey = true
			newMatches.Set(fmt.Sprintf("keyOf-%v", candidateNode.GetKey()), candidateNode)
		}
		if !prefs.DontIncludeMapValues {
			log.Debug("including value")
			candidateNode := matchingNode.CreateChildInMap(keyNode, valueNode)
			newMatches.Set(candidateNode.GetKey(), candidateNode)
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

func doTraverseMap(newMatches *orderedmap.OrderedMap, candidate *CandidateNode, wantedKey string, prefs traversePreferences, splat bool) error {
	// value.Content is a concatenated array of key, value,
	// so keys are in the even indices, values in odd.
	// merge aliases are defined first, but we only want to traverse them
	// if we don't find a match directly on this node first.

	node := candidate.Node

	var contents = node.Content
	for index := 0; index < len(contents); index = index + 2 {
		key := contents[index]
		value := contents[index+1]

		log.Debug("checking %v (%v)", key.Value, key.Tag)
		//skip the 'merge' tag, find a direct match first
		if key.Tag == "!!merge" && !prefs.DontFollowAlias {
			log.Debug("Merge anchor")
			err := traverseMergeAnchor(newMatches, candidate, value, wantedKey, prefs, splat)
			if err != nil {
				return err
			}
		} else if splat || keyMatches(key, wantedKey) {
			log.Debug("MATCHED")
			if prefs.IncludeMapKeys {
				log.Debug("including key")
				candidateNode := candidate.CreateChildInMap(key, key)
				candidateNode.IsMapKey = true
				newMatches.Set(fmt.Sprintf("keyOf-%v", candidateNode.GetKey()), candidateNode)
			}
			if !prefs.DontIncludeMapValues {
				log.Debug("including value")
				candidateNode := candidate.CreateChildInMap(key, value)
				newMatches.Set(candidateNode.GetKey(), candidateNode)
			}
		}
	}

	return nil
}

func traverseMergeAnchor(newMatches *orderedmap.OrderedMap, originalCandidate *CandidateNode, value *yaml.Node, wantedKey string, prefs traversePreferences, splat bool) error {
	switch value.Kind {
	case yaml.AliasNode:
		if value.Alias.Kind != yaml.MappingNode {
			return fmt.Errorf("can only use merge anchors with maps (!!map), but got %v", value.Alias.Tag)
		}
		candidateNode := originalCandidate.CreateReplacement(value.Alias)
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

func traverseArray(candidate *CandidateNode, operation *Operation, prefs traversePreferences) (*list.List, error) {
	log.Debug("operation Value %v", operation.Value)
	indices := []*yaml.Node{{Value: operation.StringValue}}
	return traverseArrayWithIndices(candidate, indices, prefs)
}
