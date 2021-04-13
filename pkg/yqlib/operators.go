package yqlib

import (
	"container/list"
	"fmt"

	"github.com/jinzhu/copier"
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

func resultsForRhs(d *dataTreeNavigator, context Context, lhsCandidate *CandidateNode, rhs Context, calculation crossFunctionCalculation, results *list.List) error {
	for rightEl := rhs.MatchingNodes.Front(); rightEl != nil; rightEl = rightEl.Next() {
		log.Debugf("Applying calc")
		rhsCandidate := rightEl.Value.(*CandidateNode)
		resultCandidate, err := calculation(d, context, lhsCandidate, rhsCandidate)
		if err != nil {
			return err
		}
		results.PushBack(resultCandidate)
	}
	return nil
}

func doCrossFunc(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode, calculation crossFunctionCalculation, calcWhenEmpty bool) (Context, error) {
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

	if calcWhenEmpty && lhs.MatchingNodes.Len() == 0 {
		if rhs.MatchingNodes.Len() == 0 {
			resultCandidate, err := calculation(d, context, nil, nil)
			if err != nil {
				return Context{}, err
			}
			results.PushBack(resultCandidate)
		}
		err := resultsForRhs(d, context, nil, rhs, calculation, results)
		if err != nil {
			return Context{}, err
		}
	}

	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		lhsCandidate := el.Value.(*CandidateNode)

		err := resultsForRhs(d, context, lhsCandidate, rhs, calculation, results)
		if err != nil {
			return Context{}, err
		}

	}
	return context.ChildContext(results), nil
}

func crossFunction(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode, calculation crossFunctionCalculation, calcWhenEmpty bool) (Context, error) {
	var results = list.New()

	var evaluateAllTogether = true
	for matchEl := context.MatchingNodes.Front(); matchEl != nil; matchEl = matchEl.Next() {
		evaluateAllTogether = evaluateAllTogether && matchEl.Value.(*CandidateNode).EvaluateTogether
		if !evaluateAllTogether {
			break
		}
	}
	if evaluateAllTogether {
		return doCrossFunc(d, context, expressionNode, calculation, calcWhenEmpty)
	}

	for matchEl := context.MatchingNodes.Front(); matchEl != nil; matchEl = matchEl.Next() {
		innerResults, err := doCrossFunc(d, context.SingleChildContext(matchEl.Value.(*CandidateNode)), expressionNode, calculation, calcWhenEmpty)
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

func createTraversalTree(path []interface{}, traversePrefs traversePreferences, targetKey bool) *ExpressionNode {
	if len(path) == 0 {
		return &ExpressionNode{Operation: &Operation{OperationType: selfReferenceOpType}}
	} else if len(path) == 1 {
		lastPrefs := traversePrefs
		if targetKey {
			err := copier.Copy(&lastPrefs, traversePrefs)
			if err != nil {
				panic(err)
			}
			lastPrefs.IncludeMapKeys = true
			lastPrefs.DontIncludeMapValues = true
		}
		return &ExpressionNode{Operation: &Operation{OperationType: traversePathOpType, Preferences: lastPrefs, Value: path[0], StringValue: fmt.Sprintf("%v", path[0])}}
	}

	return &ExpressionNode{
		Operation: &Operation{OperationType: shortPipeOpType},
		Lhs:       createTraversalTree(path[0:1], traversePrefs, false),
		Rhs:       createTraversalTree(path[1:], traversePrefs, targetKey),
	}
}
