package main

import (
	"fmt"
	"github.com/mikefarah/yaml/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/op/go-logging"
	"os"
	"sort"
	"testing"
)

func TestMain(m *testing.M) {
	backend.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend)
	os.Exit(m.Run())
}

func TestReadMap_simple(t *testing.T) {
	var data = parseData(`
---
b:
  c: 2
`)
	assertResult(t, 2, readMap(data, "b", []string{"c"}))
}

func TestReadMap_splat(t *testing.T) {
	var data = parseData(`
---
mapSplat:
  item1: things
  item2: whatever
`)
	var result = readMap(data, "mapSplat", []string{"*"}).([]interface{})
	var actual = []string{result[0].(string), result[1].(string)}
	sort.Strings(actual)
	assertResult(t, "[things whatever]", fmt.Sprintf("%v", actual))
}

func TestReadMap_deep_splat(t *testing.T) {
	var data = parseData(`
---
mapSplatDeep:
  item1:
    cats: bananas
  item2:
    cats: apples
`)

	var result = readMap(data, "mapSplatDeep", []string{"*", "cats"}).([]interface{})
	var actual = []string{result[0].(string), result[1].(string)}
	sort.Strings(actual)
	assertResult(t, "[apples bananas]", fmt.Sprintf("%v", actual))
}

func TestReadMap_key_doesnt_exist(t *testing.T) {
	var data = parseData(`
---
b:
  c: 2
`)
	assertResult(t, nil, readMap(data, "b.x.f", []string{"c"}))
}

func TestReadMap_recurse_against_string(t *testing.T) {
	var data = parseData(`
---
a: cat
`)
	assertResult(t, nil, readMap(data, "a", []string{"b"}))
}

func TestReadMap_with_array(t *testing.T) {
	var data = parseData(`
---
b:
  d:
    - 3
    - 4
`)
	assertResult(t, 4, readMap(data, "b", []string{"d", "1"}))
}

func TestReadMap_with_array_out_of_bounds(t *testing.T) {
	var data = parseData(`
---
b:
  d:
    - 3
    - 4
`)
	assertResult(t, nil, readMap(data, "b", []string{"d", "3"}))
}

func TestReadMap_with_array_out_of_bounds_by_1(t *testing.T) {
	var data = parseData(`
---
b:
  d:
    - 3
    - 4
`)
	assertResult(t, nil, readMap(data, "b", []string{"d", "2"}))
}

func TestReadMap_with_array_splat(t *testing.T) {
	var data = parseData(`
e:
  -
    name: Fred
    thing: cat
  -
    name: Sam
    thing: dog
`)
	assertResult(t, "[Fred Sam]", fmt.Sprintf("%v", readMap(data, "e", []string{"*", "name"})))
}

func TestWrite_really_simple(t *testing.T) {
	var data = parseData(`
    b: 2
`)

	updated := writeMap(data, []string{"b"}, "4")
	b := entryInSlice(updated, "b").Value
	assertResult(t, "4", b)
}

func TestWrite_simple(t *testing.T) {
	var data = parseData(`
b:
  c: 2
`)

	updated := writeMap(data, []string{"b", "c"}, "4")
	b := entryInSlice(updated, "b").Value.(yaml.MapSlice)
	c := entryInSlice(b, "c").Value
	assertResult(t, "4", c)
}

func TestWrite_new(t *testing.T) {
	var data = parseData(`
b:
  c: 2
`)

	updated := writeMap(data, []string{"b", "d"}, "4")
	b := entryInSlice(updated, "b").Value.(yaml.MapSlice)
	d := entryInSlice(b, "d").Value
	assertResult(t, "4", d)
}

func TestWrite_new_deep(t *testing.T) {
	var data = parseData(`
b:
  c: 2
`)

	updated := writeMap(data, []string{"b", "d", "f"}, "4")
	assertResult(t, "4", readMap(updated, "b", []string{"d", "f"}))
}

func TestWrite_array(t *testing.T) {
	var data = parseData(`
b:
  - aa
`)

	updated := writeMap(data, []string{"b", "0"}, "bb")

	b := entryInSlice(updated, "b").Value.([]interface{})
	assertResult(t, "bb", b[0].(string))
}

func TestWrite_new_array(t *testing.T) {
	var data = parseData(`
b:
  c: 2
`)

	updated := writeMap(data, []string{"b", "0"}, "4")
	assertResult(t, "4", readMap(updated, "b", []string{"0"}))
}

func TestWrite_new_array_deep(t *testing.T) {
	var data = parseData(`
b:
  c: 2
`)

	var expected = `b:
- c: "4"`

	updated := writeMap(data, []string{"b", "0", "c"}, "4")
	assertResult(t, expected, yamlToString(updated))
}

func TestWrite_new_map_array_deep(t *testing.T) {
	var data = parseData(`
b:
  c: 2
`)

	updated := writeMap(data, []string{"b", "d", "0"}, "4")
	assertResult(t, "4", readMap(updated, "b", []string{"d", "0"}))
}

func TestWrite_with_no_tail(t *testing.T) {
	var data = parseData(`
b:
  c: 2
`)
	updated := writeMap(data, []string{"b"}, "4")

	b := entryInSlice(updated, "b").Value
	assertResult(t, "4", fmt.Sprintf("%v", b))
}
