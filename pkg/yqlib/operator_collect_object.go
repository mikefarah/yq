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

func collectObjectOperator(d *dataTreeNavigator, originalContext Context, _ *ExpressionNode) (Context, error) {
	log.Debugf("collectObjectOperation")

	context := originalContext.WritableClone()

	if context.MatchingNodes.Len() == 0 {
		candidate := &CandidateNode{Kind: MappingNode, Tag: "!!map", Value: "{}"}
		log.Debugf("collectObjectOperation - starting with empty map")
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
			log.Debugf("rotate[%v] = %v", i, NodeToString(candidateNode.Content[i]))
			log.Debugf("children:\n%v", NodeContentToString(candidateNode.Content[i], 0))
			rotated[i].PushBack(candidateNode.Content[i])
		}
	}
	log.Debugf("collectObjectOperation, length of rotated is %v", len(rotated))

	newObject := list.New()
	for i := 0; i < len(first.Content); i++ {
		additions, err := collect(d, context.ChildContext(list.New()), rotated[i])
		if err != nil {
			return Context{}, err
		}
		// we should reset the parents and keys of these top level nodes,
		// as they are new
		for el := additions.MatchingNodes.Front(); el != nil; el = el.Next() {
			addition := el.Value.(*CandidateNode)
			additionCopy := addition.Copy()

			additionCopy.SetParent(nil)
			additionCopy.Key = nil

			log.Debugf("collectObjectOperation, adding result %v", NodeToString(additionCopy))

			newObject.PushBack(additionCopy)
		}
	}

	return context.ChildContext(newObject), nil

}

func collect(d *dataTreeNavigator, context Context, remainingMatches *list.List) (Context, error) {
	if remainingMatches.Len() == 0 {
		return context, nil
	}

	candidate := remainingMatches.Remove(remainingMatches.Front()).(*CandidateNode)
	log.Debugf("collectObjectOperation - collect %v", NodeToString(candidate))

	splatted, err := splat(context.SingleChildContext(candidate),
		traversePreferences{DontFollowAlias: true, IncludeMapKeys: false})

	if err != nil {
		return Context{}, err
	}

	if context.MatchingNodes.Len() == 0 {
		log.Debugf("collectObjectOperation - collect context is empty, next")
		return collect(d, splatted, remainingMatches)
	}

	newAgg := list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		aggCandidate := el.Value.(*CandidateNode)
		for splatEl := splatted.MatchingNodes.Front(); splatEl != nil; splatEl = splatEl.Next() {
			splatCandidate := splatEl.Value.(*CandidateNode)
			log.Debugf("collectObjectOperation; splatCandidate: %v", NodeToString(splatCandidate))
			newCandidate := aggCandidate.Copy()
			log.Debugf("collectObjectOperation; aggCandidate: %v", NodeToString(aggCandidate))

			newCandidate, err = multiply(multiplyPreferences{AppendArrays: false})(d, context, newCandidate, splatCandidate)

			if err != nil {
				return Context{}, err
			}
			newAgg.PushBack(newCandidate)
		}
	}
	return collect(d, context.ChildContext(newAgg), remainingMatches)

}
