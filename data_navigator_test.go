package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

func parseData(rawData string) map[interface{}]interface{} {
	var parsedData map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(rawData), &parsedData)
	if err != nil {
		fmt.Println("Error parsing yaml: %v", err)
		os.Exit(1)
	}
	return parsedData
}

func TestReadMap_simple(t *testing.T) {
	var data = parseData(`
---
b:
  c: 2
`)
	assertResult(t, 2, readMap(data, "b", []string{"c"}))
}

func TestReadMap_key_doesnt_exist(t *testing.T) {
	var data = parseData(`
---
b:
  c: 2
`)
	assertResult(t, nil, readMap(data, "b.x.f", []string{"c"}))
}

func TestReadMap_recurse_against_string(t *testing.T) {
	var data = parseData(`
---
a: cat
`)
	assertResult(t, nil, readMap(data, "a", []string{"b"}))
}

func TestReadMap_with_array(t *testing.T) {
	var data = parseData(`
---
b:
  d:
    - 3
    - 4
`)
	assertResult(t, 4, readMap(data, "b", []string{"d", "1"}))
}

func TestReadMap_with_array_out_of_bounds(t *testing.T) {
	var data = parseData(`
---
b:
  d:
    - 3
    - 4
`)
	assertResult(t, nil, readMap(data, "b", []string{"d", "3"}))
}

func TestReadMap_with_array_splat(t *testing.T) {
	var data = parseData(`
e:
  -
    name: Fred
    thing: cat
  -
    name: Sam
    thing: dog
`)
	assertResult(t, "[Fred Sam]", fmt.Sprintf("%v", readMap(data, "e", []string{"*", "name"})))
}

func TestWrite_simple(t *testing.T) {
	var data = parseData(`
b:
  c: 2
`)

	write(data, "b", []string{"c"}, "4")

	b := data["b"].(map[interface{}]interface{})
	assertResult(t, "4", b["c"].(string))
}

func TestWrite_with_no_tail(t *testing.T) {
	var data = parseData(`
b:
  c: 2
`)
	write(data, "b", []string{}, "4")

	b := data["b"]
	assertResult(t, "4", fmt.Sprintf("%v", b))
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
