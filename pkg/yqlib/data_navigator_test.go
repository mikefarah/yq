package yqlib

import (
	"fmt"
	"sort"
	"testing"

	"github.com/mikefarah/yq/v2/test"
	logging "gopkg.in/op/go-logging.v1"
)

func TestDataNavigator(t *testing.T) {
	var log = logging.MustGetLogger("yq")
	subject := NewDataNavigator(log)

	t.Run("TestReadMap_simple", func(t *testing.T) {
		var data = test.ParseData(`
---
b:
  c: 2
`)
		got, _ := subject.ReadChildValue(data, []string{"b", "c"})
		test.AssertResult(t, 2, got)
	})

	t.Run("TestReadMap_numberKey", func(t *testing.T) {
		var data = test.ParseData(`
---
200: things
`)
		got, _ := subject.ReadChildValue(data, []string{"200"})
		test.AssertResult(t, "things", got)
	})

	t.Run("TestReadMap_splat", func(t *testing.T) {
		var data = test.ParseData(`
---
mapSplat:
  item1: things
  item2: whatever
  otherThing: cat
`)
		res, _ := subject.ReadChildValue(data, []string{"mapSplat", "*"})
		test.AssertResult(t, "[things whatever cat]", fmt.Sprintf("%v", res))
	})

	t.Run("TestReadMap_prefixSplat", func(t *testing.T) {
		var data = test.ParseData(`
---
mapSplat:
  item1: things
  item2: whatever
  otherThing: cat
`)
		res, _ := subject.ReadChildValue(data, []string{"mapSplat", "item*"})
		test.AssertResult(t, "[things whatever]", fmt.Sprintf("%v", res))
	})

	t.Run("TestReadMap_deep_splat", func(t *testing.T) {
		var data = test.ParseData(`
---
mapSplatDeep:
  item1:
    cats: bananas
  item2:
    cats: apples
`)

		res, _ := subject.ReadChildValue(data, []string{"mapSplatDeep", "*", "cats"})
		result := res.([]interface{})
		var actual = []string{result[0].(string), result[1].(string)}
		sort.Strings(actual)
		test.AssertResult(t, "[apples bananas]", fmt.Sprintf("%v", actual))
	})

	t.Run("TestReadMap_key_doesnt_exist", func(t *testing.T) {
		var data = test.ParseData(`
---
b:
  c: 2
`)
		got, _ := subject.ReadChildValue(data, []string{"b", "x", "f", "c"})
		test.AssertResult(t, nil, got)
	})

	t.Run("TestReadMap_recurse_against_string", func(t *testing.T) {
		var data = test.ParseData(`
---
a: cat
`)
		got, _ := subject.ReadChildValue(data, []string{"a", "b"})
		test.AssertResult(t, nil, got)
	})

	t.Run("TestReadMap_with_array", func(t *testing.T) {
		var data = test.ParseData(`
---
b:
  d:
    - 3
    - 4
`)
		got, _ := subject.ReadChildValue(data, []string{"b", "d", "1"})
		test.AssertResult(t, 4, got)
	})

	t.Run("TestReadMap_with_array_and_bad_index", func(t *testing.T) {
		var data = test.ParseData(`
---
b:
  d:
    - 3
    - 4
`)
		_, err := subject.ReadChildValue(data, []string{"b", "d", "x"})
		if err == nil {
			t.Fatal("Expected error due to invalid path")
		}
		expectedOutput := `error accessing array: strconv.ParseInt: parsing "x": invalid syntax`
		test.AssertResult(t, expectedOutput, err.Error())
	})

	t.Run("TestReadMap_with_mapsplat_array_and_bad_index", func(t *testing.T) {
		var data = test.ParseData(`
---
b:
  d:
    e:
      - 3
      - 4
    f:
      - 1
      - 2
`)
		_, err := subject.ReadChildValue(data, []string{"b", "d", "*", "x"})
		if err == nil {
			t.Fatal("Expected error due to invalid path")
		}
		expectedOutput := `error accessing array: strconv.ParseInt: parsing "x": invalid syntax`
		test.AssertResult(t, expectedOutput, err.Error())
	})

	t.Run("TestReadMap_with_arraysplat_map_array_and_bad_index", func(t *testing.T) {
		var data = test.ParseData(`
---
b:
  d:
    - names:
        - fred
        - smith
    - names:
        - sam
        - bo
`)
		_, err := subject.ReadChildValue(data, []string{"b", "d", "*", "names", "x"})
		if err == nil {
			t.Fatal("Expected error due to invalid path")
		}
		expectedOutput := `error accessing array: strconv.ParseInt: parsing "x": invalid syntax`
		test.AssertResult(t, expectedOutput, err.Error())
	})

	t.Run("TestReadMap_with_array_out_of_bounds", func(t *testing.T) {
		var data = test.ParseData(`
---
b:
  d:
    - 3
    - 4
`)
		got, _ := subject.ReadChildValue(data, []string{"b", "d", "3"})
		test.AssertResult(t, nil, got)
	})

	t.Run("TestReadMap_with_array_out_of_bounds_by_1", func(t *testing.T) {
		var data = test.ParseData(`
---
b:
  d:
    - 3
    - 4
`)
		got, _ := subject.ReadChildValue(data, []string{"b", "d", "2"})
		test.AssertResult(t, nil, got)
	})

	t.Run("TestReadMap_with_array_splat", func(t *testing.T) {
		var data = test.ParseData(`
e:
  -
    name: Fred
    thing: cat
  -
    name: Sam
    thing: dog
`)
		got, _ := subject.ReadChildValue(data, []string{"e", "*", "name"})
		test.AssertResult(t, "[Fred Sam]", fmt.Sprintf("%v", got))
	})

	t.Run("TestWrite_really_simple", func(t *testing.T) {
		var data = test.ParseData(`
b: 2
`)

		updated := subject.UpdatedChildValue(data, []string{"b"}, "4")
		test.AssertResult(t, "[{b 4}]", fmt.Sprintf("%v", updated))
	})

	t.Run("TestWrite_simple", func(t *testing.T) {
		var data = test.ParseData(`
b:
  c: 2
`)

		updated := subject.UpdatedChildValue(data, []string{"b", "c"}, "4")
		test.AssertResult(t, "[{b [{c 4}]}]", fmt.Sprintf("%v", updated))
	})

	t.Run("TestWrite_new", func(t *testing.T) {
		var data = test.ParseData(`
b:
  c: 2
`)

		updated := subject.UpdatedChildValue(data, []string{"b", "d"}, "4")
		test.AssertResult(t, "[{b [{c 2} {d 4}]}]", fmt.Sprintf("%v", updated))
	})

	t.Run("TestWrite_new_deep", func(t *testing.T) {
		var data = test.ParseData(`
b:
  c: 2
`)

		updated := subject.UpdatedChildValue(data, []string{"b", "d", "f"}, "4")
		test.AssertResult(t, "[{b [{c 2} {d [{f 4}]}]}]", fmt.Sprintf("%v", updated))
	})

	t.Run("TestWrite_array", func(t *testing.T) {
		var data = test.ParseData(`
b:
  - aa
`)

		updated := subject.UpdatedChildValue(data, []string{"b", "0"}, "bb")

		test.AssertResult(t, "[{b [bb]}]", fmt.Sprintf("%v", updated))
	})

	t.Run("TestWrite_new_array", func(t *testing.T) {
		var data = test.ParseData(`
b:
  c: 2
`)

		updated := subject.UpdatedChildValue(data, []string{"b", "0"}, "4")
		test.AssertResult(t, "[{b [{c 2} {0 4}]}]", fmt.Sprintf("%v", updated))
	})

	t.Run("TestWrite_new_array_deep", func(t *testing.T) {
		var data = test.ParseData(`
a: apple
`)

		updated := subject.UpdatedChildValue(data, []string{"b", "+", "c"}, "4")
		test.AssertResult(t, "[{a apple} {b [[{c 4}]]}]", fmt.Sprintf("%v", updated))
	})

	t.Run("TestWrite_new_map_array_deep", func(t *testing.T) {
		var data = test.ParseData(`
b:
  c: 2
`)

		updated := subject.UpdatedChildValue(data, []string{"b", "d", "+"}, "4")
		test.AssertResult(t, "[{b [{c 2} {d [4]}]}]", fmt.Sprintf("%v", updated))
	})

	t.Run("TestWrite_add_to_array", func(t *testing.T) {
		var data = test.ParseData(`
b:
  - aa
`)

		updated := subject.UpdatedChildValue(data, []string{"b", "1"}, "bb")
		test.AssertResult(t, "[{b [aa bb]}]", fmt.Sprintf("%v", updated))
	})

	t.Run("TestWrite_with_no_tail", func(t *testing.T) {
		var data = test.ParseData(`
b:
  c: 2
`)
		updated := subject.UpdatedChildValue(data, []string{"b"}, "4")

		test.AssertResult(t, "[{b 4}]", fmt.Sprintf("%v", updated))
	})

	t.Run("TestWriteMap_no_paths", func(t *testing.T) {
		var data = test.ParseData(`
b: 5
`)
		var new = test.ParseData(`
c: 4
`)
		result := subject.UpdatedChildValue(data, []string{}, new)
		test.AssertResult(t, fmt.Sprintf("%v", new), fmt.Sprintf("%v", result))
	})

	t.Run("TestWriteArray_no_paths", func(t *testing.T) {
		var data = make([]interface{}, 1)
		data[0] = "mike"
		var new = test.ParseData(`
c: 4
`)
		result := subject.UpdatedChildValue(data, []string{}, new)
		test.AssertResult(t, fmt.Sprintf("%v", new), fmt.Sprintf("%v", result))
	})

	t.Run("TestDelete_MapItem", func(t *testing.T) {
		var data = test.ParseData(`
a: 123
b: 456
`)
		var expected = test.ParseData(`
b: 456
`)

		result, _ := subject.DeleteChildValue(data, []string{"a"})
		test.AssertResult(t, fmt.Sprintf("%v", expected), fmt.Sprintf("%v", result))
	})

	// Ensure deleting an index into a string does nothing
	t.Run("TestDelete_index_to_string", func(t *testing.T) {
		var data = test.ParseData(`
a: mystring
`)
		result, _ := subject.DeleteChildValue(data, []string{"a", "0"})
		test.AssertResult(t, fmt.Sprintf("%v", data), fmt.Sprintf("%v", result))
	})

	t.Run("TestDelete_list_index", func(t *testing.T) {
		var data = test.ParseData(`
a: [3, 4]
`)
		var expected = test.ParseData(`
a: [3]
`)
		result, _ := subject.DeleteChildValue(data, []string{"a", "1"})
		test.AssertResult(t, fmt.Sprintf("%v", expected), fmt.Sprintf("%v", result))
	})

	t.Run("TestDelete_list_index_beyond_bounds", func(t *testing.T) {
		var data = test.ParseData(`
a: [3, 4]
`)
		result, _ := subject.DeleteChildValue(data, []string{"a", "5"})
		test.AssertResult(t, fmt.Sprintf("%v", data), fmt.Sprintf("%v", result))
	})

	t.Run("TestDelete_list_index_out_of_bounds_by_1", func(t *testing.T) {
		var data = test.ParseData(`
a: [3, 4]
`)
		result, _ := subject.DeleteChildValue(data, []string{"a", "2"})
		test.AssertResult(t, fmt.Sprintf("%v", data), fmt.Sprintf("%v", result))
	})

	t.Run("TestDelete_no_paths", func(t *testing.T) {
		var data = test.ParseData(`
a: [3, 4]
b:
  - name: test
`)
		result, _ := subject.DeleteChildValue(data, []string{})
		test.AssertResult(t, fmt.Sprintf("%v", data), fmt.Sprintf("%v", result))
	})

	t.Run("TestDelete_array_map_item", func(t *testing.T) {
		var data = test.ParseData(`
b:
- name: fred
  value: blah
- name: john
  value: test
`)
		var expected = test.ParseData(`
b:
- value: blah
- name: john
  value: test
`)
		result, _ := subject.DeleteChildValue(data, []string{"b", "0", "name"})
		test.AssertResult(t, fmt.Sprintf("%v", expected), fmt.Sprintf("%v", result))
	})
}
