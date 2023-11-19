//go:build !yq_noyaml

package yqlib

import (
	"bytes"
	"errors"
	"io"
	"regexp"

	yaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

type goccyYamlDecoder struct {
	decoder      yaml.Decoder
	cm           yaml.CommentMap
	bufferRead   bytes.Buffer
	readAnything bool
	anchorMap    map[string]*CandidateNode
}

func NewGoccyYAMLDecoder() Decoder {
	return &goccyYamlDecoder{}
}

func (dec *goccyYamlDecoder) Init(reader io.Reader) error {
	dec.cm = yaml.CommentMap{}
	dec.readAnything = false
	dec.anchorMap = make(map[string]*CandidateNode)
	readerToUse := io.TeeReader(reader, &dec.bufferRead)
	dec.decoder = *yaml.NewDecoder(readerToUse, yaml.CommentToMap(dec.cm), yaml.UseOrderedMap())
	return nil
}

func (dec *goccyYamlDecoder) Decode() (*CandidateNode, error) {

	var commentLineRegEx = regexp.MustCompile(`^\s*#`)

	var ast ast.Node

	err := dec.decoder.Decode(&ast)
	if errors.Is(err, io.EOF) && !dec.readAnything {

		content := dec.bufferRead.String()
		// only null fix
		if content == "null" || content == "~" {
			dec.readAnything = true
			return createScalarNode(nil, content), nil
		} else if commentLineRegEx.MatchString(content) {
			dec.readAnything = true
			node := createScalarNode(nil, "")
			node.LeadingContent = content
			return node, nil
		}
		return nil, err
	} else if err != nil {
		return nil, err
	}

	candidateNode := &CandidateNode{}
	if err := candidateNode.UnmarshalGoccyYAML(ast, dec.cm, dec.anchorMap); err != nil {
		return nil, err
	}

	return candidateNode, nil
}
