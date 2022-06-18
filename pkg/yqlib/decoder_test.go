package yqlib

import (
	"bufio"
	"bytes"
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
	decoder        Decoder
	encoder        Encoder
}

func processFormatScenario(s formatScenario, decoder Decoder, encoder Encoder) string {

	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	if decoder == nil {
		decoder = NewYamlDecoder()
	}

	inputs, err := readDocuments(strings.NewReader(s.input), "sample.yml", 0, decoder)
	if err != nil {
		panic(err)
	}

	expression := s.expression
	if expression == "" {
		expression = "."
	}

	exp, err := getExpressionParser().ParseExpression(expression)

	if err != nil {
		panic(err)
	}

	context, err := NewDataTreeNavigator().GetMatchingNodes(Context{MatchingNodes: inputs}, exp)

	if err != nil {
		panic(err)
	}

	printer := NewPrinter(encoder, NewSinglePrinterWriter(writer))
	err = printer.PrintResults(context.MatchingNodes)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	return output.String()

}
