package yqlib

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v3"
)

func createAddOp(lhs *ExpressionNode, rhs *ExpressionNode) *ExpressionNode {
	return &ExpressionNode{Operation: &Operation{OperationType: addOpType},
		LHS: lhs,
		RHS: rhs}
}

func addAssignOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	return compoundAssignFunction(d, context, expressionNode, createAddOp)
}

func toNodes(candidate *CandidateNode, lhs *CandidateNode) ([]*yaml.Node, error) {
	if candidate.Node.Tag == "!!null" {
		return []*yaml.Node{}, nil
	}
	clone, err := candidate.Copy()
	if err != nil {
		return nil, err
	}

	switch candidate.Node.Kind {
	case yaml.SequenceNode:
		return clone.Node.Content, nil
	default:
		if len(lhs.Node.Content) > 0 {
			clone.Node.Style = lhs.Node.Content[0].Style
		}
		return []*yaml.Node{clone.Node}, nil
	}

}

func addOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("Add operator")

	return crossFunction(d, context.ReadOnlyClone(), expressionNode, add, false)
}

func add(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	lhs.Node = unwrapDoc(lhs.Node)
	rhs.Node = unwrapDoc(rhs.Node)

	lhsNode := lhs.Node

	if lhsNode.Tag == "!!null" {
		return lhs.CreateReplacement(rhs.Node), nil
	}

	target := lhs.CreateReplacement(&yaml.Node{
		Anchor: lhs.Node.Anchor,
	})

	switch lhsNode.Kind {
	case yaml.MappingNode:
		if rhs.Node.Kind != yaml.MappingNode {
			return nil, fmt.Errorf("%v (%v) cannot be added to a %v (%v)", rhs.Node.Tag, rhs.GetNicePath(), lhsNode.Tag, lhs.GetNicePath())
		}
		addMaps(target, lhs, rhs)
	case yaml.SequenceNode:
		if err := addSequences(target, lhs, rhs); err != nil {
			return nil, err
		}

	case yaml.ScalarNode:
		if rhs.Node.Kind != yaml.ScalarNode {
			return nil, fmt.Errorf("%v (%v) cannot be added to a %v (%v)", rhs.Node.Tag, rhs.GetNicePath(), lhsNode.Tag, lhs.GetNicePath())
		}
		target.Node.Kind = yaml.ScalarNode
		target.Node.Style = lhsNode.Style
		if err := addScalars(context, target, lhsNode, rhs.Node); err != nil {
			return nil, err
		}
	}
	return target, nil
}

func addScalars(context Context, target *CandidateNode, lhs *yaml.Node, rhs *yaml.Node) error {
	lhsTag := lhs.Tag
	rhsTag := guessTagFromCustomType(rhs)
	lhsIsCustom := false
	if !strings.HasPrefix(lhsTag, "!!") {
		// custom tag - we have to have a guess
		lhsTag = guessTagFromCustomType(lhs)
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
		target.Node.Tag = lhs.Tag
		target.Node.Value = lhs.Value + rhs.Value
	} else if rhsTag == "!!str" {
		target.Node.Tag = rhs.Tag
		target.Node.Value = lhs.Value + rhs.Value
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
		target.Node.Tag = lhs.Tag
		target.Node.Value = fmt.Sprintf(format, sum)
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
			target.Node.Tag = lhs.Tag
		} else {
			target.Node.Tag = "!!float"
		}
		target.Node.Value = fmt.Sprintf("%v", sum)
	} else {
		return fmt.Errorf("%v cannot be added to %v", lhsTag, rhsTag)
	}
	return nil
}

func addDateTimes(layout string, target *CandidateNode, lhs *yaml.Node, rhs *yaml.Node) error {

	duration, err := time.ParseDuration(rhs.Value)
	if err != nil {
		return fmt.Errorf("unable to parse duration [%v]: %w", rhs.Value, err)
	}

	currentTime, err := parseDateTime(layout, lhs.Value)
	if err != nil {
		return err
	}

	newTime := currentTime.Add(duration)
	target.Node.Value = newTime.Format(layout)
	return nil

}

func addSequences(target *CandidateNode, lhs *CandidateNode, rhs *CandidateNode) error {
	target.Node.Kind = yaml.SequenceNode
	if len(lhs.Node.Content) > 0 {
		target.Node.Style = lhs.Node.Style
	}
	target.Node.Tag = lhs.Node.Tag

	extraNodes, err := toNodes(rhs, lhs)
	if err != nil {
		return err
	}

	target.Node.Content = append(deepCloneContent(lhs.Node.Content), extraNodes...)
	return nil

}

func addMaps(target *CandidateNode, lhsC *CandidateNode, rhsC *CandidateNode) {
	lhs := lhsC.Node
	rhs := rhsC.Node

	target.Node.Content = make([]*yaml.Node, len(lhs.Content))
	copy(target.Node.Content, lhs.Content)

	for index := 0; index < len(rhs.Content); index = index + 2 {
		key := rhs.Content[index]
		value := rhs.Content[index+1]
		log.Debug("finding %v", key.Value)
		indexInLHS := findKeyInMap(target.Node, key)
		log.Debug("indexInLhs %v", indexInLHS)
		if indexInLHS < 0 {
			// not in there, append it
			target.Node.Content = append(target.Node.Content, key, value)
		} else {
			// it's there, replace it
			target.Node.Content[indexInLHS+1] = value
		}
	}
	target.Node.Kind = yaml.MappingNode
	if len(lhs.Content) > 0 {
		target.Node.Style = lhs.Style
	}
	target.Node.Tag = lhs.Tag
}
