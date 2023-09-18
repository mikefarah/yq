//go:build !yq_nojson

package yqlib

import (
	"fmt"
	"io"

	"github.com/goccy/go-json"
	yaml "gopkg.in/yaml.v3"
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

	var dataBucket orderedMap
	log.Debug("going to decode")
	err := dec.decoder.Decode(&dataBucket)
	if err != nil {
		return nil, err
	}
	node, err := dec.convertToYamlNode(&dataBucket)

	if err != nil {
		return nil, err
	}

	return &CandidateNode{
		Node: &yaml.Node{
			Kind:    yaml.DocumentNode,
			Content: []*yaml.Node{node},
		},
	}, nil
}

func (dec *jsonDecoder) convertToYamlNode(data *orderedMap) (*yaml.Node, error) {
	if data == nil {
		return createScalarNode(nil, "null"), nil
	}
	if data.kv == nil {
		switch rawData := data.altVal.(type) {
		case nil:
			return createScalarNode(nil, "null"), nil
		case float64, float32:
			// json decoder returns ints as float.'
			intNum := int(rawData.(float64))

			// if the integer representation is the same as the original
			// then its an int.
			if float64(intNum) == rawData.(float64) {
				return createScalarNode(intNum, fmt.Sprintf("%v", intNum)), nil
			}

			return createScalarNode(rawData, fmt.Sprintf("%v", rawData)), nil
		case int, int64, int32, string, bool:
			return createScalarNode(rawData, fmt.Sprintf("%v", rawData)), nil
		case []*orderedMap:
			return dec.parseArray(rawData)
		default:
			return nil, fmt.Errorf("unrecognised type :( %v", rawData)
		}
	}

	var yamlMap = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	for i, keyValuePair := range data.kv {
		yamlValue, err := dec.convertToYamlNode(&data.kv[i].V)
		if err != nil {
			return nil, err
		}
		yamlMap.Content = append(yamlMap.Content, createScalarNode(keyValuePair.K, keyValuePair.K), yamlValue)
	}
	return yamlMap, nil

}

func (dec *jsonDecoder) parseArray(dataArray []*orderedMap) (*yaml.Node, error) {

	var yamlMap = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}

	for _, value := range dataArray {
		yamlValue, err := dec.convertToYamlNode(value)
		if err != nil {
			return nil, err
		}
		yamlMap.Content = append(yamlMap.Content, yamlValue)
	}
	return yamlMap, nil
}
