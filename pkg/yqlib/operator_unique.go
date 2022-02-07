package yqlib

import (
	"container/list"
	"fmt"

	"github.com/elliotchance/orderedmap"
	yaml "gopkg.in/yaml.v3"
)

func unique(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	selfExpression := &ExpressionNode{Operation: &Operation{OperationType: selfReferenceOpType}}
	uniqueByExpression := &ExpressionNode{Operation: &Operation{OperationType: uniqueByOpType}, RHS: selfExpression}
	return uniqueBy(d, context, uniqueByExpression)

}

func uniqueBy(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("-- uniqueBy Operator")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		candidateNode := unwrapDoc(candidate.Node)

		if candidateNode.Kind != yaml.SequenceNode {
			return Context{}, fmt.Errorf("Only arrays are supported for unique")
		}

		var newMatches = orderedmap.NewOrderedMap()
		for _, node := range candidateNode.Content {
			child := &CandidateNode{Node: node}
			rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(child), expressionNode.RHS)

			if err != nil {
				return Context{}, err
			}

			keyValue := "null"

			if rhs.MatchingNodes.Len() > 0 {
				first := rhs.MatchingNodes.Front()
				keyCandidate := first.Value.(*CandidateNode)
				keyValue = keyCandidate.Node.Value
			}

			_, exists := newMatches.Get(keyValue)

			if !exists {
				newMatches.Set(keyValue, child.Node)
			}
		}
		resultNode := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for el := newMatches.Front(); el != nil; el = el.Next() {
			resultNode.Content = append(resultNode.Content, el.Value.(*yaml.Node))
		}

		results.PushBack(candidate.CreateReplacement(resultNode))
	}

	return context.ChildContext(results), nil

}
