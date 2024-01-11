package yqlib

import (
	"container/list"
	"fmt"
	"strconv"
)

func tryConvertToNumber(value string) (string, bool) {
	// try an int first
	_, _, err := parseInt64(value)
	if err == nil {
		return "!!int", true
	}
	// try float
	_, floatErr := strconv.ParseFloat(value, 64)

	if floatErr == nil {
		return "!!float", true
	}
	return "", false

}

func toNumberOperator(_ *dataTreeNavigator, context Context, _ *ExpressionNode) (Context, error) {
	log.Debugf("ToNumberOperator")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		if candidate.Kind != ScalarNode {
			return Context{}, fmt.Errorf("cannot convert node at path %v of tag %v to number", candidate.GetNicePath(), candidate.Tag)
		}

		if candidate.Tag == "!!int" || candidate.Tag == "!!float" {
			// it already is a number!
			results.PushBack(candidate)
		} else {
			tag, converted := tryConvertToNumber(candidate.Value)
			if converted {
				result := candidate.CreateReplacement(ScalarNode, tag, candidate.Value)
				results.PushBack(result)
			} else {
				return Context{}, fmt.Errorf("cannot convert node value [%v] at path %v of tag %v to number", candidate.Value, candidate.GetNicePath(), candidate.Tag)
			}

		}
	}

	return context.ChildContext(results), nil
}
