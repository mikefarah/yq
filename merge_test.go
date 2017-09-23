package main

import (
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestMerge(t *testing.T) {
	result, _ := mergeYaml([]string{"examples/data1.yaml", "examples/data2.yaml", "examples/data3.yaml"})
	expected := yaml.MapSlice{
		yaml.MapItem{Key: "a", Value: "simple"},
		yaml.MapItem{Key: "b", Value: []interface{}{1, 2}},
		yaml.MapItem{Key: "c", Value: yaml.MapSlice{yaml.MapItem{Key: "other", Value: true}, yaml.MapItem{Key: "test", Value: 1}}},
		yaml.MapItem{Key: "d", Value: false},
	}
	assertResultComplex(t, expected, result)
}

func TestMergeWithOverwrite(t *testing.T) {
	overwriteFlag = true
	result, _ := mergeYaml([]string{"examples/data1.yaml", "examples/data2.yaml", "examples/data3.yaml"})
	expected := yaml.MapSlice{
		yaml.MapItem{Key: "a", Value: "other"},
		yaml.MapItem{Key: "b", Value: []interface{}{2, 3, 4}},
		yaml.MapItem{Key: "c", Value: yaml.MapSlice{yaml.MapItem{Key: "other", Value: true}, yaml.MapItem{Key: "test", Value: 2}}},
		yaml.MapItem{Key: "d", Value: false},
	}
	assertResultComplex(t, expected, result)
}
