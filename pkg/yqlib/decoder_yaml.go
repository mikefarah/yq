package yqlib

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type yamlDecoder struct {
	decoder yaml.Decoder
	// work around of various parsing issues by yaml.v3 with document headers
	prefs          YamlPreferences
	leadingContent string
}

func NewYamlDecoder(prefs YamlPreferences) Decoder {
	return &yamlDecoder{prefs: prefs}
}

func (dec *yamlDecoder) processReadStream(reader *bufio.Reader) (io.Reader, string, error) {
	var commentLineRegEx = regexp.MustCompile(`^\s*#`)
	var sb strings.Builder
	for {
		peekBytes, err := reader.Peek(3)
		if errors.Is(err, io.EOF) {
			// EOF are handled else where..
			return reader, sb.String(), nil
		} else if err != nil {
			return reader, sb.String(), err
		} else if string(peekBytes) == "---" {
			_, err := reader.ReadString('\n')
			sb.WriteString("$yqDocSeperator$\n")
			if errors.Is(err, io.EOF) {
				return reader, sb.String(), nil
			} else if err != nil {
				return reader, sb.String(), err
			}
		} else if commentLineRegEx.MatchString(string(peekBytes)) {
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

func (dec *yamlDecoder) Init(reader io.Reader) error {
	readerToUse := reader
	leadingContent := ""
	var err error
	if dec.prefs.LeadingContentPreProcessing {
		readerToUse, leadingContent, err = dec.processReadStream(bufio.NewReader(reader))
		if err != nil {
			return err
		}
	}
	dec.leadingContent = leadingContent
	dec.decoder = *yaml.NewDecoder(readerToUse)
	return nil
}

func (dec *yamlDecoder) Decode() (*CandidateNode, error) {
	var dataBucket yaml.Node

	err := dec.decoder.Decode(&dataBucket)
	if err != nil {
		return nil, err
	}

	candidateNode := &CandidateNode{
		Node: &dataBucket,
	}

	if dec.leadingContent != "" {
		candidateNode.LeadingContent = dec.leadingContent
		dec.leadingContent = ""
	}
	// move document comments into candidate node
	// otherwise unwrap drops them.
	candidateNode.TrailingContent = dataBucket.FootComment
	dataBucket.FootComment = ""
	return candidateNode, nil
}
