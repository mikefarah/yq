package yqlib

import (
	"container/list"
	"fmt"

	"gopkg.in/yaml.v3"
)

type operatorHandler func(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error)

func unwrapDoc(node *yaml.Node) *yaml.Node {
	if node.Kind == yaml.DocumentNode {
		return node.Content[0]
	}
	return node
}

func emptyOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	context.MatchingNodes = list.New()
	return context, nil
}

type crossFunctionCalculation func(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error)

func doCrossFunc(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode, calculation crossFunctionCalculation) (Context, error) {
	var results = list.New()
	lhs, err := d.GetMatchingNodes(context, expressionNode.Lhs)
	if err != nil {
		return Context{}, err
	}
	log.Debugf("crossFunction LHS len: %v", lhs.MatchingNodes.Len())

	rhs, err := d.GetMatchingNodes(context, expressionNode.Rhs)

	if err != nil {
		return Context{}, err
	}

	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		lhsCandidate := el.Value.(*CandidateNode)

		for rightEl := rhs.MatchingNodes.Front(); rightEl != nil; rightEl = rightEl.Next() {
			log.Debugf("Applying calc")
			rhsCandidate := rightEl.Value.(*CandidateNode)
			resultCandidate, err := calculation(d, context, lhsCandidate, rhsCandidate)
			if err != nil {
				return Context{}, err
			}
			results.PushBack(resultCandidate)
		}

	}
	return context.ChildContext(results), nil
}

func crossFunction(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode, calculation crossFunctionCalculation) (Context, error) {
	var results = list.New()

	var evaluateAllTogether = true
	for matchEl := context.MatchingNodes.Front(); matchEl != nil; matchEl = matchEl.Next() {
		evaluateAllTogether = evaluateAllTogether && matchEl.Value.(*CandidateNode).EvaluateTogether
		if !evaluateAllTogether {
			break
		}
	}
	if evaluateAllTogether {
		return doCrossFunc(d, context, expressionNode, calculation)
	}

	for matchEl := context.MatchingNodes.Front(); matchEl != nil; matchEl = matchEl.Next() {
		innerResults, err := doCrossFunc(d, context.SingleChildContext(matchEl.Value.(*CandidateNode)), expressionNode, calculation)
		if err != nil {
			return Context{}, err
		}
		results.PushBackList(innerResults.MatchingNodes)
	}

	return context.ChildContext(results), nil
}

func createBooleanCandidate(owner *CandidateNode, value bool) *CandidateNode {
	valString := "true"
	if !value {
		valString = "false"
	}
	node := &yaml.Node{Kind: yaml.ScalarNode, Value: valString, Tag: "!!bool"}
	return owner.CreateChild(nil, node)
}

func createTraversalTree(path []interface{}, traversePrefs traversePreferences) *ExpressionNode {
	if len(path) == 0 {
		return &ExpressionNode{Operation: &Operation{OperationType: selfReferenceOpType}}
	} else if len(path) == 1 {
		return &ExpressionNode{Operation: &Operation{OperationType: traversePathOpType, Preferences: traversePrefs, Value: path[0], StringValue: fmt.Sprintf("%v", path[0])}}
	}

	return &ExpressionNode{
		Operation: &Operation{OperationType: shortPipeOpType},
		Lhs:       createTraversalTree(path[0:1], traversePrefs),
		Rhs:       createTraversalTree(path[1:], traversePrefs),
	}
}
