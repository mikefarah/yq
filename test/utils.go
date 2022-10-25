package test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/pkg/diff"
	"github.com/pkg/diff/write"
	yaml "gopkg.in/yaml.v3"
)

func ParseData(rawData string) yaml.Node {
	var parsedData yaml.Node
	err := yaml.Unmarshal([]byte(rawData), &parsedData)
	if err != nil {
		fmt.Printf("Error parsing yaml: %v\n", err)
		os.Exit(1)
	}
	return parsedData
}

func printDifference(t *testing.T, expectedValue interface{}, actualValue interface{}) {
	opts := []write.Option{write.TerminalColor()}
	var differenceBuffer bytes.Buffer
	expectedString := fmt.Sprintf("%v", expectedValue)
	actualString := fmt.Sprintf("%v", actualValue)
	if err := diff.Text("expected", "actual", expectedString, actualString, bufio.NewWriter(&differenceBuffer), opts...); err != nil {
		t.Error(err)
	} else {
		t.Error(differenceBuffer.String())
	}
}

func AssertResult(t *testing.T, expectedValue interface{}, actualValue interface{}) {
	t.Helper()
	if expectedValue != actualValue {
		printDifference(t, expectedValue, actualValue)
	}
}

func AssertResultComplex(t *testing.T, expectedValue interface{}, actualValue interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expectedValue, actualValue) {
		printDifference(t, expectedValue, actualValue)
	}
}

func AssertResultComplexWithContext(t *testing.T, expectedValue interface{}, actualValue interface{}, context interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expectedValue, actualValue) {
		t.Error(context)
		printDifference(t, expectedValue, actualValue)
	}
}

func AssertResultWithContext(t *testing.T, expectedValue interface{}, actualValue interface{}, context interface{}) {
	t.Helper()
	opts := []write.Option{write.TerminalColor()}
	if expectedValue != actualValue {
		t.Error(context)
		var differenceBuffer bytes.Buffer
		if err := diff.Text("expected", "actual", expectedValue, actualValue, bufio.NewWriter(&differenceBuffer), opts...); err != nil {
			t.Error(err)
		} else {
			t.Error(differenceBuffer.String())
		}
	}
}
