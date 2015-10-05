package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

var rawData = `
---
a: Easy!
b:
  c: 2
  d:
    - 3
    - 4
e:
  -
    name: Fred
    thing: cat
  -
    name: Sam
    thing: dog
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

func TestReadMap_key_doesnt_exist(t *testing.T) {
	assertResult(t, nil, readMap(parsedData, "b.x.f", []string{"c"}))
}

func TestReadMap_recurse_against_string(t *testing.T) {
	assertResult(t, nil, readMap(parsedData, "a", []string{"b"}))
}

func TestReadMap_with_array(t *testing.T) {
	assertResult(t, 4, readMap(parsedData, "b", []string{"d", "1"}))
}

func TestReadMap_with_array_out_of_bounds(t *testing.T) {
	assertResult(t, nil, readMap(parsedData, "b", []string{"d", "3"}))
}

func TestReadMap_with_array_splat(t *testing.T) {
	assertResult(t, "[Fred Sam]", fmt.Sprintf("%v", readMap(parsedData, "e", []string{"*", "name"})))
}

func TestWrite_simple(t *testing.T) {

	write(parsedData, "b", []string{"c"}, "4")

	b := parsedData["b"].(map[interface{}]interface{})
	assertResult(t, "4", b["c"].(string))
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
