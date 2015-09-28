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
	result := readMap(parsedData, "b", []string{"c"})
	if result != 2 {
		t.Error("Excpted 2 but got ", result)
	}
}

func TestReadMap_array(t *testing.T) {
	result := readMap(parsedData, "b", []string{"d", "1"})
	if result != 4 {
		t.Error("Excpted 4 but got ", result)
	}
}
