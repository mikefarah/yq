package yqlib

import (
	"fmt"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

func createDivideOp(lhs *ExpressionNode, rhs *ExpressionNode) *ExpressionNode {
	return &ExpressionNode{Operation: &Operation{OperationType: divideOpType},
		LHS: lhs,
		RHS: rhs}
}

// func divideAssignOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
// 	return compoundAssignFunction(d, context, expressionNode, createDivideOp)
// }

func divideOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("Divide operator")

	return crossFunction(d, context.ReadOnlyClone(), expressionNode, divide, false)
}

func divide(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	lhs.Node = unwrapDoc(lhs.Node)
	rhs.Node = unwrapDoc(rhs.Node)

	lhsNode := lhs.Node

	if lhsNode.Tag == "!!null" {
		return nil, fmt.Errorf("%v (%v) cannot be divided by %v (%v)", lhsNode.Tag, lhs.GetNicePath(), rhs.Node.Tag, rhs.GetNicePath())
	}

	target := lhs.CreateReplacement(&yaml.Node{
		Anchor: lhs.Node.Anchor,
	})

	if lhsNode.Kind == yaml.ScalarNode && rhs.Node.Kind == yaml.ScalarNode {
		if err := divideScalars(context, target, lhsNode, rhs.Node); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("%v (%v) cannot be divided by %v (%v)", lhsNode.Tag, lhs.GetNicePath(), rhs.Node.Tag, rhs.GetNicePath())
	}

	return target, nil
}

func divideScalars(context Context, target *CandidateNode, lhs *yaml.Node, rhs *yaml.Node) error {
	lhsTag := lhs.Tag
	rhsTag := guessTagFromCustomType(rhs)
	lhsIsCustom := false
	if !strings.HasPrefix(lhsTag, "!!") {
		// custom tag - we have to have a guess
		lhsTag = guessTagFromCustomType(lhs)
		lhsIsCustom = true
	}

	if lhsTag == "!!str" && rhsTag == "!!str" {
		target.Node = split(lhs.Value, rhs.Value)
		target.Node.Anchor = lhs.Anchor
	} else if (lhsTag == "!!int" || lhsTag == "!!float") && (rhsTag == "!!int" || rhsTag == "!!float") {
		target.Node.Kind = yaml.ScalarNode
		target.Node.Style = lhs.Style

		lhsNum, err := strconv.ParseFloat(lhs.Value, 64)
		if err != nil {
			return err
		}
		rhsNum, err := strconv.ParseFloat(rhs.Value, 64)
		if err != nil {
			return err
		}
		quotient := lhsNum / rhsNum
		if lhsIsCustom {
			target.Node.Tag = lhs.Tag
		} else {
			target.Node.Tag = "!!float"
		}
		target.Node.Value = fmt.Sprintf("%v", quotient)
	} else {
		return fmt.Errorf("%v cannot be divided by %v", lhsTag, rhsTag)
	}
	return nil
}
