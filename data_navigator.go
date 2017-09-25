package main

import (
	"fmt"
	"sort"
	"strconv"

	"gopkg.in/yaml.v2"
)

func entryInSlice(context yaml.MapSlice, key interface{}) *yaml.MapItem {
	for idx := range context {
		var entry = &context[idx]
		if entry.Key == key {
			return entry
		}
	}
	return nil
}

func writeMap(context interface{}, paths []string, value interface{}) yaml.MapSlice {
	log.Debugf("writeMap for %v for %v with value %v\n", paths, context, value)

	var mapSlice yaml.MapSlice
	switch context.(type) {
	case yaml.MapSlice:
		mapSlice = context.(yaml.MapSlice)
	default:
		mapSlice = make(yaml.MapSlice, 0)
	}

	if len(paths) == 0 {
		return mapSlice
	}

	child := entryInSlice(mapSlice, paths[0])
	if child == nil {
		newChild := yaml.MapItem{Key: paths[0]}
		mapSlice = append(mapSlice, newChild)
		child = entryInSlice(mapSlice, paths[0])
		log.Debugf("\tAppended child at %v for mapSlice %v\n", paths[0], mapSlice)
	}

	log.Debugf("\tchild.Value %v\n", child.Value)

	remainingPaths := paths[1:]
	child.Value = updatedChildValue(child.Value, remainingPaths, value)
	log.Debugf("\tReturning mapSlice %v\n", mapSlice)
	return mapSlice
}

func updatedChildValue(child interface{}, remainingPaths []string, value interface{}) interface{} {
	if len(remainingPaths) == 0 {
		return value
	}

	_, nextIndexErr := strconv.ParseInt(remainingPaths[0], 10, 64)
	if nextIndexErr != nil && remainingPaths[0] != "+" {
		// must be a map
		return writeMap(child, remainingPaths, value)
	}

	// must be an array
	return writeArray(child, remainingPaths, value)
}

func writeArray(context interface{}, paths []string, value interface{}) []interface{} {
	log.Debugf("writeArray for %v for %v with value %v\n", paths, context, value)
	var array []interface{}
	switch context.(type) {
	case []interface{}:
		array = context.([]interface{})
	default:
		array = make([]interface{}, 0)
	}

	if len(paths) == 0 {
		return array
	}

	log.Debugf("\tarray %v\n", array)

	rawIndex := paths[0]
	var index int64
	// the append array indicator
	if rawIndex == "+" {
		index = int64(len(array))
	} else {
		index, _ = strconv.ParseInt(rawIndex, 10, 64)
	}
	// writeArray is only called by updatedChildValue which handles parsing the
	// index, as such this renders this dead code.
	// if err != nil {
	// 	return array, fmt.Errorf("Error accessing array: %v", err)
	// }
	for index >= int64(len(array)) {
		array = append(array, nil)
	}
	currentChild := array[index]

	log.Debugf("\tcurrentChild %v\n", currentChild)

	remainingPaths := paths[1:]
	array[index] = updatedChildValue(currentChild, remainingPaths, value)
	log.Debugf("\tReturning array %v\n", array)
	return array
}

func readMap(context yaml.MapSlice, head string, tail []string) (interface{}, error) {
	if head == "*" {
		return readMapSplat(context, tail)
	}
	var value interface{}

	entry := entryInSlice(context, head)
	if entry != nil {
		value = entry.Value
	}
	return calculateValue(value, tail)
}

func readMapSplat(context yaml.MapSlice, tail []string) (interface{}, error) {
	var newArray = make([]interface{}, len(context))
	var i = 0
	for _, entry := range context {
		if len(tail) > 0 {
			val, err := recurse(entry.Value, tail[0], tail[1:])
			if err != nil {
				return nil, err
			}
			newArray[i] = val
		} else {
			newArray[i] = entry.Value
		}
		i++
	}
	return newArray, nil
}

func recurse(value interface{}, head string, tail []string) (interface{}, error) {
	switch value.(type) {
	case []interface{}:
		if head == "*" {
			return readArraySplat(value.([]interface{}), tail)
		}
		index, err := strconv.ParseInt(head, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Error accessing array: %v", err)
		}
		return readArray(value.([]interface{}), index, tail)
	case yaml.MapSlice:
		return readMap(value.(yaml.MapSlice), head, tail)
	default:
		return nil, nil
	}
}

func readArray(array []interface{}, head int64, tail []string) (interface{}, error) {
	if head >= int64(len(array)) {
		return nil, nil
	}

	value := array[head]
	return calculateValue(value, tail)
}

func readArraySplat(array []interface{}, tail []string) (interface{}, error) {
	var newArray = make([]interface{}, len(array))
	for index, value := range array {
		val, err := calculateValue(value, tail)
		if err != nil {
			return nil, err
		}
		newArray[index] = val
	}
	return newArray, nil
}

func calculateValue(value interface{}, tail []string) (interface{}, error) {
	if len(tail) > 0 {
		return recurse(value, tail[0], tail[1:])
	}
	return value, nil
}

func mapToMapSlice(data map[interface{}]interface{}) yaml.MapSlice {
	var mapSlice yaml.MapSlice

	for k, v := range data {
		if mv, ok := v.(map[interface{}]interface{}); ok {
			v = mapToMapSlice(mv)
		}
		item := yaml.MapItem{Key: k, Value: v}
		mapSlice = append(mapSlice, item)
	}

	// because the parsing of the yaml was done via a map the order will be inconsistent
	// apply order to allow a consistent output
	// Only available in Go 1.8+
	// sort.SliceStable(mapSlice, func(i, j int) bool { return mapSlice[i].Key < mapSlice[j].Key })
	sort.Sort(sortMap{mapSlice})
	return mapSlice
}

type sortMap struct {
	ms yaml.MapSlice
}

func (m sortMap) Len() int {
	return len(m.ms)
}

func (m sortMap) Less(i, j int) bool {
	return m.ms[i].Key.(string) < m.ms[j].Key.(string)
}

func (m sortMap) Swap(i, j int) {
	m.ms[i], m.ms[j] = m.ms[j], m.ms[i]
}
