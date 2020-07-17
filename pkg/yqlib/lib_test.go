package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v3/test"
)

func TestLib(t *testing.T) {

	subject := NewYqLib()

	t.Run("PathStackToString_Empty", func(t *testing.T) {
		emptyArray := make([]interface{}, 0)
		got := subject.PathStackToString(emptyArray)
		test.AssertResult(t, ``, got)
	})

	t.Run("PathStackToString", func(t *testing.T) {
		array := make([]interface{}, 3)
		array[0] = "a"
		array[1] = 0
		array[2] = "b"
		got := subject.PathStackToString(array)
		test.AssertResult(t, `a.[0].b`, got)
	})

	t.Run("MergePathStackToString", func(t *testing.T) {
		array := make([]interface{}, 3)
		array[0] = "a"
		array[1] = 0
		array[2] = "b"
		got := subject.MergePathStackToString(array, AppendArrayMergeStrategy)
		test.AssertResult(t, `a.[+].b`, got)
	})

	// 	t.Run("TestReadPath_WithError", func(t *testing.T) {
	// 		var data = test.ParseData(`
	// ---
	// b:
	//   - c
	// `)

	// 		_, err := subject.ReadPath(data, "b.[a]")
	// 		if err == nil {
	// 			t.Fatal("Expected error due to invalid path")
	// 		}
	// 	})

	// 	t.Run("TestWritePath", func(t *testing.T) {
	// 		var data = test.ParseData(`
	// ---
	// b:
	//   2: c
	// `)

	// 		got := subject.WritePath(data, "b.3", "a")
	// 		test.AssertResult(t, `[{b [{2 c} {3 a}]}]`, fmt.Sprintf("%v", got))
	// 	})

	// 	t.Run("TestPrefixPath", func(t *testing.T) {
	// 		var data = test.ParseData(`
	// ---
	// b:
	//   2: c
	// `)

	// 		got := subject.PrefixPath(data, "a.d")
	// 		test.AssertResult(t, `[{a [{d [{b [{2 c}]}]}]}]`, fmt.Sprintf("%v", got))
	// 	})

	// 	t.Run("TestDeletePath", func(t *testing.T) {
	// 		var data = test.ParseData(`
	// ---
	// b:
	//   2: c
	//   3: a
	// `)

	// 		got, _ := subject.DeletePath(data, "b.2")
	// 		test.AssertResult(t, `[{b [{3 a}]}]`, fmt.Sprintf("%v", got))
	// 	})

	// 	t.Run("TestDeletePath_WithError", func(t *testing.T) {
	// 		var data = test.ParseData(`
	// ---
	// b:
	//   - c
	// `)

	// 		_, err := subject.DeletePath(data, "b.[a]")
	// 		if err == nil {
	// 			t.Fatal("Expected error due to invalid path")
	// 		}
	// 	})

	// 	t.Run("TestMerge", func(t *testing.T) {
	// 		var dst = test.ParseData(`
	// ---
	// a: b
	// c: d
	// `)
	// 		var src = test.ParseData(`
	// ---
	// a: 1
	// b: 2
	// `)

	// 		var mergedData = make(map[interface{}]interface{})
	// 		mergedData["root"] = dst
	// 		var mapDataBucket = make(map[interface{}]interface{})
	// 		mapDataBucket["root"] = src

	// 		err := subject.Merge(&mergedData, mapDataBucket, false, false)
	// 		if err != nil {
	// 			t.Fatal("Unexpected error")
	// 		}
	// 		test.AssertResult(t, `[{a b} {c d}]`, fmt.Sprintf("%v", mergedData["root"]))
	// 	})

	// 	t.Run("TestMerge_WithOverwrite", func(t *testing.T) {
	// 		var dst = test.ParseData(`
	// ---
	// a: b
	// c: d
	// `)
	// 		var src = test.ParseData(`
	// ---
	// a: 1
	// b: 2
	// `)

	// 		var mergedData = make(map[interface{}]interface{})
	// 		mergedData["root"] = dst
	// 		var mapDataBucket = make(map[interface{}]interface{})
	// 		mapDataBucket["root"] = src

	// 		err := subject.Merge(&mergedData, mapDataBucket, true, false)
	// 		if err != nil {
	// 			t.Fatal("Unexpected error")
	// 		}
	// 		test.AssertResult(t, `[{a 1} {b 2}]`, fmt.Sprintf("%v", mergedData["root"]))
	// 	})

	// 	t.Run("TestMerge_WithAppend", func(t *testing.T) {
	// 		var dst = test.ParseData(`
	// ---
	// a: b
	// c: d
	// `)
	// 		var src = test.ParseData(`
	// ---
	// a: 1
	// b: 2
	// `)

	// 		var mergedData = make(map[interface{}]interface{})
	// 		mergedData["root"] = dst
	// 		var mapDataBucket = make(map[interface{}]interface{})
	// 		mapDataBucket["root"] = src

	// 		err := subject.Merge(&mergedData, mapDataBucket, false, true)
	// 		if err != nil {
	// 			t.Fatal("Unexpected error")
	// 		}
	// 		test.AssertResult(t, `[{a b} {c d} {a 1} {b 2}]`, fmt.Sprintf("%v", mergedData["root"]))
	// 	})

	// 	t.Run("TestMerge_WithError", func(t *testing.T) {
	// 		err := subject.Merge(nil, nil, false, false)
	// 		if err == nil {
	// 			t.Fatal("Expected error due to nil")
	// 		}
	// 	})

}
