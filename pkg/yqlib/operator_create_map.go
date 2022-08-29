package yqlib

import (
	"container/list"

	"gopkg.in/yaml.v3"
)

func createMapOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- createMapOperation")

	//each matchingNodes entry should turn into a sequence of keys to create.
	//then collect object should do a cross function of the same index sequence for all matches.

	var path []interface{}

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

	return context.SingleChildContext(&CandidateNode{Node: listToNodeSeq(sequences), Document: document, Path: path}), nil

}

func sequenceFor(d *dataTreeNavigator, context Context, matchingNode *CandidateNode, expressionNode *ExpressionNode) (*CandidateNode, error) {
	var path []interface{}
	var document uint
	var matches = list.New()

	if matchingNode != nil {
		path = matchingNode.Path
		document = matchingNode.Document
		matches.PushBack(matchingNode)
	}

	mapPairs, err := crossFunction(d, context.ChildContext(matches), expressionNode,
		func(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
			node := yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
			log.Debugf("LHS:", NodeToString(lhs))
			log.Debugf("RHS:", NodeToString(rhs))
			node.Content = []*yaml.Node{
				unwrapDoc(lhs.Node),
				unwrapDoc(rhs.Node),
			}

			return &CandidateNode{Node: &node, Document: document, Path: path}, nil
		}, false)

	if err != nil {
		return nil, err
	}
	innerList := listToNodeSeq(mapPairs.MatchingNodes)
	innerList.Style = yaml.FlowStyle
	return &CandidateNode{Node: innerList, Document: document, Path: path}, nil
}

// NOTE: here the document index gets dropped so we
// no longer know where the node originates from.
func listToNodeSeq(list *list.List) *yaml.Node {
	node := yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for entry := list.Front(); entry != nil; entry = entry.Next() {
		entryCandidate := entry.Value.(*CandidateNode)
		log.Debugf("Collecting %v into sequence", NodeToString(entryCandidate))
		node.Content = append(node.Content, entryCandidate.Node)
	}
	return &node
}
