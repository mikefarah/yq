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

func AssertResult(t *testing.T, expectedValue interface{}, actualValue interface{}) {
	t.Helper()
	if expectedValue != actualValue {
		t.Error("Expected <", expectedValue, "> but got <", actualValue, ">", fmt.Sprintf("%T", actualValue))
	}
}

func AssertResultComplex(t *testing.T, expectedValue interface{}, actualValue interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expectedValue, actualValue) {
		t.Error("\nExpected <", expectedValue, ">\nbut got  <", actualValue, ">", fmt.Sprintf("%T", actualValue))
	}
}

func AssertResultComplexWithContext(t *testing.T, expectedValue interface{}, actualValue interface{}, context interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expectedValue, actualValue) {
		t.Error(context)
		t.Error("\nExpected <", expectedValue, ">\nbut got  <", actualValue, ">", fmt.Sprintf("%T", actualValue))
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
