package yqlib

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func moduloOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("Modulo operator")

	return crossFunction(d, context.ReadOnlyClone(), expressionNode, modulo, false)
}

func modulo(_ *dataTreeNavigator, _ Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	if lhs.Tag == "!!null" {
		return nil, fmt.Errorf("%v (%v) cannot modulo by %v (%v)", lhs.Tag, lhs.GetNicePath(), rhs.Tag, rhs.GetNicePath())
	}

	target := lhs.CopyWithoutContent()

	if lhs.Kind == ScalarNode && rhs.Kind == ScalarNode {
		if err := moduloScalars(target, lhs, rhs); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("%v (%v) cannot modulo by %v (%v)", lhs.Tag, lhs.GetNicePath(), rhs.Tag, rhs.GetNicePath())
	}

	return target, nil
}

func moduloScalars(target *CandidateNode, lhs *CandidateNode, rhs *CandidateNode) error {
	lhsTag := lhs.Tag
	rhsTag := rhs.guessTagFromCustomType()
	lhsIsCustom := false
	if !strings.HasPrefix(lhsTag, "!!") {
		// custom tag - we have to have a guess
		lhsTag = lhs.guessTagFromCustomType()
		lhsIsCustom = true
	}

	if lhsTag == "!!int" && rhsTag == "!!int" {
		target.Kind = ScalarNode
		target.Style = lhs.Style

		format, lhsNum, err := parseInt64(lhs.Value)
		if err != nil {
			return err
		}
		_, rhsNum, err := parseInt64(rhs.Value)
		if err != nil {
			return err
		}
		if rhsNum == 0 {
			return fmt.Errorf("cannot modulo by 0")
		}
		remainder := lhsNum % rhsNum

		target.Tag = lhs.Tag
		target.Value = fmt.Sprintf(format, remainder)
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
		remainder := math.Mod(lhsNum, rhsNum)
		if lhsIsCustom {
			target.Tag = lhs.Tag
		} else {
			target.Tag = "!!float"
		}
		target.Value = fmt.Sprintf("%v", remainder)
	} else {
		return fmt.Errorf("%v cannot modulo by %v", lhsTag, rhsTag)
	}
	return nil
}
