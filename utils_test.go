package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

func parseData(rawData string) yaml.MapSlice {
	var parsedData yaml.MapSlice
	err := yaml.Unmarshal([]byte(rawData), &parsedData)
	if err != nil {
		fmt.Println("Error parsing yaml: %v", err)
		os.Exit(1)
	}
	return parsedData
}

func assertResult(t *testing.T, expectedValue interface{}, actualValue interface{}) {
	if expectedValue != actualValue {
		t.Error("Expected <", expectedValue, "> but got <", actualValue, ">", fmt.Sprintf("%T", actualValue))
	}
}

func assertResultWithContext(t *testing.T, expectedValue interface{}, actualValue interface{}, context interface{}) {

	if expectedValue != actualValue {
		t.Error(context)
		t.Error(": expected <", expectedValue, "> but got <", actualValue, ">")
	}
}
