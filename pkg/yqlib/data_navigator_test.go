package yqlib

import (
	"fmt"
	"sort"
	"testing"
	"github.com/mikefarah/yq/test"
)

func TestReadMap_simple(t *testing.T) {
	var data = test.ParseData(`
---
b:
  c: 2
`)
	got, _ := readMap(data, "b", []string{"c"})
	test.AssertResult(t, 2, got)
}

func TestReadMap_numberKey(t *testing.T) {
	var data = test.ParseData(`
---
200: things
`)
	got, _ := readMap(data, "200", []string{})
	test.AssertResult(t, "things", got)
}

func TestReadMap_splat(t *testing.T) {
	var data = test.ParseData(`
---
mapSplat:
  item1: things
  item2: whatever
  otherThing: cat
`)
	res, _ := readMap(data, "mapSplat", []string{"*"})
	test.AssertResult(t, "[things whatever cat]", fmt.Sprintf("%v", res))
}

func TestReadMap_prefixSplat(t *testing.T) {
	var data = test.ParseData(`
---
mapSplat:
  item1: things
  item2: whatever
  otherThing: cat
`)
	res, _ := readMap(data, "mapSplat", []string{"item*"})
	test.AssertResult(t, "[things whatever]", fmt.Sprintf("%v", res))
}

func TestReadMap_deep_splat(t *testing.T) {
	var data = test.ParseData(`
---
mapSplatDeep:
  item1:
    cats: bananas
  item2:
    cats: apples
`)

	res, _ := readMap(data, "mapSplatDeep", []string{"*", "cats"})
	result := res.([]interface{})
	var actual = []string{result[0].(string), result[1].(string)}
	sort.Strings(actual)
	test.AssertResult(t, "[apples bananas]", fmt.Sprintf("%v", actual))
}

func TestReadMap_key_doesnt_exist(t *testing.T) {
	var data = test.ParseData(`
---
b:
  c: 2
`)
	got, _ := readMap(data, "b.x.f", []string{"c"})
	test.AssertResult(t, nil, got)
}

func TestReadMap_recurse_against_string(t *testing.T) {
	var data = test.ParseData(`
---
a: cat
`)
	got, _ := readMap(data, "a", []string{"b"})
	test.AssertResult(t, nil, got)
}

func TestReadMap_with_array(t *testing.T) {
	var data = test.ParseData(`
---
b:
  d:
    - 3
    - 4
`)
	got, _ := readMap(data, "b", []string{"d", "1"})
	test.AssertResult(t, 4, got)
}

func TestReadMap_with_array_and_bad_index(t *testing.T) {
	var data = test.ParseData(`
---
b:
  d:
    - 3
    - 4
`)
	_, err := readMap(data, "b", []string{"d", "x"})
	if err == nil {
		t.Fatal("Expected error due to invalid path")
	}
	expectedOutput := `error accessing array: strconv.ParseInt: parsing "x": invalid syntax`
	test.AssertResult(t, expectedOutput, err.Error())
}

func TestReadMap_with_mapsplat_array_and_bad_index(t *testing.T) {
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
	_, err := readMap(data, "b", []string{"d", "*", "x"})
	if err == nil {
		t.Fatal("Expected error due to invalid path")
	}
	expectedOutput := `error accessing array: strconv.ParseInt: parsing "x": invalid syntax`
	test.AssertResult(t, expectedOutput, err.Error())
}

func TestReadMap_with_arraysplat_map_array_and_bad_index(t *testing.T) {
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
	_, err := readMap(data, "b", []string{"d", "*", "names", "x"})
	if err == nil {
		t.Fatal("Expected error due to invalid path")
	}
	expectedOutput := `error accessing array: strconv.ParseInt: parsing "x": invalid syntax`
	test.AssertResult(t, expectedOutput, err.Error())
}

func TestReadMap_with_array_out_of_bounds(t *testing.T) {
	var data = test.ParseData(`
---
b:
  d:
    - 3
    - 4
`)
	got, _ := readMap(data, "b", []string{"d", "3"})
	test.AssertResult(t, nil, got)
}

func TestReadMap_with_array_out_of_bounds_by_1(t *testing.T) {
	var data = test.ParseData(`
---
b:
  d:
    - 3
    - 4
`)
	got, _ := readMap(data, "b", []string{"d", "2"})
	test.AssertResult(t, nil, got)
}

func TestReadMap_with_array_splat(t *testing.T) {
	var data = test.ParseData(`
e:
  -
    name: Fred
    thing: cat
  -
    name: Sam
    thing: dog
`)
	got, _ := readMap(data, "e", []string{"*", "name"})
	test.AssertResult(t, "[Fred Sam]", fmt.Sprintf("%v", got))
}

func TestWrite_really_simple(t *testing.T) {
	var data = test.ParseData(`
    b: 2
`)

	updated := writeMap(data, []string{"b"}, "4")
	test.AssertResult(t, "[{b 4}]", fmt.Sprintf("%v", updated))
}

func TestWrite_simple(t *testing.T) {
	var data = test.ParseData(`
b:
  c: 2
`)

	updated := writeMap(data, []string{"b", "c"}, "4")
	test.AssertResult(t, "[{b [{c 4}]}]", fmt.Sprintf("%v", updated))
}

func TestWrite_new(t *testing.T) {
	var data = test.ParseData(`
b:
  c: 2
`)

	updated := writeMap(data, []string{"b", "d"}, "4")
	test.AssertResult(t, "[{b [{c 2} {d 4}]}]", fmt.Sprintf("%v", updated))
}

func TestWrite_new_deep(t *testing.T) {
	var data = test.ParseData(`
b:
  c: 2
`)

	updated := writeMap(data, []string{"b", "d", "f"}, "4")
	test.AssertResult(t, "[{b [{c 2} {d [{f 4}]}]}]", fmt.Sprintf("%v", updated))
}

func TestWrite_array(t *testing.T) {
	var data = test.ParseData(`
b:
  - aa
`)

	updated := writeMap(data, []string{"b", "0"}, "bb")

	test.AssertResult(t, "[{b [bb]}]", fmt.Sprintf("%v", updated))
}

func TestWrite_new_array(t *testing.T) {
	var data = test.ParseData(`
b:
  c: 2
`)

	updated := writeMap(data, []string{"b", "0"}, "4")
	test.AssertResult(t, "[{b [{c 2} {0 4}]}]", fmt.Sprintf("%v", updated))
}

func TestWrite_new_array_deep(t *testing.T) {
	var data = test.ParseData(`
a: apple
`)

	var expected = `a: apple
b:
- c: "4"`

	updated := writeMap(data, []string{"b", "+", "c"}, "4")
	got, _ := YamlToString(updated, true)
	test.AssertResult(t, expected, got)
}

func TestWrite_new_map_array_deep(t *testing.T) {
	var data = test.ParseData(`
b:
  c: 2
`)
	var expected = `b:
  c: 2
  d:
  - "4"`

	updated := writeMap(data, []string{"b", "d", "+"}, "4")
	got, _ := YamlToString(updated, true)
	test.AssertResult(t, expected, got)
}

func TestWrite_add_to_array(t *testing.T) {
	var data = test.ParseData(`
b:
  - aa
`)

	var expected = `b:
- aa
- bb`

	updated := writeMap(data, []string{"b", "1"}, "bb")
	got, _ := YamlToString(updated, true)
	test.AssertResult(t, expected, got)
}

func TestWrite_with_no_tail(t *testing.T) {
	var data = test.ParseData(`
b:
  c: 2
`)
	updated := writeMap(data, []string{"b"}, "4")

	test.AssertResult(t, "[{b 4}]", fmt.Sprintf("%v", updated))
}

func TestWriteMap_no_paths(t *testing.T) {
	var data = test.ParseData(`
b: 5
`)

	result := writeMap(data, []string{}, 4)
	test.AssertResult(t, fmt.Sprintf("%v", data), fmt.Sprintf("%v", result))
}

func TestWriteArray_no_paths(t *testing.T) {
	var data = make([]interface{}, 1)
	data[0] = "mike"
	result := writeArray(data, []string{}, 4)
	test.AssertResult(t, fmt.Sprintf("%v", data), fmt.Sprintf("%v", result))
}

func TestDelete_MapItem(t *testing.T) {
	var data = test.ParseData(`
a: 123
b: 456
`)
	var expected = test.ParseData(`
b: 456
`)

	result, _ := deleteMap(data, []string{"a"})
	test.AssertResult(t, fmt.Sprintf("%v", expected), fmt.Sprintf("%v", result))
}

// Ensure deleting an index into a string does nothing
func TestDelete_index_to_string(t *testing.T) {
	var data = test.ParseData(`
a: mystring
`)
	result, _ := deleteMap(data, []string{"a", "0"})
	test.AssertResult(t, fmt.Sprintf("%v", data), fmt.Sprintf("%v", result))
}

func TestDelete_list_index(t *testing.T) {
	var data = test.ParseData(`
a: [3, 4]
`)
	var expected = test.ParseData(`
a: [3]
`)
	result, _ := deleteMap(data, []string{"a", "1"})
	test.AssertResult(t, fmt.Sprintf("%v", expected), fmt.Sprintf("%v", result))
}

func TestDelete_list_index_beyond_bounds(t *testing.T) {
	var data = test.ParseData(`
a: [3, 4]
`)
	result, _ := deleteMap(data, []string{"a", "5"})
	test.AssertResult(t, fmt.Sprintf("%v", data), fmt.Sprintf("%v", result))
}

func TestDelete_list_index_out_of_bounds_by_1(t *testing.T) {
	var data = test.ParseData(`
a: [3, 4]
`)
	result, _ := deleteMap(data, []string{"a", "2"})
	test.AssertResult(t, fmt.Sprintf("%v", data), fmt.Sprintf("%v", result))
}

func TestDelete_no_paths(t *testing.T) {
	var data = test.ParseData(`
a: [3, 4]
b:
  - name: test
`)
	result, _ := deleteMap(data, []string{})
	test.AssertResult(t, fmt.Sprintf("%v", data), fmt.Sprintf("%v", result))
}

func TestDelete_array_map_item(t *testing.T) {
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
	result, _ := deleteMap(data, []string{"b", "0", "name"})
	test.AssertResult(t, fmt.Sprintf("%v", expected), fmt.Sprintf("%v", result))
}
