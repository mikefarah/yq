package yqlib

import (
	"container/list"

	yaml "gopkg.in/yaml.v3"
)

type sliceArrayPreferences struct {
	firstNumber         int
	secondNumber        int
	secondNumberDefined bool
}

func sliceArrayOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	lhs, err := d.GetMatchingNodes(context, expressionNode.LHS)
	if err != nil {
		return Context{}, err
	}
	prefs := expressionNode.Operation.Preferences.(sliceArrayPreferences)
	firstNumber := prefs.firstNumber
	secondNumber := prefs.secondNumber

	results := list.New()

	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		lhsNode := el.Value.(*CandidateNode)
		original := unwrapDoc(lhsNode.Node)

		relativeFirstNumber := firstNumber
		if relativeFirstNumber < 0 {
			relativeFirstNumber = len(original.Content) + firstNumber
		}

		relativeSecondNumber := len(original.Content)
		if prefs.secondNumberDefined {
			relativeSecondNumber = secondNumber
			if relativeSecondNumber < 0 {
				relativeSecondNumber = len(original.Content) + secondNumber
			}
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
