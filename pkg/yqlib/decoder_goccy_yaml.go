//go:build !yq_noyaml

package yqlib

import (
	"io"

	yaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

type goccyYamlDecoder struct {
	decoder yaml.Decoder
	cm      yaml.CommentMap
}

func NewGoccyYAMLDecoder() Decoder {
	return &goccyYamlDecoder{}
}

func (dec *goccyYamlDecoder) Init(reader io.Reader) error {
	dec.cm = yaml.CommentMap{}
	dec.decoder = *yaml.NewDecoder(reader, yaml.CommentToMap(dec.cm))
	return nil
}

func (dec *goccyYamlDecoder) Decode() (*CandidateNode, error) {

	var ast ast.Node

	err := dec.decoder.Decode(&ast)
	if err != nil {
		return nil, err
	}

	candidateNode := &CandidateNode{}
	candidateNode.UnmarshalGoccyYAML(ast, dec.cm)

	return candidateNode, nil
}
