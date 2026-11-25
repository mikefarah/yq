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

		sortableArray := &sortableNodeArray{}

		if candidate.CanVisitValues() {
			visitor := func(valueNode *CandidateNode) error {
				compareContext, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(valueNode), expressionNode.RHS)
				if err != nil {
					return err
				}
				sortableNode := sortableNode{Node: valueNode, CompareContext: compareContext, dateTimeLayout: context.GetDateTimeLayout()}
				sortableArray.nodes = append(sortableArray.nodes, sortableNode)
				return nil
			}
			if err := candidate.VisitValues(visitor); err != nil {
				return context, err
			}
		} else {
			return context, fmt.Errorf("node at path [%v] is not an array or map (it's a %v)", candidate.GetNicePath(), candidate.Tag)
		}

		sort.Stable(sortableArray)
		if sortableArray.err != nil {
			return context, sortableArray.err
		}

		sortedList := candidate.CopyWithoutContent()
		switch candidate.Kind {
		case MappingNode:
			for _, sortedNode := range sortableArray.nodes {
				sortedList.AddKeyValueChild(sortedNode.Node.Key, sortedNode.Node)
			}
		case SequenceNode:
			for _, sortedNode := range sortableArray.nodes {
				sortedList.AddChild(sortedNode.Node)
			}
		}

		// convert array of value nodes back to map
		results.PushBack(sortedList)
	}
	return context.ChildContext(results), nil
}

type sortableNode struct {
	Node           *CandidateNode
	CompareContext Context
	dateTimeLayout string
}

type sortableNodeArray struct {
	nodes []sortableNode
	err   error
}

func (a *sortableNodeArray) Len() int      { return len(a.nodes) }
func (a *sortableNodeArray) Swap(i, j int) { a.nodes[i], a.nodes[j] = a.nodes[j], a.nodes[i] }

func (a *sortableNodeArray) Less(i, j int) bool {
	lhsContext := a.nodes[i].CompareContext
	rhsContext := a.nodes[j].CompareContext

	rhsEl := rhsContext.MatchingNodes.Front()

	for lhsEl := lhsContext.MatchingNodes.Front(); lhsEl != nil && rhsEl != nil; lhsEl = lhsEl.Next() {
		lhs := lhsEl.Value.(*CandidateNode)
		rhs := rhsEl.Value.(*CandidateNode)

		result := a.compare(lhs, rhs, a.nodes[i].dateTimeLayout)

		if result < 0 {
			return true
		} else if result > 0 {
			return false
		}

		rhsEl = rhsEl.Next()
	}
	return lhsContext.MatchingNodes.Len() < rhsContext.MatchingNodes.Len()
}

func (a *sortableNodeArray) compare(lhs *CandidateNode, rhs *CandidateNode, dateTimeLayout string) int {
	if a.err != nil {
		return 0
	}

	if lhs.Kind == MappingNode || rhs.Kind == MappingNode {
		a.err = fmt.Errorf("maps not yet supported for comparison")
		return 0
	}
	if lhs.Kind == SequenceNode || rhs.Kind == SequenceNode {
		a.err = fmt.Errorf("arrays not yet supported for comparison")
		return 0
	}

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
			log.Warningf("Could not parse time %v with layout %v for sort, sorting by string instead: %v", lhs.Value, layout, err)
			return strings.Compare(lhs.Value, rhs.Value)
		}
		rhsTime, err := parseDateTime(layout, rhs.Value)
		if err != nil {
			log.Warningf("Could not parse time %v with layout %v for sort, sorting by string instead: %v", rhs.Value, layout, err)
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
		if lhsNum < rhsNum {
			return -1
		} else if lhsNum > rhsNum {
			return 1
		}
		return 0
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
