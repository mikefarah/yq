package yqlib

import (
	"container/list"
	"fmt"
)

func createPathNodeFor(pathElement interface{}) *CandidateNode {
	switch pathElement := pathElement.(type) {
	case string:
		return &CandidateNode{Kind: ScalarNode, Value: pathElement, Tag: "!!str"}
	default:
		return &CandidateNode{Kind: ScalarNode, Value: fmt.Sprintf("%v", pathElement), Tag: "!!int"}
	}
}

func getPathArrayFromNode(funcName string, node *CandidateNode) ([]interface{}, error) {
	if node.Kind != SequenceNode {
		return nil, fmt.Errorf("%v: expected path array, but got %v instead", funcName, node.Tag)
	}

	path := make([]interface{}, len(node.Content))

	for i, childNode := range node.Content {
		switch childNode.Tag {
		case "!!str":
			path[i] = childNode.Value
		case "!!int":
			number, err := parseInt(childNode.Value)
			if err != nil {
				return nil, fmt.Errorf("%v: could not parse %v as an int: %w", funcName, childNode.Value, err)
			}
			path[i] = number
		default:
			return nil, fmt.Errorf("%v: expected either a !!str or !!int in the path, found %v instead", funcName, childNode.Tag)
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

	lhsPathContext, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.RHS.LHS)

	if err != nil {
		return Context{}, err
	}

	if lhsPathContext.MatchingNodes.Len() != 1 {
		return Context{}, fmt.Errorf("SETPATH: expected single path but found %v results instead", lhsPathContext.MatchingNodes.Len())
	}
	lhsValue := lhsPathContext.MatchingNodes.Front().Value.(*CandidateNode)

	lhsPath, err := getPathArrayFromNode("SETPATH", lhsValue)

	if err != nil {
		return Context{}, err
	}

	lhsTraversalTree := createTraversalTree(lhsPath, traversePreferences{}, false)

	assignmentOp := &Operation{OperationType: assignOpType}

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		targetContextValue, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(candidate), expressionNode.RHS.RHS)
		if err != nil {
			return Context{}, err
		}

		if targetContextValue.MatchingNodes.Len() != 1 {
			return Context{}, fmt.Errorf("SETPATH: expected single value on RHS but found %v", targetContextValue.MatchingNodes.Len())
		}

		rhsOp := &Operation{OperationType: referenceOpType, CandidateNode: targetContextValue.MatchingNodes.Front().Value.(*CandidateNode)}

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

func delPathsOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("delPaths")
	// single RHS expression that returns an array of paths (array of arrays)

	pathArraysContext, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.RHS)
	if err != nil {
		return Context{}, err
	}
	if pathArraysContext.MatchingNodes.Len() != 1 {
		return Context{}, fmt.Errorf("DELPATHS: expected single value but found %v", pathArraysContext.MatchingNodes.Len())
	}
	pathArraysNode := pathArraysContext.MatchingNodes.Front().Value.(*CandidateNode)

	if pathArraysNode.Tag != "!!seq" {
		return Context{}, fmt.Errorf("DELPATHS: expected a sequence of sequences, but found %v", pathArraysNode.Tag)
	}

	updatedContext := context

	for i, child := range pathArraysNode.Content {

		if child.Tag != "!!seq" {
			return Context{}, fmt.Errorf("DELPATHS: expected entry [%v] to be a sequence, but its a %v. Note that delpaths takes an array of path arrays, e.g. [[\"a\", \"b\"]]", i, child.Tag)
		}
		childPath, err := getPathArrayFromNode("DELPATHS", child)

		if err != nil {
			return Context{}, err
		}

		childTraversalExp := createTraversalTree(childPath, traversePreferences{}, false)
		deleteChildOp := &Operation{OperationType: deleteChildOpType}

		deleteChildOpNode := &ExpressionNode{
			Operation: deleteChildOp,
			RHS:       childTraversalExp,
		}

		updatedContext, err = d.GetMatchingNodes(updatedContext, deleteChildOpNode)

		if err != nil {
			return Context{}, err
		}

	}

	return updatedContext, nil

}

func getPathOperator(_ *dataTreeNavigator, context Context, _ *ExpressionNode) (Context, error) {
	log.Debugf("GetPath")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := candidate.CreateReplacement(SequenceNode, "!!seq", "")

		path := candidate.GetPath()

		content := make([]*CandidateNode, len(path))
		for pathIndex := 0; pathIndex < len(path); pathIndex++ {
			path := path[pathIndex]
			content[pathIndex] = createPathNodeFor(path)
		}
		node.AddChildren(content)
		results.PushBack(node)
	}

	return context.ChildContext(results), nil
}
