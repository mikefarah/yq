package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	yaml "gopkg.in/mikefarah/yaml.v2"
	cobra "gopkg.in/spf13/cobra.v0"
)

type resulter struct {
	Error   error
	Output  string
	Command *cobra.Command
}

func runCmd(c *cobra.Command, input string) resulter {
	buf := new(bytes.Buffer)
	c.SetOutput(buf)
	c.SetArgs(strings.Split(input, " "))

	err := c.Execute()
	output := buf.String()

	return resulter{err, output, c}
}

func parseData(rawData string) yaml.MapSlice {
	var parsedData yaml.MapSlice
	err := yaml.Unmarshal([]byte(rawData), &parsedData)
	if err != nil {
		fmt.Printf("Error parsing yaml: %v\n", err)
		os.Exit(1)
	}
	return parsedData
}

func assertResult(t *testing.T, expectedValue interface{}, actualValue interface{}) {
	t.Helper()
	if expectedValue != actualValue {
		t.Error("Expected <", expectedValue, "> but got <", actualValue, ">", fmt.Sprintf("%T", actualValue))
	}
}

func assertResultComplex(t *testing.T, expectedValue interface{}, actualValue interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expectedValue, actualValue) {
		t.Error("Expected <", expectedValue, "> but got <", actualValue, ">", fmt.Sprintf("%T", actualValue))
	}
}

func assertResultWithContext(t *testing.T, expectedValue interface{}, actualValue interface{}, context interface{}) {
	t.Helper()
	if expectedValue != actualValue {
		t.Error(context)
		t.Error(": expected <", expectedValue, "> but got <", actualValue, ">")
	}
}

func writeTempYamlFile(content string) string {
	tmpfile, _ := ioutil.TempFile("", "testyaml")
	defer func() {
		_ = tmpfile.Close()
	}()

	_, _ = tmpfile.Write([]byte(content))
	return tmpfile.Name()
}

func readTempYamlFile(name string) string {
	content, _ := ioutil.ReadFile(name)
	return string(content)
}

func removeTempYamlFile(name string) {
	_ = os.Remove(name)
}
