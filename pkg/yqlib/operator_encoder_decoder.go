package yqlib

import (
	"bufio"
	"bytes"
	"container/list"
	"errors"
	"regexp"
	"strings"
)

func configureEncoder(format *Format, indent int) Encoder {

	switch format {
	case JSONFormat:
		prefs := ConfiguredJSONPreferences.Copy()
		prefs.Indent = indent
		prefs.ColorsEnabled = false
		prefs.UnwrapScalar = false
		return NewJSONEncoder(prefs)
	case YamlFormat:
		var prefs = ConfiguredYamlPreferences.Copy()
		prefs.Indent = indent
		prefs.ColorsEnabled = false
		return NewYamlEncoder(prefs)
	case XMLFormat:
		var xmlPrefs = ConfiguredXMLPreferences.Copy()
		xmlPrefs.Indent = indent
		return NewXMLEncoder(xmlPrefs)
	}
	return format.EncoderFactory()
}

func encodeToString(candidate *CandidateNode, prefs encoderPreferences) (string, error) {
	var output bytes.Buffer
	log.Debug("printing with indent: %v", prefs.indent)

	encoder := configureEncoder(prefs.format, prefs.indent)
	if encoder == nil {
		return "", errors.New("no support for output format")
	}

	printer := NewPrinter(encoder, NewSinglePrinterWriter(bufio.NewWriter(&output)))
	err := printer.PrintResults(candidate.AsList())
	return output.String(), err
}

type encoderPreferences struct {
	format *Format
	indent int
}

/* encodes object as yaml string */
var chomper = regexp.MustCompile("\n+$")

func encodeOperator(_ *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	preferences := expressionNode.Operation.Preferences.(encoderPreferences)
	var results = list.New()

	hasOnlyOneNewLine := regexp.MustCompile("[^\n].*\n$")
	endWithNewLine := regexp.MustCompile(".*\n$")

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		stringValue, err := encodeToString(candidate, preferences)

		if err != nil {
			return Context{}, err
		}

		// remove trailing newlines if needed.
		// check if we originally decoded this path, and the original thing had a single line.
		originalList := context.GetVariable("decoded: " + candidate.GetKey())
		if originalList != nil && originalList.Len() > 0 && hasOnlyOneNewLine.MatchString(stringValue) {

			original := originalList.Front().Value.(*CandidateNode)
			// original block did not have a newline at the end, get rid of this one too
			if !endWithNewLine.MatchString(original.Value) {
				stringValue = chomper.ReplaceAllString(stringValue, "")
			}
		}

		// dont print a newline when printing json on a single line.
		if (preferences.format == JSONFormat && preferences.indent == 0) ||
			preferences.format == CSVFormat ||
			preferences.format == TSVFormat {
			stringValue = chomper.ReplaceAllString(stringValue, "")
		}

		results.PushBack(candidate.CreateReplacement(ScalarNode, "!!str", stringValue))
	}
	return context.ChildContext(results), nil
}

type decoderPreferences struct {
	format *Format
}

/* takes a string and decodes it back into an object */
func decodeOperator(_ *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	preferences := expressionNode.Operation.Preferences.(decoderPreferences)

	decoder := preferences.format.DecoderFactory()
	if decoder == nil {
		return Context{}, errors.New("no support for input format")
	}

	var results = list.New()
	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		context.SetVariable("decoded: "+candidate.GetKey(), candidate.AsList())

		log.Debugf("got: [%v]", candidate.Value)

		err := decoder.Init(strings.NewReader(candidate.Value))
		if err != nil {
			return Context{}, err
		}

		node, errorReading := decoder.Decode()
		if errorReading != nil {
			return Context{}, errorReading
		}
		node.Key = candidate.Key
		node.Parent = candidate.Parent

		results.PushBack(node)
	}
	return context.ChildContext(results), nil
}
