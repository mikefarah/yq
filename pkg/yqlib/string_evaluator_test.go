package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

func TestStringEvaluator_MultipleDocumentMerge(t *testing.T) {
	yamlString := "a: Hello\n---\na: Goodbye\n"
	expected_output := "a: Goodbye\n"

	encoder := NewYamlEncoder(ConfiguredYamlPreferences)
	decoder := NewYamlDecoder(ConfiguredYamlPreferences)
	result, err := NewStringEvaluator().EvaluateAll("select(di==0) * select(di==1)", yamlString, encoder, decoder)
	if err != nil {
		t.Error(err)
	} else {
		test.AssertResult(t, expected_output, result)
	}
}

func TestStringEvaluator_Evaluate_Nominal(t *testing.T) {
	expected_output := `` +
		`yq` + "\n" +
		`---` + "\n" +
		`jq` + "\n"
	expression := ".[].name"
	input := `` +
		` - name: yq` + "\n" +
		`   description: yq is a portable command-line YAML, JSON and XML processor` + "\n" +
		`---` + "\n" +
		` - name: jq` + "\n" +
		`   description: Command-line JSON processor` + "\n"
	encoder := NewYamlEncoder(ConfiguredYamlPreferences)
	decoder := NewYamlDecoder(ConfiguredYamlPreferences)

	result, err := NewStringEvaluator().Evaluate(expression, input, encoder, decoder)
	if err != nil {
		t.Error(err)
	} else {
		test.AssertResult(t, expected_output, result)
	}
}
