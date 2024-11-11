package yqlib

import (
	"container/list"
	"fmt"
	"strconv"
)

type compareTypePref struct {
	OrEqual bool
	Greater bool
}

func compareOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("compareOperator")
	prefs := expressionNode.Operation.Preferences.(compareTypePref)
	return crossFunction(d, context, expressionNode, compare(prefs), true)
}

func compare(prefs compareTypePref) func(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	return func(_ *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
		log.Debugf("compare cross function")
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

		switch lhs.Kind {
		case MappingNode:
			return nil, fmt.Errorf("maps not yet supported for comparison")
		case SequenceNode:
			return nil, fmt.Errorf("arrays not yet supported for comparison")
		default:
			if rhs.Kind != ScalarNode {
				return nil, fmt.Errorf("%v (%v) cannot be subtracted from %v", rhs.Tag, rhs.GetNicePath(), lhs.Tag)
			}
			target := lhs.CopyWithoutContent()
			boolV, err := compareScalars(context, prefs, lhs, rhs)

			return createBooleanCandidate(target, boolV), err
		}
	}
}

func compareDateTime(layout string, prefs compareTypePref, lhs *CandidateNode, rhs *CandidateNode) (bool, error) {
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

func compareScalars(context Context, prefs compareTypePref, lhs *CandidateNode, rhs *CandidateNode) (bool, error) {
	lhsTag := lhs.guessTagFromCustomType()
	rhsTag := rhs.guessTagFromCustomType()

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

func superlativeByComparison(d *dataTreeNavigator, context Context, prefs compareTypePref) (Context, error) {
	fn := compare(prefs)

	var results = list.New()

	for seq := context.MatchingNodes.Front(); seq != nil; seq = seq.Next() {
		splatted, err := splat(context.SingleChildContext(seq.Value.(*CandidateNode)), traversePreferences{})
		if err != nil {
			return Context{}, err
		}
		result := splatted.MatchingNodes.Front()
		if result != nil {
			for el := result.Next(); el != nil; el = el.Next() {
				cmp, err := fn(d, context, el.Value.(*CandidateNode), result.Value.(*CandidateNode))
				if err != nil {
					return Context{}, err
				}
				if isTruthyNode(cmp) {
					result = el
				}
			}
			results.PushBack(result.Value)
		}
	}
	return context.ChildContext(results), nil
}

func minOperator(d *dataTreeNavigator, context Context, _ *ExpressionNode) (Context, error) {
	log.Debug(("Min"))
	return superlativeByComparison(d, context, compareTypePref{Greater: false})
}

func maxOperator(d *dataTreeNavigator, context Context, _ *ExpressionNode) (Context, error) {
	log.Debug(("Max"))
	return superlativeByComparison(d, context, compareTypePref{Greater: true})
}
