//go:build !yq_nojson

package yqlib

import (
	"fmt"
	"io"

	"github.com/goccy/go-json"
)

type jsonDecoder struct {
	decoder json.Decoder
}

func NewJSONDecoder() Decoder {
	return &jsonDecoder{}
}

func (dec *jsonDecoder) Init(reader io.Reader) error {
	dec.decoder = *json.NewDecoder(reader)
	return nil
}

func (dec *jsonDecoder) Decode() (*CandidateNode, error) {

	var dataBucket CandidateNode
	err := dec.decoder.Decode(&dataBucket)
	if err != nil {
		return nil, err
	}
	// log.Debug("convert to yaml")
	// node, err := dec.convertToYamlNode(&dataBucket)
	// log.Debug("done, %w", err)
	// if err != nil {
	// 	return nil, err
	// }

	return &CandidateNode{
		Kind:    DocumentNode,
		Content: []*CandidateNode{&dataBucket},
	}, nil
}

func (dec *jsonDecoder) convertToYamlNode(data *orderedMap) (*CandidateNode, error) {
	if data.kv == nil {
		switch rawData := data.altVal.(type) {
		case nil:
			return createScalarNode(nil, "null"), nil
		case float64, float32:
			// json decoder returns ints as float.
			return parseSnippet(fmt.Sprintf("%v", rawData))
		case int, int64, int32, string, bool:
			return createScalarNode(rawData, fmt.Sprintf("%v", rawData)), nil
		case []*orderedMap:
			return dec.parseArray(rawData)
		default:
			return nil, fmt.Errorf("unrecognised type :( %v", rawData)
		}
	}

	var yamlMap = &CandidateNode{Kind: MappingNode, Tag: "!!map"}
	for _, keyValuePair := range data.kv {
		yamlValue, err := dec.convertToYamlNode(&keyValuePair.V)
		if err != nil {
			return nil, err
		}
		yamlMap.Content = append(yamlMap.Content, createScalarNode(keyValuePair.K, keyValuePair.K), yamlValue)
	}
	return yamlMap, nil

}

func (dec *jsonDecoder) parseArray(dataArray []*orderedMap) (*CandidateNode, error) {

	var yamlMap = &CandidateNode{Kind: SequenceNode, Tag: "!!seq"}

	for _, value := range dataArray {
		yamlValue, err := dec.convertToYamlNode(value)
		if err != nil {
			return nil, err
		}
		yamlMap.Content = append(yamlMap.Content, yamlValue)
	}
	return yamlMap, nil
}
