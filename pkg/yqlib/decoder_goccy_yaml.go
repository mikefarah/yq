//go:build !yq_noyaml

//
// NOTE this is still a WIP - not yet ready.
//

package yqlib

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"regexp"
	"strings"

	yaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

type goccyYamlDecoder struct {
	decoder yaml.Decoder
	cm      yaml.CommentMap

	prefs YamlPreferences

	// work around of various parsing issues by handling document headers
	leadingContent string
	bufferRead     bytes.Buffer

	// anchor map persists over multiple documents for convenience.
	anchorMap map[string]*CandidateNode

	readAnything  bool
	firstFile     bool
	documentIndex uint
}

func NewGoccyYAMLDecoder(prefs YamlPreferences) Decoder {
	return &goccyYamlDecoder{prefs: prefs, firstFile: true}
}

func (dec *goccyYamlDecoder) processReadStream(reader *bufio.Reader) (io.Reader, string, error) {
	var commentLineRegEx = regexp.MustCompile(`^\s*#`)
	var yamlDirectiveLineRegEx = regexp.MustCompile(`^\s*%YA`)
	var sb strings.Builder
	for {
		peekBytes, err := reader.Peek(4)
		if errors.Is(err, io.EOF) {
			// EOF are handled else where..
			return reader, sb.String(), nil
		} else if err != nil {
			return reader, sb.String(), err
		} else if string(peekBytes[0]) == "\n" {
			_, err := reader.ReadString('\n')
			sb.WriteString("\n")
			if errors.Is(err, io.EOF) {
				return reader, sb.String(), nil
			} else if err != nil {
				return reader, sb.String(), err
			}
		} else if string(peekBytes) == "--- " {
			_, err := reader.ReadString(' ')
			sb.WriteString("$yqDocSeparator$\n")
			if errors.Is(err, io.EOF) {
				return reader, sb.String(), nil
			} else if err != nil {
				return reader, sb.String(), err
			}
		} else if string(peekBytes) == "---\n" {
			_, err := reader.ReadString('\n')
			sb.WriteString("$yqDocSeparator$\n")
			if errors.Is(err, io.EOF) {
				return reader, sb.String(), nil
			} else if err != nil {
				return reader, sb.String(), err
			}
		} else if commentLineRegEx.MatchString(string(peekBytes)) || yamlDirectiveLineRegEx.MatchString(string(peekBytes)) {
			line, err := reader.ReadString('\n')
			sb.WriteString(line)
			if errors.Is(err, io.EOF) {
				return reader, sb.String(), nil
			} else if err != nil {
				return reader, sb.String(), err
			}
		} else {
			return reader, sb.String(), nil
		}
	}
}

func (dec *goccyYamlDecoder) Init(reader io.Reader) error {
	readerToUse := reader
	leadingContent := ""
	dec.bufferRead = bytes.Buffer{}
	var err error
	// if we 'evaluating together' - we only process the leading content
	// of the first file - this ensures comments from subsequent files are
	// merged together correctly.
	if dec.prefs.LeadingContentPreProcessing && (!dec.prefs.EvaluateTogether || dec.firstFile) {
		readerToUse, leadingContent, err = dec.processReadStream(bufio.NewReader(reader))
		if err != nil {
			return err
		}
	} else if !dec.prefs.LeadingContentPreProcessing {
		// if we're not process the leading content
		// keep a copy of what we've read. This is incase its a
		// doc with only comments - the decoder will return nothing
		// then we can read the comments from bufferRead
		readerToUse = io.TeeReader(reader, &dec.bufferRead)
	}
	dec.leadingContent = leadingContent
	dec.readAnything = false
	dec.cm = yaml.CommentMap{}
	dec.decoder = *yaml.NewDecoder(readerToUse, yaml.CommentToMap(dec.cm), yaml.UseOrderedMap())
	dec.firstFile = false
	dec.documentIndex = 0
	dec.anchorMap = make(map[string]*CandidateNode)
	return nil
}

func (dec *goccyYamlDecoder) Decode() (*CandidateNode, error) {
	var astNode ast.Node
	err := dec.decoder.Decode(&astNode)

	if errors.Is(err, io.EOF) && dec.leadingContent != "" && !dec.readAnything {
		// force returning an empty node with a comment.
		dec.readAnything = true
		return dec.blankNodeWithComment(), nil
	} else if errors.Is(err, io.EOF) && !dec.prefs.LeadingContentPreProcessing && !dec.readAnything {
		// didn't find any yaml,
		// check the tee buffer, maybe there were comments
		dec.readAnything = true
		dec.leadingContent = dec.bufferRead.String()
		if dec.leadingContent != "" {
			return dec.blankNodeWithComment(), nil
		}
		return nil, err
	} else if err != nil {
		return nil, err
	}

	candidateNode := &CandidateNode{document: dec.documentIndex}
	if err := candidateNode.UnmarshalGoccyYAML(astNode, dec.cm, dec.anchorMap); err != nil {
		return nil, err
	}

	if dec.leadingContent != "" {
		candidateNode.LeadingContent = dec.leadingContent
		dec.leadingContent = ""
	}
	dec.readAnything = true
	dec.documentIndex++
	return candidateNode, nil
}

func (dec *goccyYamlDecoder) blankNodeWithComment() *CandidateNode {
	node := createScalarNode(nil, "")
	node.LeadingContent = dec.leadingContent
	return node
}
