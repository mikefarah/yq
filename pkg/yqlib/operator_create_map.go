package yqlib

import (
	"container/list"
)

func createMapOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("createMapOperation")

	//each matchingNodes entry should turn into a sequence of keys to create.
	//then collect object should do a cross function of the same index sequence for all matches.

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

	return context.SingleChildContext(node), nil

}

func sequenceFor(d *dataTreeNavigator, context Context, matchingNode *CandidateNode, expressionNode *ExpressionNode) (*CandidateNode, error) {
	var document uint
	var filename string
	var fileIndex int

	var matches = list.New()

	if matchingNode != nil {
		document = matchingNode.GetDocument()
		filename = matchingNode.GetFilename()
		fileIndex = matchingNode.GetFileIndex()
		matches.PushBack(matchingNode)
	}

	log.Debugf("**********sequenceFor %v", NodeToString(matchingNode))

	mapPairs, err := crossFunction(d, context.ChildContext(matches), expressionNode,
		func(_ *dataTreeNavigator, _ Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
			node := &CandidateNode{Kind: MappingNode, Tag: "!!map"}

			log.Debugf("**********adding key %v and value %v", NodeToString(lhs), NodeToString(rhs))

			node.AddKeyValueChild(lhs, rhs)

			node.document = document
			node.fileIndex = fileIndex
			node.filename = filename

			return node, nil
		}, false)

	if err != nil {
		return nil, err
	}
	innerList := listToNodeSeq(mapPairs.MatchingNodes)
	innerList.Style = FlowStyle
	innerList.document = document
	innerList.fileIndex = fileIndex
	innerList.filename = filename
	return innerList, nil
}

// NOTE: here the document index gets dropped so we
// no longer know where the node originates from.
func listToNodeSeq(list *list.List) *CandidateNode {
	node := CandidateNode{Kind: SequenceNode, Tag: "!!seq"}
	for entry := list.Front(); entry != nil; entry = entry.Next() {
		entryCandidate := entry.Value.(*CandidateNode)
		log.Debugf("Collecting %v into sequence", NodeToString(entryCandidate))
		node.AddChild(entryCandidate)
	}
	return &node
}
