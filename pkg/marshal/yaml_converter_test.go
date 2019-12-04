package marshal

import (
	"testing"

	"github.com/mikefarah/yq/v2/test"
)

func TestYamlToString(t *testing.T) {
	var raw = `b:
  c: 2
`
	var data = test.ParseData(raw)
	got, _ := NewYamlConverter().YamlToString(data, false)
	test.AssertResult(t, raw, got)
}

func TestYamlToString_withTrim(t *testing.T) {
	var raw = `b:
  c: 2`
	var data = test.ParseData(raw)
	got, _ := NewYamlConverter().YamlToString(data, true)
	test.AssertResult(t, raw, got)
}

func TestYamlToString_withIntKey(t *testing.T) {
	var raw = `b:
  2: c
`
	var data = test.ParseData(raw)
	got, _ := NewYamlConverter().YamlToString(data, false)
	test.AssertResult(t, raw, got)
}

func TestYamlToString_withBoolKey(t *testing.T) {
	var raw = `b:
  false: c
`
	var data = test.ParseData(raw)
	got, _ := NewYamlConverter().YamlToString(data, false)
	test.AssertResult(t, raw, got)
}

func TestYamlToString_withArray(t *testing.T) {
	var raw = `b:
- item: one
- item: two
`
	var data = test.ParseData(raw)
	got, _ := NewYamlConverter().YamlToString(data, false)
	test.AssertResult(t, raw, got)
}
