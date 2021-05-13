package yqlib

import (
	"github.com/elliotchance/orderedmap"
	"container/list"
	yaml "gopkg.in/yaml.v3"
	"fmt"
)

func unique(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	selfExpression := &ExpressionNode{Operation: &Operation{OperationType: selfReferenceOpType}}
	uniqueByExpression := &ExpressionNode{Operation: &Operation{OperationType: uniqueByOpType}, Rhs: selfExpression}
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
			rhs, err := d.GetMatchingNodes(context.SingleChildContext(child), expressionNode.Rhs)

			if err != nil {
				return Context{}, err
			}

			first := rhs.MatchingNodes.Front()
			keyCandidate := first.Value.(*CandidateNode)
			keyValue := keyCandidate.Node.Value
			_, exists := newMatches.Get(keyValue)

			if !exists {
				newMatches.Set(keyValue, child.Node)
			}
		}
		resultNode := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for el := newMatches.Front(); el != nil; el = el.Next() {
			resultNode.Content = append(resultNode.Content, el.Value.(*yaml.Node))
		}

		results.PushBack(candidate.CreateChild(nil, resultNode))
	}

	return context.ChildContext(results), nil

}