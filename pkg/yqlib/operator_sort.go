package yqlib

import (
	"container/list"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func sortOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	selfExpression := &ExpressionNode{Operation: &Operation{OperationType: selfReferenceOpType}}
	expressionNode.RHS = selfExpression
	return sortByOperator(d, context, expressionNode)
}

// context represents the current matching nodes in the expression pipeline
// expressionNode is your current expression (sort_by)
func sortByOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	results := list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		if candidate.Kind != SequenceNode {
			return context, fmt.Errorf("node at path [%v] is not an array (it's a %v)", candidate.GetNicePath(), candidate.Tag)
		}

		sortableArray := make(sortableNodeArray, len(candidate.Content))

		for i, originalNode := range candidate.Content {

			compareContext, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(originalNode), expressionNode.RHS)
			if err != nil {
				return Context{}, err
			}

			sortableArray[i] = sortableNode{Node: originalNode, CompareContext: compareContext, dateTimeLayout: context.GetDateTimeLayout()}

		}

		sort.Stable(sortableArray)

		sortedList := candidate.CreateReplacementWithComments(SequenceNode, "!!seq", candidate.Style)

		for _, sortedNode := range sortableArray {
			sortedList.AddChild(sortedNode.Node)
		}
		results.PushBack(sortedList)
	}
	return context.ChildContext(results), nil
}

type sortableNode struct {
	Node           *CandidateNode
	CompareContext Context
	dateTimeLayout string
}

type sortableNodeArray []sortableNode

func (a sortableNodeArray) Len() int      { return len(a) }
func (a sortableNodeArray) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a sortableNodeArray) Less(i, j int) bool {
	lhsContext := a[i].CompareContext
	rhsContext := a[j].CompareContext

	rhsEl := rhsContext.MatchingNodes.Front()
	for lhsEl := lhsContext.MatchingNodes.Front(); lhsEl != nil && rhsEl != nil; lhsEl = lhsEl.Next() {
		lhs := lhsEl.Value.(*CandidateNode)
		rhs := rhsEl.Value.(*CandidateNode)

		result := a.compare(lhs, rhs, a[i].dateTimeLayout)

		if result < 0 {
			return true
		} else if result > 0 {
			return false
		}

		rhsEl = rhsEl.Next()
	}
	return false
}

func (a sortableNodeArray) compare(lhs *CandidateNode, rhs *CandidateNode, dateTimeLayout string) int {
	lhsTag := lhs.Tag
	rhsTag := rhs.Tag

	if !strings.HasPrefix(lhsTag, "!!") {
		// custom tag - we have to have a guess
		lhsTag = lhs.guessTagFromCustomType()
	}

	if !strings.HasPrefix(rhsTag, "!!") {
		// custom tag - we have to have a guess
		rhsTag = rhs.guessTagFromCustomType()
	}

	isDateTime := lhsTag == "!!timestamp" && rhsTag == "!!timestamp"
	layout := dateTimeLayout
	// if the lhs is a string, it might be a timestamp in a custom format.
	if lhsTag == "!!str" && layout != time.RFC3339 {
		_, errLhs := parseDateTime(layout, lhs.Value)
		_, errRhs := parseDateTime(layout, rhs.Value)
		isDateTime = errLhs == nil && errRhs == nil
	}

	if lhsTag == "!!null" && rhsTag != "!!null" {
		return -1
	} else if lhsTag != "!!null" && rhsTag == "!!null" {
		return 1
	} else if lhsTag == "!!bool" && rhsTag != "!!bool" {
		return -1
	} else if lhsTag != "!!bool" && rhsTag == "!!bool" {
		return 1
	} else if lhsTag == "!!bool" && rhsTag == "!!bool" {
		lhsTruthy := isTruthyNode(lhs)

		rhsTruthy := isTruthyNode(rhs)
		if lhsTruthy == rhsTruthy {
			return 0
		} else if lhsTruthy {
			return 1
		}
		return -1
	} else if isDateTime {
		lhsTime, err := parseDateTime(layout, lhs.Value)
		if err != nil {
			log.Warningf("Could not parse time %v with layout %v for sort, sorting by string instead: %w", lhs.Value, layout, err)
			return strings.Compare(lhs.Value, rhs.Value)
		}
		rhsTime, err := parseDateTime(layout, rhs.Value)
		if err != nil {
			log.Warningf("Could not parse time %v with layout %v for sort, sorting by string instead: %w", rhs.Value, layout, err)
			return strings.Compare(lhs.Value, rhs.Value)
		}
		if lhsTime.Equal(rhsTime) {
			return 0
		} else if lhsTime.Before(rhsTime) {
			return -1
		}

		return 1
	} else if lhsTag == "!!int" && rhsTag == "!!int" {
		_, lhsNum, err := parseInt64(lhs.Value)
		if err != nil {
			panic(err)
		}
		_, rhsNum, err := parseInt64(rhs.Value)
		if err != nil {
			panic(err)
		}
		return int(lhsNum - rhsNum)
	} else if (lhsTag == "!!int" || lhsTag == "!!float") && (rhsTag == "!!int" || rhsTag == "!!float") {
		lhsNum, err := strconv.ParseFloat(lhs.Value, 64)
		if err != nil {
			panic(err)
		}
		rhsNum, err := strconv.ParseFloat(rhs.Value, 64)
		if err != nil {
			panic(err)
		}
		if lhsNum == rhsNum {
			return 0
		} else if lhsNum < rhsNum {
			return -1
		}

		return 1
	}

	return strings.Compare(lhs.Value, rhs.Value)
}
