package yqlib

import (
	"fmt"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

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

	target := &yaml.Node{}

	if lhsNode.Kind == yaml.ScalarNode && rhs.Node.Kind == yaml.ScalarNode {
		if err := divideScalars(target, lhsNode, rhs.Node); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("%v (%v) cannot be divided by %v (%v)", lhsNode.Tag, lhs.GetNicePath(), rhs.Node.Tag, rhs.GetNicePath())
	}

	return lhs.CreateReplacement(target), nil
}

func divideScalars(target *yaml.Node, lhs *yaml.Node, rhs *yaml.Node) error {
	lhsTag := lhs.Tag
	rhsTag := guessTagFromCustomType(rhs)
	lhsIsCustom := false
	if !strings.HasPrefix(lhsTag, "!!") {
		// custom tag - we have to have a guess
		lhsTag = guessTagFromCustomType(lhs)
		lhsIsCustom = true
	}

	if lhsTag == "!!str" && rhsTag == "!!str" {
		res := split(lhs.Value, rhs.Value)
		target.Kind = res.Kind
		target.Tag = res.Tag
		target.Content = res.Content
	} else if (lhsTag == "!!int" || lhsTag == "!!float") && (rhsTag == "!!int" || rhsTag == "!!float") {
		target.Kind = yaml.ScalarNode
		target.Style = lhs.Style

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
			target.Tag = lhs.Tag
		} else {
			target.Tag = "!!float"
		}
		target.Value = fmt.Sprintf("%v", quotient)
	} else {
		return fmt.Errorf("%v cannot be divided by %v", lhsTag, rhsTag)
	}
	return nil
}
