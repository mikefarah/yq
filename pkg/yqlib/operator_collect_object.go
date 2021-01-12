package yqlib

import (
	"container/list"

	yaml "gopkg.in/yaml.v3"
)

/*
[Mike: cat, Bob: dog]
[Thing: rabbit, peter: sam]

==> cross multiply

{Mike: cat, Thing: rabbit}
{Mike: cat, peter: sam}
...
*/

func collectObjectOperator(d *dataTreeNavigator, matchMap *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("-- collectObjectOperation")

	if matchMap.Len() == 0 {
		node := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map", Value: "{}"}
		candidate := &CandidateNode{Node: node}
		return nodeToMap(candidate), nil
	}
	first := matchMap.Front().Value.(*CandidateNode)
	var rotated []*list.List = make([]*list.List, len(first.Node.Content))

	for i := 0; i < len(first.Node.Content); i++ {
		rotated[i] = list.New()
	}

	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidateNode := el.Value.(*CandidateNode)
		for i := 0; i < len(first.Node.Content); i++ {
			rotated[i].PushBack(candidateNode.CreateChild(i, candidateNode.Node.Content[i]))
		}
	}

	newObject := list.New()
	for i := 0; i < len(first.Node.Content); i++ {
		additions, err := collect(d, list.New(), rotated[i])
		if err != nil {
			return nil, err
		}
		newObject.PushBackList(additions)
	}

	return newObject, nil

}

func collect(d *dataTreeNavigator, aggregate *list.List, remainingMatches *list.List) (*list.List, error) {
	if remainingMatches.Len() == 0 {
		return aggregate, nil
	}

	candidate := remainingMatches.Remove(remainingMatches.Front()).(*CandidateNode)

	splatted, err := splat(d, nodeToMap(candidate),
		&traversePreferences{FollowAlias: false, IncludeMapKeys: false})

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
			newCandidate, err := aggCandidate.Copy()
			if err != nil {
				return nil, err
			}

			newCandidate.Path = nil

			newCandidate, err = multiply(&multiplyPreferences{AppendArrays: false})(d, newCandidate, splatCandidate)
			if err != nil {
				return nil, err
			}
			newAgg.PushBack(newCandidate)
		}
	}
	return collect(d, newAgg, remainingMatches)

}
