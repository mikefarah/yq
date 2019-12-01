package marshal

import (
	"testing"
	"github.com/mikefarah/yq/test"
)

func TestJsonToString(t *testing.T) {
	var data = test.ParseData(`
---
b:
  c: 2
`)
	got, _ := NewJsonConverter().JsonToString(data)
	test.AssertResult(t, "{\"b\":{\"c\":2}}", got)
}

func TestJsonToString_withIntKey(t *testing.T) {
	var data = test.ParseData(`
---
b:
  2: c
`)
	got, _ := NewJsonConverter().JsonToString(data)
	test.AssertResult(t, `{"b":{"2":"c"}}`, got)
}

func TestJsonToString_withBoolKey(t *testing.T) {
	var data = test.ParseData(`
---
b:
  false: c
`)
	got, _ := NewJsonConverter().JsonToString(data)
	test.AssertResult(t, `{"b":{"false":"c"}}`, got)
}

func TestJsonToString_withArray(t *testing.T) {
	var data = test.ParseData(`
---
b:
  - item: one
  - item: two
`)
	got, _ := NewJsonConverter().JsonToString(data)
	test.AssertResult(t, "{\"b\":[{\"item\":\"one\"},{\"item\":\"two\"}]}", got)
}
