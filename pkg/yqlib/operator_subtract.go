package yqlib

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func createSubtractOp(lhs *ExpressionNode, rhs *ExpressionNode) *ExpressionNode {
	return &ExpressionNode{Operation: &Operation{OperationType: subtractOpType},
		LHS: lhs,
		RHS: rhs}
}

func subtractAssignOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	return compoundAssignFunction(d, context, expressionNode, createSubtractOp)
}

func subtractOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("Subtract operator")

	return crossFunction(d, context.ReadOnlyClone(), expressionNode, subtract, false)
}

func subtractArray(lhs *CandidateNode, rhs *CandidateNode) []*CandidateNode {
	newLHSArray := make([]*CandidateNode, 0)

	for lindex := 0; lindex < len(lhs.Content); lindex = lindex + 1 {
		shouldInclude := true
		for rindex := 0; rindex < len(rhs.Content) && shouldInclude; rindex = rindex + 1 {
			if recursiveNodeEqual(lhs.Content[lindex], rhs.Content[rindex]) {
				shouldInclude = false
			}
		}
		if shouldInclude {
			newLHSArray = append(newLHSArray, lhs.Content[lindex])
		}
	}
	return newLHSArray
}

func subtract(_ *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	if lhs.Tag == "!!null" {
		return lhs.CopyAsReplacement(rhs), nil
	}

	target := lhs.CopyWithoutContent()

	switch lhs.Kind {
	case MappingNode:
		return nil, fmt.Errorf("maps not yet supported for subtraction")
	case SequenceNode:
		if rhs.Kind != SequenceNode {
			return nil, fmt.Errorf("%v (%v) cannot be subtracted from %v", rhs.Tag, rhs.GetNicePath(), lhs.Tag)
		}
		target.Content = subtractArray(lhs, rhs)
	case ScalarNode:
		if rhs.Kind != ScalarNode {
			return nil, fmt.Errorf("%v (%v) cannot be subtracted from %v", rhs.Tag, rhs.GetNicePath(), lhs.Tag)
		}
		target.Kind = ScalarNode
		target.Style = lhs.Style
		if err := subtractScalars(context, target, lhs, rhs); err != nil {
			return nil, err
		}
	}

	return target, nil
}

func subtractScalars(context Context, target *CandidateNode, lhs *CandidateNode, rhs *CandidateNode) error {
	lhsTag := lhs.Tag
	rhsTag := rhs.Tag
	lhsIsCustom := false
	if !strings.HasPrefix(lhsTag, "!!") {
		// custom tag - we have to have a guess
		lhsTag = lhs.guessTagFromCustomType()
		lhsIsCustom = true
	}

	if !strings.HasPrefix(rhsTag, "!!") {
		// custom tag - we have to have a guess
		rhsTag = rhs.guessTagFromCustomType()
	}

	isDateTime := lhsTag == "!!timestamp"
	// if the lhs is a string, it might be a timestamp in a custom format.
	if lhsTag == "!!str" && context.GetDateTimeLayout() != time.RFC3339 {
		_, err := parseDateTime(context.GetDateTimeLayout(), lhs.Value)
		isDateTime = err == nil
	}

	if isDateTime {
		return subtractDateTime(context.GetDateTimeLayout(), target, lhs, rhs)
	} else if lhsTag == "!!str" {
		return fmt.Errorf("strings cannot be subtracted")
	} else if lhsTag == "!!int" && rhsTag == "!!int" {
		format, lhsNum, err := parseInt64(lhs.Value)
		if err != nil {
			return err
		}
		_, rhsNum, err := parseInt64(rhs.Value)
		if err != nil {
			return err
		}
		result := lhsNum - rhsNum
		target.Tag = lhs.Tag
		target.Value = fmt.Sprintf(format, result)
	} else if (lhsTag == "!!int" || lhsTag == "!!float") && (rhsTag == "!!int" || rhsTag == "!!float") {
		lhsNum, err := strconv.ParseFloat(lhs.Value, 64)
		if err != nil {
			return err
		}
		rhsNum, err := strconv.ParseFloat(rhs.Value, 64)
		if err != nil {
			return err
		}
		result := lhsNum - rhsNum
		if lhsIsCustom {
			target.Tag = lhs.Tag
		} else {
			target.Tag = "!!float"
		}
		target.Value = fmt.Sprintf("%v", result)
	} else {
		return fmt.Errorf("%v cannot be added to %v", lhs.Tag, rhs.Tag)
	}

	return nil
}

func subtractDateTime(layout string, target *CandidateNode, lhs *CandidateNode, rhs *CandidateNode) error {
	var durationStr string
	if strings.HasPrefix(rhs.Value, "-") {
		durationStr = rhs.Value[1:]
	} else {
		durationStr = "-" + rhs.Value
	}
	duration, err := time.ParseDuration(durationStr)

	if err != nil {
		return fmt.Errorf("unable to parse duration [%v]: %w", rhs.Value, err)
	}

	currentTime, err := parseDateTime(layout, lhs.Value)
	if err != nil {
		return err
	}

	newTime := currentTime.Add(duration)
	target.Value = newTime.Format(layout)
	return nil
}
