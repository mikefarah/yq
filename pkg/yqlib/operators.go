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
	lhs, err := d.GetMatchingNodes(context, expressionNode.LHS)
	if err != nil {
		return Context{}, err
	}

	assignmentOp := &Operation{OperationType: assignOpType, Preferences: expressionNode.Operation.Preferences}

	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		clone, err := candidate.Copy()
		if err != nil {
			return Context{}, err
		}
		valueCopyExp := &ExpressionNode{Operation: &Operation{OperationType: valueOpType, CandidateNode: clone}}

		valueExpression := &ExpressionNode{Operation: &Operation{OperationType: valueOpType, CandidateNode: candidate}}

		assignmentOpNode := &ExpressionNode{Operation: assignmentOp, LHS: valueExpression, RHS: calculation(valueCopyExp, expressionNode.RHS)}

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

func resultsForRHS(d *dataTreeNavigator, context Context, lhsCandidate *CandidateNode, prefs crossFunctionPreferences, rhsExp *ExpressionNode, results *list.List) error {

	if prefs.LhsResultValue != nil {
		result, err := prefs.LhsResultValue(lhsCandidate)
		if err != nil {
			return err
		} else if result != nil {
			results.PushBack(result)
			return nil
		}
	}

	rhs, err := d.GetMatchingNodes(context, rhsExp)
	if err != nil {
		return err
	}

	if prefs.CalcWhenEmpty && rhs.MatchingNodes.Len() == 0 {
		resultCandidate, err := prefs.Calculation(d, context, lhsCandidate, nil)
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
		resultCandidate, err := prefs.Calculation(d, context, lhsCandidate, rhsCandidate)
		if err != nil {
			return err
		}
		if resultCandidate != nil {
			results.PushBack(resultCandidate)
		}
	}
	return nil
}

type crossFunctionPreferences struct {
	CalcWhenEmpty bool
	// if this returns a result node,
	// we wont bother calculating the RHS
	LhsResultValue func(*CandidateNode) (*CandidateNode, error)
	Calculation    crossFunctionCalculation
}

func doCrossFunc(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode, prefs crossFunctionPreferences) (Context, error) {
	var results = list.New()
	lhs, err := d.GetMatchingNodes(context, expressionNode.LHS)
	if err != nil {
		return Context{}, err
	}
	log.Debugf("crossFunction LHS len: %v", lhs.MatchingNodes.Len())

	if prefs.CalcWhenEmpty && lhs.MatchingNodes.Len() == 0 {
		err := resultsForRHS(d, context, nil, prefs, expressionNode.RHS, results)
		if err != nil {
			return Context{}, err
		}
	}

	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		lhsCandidate := el.Value.(*CandidateNode)

		err = resultsForRHS(d, context, lhsCandidate, prefs, expressionNode.RHS, results)
		if err != nil {
			return Context{}, err
		}

	}
	return context.ChildContext(results), nil
}

func crossFunction(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode, calculation crossFunctionCalculation, calcWhenEmpty bool) (Context, error) {
	prefs := crossFunctionPreferences{CalcWhenEmpty: calcWhenEmpty, Calculation: calculation}
	return crossFunctionWithPrefs(d, context, expressionNode, prefs)
}

func crossFunctionWithPrefs(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode, prefs crossFunctionPreferences) (Context, error) {
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
		return doCrossFunc(d, context, expressionNode, prefs)
	}

	log.Debug("crossFunction evaluate apart!")

	for matchEl := context.MatchingNodes.Front(); matchEl != nil; matchEl = matchEl.Next() {
		innerResults, err := doCrossFunc(d, context.SingleChildContext(matchEl.Value.(*CandidateNode)), expressionNode, prefs)
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
	return owner.CreateReplacement(node)
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
		LHS:       createTraversalTree(path[0:1], traversePrefs, false),
		RHS:       createTraversalTree(path[1:], traversePrefs, targetKey),
	}
}
