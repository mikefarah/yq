package yqlib

import (
	"container/list"
	"fmt"
)

func getSliceNumber(d *dataTreeNavigator, context Context, node *CandidateNode, expressionNode *ExpressionNode) (int, error) {
	result, err := d.GetMatchingNodes(context.SingleChildContext(node), expressionNode)
	if err != nil {
		return 0, err
	}
	if result.MatchingNodes.Len() != 1 {
		return 0, fmt.Errorf("expected to find 1 number, got %v instead", result.MatchingNodes.Len())
	}
	return parseInt(result.MatchingNodes.Front().Value.(*CandidateNode).Value)
}

// clampSliceIndex resolves a possibly-negative slice index against
// length and clamps the result to [0, length].
func clampSliceIndex(index, length int) int {
	if index < 0 {
		index += length
	}
	if index < 0 {
		return 0
	}
	if index > length {
		return length
	}
	return index
}

func sliceStringNode(lhsNode *CandidateNode, firstNumber int, secondNumber int) *CandidateNode {
	runes := []rune(lhsNode.Value)
	length := len(runes)

	relativeFirstNumber := clampSliceIndex(firstNumber, length)
	relativeSecondNumber := clampSliceIndex(secondNumber, length)
	if relativeSecondNumber < relativeFirstNumber {
		relativeSecondNumber = relativeFirstNumber
	}

	log.Debugf("sliceStringNode: slice from %v to %v", relativeFirstNumber, relativeSecondNumber)

	slicedString := string(runes[relativeFirstNumber:relativeSecondNumber])
	replacement := lhsNode.CreateReplacement(ScalarNode, lhsNode.Tag, slicedString)
	replacement.Style = lhsNode.Style
	return replacement
}

func sliceArrayOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debug("slice array operator!")
	log.Debugf("lhs: %v", expressionNode.LHS.Operation.toString())
	log.Debugf("rhs: %v", expressionNode.RHS.Operation.toString())

	results := list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		lhsNode := el.Value.(*CandidateNode)

		firstNumber, err := getSliceNumber(d, context, lhsNode, expressionNode.LHS)
		if err != nil {
			return Context{}, err
		}

		secondNumber, err := getSliceNumber(d, context, lhsNode, expressionNode.RHS)
		if err != nil {
			return Context{}, err
		}

		if lhsNode.Kind == ScalarNode && lhsNode.guessTagFromCustomType() == "!!str" {
			results.PushBack(sliceStringNode(lhsNode, firstNumber, secondNumber))
			continue
		}

		relativeFirstNumber := clampSliceIndex(firstNumber, len(lhsNode.Content))
		relativeSecondNumber := clampSliceIndex(secondNumber, len(lhsNode.Content))

		log.Debugf("calculateIndicesToTraverse: slice from %v to %v", relativeFirstNumber, relativeSecondNumber)

		var newResults []*CandidateNode
		for i := relativeFirstNumber; i < relativeSecondNumber; i++ {
			newResults = append(newResults, lhsNode.Content[i])
		}

		sliceArrayNode := lhsNode.CreateReplacement(SequenceNode, lhsNode.Tag, "")
		sliceArrayNode.AddChildren(newResults)
		results.PushBack(sliceArrayNode)

	}

	// result is now the context that has the nodes we need to put back into a sequence.
	//what about multiple arrays in the context? I think we need to create an array for each one
	return context.ChildContext(results), nil
}
