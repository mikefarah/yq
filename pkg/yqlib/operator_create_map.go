package yqlib

import (
	"container/list"
)

func createMapOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- createMapOperation")

	//each matchingNodes entry should turn into a sequence of keys to create.
	//then collect object should do a cross function of the same index sequence for all matches.

	var document uint

	sequences := list.New()

	if context.MatchingNodes.Len() > 0 {

		for matchingNodeEl := context.MatchingNodes.Front(); matchingNodeEl != nil; matchingNodeEl = matchingNodeEl.Next() {
			matchingNode := matchingNodeEl.Value.(*CandidateNode)
			sequenceNode, err := sequenceFor(d, context, matchingNode, expressionNode)
			if err != nil {
				return Context{}, err
			}
			sequences.PushBack(sequenceNode)
		}
	} else {
		sequenceNode, err := sequenceFor(d, context, nil, expressionNode)
		if err != nil {
			return Context{}, err
		}
		sequences.PushBack(sequenceNode)
	}

	node := listToNodeSeq(sequences)
	node.Document = document

	return context.SingleChildContext(node), nil

}

func sequenceFor(d *dataTreeNavigator, context Context, matchingNode *CandidateNode, expressionNode *ExpressionNode) (*CandidateNode, error) {
	var document uint
	var matches = list.New()

	if matchingNode != nil {
		document = matchingNode.Document
		matches.PushBack(matchingNode)
	}

	mapPairs, err := crossFunction(d, context.ChildContext(matches), expressionNode,
		func(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
			node := CandidateNode{Kind: MappingNode, Tag: "!!map"}

			node.AddKeyValueChild(lhs, rhs)

			node.Document = document

			return &node, nil
		}, false)

	if err != nil {
		return nil, err
	}
	innerList := listToNodeSeq(mapPairs.MatchingNodes)
	innerList.Style = FlowStyle
	innerList.Document = document
	return innerList, nil
}

// NOTE: here the document index gets dropped so we
// no longer know where the node originates from.
func listToNodeSeq(list *list.List) *CandidateNode {
	node := CandidateNode{Kind: SequenceNode, Tag: "!!seq"}
	for entry := list.Front(); entry != nil; entry = entry.Next() {
		entryCandidate := entry.Value.(*CandidateNode)
		log.Debugf("Collecting %v into sequence", NodeToString(entryCandidate))
		node.Content = append(node.Content, entryCandidate)
	}
	return &node
}
