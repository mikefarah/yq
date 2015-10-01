package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

var rawData = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

var parsedData map[interface{}]interface{}

func TestMain(m *testing.M) {
	err := yaml.Unmarshal([]byte(rawData), &parsedData)
	if err != nil {
		fmt.Println("Error parsing yaml: %v", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestReadMap_simple(t *testing.T) {
	assertResult(t, 2, readMap(parsedData, "b", []string{"c"}))
}

func TestReadMap_array(t *testing.T) {
	assertResult(t, 4, readMap(parsedData, "b", []string{"d", "1"}))
}

func TestWrite_simple(t *testing.T) {

	write(parsedData, "b", []string{"c"}, "4")

	b := parsedData["b"].(map[interface{}]interface{})
	assertResult(t, "4", b["c"].(string))
}


func assertResult(t *testing.T, expectedValue interface{}, actualValue interface{}) {
	if (expectedValue != actualValue) {
		t.Error("Expected <", expectedValue, "> but got <", actualValue, ">")
	}
}
