package marshal

import (
	"encoding/json"
	"fmt"
	"strconv"

	yaml "github.com/mikefarah/yaml/v2"
)

type JsonConverter interface {
	JsonToString(context interface{}) (string, error)
}

type jsonConverter struct {}

func NewJsonConverter() JsonConverter {
	return &jsonConverter{}
}

func (j *jsonConverter) JsonToString(context interface{}) (string, error) {
	out, err := json.Marshal(j.toJSON(context))
	if err != nil {
		return "", fmt.Errorf("error printing yaml as json: %v", err)
	}
	return string(out), nil
}

func (j *jsonConverter) toJSON(context interface{}) interface{} {
	switch context := context.(type) {
	case []interface{}:
		oldArray := context
		newArray := make([]interface{}, len(oldArray))
		for index, value := range oldArray {
			newArray[index] = j.toJSON(value)
		}
		return newArray
	case yaml.MapSlice:
		oldMap := context
		newMap := make(map[string]interface{})
		for _, entry := range oldMap {
			if str, ok := entry.Key.(string); ok {
				newMap[str] = j.toJSON(entry.Value)
			} else if i, ok := entry.Key.(int); ok {
				newMap[strconv.Itoa(i)] = j.toJSON(entry.Value)
			} else if b, ok := entry.Key.(bool); ok {
				newMap[strconv.FormatBool(b)] = j.toJSON(entry.Value)
			}
		}
		return newMap
	default:
		return context
	}
}
