package yqlib

import (
	"container/list"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v3"
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

		candidateNode := unwrapDoc(candidate.Node)

		if candidateNode.Kind != yaml.SequenceNode {
			return context, fmt.Errorf("node at path [%v] is not an array (it's a %v)", candidate.GetNicePath(), candidate.GetNiceTag())
		}

		sortableArray := make(sortableNodeArray, len(candidateNode.Content))

		for i, originalNode := range candidateNode.Content {

			childCandidate := candidate.CreateChildInArray(i, originalNode)
			compareContext, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(childCandidate), expressionNode.RHS)
			if err != nil {
				return Context{}, err
			}

			nodeToCompare := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null"}
			if compareContext.MatchingNodes.Len() > 0 {
				nodeToCompare = compareContext.MatchingNodes.Front().Value.(*CandidateNode).Node
			}

			log.Debug("going to compare %v by %v", NodeToString(candidate.CreateReplacement(originalNode)), NodeToString(candidate.CreateReplacement(nodeToCompare)))

			sortableArray[i] = sortableNode{Node: originalNode, NodeToCompare: nodeToCompare, dateTimeLayout: context.GetDateTimeLayout()}

			if nodeToCompare.Kind != yaml.ScalarNode {
				return Context{}, fmt.Errorf("sort only works for scalars, got %v", nodeToCompare.Tag)
			}

		}

		sort.Stable(sortableArray)

		sortedList := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Style: candidateNode.Style}
		sortedList.Content = make([]*yaml.Node, len(candidateNode.Content))

		for i, sortedNode := range sortableArray {
			sortedList.Content[i] = sortedNode.Node
		}
		results.PushBack(candidate.CreateReplacementWithDocWrappers(sortedList))
	}
	return context.ChildContext(results), nil
}

type sortableNode struct {
	Node           *yaml.Node
	NodeToCompare  *yaml.Node
	dateTimeLayout string
}

type sortableNodeArray []sortableNode

func (a sortableNodeArray) Len() int      { return len(a) }
func (a sortableNodeArray) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a sortableNodeArray) Less(i, j int) bool {
	lhs := a[i].NodeToCompare
	rhs := a[j].NodeToCompare

	lhsTag := lhs.Tag
	rhsTag := rhs.Tag

	if !strings.HasPrefix(lhsTag, "!!") {
		// custom tag - we have to have a guess
		lhsTag = guessTagFromCustomType(lhs)
	}

	if !strings.HasPrefix(rhsTag, "!!") {
		// custom tag - we have to have a guess
		rhsTag = guessTagFromCustomType(rhs)
	}

	isDateTime := lhsTag == "!!timestamp" && rhsTag == "!!timestamp"
	layout := a[i].dateTimeLayout
	// if the lhs is a string, it might be a timestamp in a custom format.
	if lhsTag == "!!str" && layout != time.RFC3339 {
		_, errLhs := parseDateTime(layout, lhs.Value)
		_, errRhs := parseDateTime(layout, rhs.Value)
		isDateTime = errLhs == nil && errRhs == nil
	}

	if lhsTag == "!!null" && rhsTag != "!!null" {
		return true
	} else if lhsTag != "!!null" && rhsTag == "!!null" {
		return false
	} else if lhsTag == "!!bool" && rhsTag != "!!bool" {
		return true
	} else if lhsTag != "!!bool" && rhsTag == "!!bool" {
		return false
	} else if lhsTag == "!!bool" && rhsTag == "!!bool" {
		lhsTruthy, err := isTruthyNode(lhs)
		if err != nil {
			panic(fmt.Errorf("could not parse %v as boolean: %w", lhs.Value, err))
		}

		rhsTruthy, err := isTruthyNode(rhs)
		if err != nil {
			panic(fmt.Errorf("could not parse %v as boolean: %w", rhs.Value, err))
		}

		return !lhsTruthy && rhsTruthy
	} else if isDateTime {
		lhsTime, err := parseDateTime(layout, lhs.Value)
		if err != nil {
			log.Warningf("Could not parse time %v with layout %v for sort, sorting by string instead: %w", lhs.Value, layout, err)
			return strings.Compare(lhs.Value, rhs.Value) < 0
		}
		rhsTime, err := parseDateTime(layout, rhs.Value)
		if err != nil {
			log.Warningf("Could not parse time %v with layout %v for sort, sorting by string instead: %w", rhs.Value, layout, err)
			return strings.Compare(lhs.Value, rhs.Value) < 0
		}
		return lhsTime.Before(rhsTime)
	} else if lhsTag == "!!int" && rhsTag == "!!int" {
		_, lhsNum, err := parseInt64(lhs.Value)
		if err != nil {
			panic(err)
		}
		_, rhsNum, err := parseInt64(rhs.Value)
		if err != nil {
			panic(err)
		}
		return lhsNum < rhsNum
	} else if (lhsTag == "!!int" || lhsTag == "!!float") && (rhsTag == "!!int" || rhsTag == "!!float") {
		lhsNum, err := strconv.ParseFloat(lhs.Value, 64)
		if err != nil {
			panic(err)
		}
		rhsNum, err := strconv.ParseFloat(rhs.Value, 64)
		if err != nil {
			panic(err)
		}
		return lhsNum < rhsNum
	}

	return strings.Compare(lhs.Value, rhs.Value) < 0
}
