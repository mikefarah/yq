package yqlib

import (
	"container/list"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

func getSliceNumber(d *dataTreeNavigator, context Context, node *CandidateNode, expressionNode *ExpressionNode) (int, error) {
	result, err := d.GetMatchingNodes(context.SingleChildContext(node), expressionNode)
	if err != nil {
		return 0, err
	}
	if result.MatchingNodes.Len() != 1 {
		return 0, fmt.Errorf("expected to find 1 number, got %v instead", result.MatchingNodes.Len())
	}
	return parseInt(result.MatchingNodes.Front().Value.(*CandidateNode).Node.Value)
}

func sliceArrayOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debug("slice array operator!")
	log.Debug("lhs: %v", expressionNode.LHS.Operation.toString())
	log.Debug("rhs: %v", expressionNode.RHS.Operation.toString())

	results := list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		lhsNode := el.Value.(*CandidateNode)
		original := unwrapDoc(lhsNode.Node)

		firstNumber, err := getSliceNumber(d, context, lhsNode, expressionNode.LHS)

		if err != nil {
			return Context{}, err
		}
		relativeFirstNumber := firstNumber
		if relativeFirstNumber < 0 {
			relativeFirstNumber = len(original.Content) + firstNumber
		}

		secondNumber, err := getSliceNumber(d, context, lhsNode, expressionNode.RHS)
		if err != nil {
			return Context{}, err
		}

		relativeSecondNumber := secondNumber
		if relativeSecondNumber < 0 {
			relativeSecondNumber = len(original.Content) + secondNumber
		} else if relativeSecondNumber > len(original.Content) {
			relativeSecondNumber = len(original.Content)
		}

		log.Debug("calculateIndicesToTraverse: slice from %v to %v", relativeFirstNumber, relativeSecondNumber)

		var newResults []*yaml.Node
		for i := relativeFirstNumber; i < relativeSecondNumber; i++ {
			newResults = append(newResults, original.Content[i])
		}

		slicedArrayNode := &yaml.Node{
			Kind:    yaml.SequenceNode,
			Tag:     original.Tag,
			Content: newResults,
		}
		results.PushBack(lhsNode.CreateReplacement(slicedArrayNode))

	}

	// result is now the context that has the nodes we need to put back into a sequence.
	//what about multiple arrays in the context? I think we need to create an array for each one
	return context.ChildContext(results), nil
}
