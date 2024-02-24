package yqlib

import (
	"encoding/base64"
	"fmt"
	"io"
)

type base64Encoder struct {
	encoding base64.Encoding
}

func NewBase64Encoder() Encoder {
	return &base64Encoder{encoding: *base64.StdEncoding}
}

func (e *base64Encoder) CanHandleAliases() bool {
	return false
}

func (e *base64Encoder) PrintDocumentSeparator(_ io.Writer) error {
	return nil
}

func (e *base64Encoder) PrintLeadingContent(_ io.Writer, _ string) error {
	return nil
}

func (e *base64Encoder) Encode(writer io.Writer, node *CandidateNode) error {
	if node.guessTagFromCustomType() != "!!str" {
		return fmt.Errorf("cannot encode %v as base64, can only operate on strings", node.Tag)
	}
	_, err := writer.Write([]byte(e.encoding.EncodeToString([]byte(node.Value))))
	return err
}
