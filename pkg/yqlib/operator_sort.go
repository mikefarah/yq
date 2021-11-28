package yqlib

import (
	"container/list"
	"fmt"
	"sort"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// context represents the current matching nodes in the expression pipeline
//expressionNode is your current expression (sort_by)
func sortByOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	results := list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		candidateNode := unwrapDoc(candidate.Node)

		if candidateNode.Kind != yaml.SequenceNode {
			return context, fmt.Errorf("%v is not an array", candidate.GetKey())
		}

		sortableArray := make(sortableNodeArray, len(candidateNode.Content))

		for i, originalNode := range candidateNode.Content {

			childCandidate := candidate.CreateChildInArray(i, originalNode)
			compareContext, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(childCandidate), expressionNode.Rhs)
			if err != nil {
				return Context{}, err
			}

			nodeToCompare := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null"}
			if compareContext.MatchingNodes.Len() > 0 {
				nodeToCompare = compareContext.MatchingNodes.Front().Value.(*CandidateNode).Node
			}

			log.Debug("going to compare %v by %v", NodeToString(candidate.CreateReplacement(originalNode)), NodeToString(candidate.CreateReplacement(nodeToCompare)))

			sortableArray[i] = sortableNode{Node: originalNode, NodeToCompare: nodeToCompare}

		}

		sort.Sort(sortableArray)

		sortedList := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Style: candidateNode.Style}
		sortedList.Content = make([]*yaml.Node, len(candidateNode.Content))

		for i, sortedNode := range sortableArray {
			sortedList.Content[i] = sortedNode.Node
		}
		results.PushBack(candidate.CreateReplacement(sortedList))
	}
	return context.ChildContext(results), nil
}

type sortableNode struct {
	Node          *yaml.Node
	NodeToCompare *yaml.Node
}

type sortableNodeArray []sortableNode

func (a sortableNodeArray) Len() int      { return len(a) }
func (a sortableNodeArray) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a sortableNodeArray) Less(i, j int) bool {
	lhs := a[i].NodeToCompare
	rhs := a[j].NodeToCompare

	if lhs.Tag != rhs.Tag || lhs.Tag == "!!str" {
		return strings.Compare(lhs.Value, rhs.Value) < 0
	} else if lhs.Tag == "!!int" && rhs.Tag == "!!int" {
		_, lhsNum, err := parseInt(lhs.Value)
		if err != nil {
			panic(err)
		}
		_, rhsNum, err := parseInt(rhs.Value)
		if err != nil {
			panic(err)
		}
		return lhsNum < rhsNum
	} else if (lhs.Tag == "!!int" || lhs.Tag == "!!float") && (rhs.Tag == "!!int" || rhs.Tag == "!!float") {
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

	return true
}
