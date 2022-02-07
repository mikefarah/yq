package yqlib

import (
	"container/list"
	"fmt"

	"github.com/elliotchance/orderedmap"
	yaml "gopkg.in/yaml.v3"
)

func processIntoGroups(d *dataTreeNavigator, context Context, rhsExp *ExpressionNode, node *yaml.Node) (*orderedmap.OrderedMap, error) {
	var newMatches = orderedmap.NewOrderedMap()
	for _, node := range node.Content {
		child := &CandidateNode{Node: node}
		rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(child), rhsExp)

		if err != nil {
			return nil, err
		}

		keyValue := "null"

		if rhs.MatchingNodes.Len() > 0 {
			first := rhs.MatchingNodes.Front()
			keyCandidate := first.Value.(*CandidateNode)
			keyValue = keyCandidate.Node.Value
		}

		groupList, exists := newMatches.Get(keyValue)

		if !exists {
			groupList = list.New()
			newMatches.Set(keyValue, groupList)
		}
		groupList.(*list.List).PushBack(node)
	}
	return newMatches, nil
}

func groupBy(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("-- groupBy Operator")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		candidateNode := unwrapDoc(candidate.Node)

		if candidateNode.Kind != yaml.SequenceNode {
			return Context{}, fmt.Errorf("Only arrays are supported for group by")
		}

		newMatches, err := processIntoGroups(d, context, expressionNode.RHS, candidateNode)

		if err != nil {
			return Context{}, err
		}

		resultNode := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for groupEl := newMatches.Front(); groupEl != nil; groupEl = groupEl.Next() {
			groupResultNode := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
			groupList := groupEl.Value.(*list.List)
			for groupItem := groupList.Front(); groupItem != nil; groupItem = groupItem.Next() {
				groupResultNode.Content = append(groupResultNode.Content, groupItem.Value.(*yaml.Node))
			}

			resultNode.Content = append(resultNode.Content, groupResultNode)
		}

		results.PushBack(candidate.CreateReplacement(resultNode))

	}

	return context.ChildContext(results), nil

}
