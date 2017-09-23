package main

import (
	"testing"
)

func TestJsonToString(t *testing.T) {
	var data = parseData(`
---
b:
  c: 2
`)
	got, _ := jsonToString(data)
	assertResult(t, "{\"b\":{\"c\":2}}", got)
}

func TestJsonToString_withIntKey(t *testing.T) {
	var data = parseData(`
---
b:
  2: c
`)
	got, _ := jsonToString(data)
	assertResult(t, `{"b":{"2":"c"}}`, got)
}

func TestJsonToString_withBoolKey(t *testing.T) {
	var data = parseData(`
---
b:
  false: c
`)
	got, _ := jsonToString(data)
	assertResult(t, `{"b":{"false":"c"}}`, got)
}

func TestJsonToString_withArray(t *testing.T) {
	var data = parseData(`
---
b:
  - item: one
  - item: two
`)
	got, _ := jsonToString(data)
	assertResult(t, "{\"b\":[{\"item\":\"one\"},{\"item\":\"two\"}]}", got)
}
