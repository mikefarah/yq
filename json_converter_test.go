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
	assertResult(t, "{\"b\":{\"c\":2}}", jsonToString(data))
}

func TestJsonToString_withIntKey(t *testing.T) {
	var data = parseData(`
---
b:
  2: c
`)
	assertResult(t, `{"b":{"2":"c"}}`, jsonToString(data))
}

func TestJsonToString_withBoolKey(t *testing.T) {
	var data = parseData(`
---
b:
  false: c
`)
	assertResult(t, `{"b":{"false":"c"}}`, jsonToString(data))
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
