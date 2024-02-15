package yqlib

import (
	"container/list"
	"fmt"
)

func entrySeqFor(key *CandidateNode, value *CandidateNode) *CandidateNode {
	var keyKey = &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "key"}
	var valueKey = &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "value"}
	candidate := &CandidateNode{Kind: MappingNode, Tag: "!!map"}
	candidate.AddKeyValueChild(keyKey, key)
	candidate.AddKeyValueChild(valueKey, value)
	return candidate
}

func toEntriesFromMap(candidateNode *CandidateNode) *CandidateNode {
	var sequence = candidateNode.CreateReplacementWithComments(SequenceNode, "!!seq", 0)

	var contents = candidateNode.Content
	for index := 0; index < len(contents); index = index + 2 {
		key := contents[index]
		value := contents[index+1]

		sequence.AddChild(entrySeqFor(key, value))
	}
	return sequence
}

func toEntriesfromSeq(candidateNode *CandidateNode) *CandidateNode {
	var sequence = candidateNode.CreateReplacementWithComments(SequenceNode, "!!seq", 0)

	var contents = candidateNode.Content
	for index := 0; index < len(contents); index = index + 1 {
		key := &CandidateNode{Kind: ScalarNode, Tag: "!!int", Value: fmt.Sprintf("%v", index)}
		value := contents[index]

		sequence.AddChild(entrySeqFor(key, value))
	}
	return sequence
}

func toEntriesOperator(_ *dataTreeNavigator, context Context, _ *ExpressionNode) (Context, error) {
	var results = list.New()
	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		switch candidate.Kind {
		case MappingNode:
			results.PushBack(toEntriesFromMap(candidate))

		case SequenceNode:
			results.PushBack(toEntriesfromSeq(candidate))
		default:
			if candidate.Tag != "!!null" {
				return Context{}, fmt.Errorf("%v has no keys", candidate.Tag)
			}
		}
	}

	return context.ChildContext(results), nil
}

func parseEntry(candidateNode *CandidateNode, position int) (*CandidateNode, *CandidateNode, error) {
	prefs := traversePreferences{DontAutoCreate: true}

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

	return keyResults.Front().Value.(*CandidateNode), valueResults.Front().Value.(*CandidateNode), nil

}

func fromEntries(candidateNode *CandidateNode) (*CandidateNode, error) {
	var node = candidateNode.CopyWithoutContent()

	var contents = candidateNode.Content

	for index := 0; index < len(contents); index = index + 1 {
		key, value, err := parseEntry(contents[index], index)
		if err != nil {
			return nil, err
		}

		node.AddKeyValueChild(key, value)
	}
	node.Kind = MappingNode
	node.Tag = "!!map"
	return node, nil
}

func fromEntriesOperator(_ *dataTreeNavigator, context Context, _ *ExpressionNode) (Context, error) {
	var results = list.New()
	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		switch candidate.Kind {
		case SequenceNode:
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

		newResults := list.New()

		for itemEl := splatted.MatchingNodes.Front(); itemEl != nil; itemEl = itemEl.Next() {
			result, err := d.GetMatchingNodes(splatted.SingleChildContext(itemEl.Value.(*CandidateNode)), expressionNode.RHS)
			if err != nil {
				return Context{}, err
			}
			newResults.PushBackList(result.MatchingNodes)
		}

		selfExpression := &ExpressionNode{Operation: &Operation{OperationType: selfReferenceOpType}}
		collected, err := collectTogether(d, splatted.ChildContext(newResults), selfExpression)
		if err != nil {
			return Context{}, err
		}
		log.Debug("candidate %v", NodeToString(candidate))
		log.Debug("candidate leading content: %v", candidate.LeadingContent)
		collected.LeadingContent = candidate.LeadingContent
		log.Debug("candidate FootComment: [%v]", candidate.FootComment)

		collected.HeadComment = candidate.HeadComment
		collected.FootComment = candidate.FootComment

		log.Debugf("collected %v", collected.LeadingContent)

		fromEntries, err := fromEntriesOperator(d, context.SingleChildContext(collected), expressionNode)
		if err != nil {
			return Context{}, err
		}
		results.PushBackList(fromEntries.MatchingNodes)

	}

	//from_entries on the result
	return context.ChildContext(results), nil
}
