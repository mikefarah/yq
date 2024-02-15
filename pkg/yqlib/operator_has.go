package yqlib

import (
	"container/list"
	"strconv"
)

func hasOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("hasOperation")
	var results = list.New()

	rhs, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.RHS)

	if err != nil {
		return Context{}, err
	}

	wantedKey := "null"
	wanted := &CandidateNode{Tag: "!!null"}
	if rhs.MatchingNodes.Len() != 0 {
		wanted = rhs.MatchingNodes.Front().Value.(*CandidateNode)
		wantedKey = wanted.Value
	}

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		var contents = candidate.Content
		switch candidate.Kind {
		case MappingNode:
			candidateHasKey := false
			for index := 0; index < len(contents) && !candidateHasKey; index = index + 2 {
				key := contents[index]
				if key.Value == wantedKey {
					candidateHasKey = true
				}
			}
			results.PushBack(createBooleanCandidate(candidate, candidateHasKey))
		case SequenceNode:
			candidateHasKey := false
			if wanted.Tag == "!!int" {
				var number, errParsingInt = strconv.ParseInt(wantedKey, 10, 64)
				if errParsingInt != nil {
					return Context{}, errParsingInt
				}
				candidateHasKey = int64(len(contents)) > number
			}
			results.PushBack(createBooleanCandidate(candidate, candidateHasKey))
		default:
			results.PushBack(createBooleanCandidate(candidate, false))
		}
	}
	return context.ChildContext(results), nil
}
