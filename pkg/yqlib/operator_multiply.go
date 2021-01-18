package yqlib

import (
	"fmt"
	"strconv"

	"container/list"

	yaml "gopkg.in/yaml.v3"
)

type crossFunctionCalculation func(d *dataTreeNavigator, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error)

func doCrossFunc(d *dataTreeNavigator, contextList *list.List, expressionNode *ExpressionNode, calculation crossFunctionCalculation) (*list.List, error) {
	var results = list.New()
	lhs, err := d.GetMatchingNodes(contextList, expressionNode.Lhs)
	if err != nil {
		return nil, err
	}
	log.Debugf("crossFunction LHS len: %v", lhs.Len())

	rhs, err := d.GetMatchingNodes(contextList, expressionNode.Rhs)

	if err != nil {
		return nil, err
	}

	for el := lhs.Front(); el != nil; el = el.Next() {
		lhsCandidate := el.Value.(*CandidateNode)

		for rightEl := rhs.Front(); rightEl != nil; rightEl = rightEl.Next() {
			log.Debugf("Applying calc")
			rhsCandidate := rightEl.Value.(*CandidateNode)
			resultCandidate, err := calculation(d, lhsCandidate, rhsCandidate)
			if err != nil {
				return nil, err
			}
			results.PushBack(resultCandidate)
		}

	}
	return results, nil
}

func crossFunction(d *dataTreeNavigator, matchingNodes *list.List, expressionNode *ExpressionNode, calculation crossFunctionCalculation) (*list.List, error) {
	var results = list.New()

	var evaluateAllTogether = true
	for matchEl := matchingNodes.Front(); matchEl != nil; matchEl = matchEl.Next() {
		evaluateAllTogether = evaluateAllTogether && matchEl.Value.(*CandidateNode).EvaluateTogether
		if !evaluateAllTogether {
			break
		}
	}
	if evaluateAllTogether {
		return doCrossFunc(d, matchingNodes, expressionNode, calculation)
	}

	for matchEl := matchingNodes.Front(); matchEl != nil; matchEl = matchEl.Next() {
		contextList := nodeToMap(matchEl.Value.(*CandidateNode))
		innerResults, err := doCrossFunc(d, contextList, expressionNode, calculation)
		if err != nil {
			return nil, err
		}
		results.PushBackList(innerResults)
	}

	return results, nil
}

type multiplyPreferences struct {
	AppendArrays  bool
	TraversePrefs traversePreferences
}

func multiplyOperator(d *dataTreeNavigator, matchingNodes *list.List, expressionNode *ExpressionNode) (*list.List, error) {
	log.Debugf("-- MultiplyOperator")
	return crossFunction(d, matchingNodes, expressionNode, multiply(expressionNode.Operation.Preferences.(multiplyPreferences)))
}

func multiply(preferences multiplyPreferences) func(d *dataTreeNavigator, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	return func(d *dataTreeNavigator, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
		lhs.Node = unwrapDoc(lhs.Node)
		rhs.Node = unwrapDoc(rhs.Node)
		log.Debugf("Multipling LHS: %v", lhs.Node.Tag)
		log.Debugf("-          RHS: %v", rhs.Node.Tag)

		if lhs.Node.Kind == yaml.MappingNode && rhs.Node.Kind == yaml.MappingNode ||
			(lhs.Node.Kind == yaml.SequenceNode && rhs.Node.Kind == yaml.SequenceNode) {

			var newBlank = lhs.CreateChild(nil, &yaml.Node{})
			var newThing, err = mergeObjects(d, newBlank, lhs, multiplyPreferences{})
			if err != nil {
				return nil, err
			}
			return mergeObjects(d, newThing, rhs, preferences)
		} else if lhs.Node.Tag == "!!int" && rhs.Node.Tag == "!!int" {
			return multiplyIntegers(lhs, rhs)
		}
		return nil, fmt.Errorf("Cannot multiply %v with %v", lhs.Node.Tag, rhs.Node.Tag)
	}
}

func multiplyIntegers(lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	target := lhs.CreateChild(nil, &yaml.Node{})
	target.Node.Kind = yaml.ScalarNode
	target.Node.Style = lhs.Node.Style
	target.Node.Tag = "!!int"

	lhsNum, err := strconv.Atoi(lhs.Node.Value)
	if err != nil {
		return nil, err
	}
	rhsNum, err := strconv.Atoi(rhs.Node.Value)
	if err != nil {
		return nil, err
	}
	target.Node.Value = fmt.Sprintf("%v", lhsNum*rhsNum)
	return target, nil
}

func mergeObjects(d *dataTreeNavigator, lhs *CandidateNode, rhs *CandidateNode, preferences multiplyPreferences) (*CandidateNode, error) {
	shouldAppendArrays := preferences.AppendArrays
	var results = list.New()

	// shouldn't recurse arrays if appending
	prefs := recursiveDescentPreferences{RecurseArray: !shouldAppendArrays,
		TraversePreferences: traversePreferences{DontFollowAlias: true}}
	err := recursiveDecent(d, results, nodeToMap(rhs), prefs)
	if err != nil {
		return nil, err
	}

	var pathIndexToStartFrom int = 0
	if results.Front() != nil {
		pathIndexToStartFrom = len(results.Front().Value.(*CandidateNode).Path)
	}

	for el := results.Front(); el != nil; el = el.Next() {
		err := applyAssignment(d, pathIndexToStartFrom, lhs, el.Value.(*CandidateNode), preferences)
		if err != nil {
			return nil, err
		}
	}
	return lhs, nil
}

func applyAssignment(d *dataTreeNavigator, pathIndexToStartFrom int, lhs *CandidateNode, rhs *CandidateNode, preferences multiplyPreferences) error {
	shouldAppendArrays := preferences.AppendArrays
	log.Debugf("merge - applyAssignment lhs %v, rhs: %v", NodeToString(lhs), NodeToString(rhs))

	lhsPath := rhs.Path[pathIndexToStartFrom:]

	assignmentOp := &Operation{OperationType: assignAttributesOpType}
	if rhs.Node.Kind == yaml.ScalarNode || rhs.Node.Kind == yaml.AliasNode {
		assignmentOp.OperationType = assignOpType
		assignmentOp.UpdateAssign = false
	} else if shouldAppendArrays && rhs.Node.Kind == yaml.SequenceNode {
		assignmentOp.OperationType = addAssignOpType
	}
	rhsOp := &Operation{OperationType: valueOpType, CandidateNode: rhs}

	assignmentOpNode := &ExpressionNode{Operation: assignmentOp, Lhs: createTraversalTree(lhsPath, preferences.TraversePrefs), Rhs: &ExpressionNode{Operation: rhsOp}}

	_, err := d.GetMatchingNodes(nodeToMap(lhs), assignmentOpNode)

	return err
}
