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
}

func TestParseValue(t *testing.T) {
	for _, tt := range parseValueTests {
		assertResultWithContext(t, tt.expectedResult, parseValue(tt.argument), tt.testDescription)
	}
}

func TestRead(t *testing.T) {
	result := read([]string{"sample.yaml", "b.c"})
	assertResult(t, 2, result)
}

func TestOrder(t *testing.T) {
	result := read([]string{"order.yaml"})
	formattedResult := yamlToString(result)
	assertResult(t,
		`version: 3
application: MyApp`,
		formattedResult)
}

func TestNewYaml(t *testing.T) {
	result := newYaml([]string{"b.c", "3"})
	formattedResult := fmt.Sprintf("%v", result)
	assertResult(t,
		"[{b [{c 3}]}]",
		formattedResult)
}

func TestUpdateYaml(t *testing.T) {
	result := updateYaml([]string{"sample.yaml", "b.c", "3"})
	formattedResult := fmt.Sprintf("%v", result)
	assertResult(t,
		"[{a Easy! as one two three} {b [{c 3} {d [3 4]} {e [[{name fred} {value 3}] [{name sam} {value 4}]]}]}]",
		formattedResult)
}

func TestUpdateYaml_WithScript(t *testing.T) {
	writeScript = "instruction_sample.yaml"
	updateYaml([]string{"sample.yaml"})
}

func TestNewYaml_WithScript(t *testing.T) {
	writeScript = "instruction_sample.yaml"
	result := newYaml([]string{""})
	formattedResult := fmt.Sprintf("%v", result)
	assertResult(t,
		"[{b [{c cat} {e [[{name Mike Farah}]]}]}]",
		formattedResult)
}
