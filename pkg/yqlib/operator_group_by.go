package yqlib

import (
	"container/list"
	"fmt"

	"github.com/elliotchance/orderedmap"
)

func processIntoGroups(d *dataTreeNavigator, context Context, rhsExp *ExpressionNode, node *CandidateNode) (*orderedmap.OrderedMap, error) {
	var newMatches = orderedmap.NewOrderedMap()
	for _, child := range node.Content {
		rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(child), rhsExp)

		if err != nil {
			return nil, err
		}

		keyValue := "null"

		if rhs.MatchingNodes.Len() > 0 {
			first := rhs.MatchingNodes.Front()
			keyCandidate := first.Value.(*CandidateNode)
			keyValue = keyCandidate.Value
		}

		groupList, exists := newMatches.Get(keyValue)

		if !exists {
			groupList = list.New()
			newMatches.Set(keyValue, groupList)
		}
		groupList.(*list.List).PushBack(child)
	}
	return newMatches, nil
}

func groupBy(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("groupBy Operator")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		if candidate.Kind != SequenceNode {
			return Context{}, fmt.Errorf("only arrays are supported for group by")
		}

		newMatches, err := processIntoGroups(d, context, expressionNode.RHS, candidate)

		if err != nil {
			return Context{}, err
		}

		resultNode := candidate.CreateReplacement(SequenceNode, "!!seq", "")
		for groupEl := newMatches.Front(); groupEl != nil; groupEl = groupEl.Next() {
			groupResultNode := &CandidateNode{Kind: SequenceNode, Tag: "!!seq"}
			groupList := groupEl.Value.(*list.List)
			for groupItem := groupList.Front(); groupItem != nil; groupItem = groupItem.Next() {
				groupResultNode.AddChild(groupItem.Value.(*CandidateNode))
			}

			resultNode.AddChild(groupResultNode)
		}

		results.PushBack(resultNode)

	}

	return context.ChildContext(results), nil

}
