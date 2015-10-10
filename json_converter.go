package main

import (
	"encoding/json"
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
	case map[interface{}]interface{}:
		oldMap := context.(map[interface{}]interface{})
		newMap := make(map[string]interface{})
		for key, value := range oldMap {
			newMap[key.(string)] = toJSON(value)
		}
		return newMap
	default:
		return context
	}
}
