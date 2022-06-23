package yqlib

import (
	"container/list"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

func entrySeqFor(key *yaml.Node, value *yaml.Node) *yaml.Node {
	var keyKey = &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "key"}
	var valueKey = &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "value"}

	return &yaml.Node{
		Kind:    yaml.MappingNode,
		Tag:     "!!map",
		Content: []*yaml.Node{keyKey, key, valueKey, value},
	}
}

func toEntriesFromMap(candidateNode *CandidateNode) *CandidateNode {
	var sequence = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	var entriesNode = candidateNode.CreateReplacementWithDocWrappers(sequence)

	var contents = unwrapDoc(candidateNode.Node).Content
	for index := 0; index < len(contents); index = index + 2 {
		key := contents[index]
		value := contents[index+1]

		sequence.Content = append(sequence.Content, entrySeqFor(key, value))
	}
	return entriesNode
}

func toEntriesfromSeq(candidateNode *CandidateNode) *CandidateNode {
	var sequence = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	var entriesNode = candidateNode.CreateReplacementWithDocWrappers(sequence)

	var contents = unwrapDoc(candidateNode.Node).Content
	for index := 0; index < len(contents); index = index + 1 {
		key := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: fmt.Sprintf("%v", index)}
		value := contents[index]

		sequence.Content = append(sequence.Content, entrySeqFor(key, value))
	}
	return entriesNode
}

func toEntriesOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	var results = list.New()
	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		candidateNode := unwrapDoc(candidate.Node)

		switch candidateNode.Kind {
		case yaml.MappingNode:
			results.PushBack(toEntriesFromMap(candidate))

		case yaml.SequenceNode:
			results.PushBack(toEntriesfromSeq(candidate))
		default:
			if candidateNode.Tag != "!!null" {
				return Context{}, fmt.Errorf("%v has no keys", candidate.Node.Tag)
			}
		}
	}

	return context.ChildContext(results), nil
}

func parseEntry(entry *yaml.Node, position int) (*yaml.Node, *yaml.Node, error) {
	prefs := traversePreferences{DontAutoCreate: true}
	candidateNode := &CandidateNode{Node: entry}

	keyResults, err := traverseMap(Context{}, candidateNode, createStringScalarNode("key"), prefs, false)

	if err != nil {
		return nil, nil, err
	} else if keyResults.Len() != 1 {
		return nil, nil, fmt.Errorf("expected to find one 'key' entry but found %v in position %v", keyResults.Len(), position)
	}

	valueResults, err := traverseMap(Context{}, candidateNode, createStringScalarNode("value"), prefs, false)

	if err != nil {
		return nil, nil, err
	} else if valueResults.Len() != 1 {
		return nil, nil, fmt.Errorf("expected to find one 'value' entry but found %v in position %v", valueResults.Len(), position)
	}

	return keyResults.Front().Value.(*CandidateNode).Node, valueResults.Front().Value.(*CandidateNode).Node, nil

}

func fromEntries(candidateNode *CandidateNode) (*CandidateNode, error) {
	var node = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	var mapCandidateNode = candidateNode.CreateReplacementWithDocWrappers(node)

	var contents = unwrapDoc(candidateNode.Node).Content

	for index := 0; index < len(contents); index = index + 1 {
		key, value, err := parseEntry(contents[index], index)
		if err != nil {
			return nil, err
		}

		node.Content = append(node.Content, key, value)
	}
	return mapCandidateNode, nil
}

func fromEntriesOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	var results = list.New()
	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		candidateNode := unwrapDoc(candidate.Node)

		switch candidateNode.Kind {
		case yaml.SequenceNode:
			mapResult, err := fromEntries(candidate)
			if err != nil {
				return Context{}, err
			}
			results.PushBack(mapResult)
		default:
			return Context{}, fmt.Errorf("from entries only runs against arrays")
		}
	}

	return context.ChildContext(results), nil
}

func withEntriesOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	//to_entries on the context
	toEntries, err := toEntriesOperator(d, context, expressionNode)
	if err != nil {
		return Context{}, err
	}

	var results = list.New()

	for el := toEntries.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		//run expression against entries
		// splat toEntries and pipe it into Rhs
		splatted, err := splat(context.SingleChildContext(candidate), traversePreferences{})
		if err != nil {
			return Context{}, err
		}

		result, err := d.GetMatchingNodes(splatted, expressionNode.RHS)
		log.Debug("expressionNode.Rhs %v", expressionNode.RHS.Operation.OperationType)
		log.Debug("result %v", result)
		if err != nil {
			return Context{}, err
		}

		selfExpression := &ExpressionNode{Operation: &Operation{OperationType: selfReferenceOpType}}
		collected, err := collectTogether(d, result, selfExpression)
		if err != nil {
			return Context{}, err
		}
		collected.LeadingContent = candidate.LeadingContent
		collected.TrailingContent = candidate.TrailingContent

		log.Debugf("**** collected %v", collected.LeadingContent)

		fromEntries, err := fromEntriesOperator(d, context.SingleChildContext(collected), expressionNode)
		if err != nil {
			return Context{}, err
		}
		results.PushBackList(fromEntries.MatchingNodes)

	}

	//from_entries on the result
	return context.ChildContext(results), nil
}
