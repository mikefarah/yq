package yqlib

import (
	"fmt"
	"strconv"

	yaml "gopkg.in/yaml.v3"
)

type compareTypePref struct {
	OrEqual bool
	Greater bool
}

func compareOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- compareOperator")
	prefs := expressionNode.Operation.Preferences.(compareTypePref)
	return crossFunction(d, context, expressionNode, compare(prefs), true)
}

func compare(prefs compareTypePref) func(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	return func(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
		log.Debugf("-- compare cross function")
		if lhs == nil && rhs == nil {
			owner := &CandidateNode{}
			return createBooleanCandidate(owner, prefs.OrEqual), nil
		} else if lhs == nil {
			log.Debugf("lhs nil, but rhs is not")
			return createBooleanCandidate(rhs, false), nil
		} else if rhs == nil {
			log.Debugf("rhs nil, but rhs is not")
			return createBooleanCandidate(lhs, false), nil
		}

		lhs.Node = unwrapDoc(lhs.Node)
		rhs.Node = unwrapDoc(rhs.Node)

		switch lhs.Node.Kind {
		case yaml.MappingNode:
			return nil, fmt.Errorf("maps not yet supported for comparison")
		case yaml.SequenceNode:
			return nil, fmt.Errorf("arrays not yet supported for comparison")
		default:
			if rhs.Node.Kind != yaml.ScalarNode {
				return nil, fmt.Errorf("%v (%v) cannot be subtracted from %v", rhs.Node.Tag, rhs.Path, lhs.Node.Tag)
			}
			target := lhs.CreateReplacement(&yaml.Node{})
			boolV, err := compareScalars(context, prefs, lhs.Node, rhs.Node)

			return createBooleanCandidate(target, boolV), err
		}
	}
}

func compareDateTime(layout string, prefs compareTypePref, lhs *yaml.Node, rhs *yaml.Node) (bool, error) {
	lhsTime, err := parseDateTime(layout, lhs.Value)
	if err != nil {
		return false, err
	}

	rhsTime, err := parseDateTime(layout, rhs.Value)
	if err != nil {
		return false, err
	}

	if prefs.OrEqual && lhsTime.Equal(rhsTime) {
		return true, nil
	}
	if prefs.Greater {
		return lhsTime.After(rhsTime), nil
	}
	return lhsTime.Before(rhsTime), nil

}

func compareScalars(context Context, prefs compareTypePref, lhs *yaml.Node, rhs *yaml.Node) (bool, error) {
	lhsTag := guessTagFromCustomType(lhs)
	rhsTag := guessTagFromCustomType(rhs)

	isDateTime := lhs.Tag == "!!timestamp"
	// if the lhs is a string, it might be a timestamp in a custom format.
	if lhsTag == "!!str" {
		_, err := parseDateTime(context.GetDateTimeLayout(), lhs.Value)
		isDateTime = err == nil
	}
	if isDateTime {
		return compareDateTime(context.GetDateTimeLayout(), prefs, lhs, rhs)
	} else if lhsTag == "!!int" && rhsTag == "!!int" {
		_, lhsNum, err := parseInt64(lhs.Value)
		if err != nil {
			return false, err
		}
		_, rhsNum, err := parseInt64(rhs.Value)
		if err != nil {
			return false, err
		}

		if prefs.OrEqual && lhsNum == rhsNum {
			return true, nil
		}
		if prefs.Greater {
			return lhsNum > rhsNum, nil
		}
		return lhsNum < rhsNum, nil
	} else if (lhsTag == "!!int" || lhsTag == "!!float") && (rhsTag == "!!int" || rhsTag == "!!float") {
		lhsNum, err := strconv.ParseFloat(lhs.Value, 64)
		if err != nil {
			return false, err
		}
		rhsNum, err := strconv.ParseFloat(rhs.Value, 64)
		if err != nil {
			return false, err
		}
		if prefs.OrEqual && lhsNum == rhsNum {
			return true, nil
		}
		if prefs.Greater {
			return lhsNum > rhsNum, nil
		}
		return lhsNum < rhsNum, nil
	} else if lhsTag == "!!str" && rhsTag == "!!str" {
		if prefs.OrEqual && lhs.Value == rhs.Value {
			return true, nil
		}
		if prefs.Greater {
			return lhs.Value > rhs.Value, nil
		}
		return lhs.Value < rhs.Value, nil
	} else if lhsTag == "!!null" && rhsTag == "!!null" && prefs.OrEqual {
		return true, nil
	} else if lhsTag == "!!null" || rhsTag == "!!null" {
		return false, nil
	}

	return false, fmt.Errorf("%v not yet supported for comparison", lhs.Tag)
}
