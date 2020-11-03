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

func CollectObjectOperator(d *dataTreeNavigator, matchMap *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("-- collectObjectOperation")
	return collect(d, list.New(), matchMap)

}

func collect(d *dataTreeNavigator, aggregate *list.List, remainingMatches *list.List) (*list.List, error) {
	if remainingMatches.Len() == 0 {
		return aggregate, nil
	}

	candidate := remainingMatches.Remove(remainingMatches.Front()).(*CandidateNode)
	splatted, err := Splat(d, nodeToMap(candidate))

	for splatEl := splatted.Front(); splatEl != nil; splatEl = splatEl.Next() {
		splatEl.Value.(*CandidateNode).Path = nil
	}

	if err != nil {
		return nil, err
	}

	if aggregate.Len() == 0 {
		return collect(d, splatted, remainingMatches)
	}

	newAgg := list.New()

	for el := aggregate.Front(); el != nil; el = el.Next() {
		aggCandidate := el.Value.(*CandidateNode)
		for splatEl := splatted.Front(); splatEl != nil; splatEl = splatEl.Next() {
			splatCandidate := splatEl.Value.(*CandidateNode)
			newCandidate := aggCandidate.Copy()
			newCandidate.Path = nil

			newCandidate, err := multiply(d, newCandidate, splatCandidate)
			if err != nil {
				return nil, err
			}
			newAgg.PushBack(newCandidate)
		}
	}
	return collect(d, newAgg, remainingMatches)

}
