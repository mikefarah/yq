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

func TestStringEvaluator_InterpolationAfterFilter(t *testing.T) {
	previous := StringInterpolationEnabled
	t.Cleanup(func() { StringInterpolationEnabled = previous })
	StringInterpolationEnabled = true

	const documents = "name: Alice\n---\nname: Bob\n"
	scenarios := []struct {
		description string
		input       string
		expression  string
		expected    string
		evaluateAll bool
	}{
		{
			description: "first document",
			input:       documents,
			expression:  `select(document_index == 0) | "My name is \(.name)"`,
			expected:    "My name is Alice\n",
		},
		{
			description: "last document",
			input:       documents,
			expression:  `select(document_index == 1) | "My name is \(.name)"`,
			expected:    "My name is Bob\n",
		},
		{
			description: "no documents selected",
			input:       documents,
			expression:  `select(document_index == 2) | "My name is \(.name)"`,
			expected:    "",
		},
		{
			description: "non-existent document index produces no blank lines",
			input:       documents,
			expression:  `select(document_index == 99) | "\(.name)"`,
			expected:    "",
		},
		{
			description: "all documents selected",
			input:       documents,
			expression:  `select(true) | "My name is \(.name)"`,
			expected:    "My name is Alice\nMy name is Bob\n",
		},
		{
			description: "nested string interpolation",
			input:       documents,
			expression:  `select(document_index == 0) | "My name is \(\"Dr \(.name)\")"`,
			expected:    "My name is Dr Alice\n",
		},
		{
			description: "multiple interpolations",
			input:       documents,
			expression:  `select(document_index == 0) | "\(.name) has \(.name | length) letters"`,
			expected:    "Alice has 5 letters\n",
		},
		{
			description: "upcase before interpolation",
			input:       documents,
			expression:  `select(document_index == 0) | .name | upcase | "My name is \(.)"`,
			expected:    "My name is ALICE\n",
		},
		{
			description: "upcase after interpolation",
			input:       documents,
			expression:  `select(document_index == 0) | "My name is \(.name)" | upcase`,
			expected:    "MY NAME IS ALICE\n",
		},
		{
			description: "length after interpolation",
			input:       documents,
			expression:  `select(document_index == 0) | "My name is \(.name)" | length`,
			expected:    "16\n",
		},
		{
			description: "concatenation after interpolation",
			input:       documents,
			expression:  `select(document_index == 0) | ("My name is \(.name)" + "!")`,
			expected:    "My name is Alice!\n",
		},
		{
			description: "plain string after rejected select",
			input:       documents,
			expression:  `select(false) | "hello"`,
			expected:    "",
		},
		{
			description: "explicit empty array filter",
			input:       documents,
			expression:  `[] | .[] | "My name is \(.name)"`,
			expected:    "",
		},
		{
			description: "empty input array filter",
			input:       "[]\n",
			expression:  `.[] | "My name is \(.name)"`,
			expected:    "",
		},
		{
			description: "selected missing name",
			input:       "{}\n",
			expression:  `select(true) | "\(.name)"`,
			expected:    "null\n",
		},
		{
			description: "selected null name",
			input:       "name: null\n",
			expression:  `select(true) | "\(.name)"`,
			expected:    "null\n",
		},
		{
			description: "selected empty name",
			input:       "name: \"\"\n",
			expression:  `select(true) | "\(.name)"`,
			expected:    "\n",
		},
		{
			description: "empty expression within selected interpolation",
			input:       "name: Alice\n",
			expression:  `select(true) | "My name is \(select(false))"`,
			expected:    "My name is \n",
		},
		{
			description: "evaluate all first document",
			input:       documents,
			expression:  `select(document_index == 0) | "My name is \(.name)"`,
			expected:    "My name is Alice\n",
			evaluateAll: true,
		},
		{
			description: "evaluate all last document",
			input:       documents,
			expression:  `select(document_index == 1) | "My name is \(.name)"`,
			expected:    "My name is Bob\n",
			evaluateAll: true,
		},
		{
			description: "evaluate all no documents selected",
			input:       documents,
			expression:  `select(document_index == 2) | "My name is \(.name)"`,
			expected:    "",
			evaluateAll: true,
		},
		{
			description: "evaluate all documents selected",
			input:       documents,
			expression:  `select(true) | "My name is \(.name)"`,
			expected:    "My name is Alice\nMy name is Bob\n",
			evaluateAll: true,
		},
		{
			description: "evaluate all empty input",
			input:       "",
			expression:  `"My name is \(.name)"`,
			expected:    "",
			evaluateAll: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			prefs := NewDefaultYamlPreferences()
			evaluate := NewStringEvaluator().Evaluate
			if scenario.evaluateAll {
				evaluate = NewStringEvaluator().EvaluateAll
			}
			result, err := evaluate(scenario.expression, scenario.input, NewYamlEncoder(prefs), NewYamlDecoder(prefs))
			if err != nil {
				t.Fatal(err)
			}
			test.AssertResult(t, scenario.expected, result)
		})
	}
}

func TestStringEvaluator_DisabledInterpolationAfterFilter(t *testing.T) {
	previous := StringInterpolationEnabled
	t.Cleanup(func() { StringInterpolationEnabled = previous })
	StringInterpolationEnabled = false

	const documents = "name: Alice\n---\nname: Bob\n"
	scenarios := []struct {
		description string
		input       string
		expression  string
		expected    string
		evaluateAll bool
	}{
		{
			description: "selected document keeps literal interpolation",
			input:       documents,
			expression:  `select(document_index == 0) | "My name is \(.name)"`,
			expected:    "My name is \\(.name)\n",
		},
		{
			description: "rejected documents",
			input:       documents,
			expression:  `select(false) | "My name is \(.name)"`,
			expected:    "",
		},
		{
			description: "empty input",
			input:       "",
			expression:  `"My name is \(.name)"`,
			expected:    "",
		},
		{
			description: "evaluate all selected document keeps literal interpolation",
			input:       documents,
			expression:  `select(document_index == 0) | "My name is \(.name)"`,
			expected:    "My name is \\(.name)\n",
			evaluateAll: true,
		},
		{
			description: "evaluate all rejected documents",
			input:       documents,
			expression:  `select(false) | "My name is \(.name)"`,
			expected:    "",
			evaluateAll: true,
		},
		{
			description: "evaluate all empty input",
			input:       "",
			expression:  `"My name is \(.name)"`,
			expected:    "",
			evaluateAll: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			prefs := NewDefaultYamlPreferences()
			evaluate := NewStringEvaluator().Evaluate
			if scenario.evaluateAll {
				evaluate = NewStringEvaluator().EvaluateAll
			}
			result, err := evaluate(scenario.expression, scenario.input, NewYamlEncoder(prefs), NewYamlDecoder(prefs))
			if err != nil {
				t.Fatal(err)
			}
			test.AssertResult(t, scenario.expected, result)
		})
	}
}
