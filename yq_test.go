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

func TestReadMulti(t *testing.T) {
	docIndex = 1
	result, _ := read([]string{"examples/multiple_docs.yaml", "another.document"})
	assertResult(t, "here", result)
	docIndex = 0
}

func TestReadArray(t *testing.T) {
	result, _ := read([]string{"examples/sample_array.yaml", "[1]"})
	assertResult(t, 2, result)
}

func TestReadString(t *testing.T) {
	result, _ := read([]string{"examples/sample_text.yaml"})
	assertResult(t, "hi", result)
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

func TestDeleteYaml(t *testing.T) {
	result, _ := deleteYaml([]string{"examples/sample.yaml", "b.c"})
	formattedResult := fmt.Sprintf("%v", result)
	assertResult(t,
		"[{a Easy! as one two three} {b [{d [3 4]} {e [[{name fred} {value 3}] [{name sam} {value 4}]]}]}]",
		formattedResult)
}
