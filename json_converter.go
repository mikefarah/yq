package main

import (
	"encoding/json"
	"github.com/mikefarah/yaml/Godeps/_workspace/src/gopkg.in/yaml.v2"
)

func fromJSONBytes(jsonBytes []byte, parsedData *map[interface{}]interface{}) {
	*parsedData = make(map[interface{}]interface{})
	var jsonData map[string]interface{}
	err := json.Unmarshal(jsonBytes, &jsonData)
	if err != nil {
		die("error parsing data: ", err)
	}

	for key, value := range jsonData {
		(*parsedData)[key] = fromJSON(value)
	}
}

func jsonToString(context interface{}) string {
	out, err := json.Marshal(toJSON(context))
	if err != nil {
		die("error printing yaml as json: ", err)
	}
	return string(out)
}

func fromJSON(context interface{}) interface{} {
	switch context.(type) {
	case []interface{}:
		oldArray := context.([]interface{})
		newArray := make([]interface{}, len(oldArray))
		for index, value := range oldArray {
			newArray[index] = fromJSON(value)
		}
		return newArray
	case map[string]interface{}:
		oldMap := context.(map[string]interface{})
		newMap := make(map[interface{}]interface{})
		for key, value := range oldMap {
			newMap[key] = fromJSON(value)
		}
		return newMap
	default:
		return context
	}
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
