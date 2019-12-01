package main

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/mikefarah/yq/pkg/marshal"
	"github.com/mikefarah/yq/test"
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
		test.AssertResultWithContext(t, tt.expectedResult, parseValue(tt.argument), tt.testDescription)
	}
}

func TestMultilineString(t *testing.T) {
	testString := `
	abcd
	efg`
	formattedResult, _ := marshal.NewYamlConverter().YamlToString(testString, false)
	test.AssertResult(t, testString, formattedResult)
}

func TestNewYaml(t *testing.T) {
	result, _ := newYaml([]string{"b.c", "3"})
	formattedResult := fmt.Sprintf("%v", result)
	test.AssertResult(t,
		"[{b [{c 3}]}]",
		formattedResult)
}

func TestNewYamlArray(t *testing.T) {
	result, _ := newYaml([]string{"[0].cat", "meow"})
	formattedResult := fmt.Sprintf("%v", result)
	test.AssertResult(t,
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
	actualResult, _ := marshal.NewYamlConverter().YamlToString(result, true)
	test.AssertResult(t, expectedResult, actualResult)
}

func TestNewYaml_WithUnknownScript(t *testing.T) {
	writeScript = "fake-unknown"
	_, err := newYaml([]string{""})
	if err == nil {
		t.Error("Expected error due to unknown file")
	}
	var expectedOutput string
	if runtime.GOOS == "windows" {
		expectedOutput = `open fake-unknown: The system cannot find the file specified.`
	} else {
		expectedOutput = `open fake-unknown: no such file or directory`
	}
	test.AssertResult(t, expectedOutput, err.Error())
}
