package yqlib

import (
	"fmt"
	"testing"
	"github.com/mikefarah/yq/test"
)

func TestReadPath(t *testing.T) {
	var data = test.ParseData(`
---
b:
  2: c
`)

	got, _ := ReadPath(data, "b.2")
	test.AssertResult(t, `c`, got)
}

func TestReadPath_WithError(t *testing.T) {
	var data = test.ParseData(`
---
b:
  - c
`)

	_, err := ReadPath(data, "b.[a]")
	if err == nil {
		t.Fatal("Expected error due to invalid path")
	}
}

func TestWritePath(t *testing.T) {
	var data = test.ParseData(`
---
b:
  2: c
`)

	got := WritePath(data, "b.3", "a")
	test.AssertResult(t, `[{b [{2 c} {3 a}]}]`, fmt.Sprintf("%v", got))
}

func TestPrefixPath(t *testing.T) {
	var data = test.ParseData(`
---
b:
  2: c
`)

	got := PrefixPath(data, "d")
	test.AssertResult(t, `[{d [{b [{2 c}]}]}]`, fmt.Sprintf("%v", got))
}

func TestDeletePath(t *testing.T) {
	var data = test.ParseData(`
---
b:
  2: c
  3: a
`)

	got, _ := DeletePath(data, "b.2")
	test.AssertResult(t, `[{b [{3 a}]}]`, fmt.Sprintf("%v", got))
}

func TestDeletePath_WithError(t *testing.T) {
	var data = test.ParseData(`
---
b:
  - c
`)

	_, err := DeletePath(data, "b.[a]")
	if err == nil {
		t.Fatal("Expected error due to invalid path")
	}
}

func TestMerge(t *testing.T) {
	var dst = test.ParseData(`
---
a: b
c: d
`)
	var src = test.ParseData(`
---
a: 1
b: 2
`)

	var mergedData = make(map[interface{}]interface{})
	mergedData["root"] = dst
	var mapDataBucket = make(map[interface{}]interface{})
	mapDataBucket["root"] = src

	Merge(&mergedData, mapDataBucket, false, false)
	test.AssertResult(t, `[{a b} {c d}]`, fmt.Sprintf("%v", mergedData["root"]))
}

func TestMerge_WithOverwrite(t *testing.T) {
	var dst = test.ParseData(`
---
a: b
c: d
`)
	var src = test.ParseData(`
---
a: 1
b: 2
`)

	var mergedData = make(map[interface{}]interface{})
	mergedData["root"] = dst
	var mapDataBucket = make(map[interface{}]interface{})
	mapDataBucket["root"] = src

	Merge(&mergedData, mapDataBucket, true, false)
	test.AssertResult(t, `[{a 1} {b 2}]`, fmt.Sprintf("%v", mergedData["root"]))
}

func TestMerge_WithAppend(t *testing.T) {
	var dst = test.ParseData(`
---
a: b
c: d
`)
	var src = test.ParseData(`
---
a: 1
b: 2
`)

	var mergedData = make(map[interface{}]interface{})
	mergedData["root"] = dst
	var mapDataBucket = make(map[interface{}]interface{})
	mapDataBucket["root"] = src

	Merge(&mergedData, mapDataBucket, false, true)
	test.AssertResult(t, `[{a b} {c d} {a 1} {b 2}]`, fmt.Sprintf("%v", mergedData["root"]))
}

func TestMerge_WithError(t *testing.T) {
	err := Merge(nil, nil, false, false)
	if err == nil {
		t.Fatal("Expected error due to nil")
	}
}