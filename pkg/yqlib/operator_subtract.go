package yqlib

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
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

func subtractArray(lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	newLHSArray := make([]*yaml.Node, 0)

	for lindex := 0; lindex < len(lhs.Node.Content); lindex = lindex + 1 {
		shouldInclude := true
		for rindex := 0; rindex < len(rhs.Node.Content) && shouldInclude; rindex = rindex + 1 {
			if recursiveNodeEqual(lhs.Node.Content[lindex], rhs.Node.Content[rindex]) {
				shouldInclude = false
			}
		}
		if shouldInclude {
			newLHSArray = append(newLHSArray, lhs.Node.Content[lindex])
		}
	}
	lhs.Node.Content = newLHSArray
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
		return nil, fmt.Errorf("maps not yet supported for subtraction")
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
		if err := subtractScalars(context, target, lhsNode, rhs.Node); err != nil {
			return nil, err
		}
	}

	return target, nil
}

func subtractScalars(context Context, target *CandidateNode, lhs *yaml.Node, rhs *yaml.Node) error {
	lhsTag := lhs.Tag
	rhsTag := rhs.Tag
	lhsIsCustom := false
	if !strings.HasPrefix(lhsTag, "!!") {
		// custom tag - we have to have a guess
		lhsTag = guessTagFromCustomType(lhs)
		lhsIsCustom = true
	}

	if !strings.HasPrefix(rhsTag, "!!") {
		// custom tag - we have to have a guess
		rhsTag = guessTagFromCustomType(rhs)
	}

	isDateTime := lhs.Tag == "!!timestamp"
	// if the lhs is a string, it might be a timestamp in a custom format.
	if lhsTag == "!!str" && context.GetDateTimeLayout() != time.RFC3339 {
		_, err := time.Parse(context.GetDateTimeLayout(), lhs.Value)
		isDateTime = err == nil
	}

	if isDateTime {
		return subtractDateTime(context.GetDateTimeLayout(), target, lhs, rhs)
	} else if lhsTag == "!!str" {
		return fmt.Errorf("strings cannot be subtracted")
	} else if lhsTag == "!!int" && rhsTag == "!!int" {
		format, lhsNum, err := parseInt(lhs.Value)
		if err != nil {
			return err
		}
		_, rhsNum, err := parseInt(rhs.Value)
		if err != nil {
			return err
		}
		result := lhsNum - rhsNum
		target.Node.Tag = lhs.Tag
		target.Node.Value = fmt.Sprintf(format, result)
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
			target.Node.Tag = lhs.Tag
		} else {
			target.Node.Tag = "!!float"
		}
		target.Node.Value = fmt.Sprintf("%v", result)
	} else {
		return fmt.Errorf("%v cannot be added to %v", lhs.Tag, rhs.Tag)
	}

	return nil
}

func subtractDateTime(layout string, target *CandidateNode, lhs *yaml.Node, rhs *yaml.Node) error {
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

	currentTime, err := time.Parse(layout, lhs.Value)
	if err != nil {
		return err
	}

	newTime := currentTime.Add(duration)
	target.Node.Value = newTime.Format(layout)
	return nil
}
