package yqlib

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

type formatScenario struct {
	input          string
	indent         int
	expression     string
	expected       string
	description    string
	subdescription string
	skipDoc        bool
	scenarioType   string
	expectedError  string
}

func processFormatScenario(s formatScenario, decoder Decoder, encoder Encoder) (string, error) {
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	if decoder == nil {
		decoder = NewYamlDecoder(ConfiguredYamlPreferences)
	}

	log.Debugf("reading docs")
	inputs, err := readDocuments(strings.NewReader(s.input), "sample.yml", 0, decoder)
	if err != nil {
		return "", err
	}

	log.Debugf("done reading the documents")

	expression := s.expression
	if expression == "" {
		expression = "."
	}

	exp, err := getExpressionParser().ParseExpression(expression)
	if err != nil {
		return "", err
	}

	context, err := NewDataTreeNavigator().GetMatchingNodes(Context{MatchingNodes: inputs}, exp)

	log.Debugf("Going to print: %v", NodesToString(context.MatchingNodes))

	if err != nil {
		return "", err
	}

	printer := NewPrinter(encoder, NewSinglePrinterWriter(writer))
	err = printer.PrintResults(context.MatchingNodes)
	if err != nil {
		return "", err
	}
	writer.Flush()

	return output.String(), nil
}

func mustProcessFormatScenario(s formatScenario, decoder Decoder, encoder Encoder) string {
	result, err := processFormatScenario(s, decoder, encoder)
	if err != nil {
		log.Error("Bad scenario %v: %w", s.description, err)
		return fmt.Sprintf("Bad scenario %v: %v", s.description, err.Error())
	}
	return result
}
