package yqlib

import (
	"fmt"
	"strconv"
	"strings"
)

func divideOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("Divide operator")

	return crossFunction(d, context.ReadOnlyClone(), expressionNode, divide, false)
}

func divide(_ *dataTreeNavigator, _ Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	if lhs.Tag == "!!null" {
		return nil, fmt.Errorf("%v (%v) cannot be divided by %v (%v)", lhs.Tag, lhs.GetNicePath(), rhs.Tag, rhs.GetNicePath())
	}

	target := lhs.CopyWithoutContent()

	if lhs.Kind == ScalarNode && rhs.Kind == ScalarNode {
		if err := divideScalars(target, lhs, rhs); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("%v (%v) cannot be divided by %v (%v)", lhs.Tag, lhs.GetNicePath(), rhs.Tag, rhs.GetNicePath())
	}

	return target, nil
}

func divideScalars(target *CandidateNode, lhs *CandidateNode, rhs *CandidateNode) error {
	lhsTag := lhs.Tag
	rhsTag := rhs.guessTagFromCustomType()
	lhsIsCustom := false
	if !strings.HasPrefix(lhsTag, "!!") {
		// custom tag - we have to have a guess
		lhsTag = lhs.guessTagFromCustomType()
		lhsIsCustom = true
	}

	if lhsTag == "!!str" && rhsTag == "!!str" {
		tKind, tTag, res := split(lhs.Value, rhs.Value)
		target.Kind = tKind
		target.Tag = tTag
		target.AddChildren(res)
	} else if (lhsTag == "!!int" || lhsTag == "!!float") && (rhsTag == "!!int" || rhsTag == "!!float") {
		target.Kind = ScalarNode
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
