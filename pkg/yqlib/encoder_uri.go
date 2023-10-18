package yqlib

import (
	"fmt"
	"io"
	"net/url"
)

type uriEncoder struct {
}

func NewUriEncoder() Encoder {
	return &uriEncoder{}
}

func (e *uriEncoder) CanHandleAliases() bool {
	return false
}

func (e *uriEncoder) PrintDocumentSeparator(writer io.Writer) error {
	return nil
}

func (e *uriEncoder) PrintLeadingContent(writer io.Writer, content string) error {
	return nil
}

func (e *uriEncoder) Encode(writer io.Writer, node *CandidateNode) error {
	if node.guessTagFromCustomType() != "!!str" {
		return fmt.Errorf("cannot encode %v as URI, can only operate on strings. Please first pipe through another encoding operator to convert the value to a string", node.Tag)
	}
	_, err := writer.Write([]byte(url.QueryEscape(node.Value)))
	return err
}
