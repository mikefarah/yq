package yqlib

import (
	"fmt"
	"testing"

	"github.com/mikefarah/yq/test"
	logging "gopkg.in/op/go-logging.v1"
)

func TestLib(t *testing.T) {

	var log = logging.MustGetLogger("yq")
	subject := NewYqLib(log)

	t.Run("TestReadPath", func(t *testing.T) {
		var data = test.ParseData(`
---
b:
  2: c
`)

		got, _ := subject.ReadPath(data, "b.2")
		test.AssertResult(t, `c`, got)
	})

	t.Run("TestReadPath_WithError", func(t *testing.T) {
		var data = test.ParseData(`
---
b:
  - c
`)

		_, err := subject.ReadPath(data, "b.[a]")
		if err == nil {
			t.Fatal("Expected error due to invalid path")
		}
	})

	t.Run("TestWritePath", func(t *testing.T) {
		var data = test.ParseData(`
---
b:
  2: c
`)

		got := subject.WritePath(data, "b.3", "a")
		test.AssertResult(t, `[{b [{2 c} {3 a}]}]`, fmt.Sprintf("%v", got))
	})

	t.Run("TestPrefixPath", func(t *testing.T) {
		var data = test.ParseData(`
---
b:
  2: c
`)

		got := subject.PrefixPath(data, "a.d")
		test.AssertResult(t, `[{a [{d [{b [{2 c}]}]}]}]`, fmt.Sprintf("%v", got))
	})

	t.Run("TestDeletePath", func(t *testing.T) {
		var data = test.ParseData(`
---
b:
  2: c
  3: a
`)

		got, _ := subject.DeletePath(data, "b.2")
		test.AssertResult(t, `[{b [{3 a}]}]`, fmt.Sprintf("%v", got))
	})

	t.Run("TestDeletePath_WithError", func(t *testing.T) {
		var data = test.ParseData(`
---
b:
  - c
`)

		_, err := subject.DeletePath(data, "b.[a]")
		if err == nil {
			t.Fatal("Expected error due to invalid path")
		}
	})

	t.Run("TestMerge", func(t *testing.T) {
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

		err := subject.Merge(&mergedData, mapDataBucket, false, false)
		if err != nil {
			t.Fatal("Unexpected error")
		}
		test.AssertResult(t, `[{a b} {c d}]`, fmt.Sprintf("%v", mergedData["root"]))
	})

	t.Run("TestMerge_WithOverwrite", func(t *testing.T) {
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

		err := subject.Merge(&mergedData, mapDataBucket, true, false)
		if err != nil {
			t.Fatal("Unexpected error")
		}
		test.AssertResult(t, `[{a 1} {b 2}]`, fmt.Sprintf("%v", mergedData["root"]))
	})

	t.Run("TestMerge_WithAppend", func(t *testing.T) {
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

		err := subject.Merge(&mergedData, mapDataBucket, false, true)
		if err != nil {
			t.Fatal("Unexpected error")
		}
		test.AssertResult(t, `[{a b} {c d} {a 1} {b 2}]`, fmt.Sprintf("%v", mergedData["root"]))
	})

	t.Run("TestMerge_WithError", func(t *testing.T) {
		err := subject.Merge(nil, nil, false, false)
		if err == nil {
			t.Fatal("Expected error due to nil")
		}
	})

}
