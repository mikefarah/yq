package yqlib

import (
	"fmt"
	"testing"

	"github.com/goccy/go-yaml/parser"
	"gopkg.in/yaml.v3"
)

func TestMarshalGoccyYAML(t *testing.T) {
	input := `
a:
  b: 2
  c: &anc
    d: !mytag ef
    e: 3.0
s: 
- 1
- 2
t: [5, five]
f: [6, {y: true}]
`

	goccyAst, err := parser.ParseBytes([]byte(input), parser.ParseComments)
	fmt.Println(goccyAst)

	var yamlNode yaml.Node
	err = yaml.Unmarshal([]byte(input), &yamlNode)
	if err != nil {
		t.Error(err)
		return
	}

	candidate := &CandidateNode{}
	err = candidate.UnmarshalYAML(yamlNode.Content[0], make(map[string]*CandidateNode))
	if err != nil {
		t.Error(err)
		return
	}

	goccyNode, err := candidate.MarshalGoccyYAML()
	if err != nil {
		t.Error(err)
		return
	}

	parsed, err := parser.ParseBytes([]byte(input), parser.ParseComments)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(parsed)
	fmt.Println(goccyNode)
}
