package yqlib

import (
	"container/list"
)

/*
[Mike: cat, Bob: dog]
[Thing: rabbit, peter: sam]

==> cross multiply

{Mike: cat, Thing: rabbit}
{Mike: cat, peter: sam}
...
*/

func collectObjectOperator(d *dataTreeNavigator, originalContext Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- collectObjectOperation")

	context := originalContext.WritableClone()

	if context.MatchingNodes.Len() == 0 {
		candidate := &CandidateNode{Kind: MappingNode, Tag: "!!map", Value: "{}"}
		return context.SingleChildContext(candidate), nil
	}
	first := context.MatchingNodes.Front().Value.(*CandidateNode)
	var rotated = make([]*list.List, len(first.Content))

	for i := 0; i < len(first.Content); i++ {
		rotated[i] = list.New()
	}

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidateNode := el.Value.(*CandidateNode)
		for i := 0; i < len(first.Content); i++ {

			rotated[i].PushBack(candidateNode.Content[i])
		}
	}

	newObject := list.New()
	for i := 0; i < len(first.Content); i++ {
		additions, err := collect(d, context.ChildContext(list.New()), rotated[i])
		if err != nil {
			return Context{}, err
		}
		newObject.PushBackList(additions.MatchingNodes)
	}

	return context.ChildContext(newObject), nil

}

func collect(d *dataTreeNavigator, context Context, remainingMatches *list.List) (Context, error) {
	if remainingMatches.Len() == 0 {
		return context, nil
	}

	candidate := remainingMatches.Remove(remainingMatches.Front()).(*CandidateNode)

	splatted, err := splat(context.SingleChildContext(candidate),
		traversePreferences{DontFollowAlias: true, IncludeMapKeys: false})

	if err != nil {
		return Context{}, err
	}

	if context.MatchingNodes.Len() == 0 {
		return collect(d, splatted, remainingMatches)
	}

	newAgg := list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		aggCandidate := el.Value.(*CandidateNode)
		for splatEl := splatted.MatchingNodes.Front(); splatEl != nil; splatEl = splatEl.Next() {
			splatCandidate := splatEl.Value.(*CandidateNode)
			newCandidate := aggCandidate.Copy()

			newCandidate, err = multiply(multiplyPreferences{AppendArrays: false})(d, context, newCandidate, splatCandidate)
			if err != nil {
				return Context{}, err
			}
			newAgg.PushBack(newCandidate)
		}
	}
	return collect(d, context.ChildContext(newAgg), remainingMatches)

}
