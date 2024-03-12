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
