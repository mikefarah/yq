package yqlib

import (
	"container/list"
)

func inOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("inOperation")
	var results = list.New()

	rhs, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.RHS)

	if err != nil {
		return Context{}, err
	}

	var collection *CandidateNode
	if rhs.MatchingNodes.Len() != 0 {
		collection = rhs.MatchingNodes.Front().Value.(*CandidateNode)
	} else {
		// no collection provided, return false for all
		for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
			candidate := el.Value.(*CandidateNode)
			results.PushBack(createBooleanCandidate(candidate, false))
		}
		return context.ChildContext(results), nil
	}

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		isIn := false

		switch collection.Kind {
		case MappingNode:
			// Check if candidate value exists as a key in the mapping
			for index := 0; index < len(collection.Content); index = index + 2 {
				key := collection.Content[index]
				if key.Value == candidate.Value {
					isIn = true
					break
				}
			}
		case SequenceNode:
			// Check if candidate value is present in the sequence (value membership)
			for _, element := range collection.Content {
				if element.Value == candidate.Value {
					isIn = true
					break
				}
			}
		default:
			isIn = false
		}

		results.PushBack(createBooleanCandidate(candidate, isIn))
	}
	return context.ChildContext(results), nil
}
