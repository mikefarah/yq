//go:build !yq_noini

package yqlib

import (
	"bytes"
	"fmt"
	"io"

	"github.com/go-ini/ini"
)

type iniEncoder struct {
	indentString string
}

// NewINIEncoder creates a new INI encoder
func NewINIEncoder() Encoder {
	// Hardcoded indent value of 0, meaning no additional spacing.
	return &iniEncoder{""}
}

// CanHandleAliases indicates whether the encoder supports aliases. INI does not support aliases.
func (ie *iniEncoder) CanHandleAliases() bool {
	return false
}

// PrintDocumentSeparator is a no-op since INI does not support multiple documents.
func (ie *iniEncoder) PrintDocumentSeparator(_ io.Writer) error {
	return nil
}

// PrintLeadingContent is a no-op since INI does not support leading content or comments at the encoder level.
func (ie *iniEncoder) PrintLeadingContent(_ io.Writer, _ string) error {
	return nil
}

// Encode converts a CandidateNode into INI format and writes it to the provided writer.
func (ie *iniEncoder) Encode(writer io.Writer, node *CandidateNode) error {
	log.Debugf("I need to encode %v", NodeToString(node))
	log.Debugf("kids %v", len(node.Content))

	if node.Kind == ScalarNode {
		return writeStringINI(writer, node.Value+"\n")
	}

	// Create a new INI configuration.
	cfg := ini.Empty()

	if node.Kind == MappingNode {
		// Default section for key-value pairs at the root level.
		defaultSection, err := cfg.NewSection(ini.DefaultSection)
		if err != nil {
			return err
		}

		// Process the node's content.
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]
			key := keyNode.Value

			switch valueNode.Kind {
			case ScalarNode:
				// Add key-value pair to the default section.
				_, err := defaultSection.NewKey(key, valueNode.Value)
				if err != nil {
					return err
				}
			case MappingNode:
				// Create a new section for nested MappingNode.
				section, err := cfg.NewSection(key)
				if err != nil {
					return err
				}
				// Process nested key-value pairs.
				for j := 0; j < len(valueNode.Content); j += 2 {
					nestedKeyNode := valueNode.Content[j]
					nestedValueNode := valueNode.Content[j+1]
					if nestedValueNode.Kind == ScalarNode {
						_, err := section.NewKey(nestedKeyNode.Value, nestedValueNode.Value)
						if err != nil {
							return err
						}
					} else {
						log.Debugf("Skipping nested non-scalar value for key %s: %v", nestedKeyNode.Value, nestedValueNode.Kind)
					}
				}
			default:
				log.Debugf("Skipping non-scalar value for key %s: %v", key, valueNode.Kind)
			}
		}
	} else {
		return fmt.Errorf("INI encoder supports only MappingNode at the root level, got %v", node.Kind)
	}

	// Use a buffer to store the INI output as the library doesn't support direct io.Writer with indent.
	var buffer bytes.Buffer
	_, err := cfg.WriteToIndent(&buffer, ie.indentString)
	if err != nil {
		return err
	}

	// Write the buffer content to the provided writer.
	_, err = writer.Write(buffer.Bytes())
	return err
}

// writeStringINI is a helper function to write a string to the provided writer for INI encoder.
func writeStringINI(writer io.Writer, content string) error {
	_, err := writer.Write([]byte(content))
	return err
}
