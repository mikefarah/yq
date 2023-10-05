package yqlib

import (
	"container/list"
	"fmt"
	"strconv"

	yaml "gopkg.in/yaml.v3"
)

func tryConvertToNumber(value string) (string, bool) {
	// try a int first
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

func toNumberOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("ToNumberOperator")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		if candidate.Node.Kind != yaml.ScalarNode {
			return Context{}, fmt.Errorf("cannot convert node at path %v of tag %v to number", candidate.GetNicePath(), candidate.GetNiceTag())
		}

		if candidate.Node.Tag == "!!int" || candidate.Node.Tag == "!!float" {
			// it already is a number!
			results.PushBack(candidate)
		} else {
			tag, converted := tryConvertToNumber(candidate.Node.Value)
			if converted {
				node := &yaml.Node{Kind: yaml.ScalarNode, Value: candidate.Node.Value, Tag: tag}

				result := candidate.CreateReplacement(node)
				results.PushBack(result)
			} else {
				return Context{}, fmt.Errorf("cannot convert node value [%v] at path %v of tag %v to number", candidate.Node.Value, candidate.GetNicePath(), candidate.GetNiceTag())
			}

		}
	}

	return context.ChildContext(results), nil
}
