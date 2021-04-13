package yqlib

import (
	"fmt"

	"strconv"

	yaml "gopkg.in/yaml.v3"
)

func createSubtractOp(lhs *ExpressionNode, rhs *ExpressionNode) *ExpressionNode {
	return &ExpressionNode{Operation: &Operation{OperationType: subtractOpType},
		Lhs: lhs,
		Rhs: rhs}
}

func subtractAssignOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	assignmentOp := &Operation{OperationType: assignOpType}
	assignmentOp.UpdateAssign = true
	selfExpression := &ExpressionNode{Operation: &Operation{OperationType: selfReferenceOpType}}
	assignmentOpNode := &ExpressionNode{Operation: assignmentOp, Lhs: expressionNode.Lhs, Rhs: createSubtractOp(selfExpression, expressionNode.Rhs)}
	return d.GetMatchingNodes(context, assignmentOpNode)
}

func subtractOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("Subtract operator")

	return crossFunction(d, context, expressionNode, subtract, false)
}

func subtract(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	lhs.Node = unwrapDoc(lhs.Node)
	rhs.Node = unwrapDoc(rhs.Node)

	lhsNode := lhs.Node

	if lhsNode.Tag == "!!null" {
		return lhs.CreateChild(nil, rhs.Node), nil
	}

	target := lhs.CreateChild(nil, &yaml.Node{})

	switch lhsNode.Kind {
	case yaml.MappingNode:
		return nil, fmt.Errorf("Maps not yet supported for subtraction")
	case yaml.SequenceNode:
		return nil, fmt.Errorf("Sequences not yet supported for subtraction")
		// target.Node.Kind = yaml.SequenceNode
		// target.Node.Style = lhsNode.Style
		// target.Node.Tag = "!!seq"
		// target.Node.Content = append(lhsNode.Content, toNodes(rhs)...)
	case yaml.ScalarNode:
		if rhs.Node.Kind != yaml.ScalarNode {
			return nil, fmt.Errorf("%v (%v) cannot be added to a %v", rhs.Node.Tag, rhs.Path, lhsNode.Tag)
		}
		target.Node.Kind = yaml.ScalarNode
		target.Node.Style = lhsNode.Style
		return subtractScalars(target, lhsNode, rhs.Node)
	}

	return target, nil
}

func subtractScalars(target *CandidateNode, lhs *yaml.Node, rhs *yaml.Node) (*CandidateNode, error) {

	if lhs.Tag == "!!str" {
		return nil, fmt.Errorf("strings cannot be subtracted")
	} else if lhs.Tag == "!!int" && rhs.Tag == "!!int" {
		lhsNum, err := strconv.Atoi(lhs.Value)
		if err != nil {
			return nil, err
		}
		rhsNum, err := strconv.Atoi(rhs.Value)
		if err != nil {
			return nil, err
		}
		result := lhsNum - rhsNum
		target.Node.Tag = "!!int"
		target.Node.Value = fmt.Sprintf("%v", result)
	} else if (lhs.Tag == "!!int" || lhs.Tag == "!!float") && (rhs.Tag == "!!int" || rhs.Tag == "!!float") {
		lhsNum, err := strconv.ParseFloat(lhs.Value, 64)
		if err != nil {
			return nil, err
		}
		rhsNum, err := strconv.ParseFloat(rhs.Value, 64)
		if err != nil {
			return nil, err
		}
		result := lhsNum - rhsNum
		target.Node.Tag = "!!float"
		target.Node.Value = fmt.Sprintf("%v", result)
	} else {
		return nil, fmt.Errorf("%v cannot be added to %v", lhs.Tag, rhs.Tag)
	}

	return target, nil
}
