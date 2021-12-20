package yqlib

import (
	"fmt"
	"strconv"

	"gopkg.in/yaml.v3"
)

func createSubtractOp(lhs *ExpressionNode, rhs *ExpressionNode) *ExpressionNode {
	return &ExpressionNode{Operation: &Operation{OperationType: subtractOpType},
		Lhs: lhs,
		Rhs: rhs}
}

func subtractAssignOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	return compoundAssignFunction(d, context, expressionNode, createSubtractOp)
}

func subtractOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("Subtract operator")

	return crossFunction(d, context.ReadOnlyClone(), expressionNode, subtract, false)
}

func subtractArray(lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	newLhsArray := make([]*yaml.Node, 0)

	for lindex := 0; lindex < len(lhs.Node.Content); lindex = lindex + 1 {
		shouldInclude := true
		for rindex := 0; rindex < len(rhs.Node.Content) && shouldInclude; rindex = rindex + 1 {
			if recursiveNodeEqual(lhs.Node.Content[lindex], rhs.Node.Content[rindex]) {
				shouldInclude = false
			}
		}
		if shouldInclude {
			newLhsArray = append(newLhsArray, lhs.Node.Content[lindex])
		}
	}
	lhs.Node.Content = newLhsArray
	return lhs, nil
}

func subtract(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	lhs.Node = unwrapDoc(lhs.Node)
	rhs.Node = unwrapDoc(rhs.Node)

	lhsNode := lhs.Node

	if lhsNode.Tag == "!!null" {
		return lhs.CreateReplacement(rhs.Node), nil
	}

	target := lhs.CreateReplacement(&yaml.Node{})

	switch lhsNode.Kind {
	case yaml.MappingNode:
		return nil, fmt.Errorf("Maps not yet supported for subtraction")
	case yaml.SequenceNode:
		if rhs.Node.Kind != yaml.SequenceNode {
			return nil, fmt.Errorf("%v (%v) cannot be subtracted from %v", rhs.Node.Tag, rhs.Path, lhsNode.Tag)
		}
		return subtractArray(lhs, rhs)
	case yaml.ScalarNode:
		if rhs.Node.Kind != yaml.ScalarNode {
			return nil, fmt.Errorf("%v (%v) cannot be subtracted from %v", rhs.Node.Tag, rhs.Path, lhsNode.Tag)
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
		format, lhsNum, err := parseInt(lhs.Value)
		if err != nil {
			return nil, err
		}
		_, rhsNum, err := parseInt(rhs.Value)
		if err != nil {
			return nil, err
		}
		result := lhsNum - rhsNum
		target.Node.Tag = "!!int"
		target.Node.Value = fmt.Sprintf(format, result)
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
