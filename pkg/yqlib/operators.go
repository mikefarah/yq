package yqlib

import (
	"container/list"
	"fmt"

	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v3"
)

type operatorHandler func(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error)

type compoundCalculation func(lhs *ExpressionNode, rhs *ExpressionNode) *ExpressionNode

func compoundAssignFunction(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode, calculation compoundCalculation) (Context, error) {
	lhs, err := d.GetMatchingNodes(context, expressionNode.Lhs)
	if err != nil {
		return Context{}, err
	}

	assignmentOp := &Operation{OperationType: assignOpType}
	valueOp := &Operation{OperationType: valueOpType}

	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		valueNodes := list.New()
		valueNodes.PushBack(candidate)
		valueOp.ValueNodes = valueNodes
		valueExpression := &ExpressionNode{Operation: valueOp}

		assignmentOpNode := &ExpressionNode{Operation: assignmentOp, Lhs: valueExpression, Rhs: calculation(valueExpression, expressionNode.Rhs)}

		_, err = d.GetMatchingNodes(context, assignmentOpNode)
		if err != nil {
			return Context{}, err
		}
	}
	return context, nil
}

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

func resultsForRhs(d *dataTreeNavigator, context Context, lhsCandidate *CandidateNode, rhs Context, calculation crossFunctionCalculation, results *list.List, calcWhenEmpty bool) error {

	if calcWhenEmpty && rhs.MatchingNodes.Len() == 0 {
		resultCandidate, err := calculation(d, context, lhsCandidate, nil)
		if err != nil {
			return err
		}
		if resultCandidate != nil {
			results.PushBack(resultCandidate)
		}
		return nil
	}

	for rightEl := rhs.MatchingNodes.Front(); rightEl != nil; rightEl = rightEl.Next() {
		log.Debugf("Applying calc")
		rhsCandidate := rightEl.Value.(*CandidateNode)
		resultCandidate, err := calculation(d, context, lhsCandidate, rhsCandidate)
		if err != nil {
			return err
		}
		if resultCandidate != nil {
			results.PushBack(resultCandidate)
		}
	}
	return nil
}

func doCrossFunc(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode, calculation crossFunctionCalculation, calcWhenEmpty bool) (Context, error) {
	var results = list.New()
	lhs, err := d.GetMatchingNodes(context, expressionNode.Lhs)
	if err != nil {
		return Context{}, err
	}
	log.Debugf("crossFunction LHS %v", NodesToString(lhs.MatchingNodes))

	rhs, err := d.GetMatchingNodes(context, expressionNode.Rhs)

	if err != nil {
		return Context{}, err
	}

	if calcWhenEmpty && lhs.MatchingNodes.Len() == 0 {
		err := resultsForRhs(d, context, nil, rhs, calculation, results, calcWhenEmpty)
		if err != nil {
			return Context{}, err
		}
	}

	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		lhsCandidate := el.Value.(*CandidateNode)

		err := resultsForRhs(d, context, lhsCandidate, rhs, calculation, results, calcWhenEmpty)
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
		log.Debug("crossFunction evaluateAllTogether!")
		return doCrossFunc(d, context, expressionNode, calculation, calcWhenEmpty)
	}

	log.Debug("crossFunction evaluate apart!")

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
