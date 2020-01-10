package main

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/mikefarah/yq/v2/pkg/marshal"
	"github.com/mikefarah/yq/v2/test"
)

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

func TestNewYamlBigInt(t *testing.T) {
	result, _ := newYaml([]string{"b", "1212121"})
	formattedResult := fmt.Sprintf("%v", result)
	test.AssertResult(t,
		"[{b 1212121}]",
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
