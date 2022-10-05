package yqlib

import (
	"container/list"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

func createPathNodeFor(pathElement interface{}) *yaml.Node {
	switch pathElement := pathElement.(type) {
	case string:
		return &yaml.Node{Kind: yaml.ScalarNode, Value: pathElement, Tag: "!!str"}
	default:
		return &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%v", pathElement), Tag: "!!int"}
	}
}

func getPathArrayFromExp(d *dataTreeNavigator, context Context, pathExp *ExpressionNode) ([]interface{}, error) {
	lhsPathContext, err := d.GetMatchingNodes(context.ReadOnlyClone(), pathExp)

	if err != nil {
		return nil, err
	}

	if lhsPathContext.MatchingNodes.Len() != 1 {
		return nil, fmt.Errorf("expected single path but found %v results instead", lhsPathContext.MatchingNodes.Len())
	}
	lhsValue := lhsPathContext.MatchingNodes.Front().Value.(*CandidateNode)
	if lhsValue.Node.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("expected path array, but got %v instead", lhsValue.Node.Tag)
	}

	path := make([]interface{}, len(lhsValue.Node.Content))

	for i, childNode := range lhsValue.Node.Content {
		if childNode.Tag == "!!str" {
			path[i] = childNode.Value
		} else if childNode.Tag == "!!int" {
			number, err := parseInt(childNode.Value)
			if err != nil {
				return nil, fmt.Errorf("could not parse %v as an int: %w", childNode.Value, err)
			}
			path[i] = number
		} else {
			return nil, fmt.Errorf("expected either a !!str or !!int in the path, found %v instead", childNode.Tag)
		}

	}
	return path, nil
}

// SETPATH(pathArray; value)
func setPathOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("SetPath")

	if expressionNode.RHS.Operation.OperationType != blockOpType {
		return Context{}, fmt.Errorf("SETPATH must be given a block (;), got %v instead", expressionNode.RHS.Operation.OperationType.Type)
	}

	lhsPath, err := getPathArrayFromExp(d, context, expressionNode.RHS.LHS)

	if err != nil {
		return Context{}, err
	}

	lhsTraversalTree := createTraversalTree(lhsPath, traversePreferences{}, false)

	assignmentOp := &Operation{OperationType: assignOpType}

	//TODO if context is empty, create a new one

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		targetContextValue, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(candidate), expressionNode.RHS.RHS)
		if err != nil {
			return Context{}, err
		}

		if targetContextValue.MatchingNodes.Len() != 1 {
			return Context{}, fmt.Errorf("Expected single value on RHS but found %v", targetContextValue.MatchingNodes.Len())
		}

		rhsOp := &Operation{OperationType: valueOpType, CandidateNode: targetContextValue.MatchingNodes.Front().Value.(*CandidateNode)}

		assignmentOpNode := &ExpressionNode{
			Operation: assignmentOp,
			LHS:       lhsTraversalTree,
			RHS:       &ExpressionNode{Operation: rhsOp},
		}

		_, err = d.GetMatchingNodes(context.SingleChildContext(candidate), assignmentOpNode)

		if err != nil {
			return Context{}, err
		}

	}
	return context, nil
}

func getPathOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("GetPath")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}

		content := make([]*yaml.Node, len(candidate.Path))
		for pathIndex := 0; pathIndex < len(candidate.Path); pathIndex++ {
			path := candidate.Path[pathIndex]
			content[pathIndex] = createPathNodeFor(path)
		}
		node.Content = content
		result := candidate.CreateReplacement(node)
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}
