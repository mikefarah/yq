package yqlib

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type yamlDecoder struct {
	decoder yaml.Decoder

	prefs YamlPreferences

	// work around of various parsing issues by yaml.v3 with document headers
	leadingContent string
	bufferRead     bytes.Buffer

	readAnything bool
	firstFile    bool
}

func NewYamlDecoder(prefs YamlPreferences) Decoder {
	return &yamlDecoder{prefs: prefs, firstFile: true}
}

func (dec *yamlDecoder) processReadStream(reader *bufio.Reader) (io.Reader, string, error) {
	var commentLineRegEx = regexp.MustCompile(`^\s*#`)
	var yamlDirectiveLineRegEx = regexp.MustCompile(`^\s*%YA`)
	var sb strings.Builder
	for {
		peekBytes, err := reader.Peek(3)
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
		} else if string(peekBytes) == "---" {
			_, err := reader.ReadString('\n')
			sb.WriteString("$yqDocSeperator$\n")
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
	return nil
}

func (dec *yamlDecoder) convertKind(oKind yaml.Kind) Kind {
	switch oKind {
	case yaml.DocumentNode:
		return DocumentNode
	case yaml.SequenceNode:
		return SequenceNode
	case yaml.MappingNode:
		return MappingNode
	case yaml.ScalarNode:
		return ScalarNode
	case yaml.AliasNode:
		return AliasNode
	}
	return ScalarNode
}

func (dec *yamlDecoder) convertStyle(oStyle yaml.Style) Style {
	switch oStyle {
	case yaml.TaggedStyle:
		return TaggedStyle
	case yaml.DoubleQuotedStyle:
		return DoubleQuotedStyle
	case yaml.SingleQuotedStyle:
		return SingleQuotedStyle
	case yaml.LiteralStyle:
		return LiteralStyle
	case yaml.FoldedStyle:
		return FoldedStyle
	case yaml.FlowStyle:
		return FlowStyle
	}
	return 0
}

func (dec *yamlDecoder) ConvertToCandidateNode(yamlNode *yaml.Node) *CandidateNode {
	kids := make([]*CandidateNode, len(yamlNode.Content))
	for i, v := range yamlNode.Content {
		kids[i] = dec.ConvertToCandidateNode(v)
	}

	return &CandidateNode{
		Kind:  dec.convertKind(yamlNode.Kind),
		Style: dec.convertStyle(yamlNode.Style),

		Tag:     yamlNode.Tag,
		Value:   yamlNode.Value,
		Anchor:  yamlNode.Anchor,
		Alias:   dec.ConvertToCandidateNode(yamlNode.Alias),
		Content: kids,

		HeadComment: yamlNode.HeadComment,
		LineComment: yamlNode.LineComment,
		FootComment: yamlNode.FootComment,

		// Parent: yamlNode.Parent,
		// Key:    yamlNode.Key,

		// LeadingContent:  yamlNode.LeadingContent,
		// TrailingContent: yamlNode.TrailingContent,

		// Path:     yamlNode.Path,
		// Document: yamlNode.Document,
		// Filename: yamlNode.Filename,

		Line:   yamlNode.Line,
		Column: yamlNode.Column,

		// FileIndex:        yamlNode.FileIndex,
		// EvaluateTogether: yamlNode.EvaluateTogether,
		// IsMapKey:         yamlNode.IsMapKey,
	}

}

func (dec *yamlDecoder) Decode() (*CandidateNode, error) {
	var dataBucket yaml.Node
	err := dec.decoder.Decode(&dataBucket)
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

	candidateNode := dec.ConvertToCandidateNode(&dataBucket)

	if dec.leadingContent != "" {
		candidateNode.LeadingContent = dec.leadingContent
		dec.leadingContent = ""
	}
	dec.readAnything = true
	// move document comments into candidate node
	// otherwise unwrap drops them.
	candidateNode.TrailingContent = dataBucket.FootComment
	dataBucket.FootComment = ""
	return candidateNode, nil
}

func (dec *yamlDecoder) blankNodeWithComment() *CandidateNode {
	return &CandidateNode{
		Document:       0,
		Filename:       "",
		Kind:           DocumentNode,
		Content:        []*CandidateNode{createScalarNode(nil, "")},
		FileIndex:      0,
		LeadingContent: dec.leadingContent,
	}
}
