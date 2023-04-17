package yqlib

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func createAddOp(lhs *ExpressionNode, rhs *ExpressionNode) *ExpressionNode {
	return &ExpressionNode{Operation: &Operation{OperationType: addOpType},
		LHS: lhs,
		RHS: rhs}
}

func addAssignOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	return compoundAssignFunction(d, context, expressionNode, createAddOp)
}

func toNodes(candidate *CandidateNode, lhs *CandidateNode) ([]*CandidateNode, error) {
	if candidate.Tag == "!!null" {
		return []*CandidateNode{}, nil
	}
	clone := candidate.Copy()

	switch candidate.Kind {
	case SequenceNode:
		return clone.Content, nil
	default:
		if len(lhs.Content) > 0 {
			clone.Style = lhs.Content[0].Style
		}
		return []*CandidateNode{clone}, nil
	}

}

func addOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("Add operator")

	return crossFunction(d, context.ReadOnlyClone(), expressionNode, add, false)
}

func add(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	lhs = lhs.unwrapDocument()
	rhs = rhs.unwrapDocument()

	lhsNode := lhs

	if lhsNode.Tag == "!!null" {
		return lhs.CopyAsReplacement(rhs), nil
	}

	target := lhs.CopyWithoutContent()

	switch lhsNode.Kind {
	case MappingNode:
		if rhs.Kind != MappingNode {
			return nil, fmt.Errorf("%v (%v) cannot be added to a %v (%v)", rhs.Tag, rhs.GetNicePath(), lhsNode.Tag, lhs.GetNicePath())
		}
		addMaps(target, lhs, rhs)
	case SequenceNode:
		if err := addSequences(target, lhs, rhs); err != nil {
			return nil, err
		}

	case ScalarNode:
		if rhs.Kind != ScalarNode {
			return nil, fmt.Errorf("%v (%v) cannot be added to a %v (%v)", rhs.Tag, rhs.GetNicePath(), lhsNode.Tag, lhs.GetNicePath())
		}
		target.Kind = ScalarNode
		target.Style = lhsNode.Style
		if err := addScalars(context, target, lhsNode, rhs); err != nil {
			return nil, err
		}
	}
	return target, nil
}

func addScalars(context Context, target *CandidateNode, lhs *CandidateNode, rhs *CandidateNode) error {
	lhsTag := lhs.Tag
	rhsTag := rhs.guessTagFromCustomType()
	lhsIsCustom := false
	if !strings.HasPrefix(lhsTag, "!!") {
		// custom tag - we have to have a guess
		lhsTag = lhs.guessTagFromCustomType()
		lhsIsCustom = true
	}

	isDateTime := lhs.Tag == "!!timestamp"

	// if the lhs is a string, it might be a timestamp in a custom format.
	if lhsTag == "!!str" && context.GetDateTimeLayout() != time.RFC3339 {
		_, err := parseDateTime(context.GetDateTimeLayout(), lhs.Value)
		isDateTime = err == nil
	}

	if isDateTime {
		return addDateTimes(context.GetDateTimeLayout(), target, lhs, rhs)

	} else if lhsTag == "!!str" {
		target.Tag = lhs.Tag
		target.Value = lhs.Value + rhs.Value
	} else if rhsTag == "!!str" {
		target.Tag = rhs.Tag
		target.Value = lhs.Value + rhs.Value
	} else if lhsTag == "!!int" && rhsTag == "!!int" {
		format, lhsNum, err := parseInt64(lhs.Value)
		if err != nil {
			return err
		}
		_, rhsNum, err := parseInt64(rhs.Value)
		if err != nil {
			return err
		}
		sum := lhsNum + rhsNum
		target.Tag = lhs.Tag
		target.Value = fmt.Sprintf(format, sum)
	} else if (lhsTag == "!!int" || lhsTag == "!!float") && (rhsTag == "!!int" || rhsTag == "!!float") {
		lhsNum, err := strconv.ParseFloat(lhs.Value, 64)
		if err != nil {
			return err
		}
		rhsNum, err := strconv.ParseFloat(rhs.Value, 64)
		if err != nil {
			return err
		}
		sum := lhsNum + rhsNum
		if lhsIsCustom {
			target.Tag = lhs.Tag
		} else {
			target.Tag = "!!float"
		}
		target.Value = fmt.Sprintf("%v", sum)
	} else {
		return fmt.Errorf("%v cannot be added to %v", lhsTag, rhsTag)
	}
	return nil
}

func addDateTimes(layout string, target *CandidateNode, lhs *CandidateNode, rhs *CandidateNode) error {

	duration, err := time.ParseDuration(rhs.Value)
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

func addSequences(target *CandidateNode, lhs *CandidateNode, rhs *CandidateNode) error {
	log.Debugf("adding sequences! target: %v; lhs %v; rhs: %v", NodeToString(target), NodeToString(lhs), NodeToString(rhs))
	target.Kind = SequenceNode
	if len(lhs.Content) == 0 {
		log.Debugf("dont copy lhs style")
		target.Style = 0
	}
	target.Tag = lhs.Tag

	extraNodes, err := toNodes(rhs, lhs)
	if err != nil {
		return err
	}

	target.Content = append(lhs.CopyChildren(), extraNodes...)
	return nil

}

func addMaps(target *CandidateNode, lhsC *CandidateNode, rhsC *CandidateNode) {
	lhs := lhsC
	rhs := rhsC

	if len(lhs.Content) == 0 {
		log.Debugf("dont copy lhs style")
		target.Style = 0
	}

	target.Content = make([]*CandidateNode, len(lhs.Content))
	copy(target.Content, lhs.Content)

	for index := 0; index < len(rhs.Content); index = index + 2 {
		key := rhs.Content[index]
		value := rhs.Content[index+1]
		log.Debug("finding %v", key.Value)
		indexInLHS := findKeyInMap(target, key)
		log.Debug("indexInLhs %v", indexInLHS)
		if indexInLHS < 0 {
			// not in there, append it
			target.Content = append(target.Content, key, value)
		} else {
			// it's there, replace it
			target.Content[indexInLHS+1] = value
		}
	}
	target.Kind = MappingNode
	if len(lhs.Content) > 0 {
		target.Style = lhs.Style
	}
	target.Tag = lhs.Tag
}
