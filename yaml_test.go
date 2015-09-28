package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

var raw_data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

var parsed_data map[interface{}]interface{}

func TestMain(m *testing.M) {
	err := yaml.Unmarshal([]byte(raw_data), &parsed_data)
	if err != nil {
		fmt.Println("Error parsing yaml: %v", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestRead_map_simple(t *testing.T) {
	result := read_map(parsed_data, "b", []string{"c"})
	if result != 2 {
		t.Error("Excpted 2 but got ", result)
	}
}

func TestRead_map_array(t *testing.T) {
	result := read_map(parsed_data, "b", []string{"d", "1"})
	if result != 4 {
		t.Error("Excpted 4 but got ", result)
	}
}
