package yqlib

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	yaml "github.com/mikefarah/yaml/v2"
	logging "gopkg.in/op/go-logging.v1"
)

type DataNavigator interface {
	ReadChildValue(child interface{}, remainingPaths []string) (interface{}, error)
	UpdatedChildValue(child interface{}, remainingPaths []string, value interface{}) interface{}
	DeleteChildValue(child interface{}, remainingPaths []string) (interface{}, error)
}

type navigator struct {
	log *logging.Logger
}

func NewDataNavigator(l *logging.Logger) DataNavigator {
	return &navigator {
		log: l,
	}
}

func (n *navigator) ReadChildValue(child interface{}, remainingPaths []string) (interface{}, error) {
	if len(remainingPaths) == 0 {
		return child, nil
	}
	return n.recurse(child, remainingPaths[0], remainingPaths[1:])
}

func (n *navigator) UpdatedChildValue(child interface{}, remainingPaths []string, value interface{}) interface{} {
	if len(remainingPaths) == 0 {
		return value
	}
	n.log.Debugf("UpdatedChildValue for child %v with path %v to set value %v", child, remainingPaths, value)
	n.log.Debugf("type of child is %v", reflect.TypeOf(child))

	switch child := child.(type) {
	case nil:
		if remainingPaths[0] == "+" || remainingPaths[0] == "*" {
			return n.writeArray(child, remainingPaths, value)
		}
	case []interface{}:
		_, nextIndexErr := strconv.ParseInt(remainingPaths[0], 10, 64)
		arrayCommand := nextIndexErr == nil || remainingPaths[0] == "+" || remainingPaths[0] == "*"
		if arrayCommand {
			return n.writeArray(child, remainingPaths, value)
		}
	}
	return n.writeMap(child, remainingPaths, value)
}

func (n *navigator) DeleteChildValue(child interface{}, remainingPaths []string) (interface{}, error) {
	n.log.Debugf("DeleteChildValue for %v for %v\n", remainingPaths, child)
	if len(remainingPaths) == 0 {
		return child, nil
	}
	var head = remainingPaths[0]
	var tail = remainingPaths[1:]
	switch child := child.(type) {
	case yaml.MapSlice:
		return n.deleteMap(child, remainingPaths)
	case []interface{}:
		if head == "*" {
			return n.deleteArraySplat(child, tail)
		}
		index, err := strconv.ParseInt(head, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error accessing array: %v", err)
		}
		return n.deleteArray(child, remainingPaths, index)
	}
	return child, nil
}

func (n *navigator) recurse(value interface{}, head string, tail []string) (interface{}, error) {
	switch value := value.(type) {
	case []interface{}:
		if head == "*" {
			return n.readArraySplat(value, tail)
		}
		index, err := strconv.ParseInt(head, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error accessing array: %v", err)
		}
		return n.readArray(value, index, tail)
	case yaml.MapSlice:
		return n.readMap(value, head, tail)
	default:
		return nil, nil
	}
}

func (n *navigator) matchesKey(key string, actual interface{}) bool {
	var actualString = fmt.Sprintf("%v", actual)
	var prefixMatch = strings.TrimSuffix(key, "*")
	if prefixMatch != key {
		return strings.HasPrefix(actualString, prefixMatch)
	}
	return actualString == key
}

func (n *navigator) entriesInSlice(context yaml.MapSlice, key string) []*yaml.MapItem {
	var matches = make([]*yaml.MapItem, 0)
	for idx := range context {
		var entry = &context[idx]
		if n.matchesKey(key, entry.Key) {
			matches = append(matches, entry)
		}
	}
	return matches
}

func (n *navigator) getMapSlice(context interface{}) yaml.MapSlice {
	var mapSlice yaml.MapSlice
	switch context := context.(type) {
	case yaml.MapSlice:
		mapSlice = context
	default:
		mapSlice = make(yaml.MapSlice, 0)
	}
	return mapSlice
}

func (n *navigator) getArray(context interface{}) (array []interface{}, ok bool) {
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

func (n *navigator) writeMap(context interface{}, paths []string, value interface{}) interface{} {
	n.log.Debugf("writeMap with path %v for %v to set value %v\n", paths, context, value)

	mapSlice := n.getMapSlice(context)

	if len(paths) == 0 {
		return context
	}

	children := n.entriesInSlice(mapSlice, paths[0])

	if len(children) == 0 && paths[0] == "*" {
		n.log.Debugf("\tNo matches, return map as is")
		return context
	}

	if len(children) == 0 {
		newChild := yaml.MapItem{Key: paths[0]}
		mapSlice = append(mapSlice, newChild)
		children = n.entriesInSlice(mapSlice, paths[0])
		n.log.Debugf("\tAppended child at %v for mapSlice %v\n", paths[0], mapSlice)
	}

	remainingPaths := paths[1:]
	for _, child := range children {
		child.Value = n.UpdatedChildValue(child.Value, remainingPaths, value)
	}
	n.log.Debugf("\tReturning mapSlice %v\n", mapSlice)
	return mapSlice
}

func (n *navigator) writeArray(context interface{}, paths []string, value interface{}) []interface{} {
	n.log.Debugf("writeArray with path %v for %v to set value %v\n", paths, context, value)
	array, _ := n.getArray(context)

	if len(paths) == 0 {
		return array
	}

	n.log.Debugf("\tarray %v\n", array)

	rawIndex := paths[0]
	remainingPaths := paths[1:]
	var index int64
	// the append array indicator
	if rawIndex == "+" {
		index = int64(len(array))
	} else if rawIndex == "*" {
		for index, oldChild := range array {
			array[index] = n.UpdatedChildValue(oldChild, remainingPaths, value)
		}
		return array
	} else {
		index, _ = strconv.ParseInt(rawIndex, 10, 64) // nolint
		// writeArray is only called by UpdatedChildValue which handles parsing the
		// index, as such this renders this dead code.
	}

	for index >= int64(len(array)) {
		array = append(array, nil)
	}
	currentChild := array[index]

	n.log.Debugf("\tcurrentChild %v\n", currentChild)

	array[index] = n.UpdatedChildValue(currentChild, remainingPaths, value)
	n.log.Debugf("\tReturning array %v\n", array)
	return array
}

func (n *navigator) readMap(context yaml.MapSlice, head string, tail []string) (interface{}, error) {
	n.log.Debugf("readingMap %v with key %v\n", context, head)
	if head == "*" {
		return n.readMapSplat(context, tail)
	}

	entries := n.entriesInSlice(context, head)
	if len(entries) == 1 {
		return n.calculateValue(entries[0].Value, tail)
	} else if len(entries) == 0 {
		return nil, nil
	}
	var errInIdx error
	values := make([]interface{}, len(entries))
	for idx, entry := range entries {
		values[idx], errInIdx = n.calculateValue(entry.Value, tail)
		if errInIdx != nil {
			n.log.Errorf("Error updating index %v in %v", idx, context)
			return nil, errInIdx
		}

	}
	return values, nil
}

func (n *navigator) readMapSplat(context yaml.MapSlice, tail []string) (interface{}, error) {
	var newArray = make([]interface{}, len(context))
	var i = 0
	for _, entry := range context {
		if len(tail) > 0 {
			val, err := n.recurse(entry.Value, tail[0], tail[1:])
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

func (n *navigator) readArray(array []interface{}, head int64, tail []string) (interface{}, error) {
	if head >= int64(len(array)) {
		return nil, nil
	}

	value := array[head]
	return n.calculateValue(value, tail)
}

func (n *navigator) readArraySplat(array []interface{}, tail []string) (interface{}, error) {
	var newArray = make([]interface{}, len(array))
	for index, value := range array {
		val, err := n.calculateValue(value, tail)
		if err != nil {
			return nil, err
		}
		newArray[index] = val
	}
	return newArray, nil
}

func (n *navigator) calculateValue(value interface{}, tail []string) (interface{}, error) {
	if len(tail) > 0 {
		return n.recurse(value, tail[0], tail[1:])
	}
	return value, nil
}

func (n *navigator) deleteMap(context interface{}, paths []string) (yaml.MapSlice, error) {
	n.log.Debugf("deleteMap for %v for %v\n", paths, context)

	mapSlice := n.getMapSlice(context)

	if len(paths) == 0 {
		return mapSlice, nil
	}

	var index int
	var child yaml.MapItem
	for index, child = range mapSlice {
		if n.matchesKey(paths[0], child.Key) {
			n.log.Debugf("\tMatched [%v] with [%v] at index %v", paths[0], child.Key, index)
			var badDelete error
			mapSlice, badDelete = n.deleteEntryInMap(mapSlice, child, index, paths)
			if badDelete != nil {
				return nil, badDelete
			}
		}
	}

	return mapSlice, nil

}

func (n *navigator) deleteEntryInMap(original yaml.MapSlice, child yaml.MapItem, index int, paths []string) (yaml.MapSlice, error) {
	remainingPaths := paths[1:]

	var newSlice yaml.MapSlice
	if len(remainingPaths) > 0 {
		newChild := yaml.MapItem{Key: child.Key}
		var errorDeleting error
		newChild.Value, errorDeleting = n.DeleteChildValue(child.Value, remainingPaths)
		if errorDeleting != nil {
			return nil, errorDeleting
		}

		newSlice = make(yaml.MapSlice, len(original))
		for i := range original {
			item := original[i]
			if i == index {
				item = newChild
			}
			newSlice[i] = item
		}
	} else {
		// Delete item from slice at index
		newSlice = append(original[:index], original[index+1:]...)
		n.log.Debugf("\tDeleted item index %d from original", index)
	}

	n.log.Debugf("\tReturning original %v\n", original)
	return newSlice, nil
}

func (n *navigator) deleteArraySplat(array []interface{}, tail []string) (interface{}, error) {
	n.log.Debugf("deleteArraySplat for %v for %v\n", tail, array)
	var newArray = make([]interface{}, len(array))
	for index, value := range array {
		val, err := n.DeleteChildValue(value, tail)
		if err != nil {
			return nil, err
		}
		newArray[index] = val
	}
	return newArray, nil
}

func (n *navigator) deleteArray(array []interface{}, paths []string, index int64) (interface{}, error) {
	n.log.Debugf("deleteArray for %v for %v\n", paths, array)

	if index >= int64(len(array)) {
		return array, nil
	}

	remainingPaths := paths[1:]
	if len(remainingPaths) > 0 {
		// recurse into the array element at index
		var errorDeleting error
		array[index], errorDeleting = n.deleteMap(array[index], remainingPaths)
		if errorDeleting != nil {
			return nil, errorDeleting
		}

	} else {
		// Delete the array element at index
		array = append(array[:index], array[index+1:]...)
		n.log.Debugf("\tDeleted item index %d from array, leaving %v", index, array)
	}

	n.log.Debugf("\tReturning array: %v\n", array)
	return array, nil
}
