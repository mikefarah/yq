package test

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

type resulter struct {
	Error   error
	Output  string
	Command *cobra.Command
}

func RunCmd(c *cobra.Command, input string) resulter {
	buf := new(bytes.Buffer)
	c.SetOutput(buf)
	c.SetArgs(strings.Split(input, " "))

	err := c.Execute()
	output := buf.String()

	return resulter{err, output, c}
}

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
	if expectedValue != actualValue {
		t.Error(context)
		t.Error(": expected <", expectedValue, "> but got <", actualValue, ">")
	}
}
