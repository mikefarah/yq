//go:build !yq_noyaml

package yqlib

import (
	"io"

	yaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

type goccyYamlDecoder struct {
	decoder yaml.Decoder
}

func NewGoccyYAMLDecoder() Decoder {
	return &goccyYamlDecoder{}
}

func (dec *goccyYamlDecoder) Init(reader io.Reader) error {
	dec.decoder = *yaml.NewDecoder(reader)
	return nil
}

func (dec *goccyYamlDecoder) Decode() (*CandidateNode, error) {

	var ast ast.Node

	err := dec.decoder.Decode(&ast)
	if err != nil {
		log.Debug("badasda: %v", err)

		return nil, err
	}

	log.Debug("ASTasdasdadasd: %v", ast.Type().String())

	candidateNode := &CandidateNode{}
	candidateNode.UnmarshalGoccyYAML(ast, nil)

	return candidateNode, nil
}
