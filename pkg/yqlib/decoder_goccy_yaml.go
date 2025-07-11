//go:build !yq_noyaml

// Package yqlib provides YAML processing functionality using the goccy/go-yaml parser.
// This decoder implementation provides compatibility with the legacy yaml.v3 parser
// while leveraging the actively maintained goccy/go-yaml library.
package yqlib

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	yaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

// Constants for document processing
const (
	// docSeparatorPlaceholder is used internally to mark document separators during preprocessing.
	docSeparatorPlaceholder = "$yqDocSeparator$"
)

var commentLineRegEx = regexp.MustCompile(`^\s*#`)
var yamlDirectiveLineRegEx = regexp.MustCompile(`^\s*%YA`) // e.g., %YAML

type goccyYamlDecoder struct {
	decoder yaml.Decoder
	cm      yaml.CommentMap

	prefs YamlPreferences

	// leadingContent stores comments and directives found before the actual YAML content.
	leadingContent string
	// bufferRead is used with a TeeReader when LeadingContentPreProcessing is off,
	// to capture initial bytes that might only contain comments.
	bufferRead bytes.Buffer

	// anchorMap stores anchors encountered within the current YAML stream.
	// It is reset by Init() for each new stream.
	anchorMap map[string]*CandidateNode

	// dateTimePreprocessor handles automatic timestamp tagging for datetime arithmetic compatibility
	dateTimePreprocessor *DateTimePreprocessor

	readAnything  bool // Flag to track if any actual YAML node (or synthesized comment node) has been decoded.
	firstFile     bool // Flag for 'evaluateTogether' mode to handle leading content only once.
	documentIndex uint // Index of the current document within the stream.
}

// NewGoccyYAMLDecoder creates a new YAML decoder using the goccy/go-yaml library,
// configured with the given preferences.
func NewGoccyYAMLDecoder(prefs YamlPreferences) Decoder {
	return &goccyYamlDecoder{
		prefs:                prefs,
		firstFile:            true,
		dateTimePreprocessor: NewDateTimePreprocessor(true), // Enable datetime preprocessing for Goccy
	}
}

// processReadStream pre-processes the input stream to extract leading comments,
// directives, and document separators before the main YAML content is parsed.
// It returns a reader for the remaining stream and the extracted leading content.
func (dec *goccyYamlDecoder) processReadStream(reader *bufio.Reader) (io.Reader, string, error) {
	var sb strings.Builder
	for {
		peekBytes, err := reader.Peek(4) // Peek up to 4 bytes to identify line types.

		if err != nil {
			if errors.Is(err, io.EOF) { // EOF encountered while peeking.
				// No more full lines can be definitively processed as leading content.
				return reader, sb.String(), nil
			}
			return reader, sb.String(), err // Other errors during Peek.
		}

		peekString := string(peekBytes) // Convert peeked bytes to string once per iteration.

		if strings.HasPrefix(peekString, "\n") { // Line is blank.
			_, readErr := reader.ReadString('\n')
			sb.WriteString("\n")
			if readErr != nil {
				if errors.Is(readErr, io.EOF) {
					return reader, sb.String(), nil
				}
				return reader, sb.String(), readErr
			}
		} else if strings.HasPrefix(peekString, "--- ") { // Document separator "--- "
			_, readErr := reader.ReadString(' ') // Consume "--- "
			sb.WriteString(docSeparatorPlaceholder + "\n")
			if readErr != nil {
				if errors.Is(readErr, io.EOF) {
					return reader, sb.String(), nil
				}
				return reader, sb.String(), readErr
			}
		} else if strings.HasPrefix(peekString, "---\n") { // Document separator "---\n"
			_, readErr := reader.ReadString('\n') // Consume "---\n"
			sb.WriteString(docSeparatorPlaceholder + "\n")
			if readErr != nil {
				if errors.Is(readErr, io.EOF) {
					return reader, sb.String(), nil
				}
				return reader, sb.String(), readErr
			}
		} else if commentLineRegEx.MatchString(peekString) || yamlDirectiveLineRegEx.MatchString(peekString) {
			// Line is a comment or a YAML directive.
			line, readErr := reader.ReadString('\n')
			sb.WriteString(line)
			if readErr != nil {
				if errors.Is(readErr, io.EOF) {
					return reader, sb.String(), nil
				}
				return reader, sb.String(), readErr
			}
		} else {
			// Line does not match any leading content pattern, so actual YAML content begins.
			return reader, sb.String(), nil
		}
	}
}

// Init initializes the decoder with a new reader.
// It handles leading content preprocessing according to preferences.
func (dec *goccyYamlDecoder) Init(reader io.Reader) error {
	readerToUse := reader
	processedLeadingContent := ""
	dec.bufferRead = bytes.Buffer{} // Reset buffer
	var err error

	if dec.prefs.LeadingContentPreProcessing && (!dec.prefs.EvaluateTogether || dec.firstFile) {
		// Only process leading content for the first file if 'evaluateTogether' is active.
		readerToUse, processedLeadingContent, err = dec.processReadStream(bufio.NewReader(reader))
		if err != nil {
			return err
		}

		// Apply datetime preprocessing to the remaining content
		if dec.dateTimePreprocessor != nil {
			remainingBytes, err := io.ReadAll(readerToUse)
			if err != nil && errors.Is(err, io.EOF) {
				return err
			}
			preprocessedContent := dec.dateTimePreprocessor.PreprocessDocument(string(remainingBytes))
			readerToUse = strings.NewReader(preprocessedContent)
		}
	} else if !dec.prefs.LeadingContentPreProcessing {
		// If not preprocessing, TeeReader captures initial bytes in case it's a comment-only document.
		if dec.dateTimePreprocessor != nil {
			// Read all content, apply datetime preprocessing, then create new reader
			allBytes, err := io.ReadAll(reader)
			if err != nil && errors.Is(err, io.EOF) {
				return err
			}
			preprocessedContent := dec.dateTimePreprocessor.PreprocessDocument(string(allBytes))
			readerToUse = io.TeeReader(strings.NewReader(preprocessedContent), &dec.bufferRead)
		} else {
			readerToUse = io.TeeReader(reader, &dec.bufferRead)
		}
	}

	dec.leadingContent = processedLeadingContent
	dec.readAnything = false
	dec.cm = yaml.CommentMap{} // Reset comment map for the new stream.
	dec.decoder = *yaml.NewDecoder(readerToUse, yaml.CommentToMap(dec.cm), yaml.UseOrderedMap())
	dec.firstFile = false                           // Subsequent Init calls are not for the 'firstFile' in 'evaluateTogether' context.
	dec.documentIndex = 0                           // Reset document index for the new stream.
	dec.anchorMap = make(map[string]*CandidateNode) // Reset anchor map for the new stream.
	return nil
}

// Decode reads the next YAML document from the stream and decodes it into a CandidateNode.
// It handles EOF conditions and the synthesis of nodes for comment-only content.
func (dec *goccyYamlDecoder) Decode() (*CandidateNode, error) {
	var astNode ast.Node
	err := dec.decoder.Decode(&astNode)

	if err != nil { // An error occurred (could be EOF).
		if errors.Is(err, io.EOF) {
			// Handle EOF: potentially return a comment-only node, or propagate EOF.
			// Case 1: Leading content was processed, and nothing else has been successfully decoded from this stream yet.
			if dec.leadingContent != "" && !dec.readAnything {
				dec.readAnything = true // Mark that we are consuming the leadingContent as a node.
				return dec.blankNodeWithComment(), nil
			}
			// Case 2: Leading content preprocessing was disabled, nothing successfully decoded yet,
			// so check the TeeReader buffer for comments.
			if !dec.prefs.LeadingContentPreProcessing && !dec.readAnything {
				dec.readAnything = true                      // Mark that we are attempting to consume buffered content.
				dec.leadingContent = dec.bufferRead.String() // Transfer buffered content to leadingContent.
				if dec.leadingContent != "" {
					return dec.blankNodeWithComment(), nil
				}
			}
			// If neither of the above conditions for a comment-only node is met, it's a genuine EOF.
			return nil, io.EOF
		}

		// Provide more informative error messages for known goccy limitations
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "could not find flow map content") {
			return nil, fmt.Errorf("flow maps with alias keys are not supported in this parser. Consider using block map syntax instead. Original error: %w", err)
		}
		if strings.Contains(errorMsg, "sequence was used where mapping is expected") {
			return nil, fmt.Errorf("merge anchor only supports maps, got !!seq instead")
		}

		// Non-EOF error.
		return nil, err
	}

	// No error from dec.decoder.Decode(&astNode) means a successful decode of a YAML structure.
	dec.readAnything = true

	candidateNode := &CandidateNode{document: dec.documentIndex}
	if errUnmarshal := candidateNode.UnmarshalGoccyYAML(astNode, dec.cm, dec.anchorMap); errUnmarshal != nil {
		return nil, errUnmarshal
	}

	// If there was leading content processed before this node, attach it.
	if dec.leadingContent != "" {
		candidateNode.LeadingContent = dec.leadingContent
		dec.leadingContent = "" // Clear after attaching.
	}

	dec.documentIndex++
	return candidateNode, nil
}

// blankNodeWithComment creates an empty scalar node and attaches the current leadingContent as its comment.
// This is used for documents that only contain comments or directives.
func (dec *goccyYamlDecoder) blankNodeWithComment() *CandidateNode {
	node := createScalarNode(nil, "") // Create an empty scalar node.
	node.LeadingContent = dec.leadingContent
	return node
}
