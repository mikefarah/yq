package yqlib

import (
	"io"
)

type Encoder interface {
	Encode(writer io.Writer, node *CandidateNode) error
	PrintDocumentSeparator(writer io.Writer) error
	PrintLeadingContent(writer io.Writer, content string) error
	CanHandleAliases() bool
}

func mapKeysToStrings(node *CandidateNode) {

	if node.Kind == MappingNode {
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
