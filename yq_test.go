package main

import (
	"fmt"
	"testing"
)

var parseValueTests = []struct {
	argument        string
	expectedResult  interface{}
	testDescription string
}{
	{"true", true, "boolean"},
	{"\"true\"", "true", "boolean as string"},
	{"3.4", 3.4, "number"},
	{"\"3.4\"", "3.4", "number as string"},
	{"", "", "empty string"},
}

func TestParseValue(t *testing.T) {
	for _, tt := range parseValueTests {
		assertResultWithContext(t, tt.expectedResult, parseValue(tt.argument), tt.testDescription)
	}
}

func TestRead(t *testing.T) {
	result, _ := read([]string{"examples/sample.yaml", "b.c"})
	assertResult(t, 2, result)
}

func TestReadArray(t *testing.T) {
	result, _ := read([]string{"examples/sample_array.yaml", "[1]"})
	assertResult(t, 2, result)
}

func TestReadString(t *testing.T) {
	result, _ := read([]string{"examples/sample_text.yaml"})
	assertResult(t, "hi", result)
}

func TestReadMultipleDocuments(t *testing.T) {
	t.Run("multiple documents assumes first", func(t *testing.T) {
		result, _ := read([]string{"examples/multidocument.yaml", "b.c"})
		assertResult(t, 1, result)
	})
	t.Run("multiple documents with @n path reads document n", func(t *testing.T) {
		result, _ := read([]string{"examples/multidocument.yaml", "@0"})
		actual, err := yamlToString(result)
		assertNilErr(t, err)
		assertResult(t, `a: Document One
b:
  c: 1`, actual)
	})
	t.Run("multiple documents with leading @n path reads path from document n", func(t *testing.T) {
		result, _ := read([]string{"examples/multidocument.yaml", "@0.b.c"})
		assertResult(t, 1, result)

		result, _ = read([]string{"examples/multidocument.yaml", "@1.b.c"})
		assertResult(t, 2, result)

		_, err := read([]string{"examples/multidocument.yaml", "@2.b.c"})
		assertAnyErr(t, err)
	})
}

func TestOrder(t *testing.T) {
	result, _ := read([]string{"examples/order.yaml"})
	formattedResult, _ := yamlToString(result)
	assertResult(t,
		`version: 3
application: MyApp`,
		formattedResult)
}

func TestMultilineString(t *testing.T) {
	testString := `
	abcd
	efg`
	formattedResult, _ := yamlToString(testString)
	assertResult(t, testString, formattedResult)
}

func TestNewYaml(t *testing.T) {
	result, _ := newYaml([]string{"b.c", "3"})
	formattedResult := fmt.Sprintf("%v", result)
	assertResult(t,
		"[{b [{c 3}]}]",
		formattedResult)
}

func TestNewYamlArray(t *testing.T) {
	result, _ := newYaml([]string{"[0].cat", "meow"})
	formattedResult := fmt.Sprintf("%v", result)
	assertResult(t,
		"[[{cat meow}]]",
		formattedResult)
}

func TestUpdateYaml(t *testing.T) {
	result, _ := updateYaml([]string{"examples/sample.yaml", "b.c", "3"})
	formattedResult := fmt.Sprintf("%v", result)
	assertResult(t,
		"[{a Easy! as one two three} {b [{c 3} {d [3 4]} {e [[{name fred} {value 3}] [{name sam} {value 4}]]}]}]",
		formattedResult)
}

func TestUpdateYamlArray(t *testing.T) {
	result, _ := updateYaml([]string{"examples/sample_array.yaml", "[0]", "3"})
	formattedResult := fmt.Sprintf("%v", result)
	assertResult(t,
		"[3 2 3]",
		formattedResult)
}

func TestUpdateYaml_WithScript(t *testing.T) {
	writeScript = "examples/instruction_sample.yaml"
	_, _ = updateYaml([]string{"examples/sample.yaml"})
}

func TestUpdateYaml_WithUnknownScript(t *testing.T) {
	writeScript = "fake-unknown"
	_, err := updateYaml([]string{"examples/sample.yaml"})
	if err == nil {
		t.Error("Expected error due to unknown file")
	}
	expectedOutput := `open fake-unknown: no such file or directory`
	assertResult(t, expectedOutput, err.Error())
}

func TestNewYaml_WithScript(t *testing.T) {
	writeScript = "examples/instruction_sample.yaml"
	expectedResult := `b:
  c: cat
  e:
  - name: Mike Farah`
	result, _ := newYaml([]string{""})
	actualResult, _ := yamlToString(result)
	assertResult(t, expectedResult, actualResult)
}

func TestNewYaml_WithUnknownScript(t *testing.T) {
	writeScript = "fake-unknown"
	_, err := newYaml([]string{""})
	if err == nil {
		t.Error("Expected error due to unknown file")
	}
	expectedOutput := `open fake-unknown: no such file or directory`
	assertResult(t, expectedOutput, err.Error())
}

func TestSplitData(t *testing.T) {
	t.Run("document index zero", func(t *testing.T) {
		actual, err := readSubDocument("examples/multidocument.yaml", 0)
		assertNilErr(t, err)
		expected := `a: Document One
b:
  c: 1`
		assertResult(t, expected, actual)
		// assert.Equal(t, expected, actual) // needed this to find whitespace problem
	})

	t.Run("document index one", func(t *testing.T) {
		actual, err := readSubDocument("examples/multidocument.yaml", 1)
		assertNilErr(t, err)
		expected := `a: Document Two
b:
  c: 2`
		assertResult(t, expected, actual)
	})

	t.Run("document index beyond available", func(t *testing.T) {
		actual, err := readSubDocument("examples/multidocument.yaml", 2)
		assertAnyErr(t, err)
		assertResult(t, "", actual)
	})
}
