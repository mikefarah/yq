package yqlib

import (
	"bufio"
	"bytes"
	"container/list"
	"strings"

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

/* takes a string and decodes it back into an object */
func decodeOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	var results = list.New()
	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		var dataBucket yaml.Node
		log.Debugf("got: [%v]", candidate.Node.Value)
		decoder := yaml.NewDecoder(strings.NewReader(unwrapDoc(candidate.Node).Value))
		errorReading := decoder.Decode(&dataBucket)
		if errorReading != nil {
			return Context{}, errorReading
		}
		//first node is a doc
		node := unwrapDoc(&dataBucket)

		results.PushBack(candidate.CreateChild(nil, node))
	}
	return context.ChildContext(results), nil
}
