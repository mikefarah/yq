package yqlib

import (
	"bufio"
	"bytes"
	"container/list"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

func yamlToString(candidate *CandidateNode, prefs encoderPreferences) (string, error) {
	var output bytes.Buffer
	printer := NewPrinter(bufio.NewWriter(&output), prefs.format, true, false, 2, true)
	err := printer.PrintResults(candidate.AsList())
	return output.String(), err
}

type encoderPreferences struct {
	format PrinterOutputFormat
}

/* encodes object as yaml string */

func encodeOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	preferences := expressionNode.Operation.Preferences.(encoderPreferences)
	var results = list.New()

	hasOnlyOneNewLine := regexp.MustCompile("[^\n].*\n$")
	endWithNewLine := regexp.MustCompile(".*\n$")
	chomper := regexp.MustCompile("\n+$")

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		stringValue, err := yamlToString(candidate, preferences)

		if err != nil {
			return Context{}, err
		}

		// remove trailing newlines if needed.
		// check if we originally decoded this path, and the original thing had a single line.
		originalList := context.GetVariable("decoded: " + candidate.GetKey())
		if originalList != nil && originalList.Len() > 0 && hasOnlyOneNewLine.MatchString(stringValue) {

			original := originalList.Front().Value.(*CandidateNode)
			originalNode := unwrapDoc(original.Node)
			// original block did not have a new line at the end, get rid of this one too
			if !endWithNewLine.MatchString(originalNode.Value) {
				stringValue = chomper.ReplaceAllString(stringValue, "")
			}
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

		context.SetVariable("decoded: "+candidate.GetKey(), candidate.AsList())

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
