package main

import (
	"fmt"
	"strconv"

	yaml "github.com/mikefarah/yaml"
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

func getMapSlice(context interface{}) yaml.MapSlice {
	var mapSlice yaml.MapSlice
	switch context := context.(type) {
	case yaml.MapSlice:
		mapSlice = context
	default:
		mapSlice = make(yaml.MapSlice, 0)
	}
	return mapSlice
}

func getArray(context interface{}) (array []interface{}, ok bool) {
	switch context := context.(type) {
	case []interface{}:
		array = context
		ok = true
	default:
		array = make([]interface{}, 0)
		ok = false
	}
	return
}

func writeMap(context interface{}, paths []string, value interface{}) yaml.MapSlice {
	log.Debugf("writeMap for %v for %v with value %v\n", paths, context, value)

	mapSlice := getMapSlice(context)

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
	array, _ := getArray(context)

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
		index, _ = strconv.ParseInt(rawIndex, 10, 64) // nolint
		// writeArray is only called by updatedChildValue which handles parsing the
		// index, as such this renders this dead code.
	}

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
	switch value := value.(type) {
	case []interface{}:
		if head == "*" {
			return readArraySplat(value, tail)
		}
		index, err := strconv.ParseInt(head, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error accessing array: %v", err)
		}
		return readArray(value, index, tail)
	case yaml.MapSlice:
		return readMap(value, head, tail)
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

func deleteMap(context interface{}, paths []string) yaml.MapSlice {
	log.Debugf("deleteMap for %v for %v\n", paths, context)

	mapSlice := getMapSlice(context)

	if len(paths) == 0 {
		return mapSlice
	}

	var found bool
	var index int
	var child yaml.MapItem
	for index, child = range mapSlice {
		if child.Key == paths[0] {
			found = true
			break
		}
	}

	if !found {
		return mapSlice
	}

	remainingPaths := paths[1:]

	var newSlice yaml.MapSlice
	if len(remainingPaths) > 0 {
		newChild := yaml.MapItem{Key: child.Key}
		newChild.Value = deleteChildValue(child.Value, remainingPaths)

		newSlice = make(yaml.MapSlice, len(mapSlice))
		for i := range mapSlice {
			item := mapSlice[i]
			if i == index {
				item = newChild
			}
			newSlice[i] = item
		}
	} else {
		// Delete item from slice at index
		newSlice = append(mapSlice[:index], mapSlice[index+1:]...)
		log.Debugf("\tDeleted item index %d from mapSlice", index)
	}

	log.Debugf("\t\tlen: %d\tcap: %d\tslice: %v", len(mapSlice), cap(mapSlice), mapSlice)
	log.Debugf("\tReturning mapSlice %v\n", mapSlice)
	return newSlice
}

func deleteArray(context interface{}, paths []string, index int64) interface{} {
	log.Debugf("deleteArray for %v for %v\n", paths, context)

	array, ok := getArray(context)
	if !ok {
		// did not get an array
		return context
	}

	if index >= int64(len(array)) {
		return array
	}

	remainingPaths := paths[1:]
	if len(remainingPaths) > 0 {
		// Recurse into the array element at index
		array[index] = deleteMap(array[index], remainingPaths)
	} else {
		// Delete the array element at index
		array = append(array[:index], array[index+1:]...)
		log.Debugf("\tDeleted item index %d from array, leaving %v", index, array)
	}

	log.Debugf("\tReturning array: %v\n", array)
	return array
}

func deleteChildValue(child interface{}, remainingPaths []string) interface{} {
	log.Debugf("deleteChildValue for %v for %v\n", remainingPaths, child)

	idx, nextIndexErr := strconv.ParseInt(remainingPaths[0], 10, 64)
	if nextIndexErr != nil {
		// must be a map
		log.Debugf("\tdetected a map, invoking deleteMap\n")
		return deleteMap(child, remainingPaths)
	}

	log.Debugf("\tdetected an array, so traversing element with index %d\n", idx)
	return deleteArray(child, remainingPaths, idx)
}
