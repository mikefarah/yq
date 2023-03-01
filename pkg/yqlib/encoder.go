package yqlib

import (
	"io"

	yaml "gopkg.in/yaml.v3"
)

type Encoder interface {
	Encode(writer io.Writer, node *yaml.Node) error
	PrintDocumentSeparator(writer io.Writer) error
	PrintLeadingContent(writer io.Writer, content string) error
	CanHandleAliases() bool
}

func mapKeysToStrings(node *yaml.Node) {

	if node.Kind == yaml.MappingNode {
		for index, child := range node.Content {
			if index%2 == 0 { // its a map key
				child.Tag = "!!str"
			}
		}
	}

	for _, child := range node.Content {
		mapKeysToStrings(child)
	}
}
