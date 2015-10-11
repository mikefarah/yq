package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestJsonFromString(t *testing.T) {
	var data = parseJSONData(`
  {
    "b": {
      "c": 2
    }
  }
`)
	assertResult(t, "map[b:map[c:2]]", fmt.Sprintf("%v", data))
}

func TestJsonFromString_withArray(t *testing.T) {
	var data = parseJSONData(`
  {
    "b": [
      { "c": 5 },
      { "c": 6 }
    ]
  }
`)
	assertResult(t, "map[b:[map[c:5] map[c:6]]]", fmt.Sprintf("%v", data))
}

func TestJsonToString(t *testing.T) {
	var data = parseData(`
---
b:
  c: 2
`)
	assertResult(t, "{\"b\":{\"c\":2}}", jsonToString(data))
}

func TestJsonToString_withArray(t *testing.T) {
	var data = parseData(`
---
b:
  - item: one
  - item: two
`)
	assertResult(t, "{\"b\":[{\"item\":\"one\"},{\"item\":\"two\"}]}", jsonToString(data))
}

func parseJSONData(rawData string) map[string]interface{} {
	var parsedData map[string]interface{}
	err := json.Unmarshal([]byte(rawData), &parsedData)
	if err != nil {
		fmt.Println("Error parsing json: ", err)
		os.Exit(1)
	}
	return parsedData
}
