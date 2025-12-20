package yqlib

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"regexp"
	"strings"

	yaml "go.yaml.in/yaml/v4"
)

var (
	commentLineRe       = regexp.MustCompile(`^\s*#`)
	yamlDirectiveLineRe = regexp.MustCompile(`^\s*%YAML`)
	separatorLineRe     = regexp.MustCompile(`^\s*---\s*$`)
	separatorPrefixRe   = regexp.MustCompile(`^\s*---\s+`)
)

type yamlDecoder struct {
	decoder yaml.Decoder

	prefs YamlPreferences

	// work around of various parsing issues by yaml.v3 with document headers
	leadingContent string
	bufferRead     bytes.Buffer

	// anchor map persists over multiple documents for convenience.
	anchorMap map[string]*CandidateNode

	readAnything  bool
	firstFile     bool
	documentIndex uint
}

func NewYamlDecoder(prefs YamlPreferences) Decoder {
	return &yamlDecoder{prefs: prefs, firstFile: true}
}

func (dec *yamlDecoder) processReadStream(reader *bufio.Reader) (io.Reader, string, error) {
	var sb strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if errors.Is(err, io.EOF) && line == "" {
			// no more data
			return reader, sb.String(), nil
		}
		if err != nil && !errors.Is(err, io.EOF) {
			return reader, sb.String(), err
		}

		// Determine newline style and strip it for inspection
		newline := ""
		if strings.HasSuffix(line, "\r\n") {
			newline = "\r\n"
			line = strings.TrimSuffix(line, "\r\n")
		} else if strings.HasSuffix(line, "\n") {
			newline = "\n"
			line = strings.TrimSuffix(line, "\n")
		}

		trimmed := strings.TrimSpace(line)

		// Document separator: exact line '---' or a '--- ' prefix followed by content
		if separatorLineRe.MatchString(trimmed) {
			sb.WriteString("$yqDocSeparator$")
			sb.WriteString(newline)
			if errors.Is(err, io.EOF) {
				return reader, sb.String(), nil
			}
			continue
		}

		// Handle lines that start with '--- ' followed by more content (e.g. '--- cat')
		if separatorPrefixRe.MatchString(line) {
			match := separatorPrefixRe.FindString(line)
			remainder := line[len(match):]
			// normalise separator newline: if original had none, default to LF
			sepNewline := newline
			if sepNewline == "" {
				sepNewline = "\n"
			}
			sb.WriteString("$yqDocSeparator$")
			sb.WriteString(sepNewline)
			// push the remainder back onto the reader and continue processing
			reader = bufio.NewReader(io.MultiReader(strings.NewReader(remainder), reader))
			if errors.Is(err, io.EOF) && remainder == "" {
				return reader, sb.String(), nil
			}
			continue
		}

		// Comments, YAML directives, and blank lines are leading content
		if commentLineRe.MatchString(line) || yamlDirectiveLineRe.MatchString(line) || trimmed == "" {
			sb.WriteString(line)
			sb.WriteString(newline)
			if errors.Is(err, io.EOF) {
				return reader, sb.String(), nil
			}
			continue
		}

		// First non-leading line: push it back onto a reader and return
		originalLine := line + newline
		return io.MultiReader(strings.NewReader(originalLine), reader), sb.String(), nil
	}
}

func (dec *yamlDecoder) Init(reader io.Reader) error {
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
	dec.decoder = *yaml.NewDecoder(readerToUse)
	dec.firstFile = false
	dec.documentIndex = 0
	dec.anchorMap = make(map[string]*CandidateNode)
	return nil
}

func (dec *yamlDecoder) Decode() (*CandidateNode, error) {
	var yamlNode yaml.Node
	err := dec.decoder.Decode(&yamlNode)

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
	} else if len(yamlNode.Content) == 0 {
		return nil, errors.New("yaml node has no content")
	}

	candidateNode := CandidateNode{document: dec.documentIndex}
	// don't bother with the DocumentNode
	err = candidateNode.UnmarshalYAML(yamlNode.Content[0], dec.anchorMap)
	if err != nil {
		return nil, err
	}

	candidateNode.HeadComment = yamlNode.HeadComment + candidateNode.HeadComment
	candidateNode.FootComment = yamlNode.FootComment + candidateNode.FootComment

	if dec.leadingContent != "" {
		candidateNode.LeadingContent = dec.leadingContent
		dec.leadingContent = ""
	}
	dec.readAnything = true
	dec.documentIndex++
	return &candidateNode, nil
}

func (dec *yamlDecoder) blankNodeWithComment() *CandidateNode {
	node := createScalarNode(nil, "")
	node.LeadingContent = dec.leadingContent
	return node
}
