package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

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
	encoder := NewYamlEncoder(2, true, ConfiguredYamlPreferences)
	decoder := NewYamlDecoder(ConfiguredYamlPreferences)

	result, err := NewStringEvaluator().Evaluate(expression, input, encoder, decoder)
	if err != nil {
		t.Error(err)
	}

	test.AssertResult(t, expected_output, result)
}
