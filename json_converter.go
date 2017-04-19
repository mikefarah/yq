package main

import (
	"encoding/json"
	"github.com/mikefarah/yaml/Godeps/_workspace/src/gopkg.in/yaml.v2"
)

func jsonToString(context interface{}) string {
	out, err := json.Marshal(toJSON(context))
	if err != nil {
		die("error printing yaml as json: ", err)
	}
	return string(out)
}

func toJSON(context interface{}) interface{} {
	switch context.(type) {
	case []interface{}:
		oldArray := context.([]interface{})
		newArray := make([]interface{}, len(oldArray))
		for index, value := range oldArray {
			newArray[index] = toJSON(value)
		}
		return newArray
	case yaml.MapSlice:
		oldMap := context.(yaml.MapSlice)
		newMap := make(map[string]interface{})
		for _, entry := range oldMap {
			newMap[entry.Key.(string)] = toJSON(entry.Value)
		}
		return newMap
	default:
		return context
	}
}
