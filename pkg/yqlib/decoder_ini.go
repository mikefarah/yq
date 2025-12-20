//go:build !yq_noini

package yqlib

import (
	"fmt"
	"io"

	"github.com/go-ini/ini"
)

type iniDecoder struct {
	reader   io.Reader
	finished bool // Flag to signal completion of processing
}

func NewINIDecoder() Decoder {
	return &iniDecoder{
		finished: false, // Initialise the flag as false
	}
}

func (dec *iniDecoder) Init(reader io.Reader) error {
	// Store the reader for use in Decode
	dec.reader = reader
	return nil
}

func (dec *iniDecoder) Decode() (*CandidateNode, error) {
	// If processing is already finished, return io.EOF
	if dec.finished {
		return nil, io.EOF
	}

	// Read all content from the stored reader
	content, err := io.ReadAll(dec.reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read INI content: %w", err)
	}

	// Parse the INI content
	cfg, err := ini.Load(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse INI content: %w", err)
	}

	// Create a root CandidateNode as a MappingNode (since INI is key-value based)
	root := &CandidateNode{
		Kind:  MappingNode,
		Tag:   "!!map",
		Value: "",
	}

	// Process each section in the INI file
	for _, section := range cfg.Sections() {
		sectionName := section.Name()

		if sectionName == ini.DefaultSection {
			// For the default section, add key-value pairs directly to the root node
			for _, key := range section.Keys() {
				keyName := key.Name()
				keyValue := key.String()

				// Create a key node (scalar for the key name)
				keyNode := createStringScalarNode(keyName)
				// Create a value node (scalar for the value)
				valueNode := createStringScalarNode(keyValue)

				// Add key-value pair to the root node
				root.AddKeyValueChild(keyNode, valueNode)
			}
		} else {
			// For named sections, create a nested map
			sectionNode := &CandidateNode{
				Kind:  MappingNode,
				Tag:   "!!map",
				Value: "",
			}

			// Add key-value pairs to the section node
			for _, key := range section.Keys() {
				keyName := key.Name()
				keyValue := key.String()

				// Create a key node (scalar for the key name)
				keyNode := createStringScalarNode(keyName)
				// Create a value node (scalar for the value)
				valueNode := createStringScalarNode(keyValue)

				// Add key-value pair to the section node
				sectionNode.AddKeyValueChild(keyNode, valueNode)
			}

			// Create a key node for the section name
			sectionKeyNode := createStringScalarNode(sectionName)
			// Add the section as a nested map to the root node
			root.AddKeyValueChild(sectionKeyNode, sectionNode)
		}
	}

	// Set the finished flag to true to prevent further Decode calls
	dec.finished = true

	// Return the root node
	return root, nil
}
