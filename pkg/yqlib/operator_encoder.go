package yqlib

import (
	"bufio"
	"bytes"
	"container/list"

	"gopkg.in/yaml.v3"
)

func yamlToString(candidate *CandidateNode, prefs encoderPreferences) (string, error) {
	var output bytes.Buffer
	printer := NewPrinter(bufio.NewWriter(&output), prefs.format, true, false, 2, true)
	elMap := list.New()
	elMap.PushBack(candidate)
	err := printer.PrintResults(elMap)
	return output.String(), err
}

type encoderPreferences struct {
	format PrinterOutputFormat
}

/* encodes object as yaml string */

func encodeOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	preferences := expressionNode.Operation.Preferences.(encoderPreferences)
	var results = list.New()
	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		stringValue, err := yamlToString(candidate, preferences)
		if err != nil {
			return Context{}, err
		}

		stringContentNode := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: stringValue}
		results.PushBack(candidate.CreateChild(nil, stringContentNode))
	}
	return context.ChildContext(results), nil
}
