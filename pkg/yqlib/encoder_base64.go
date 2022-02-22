package yqlib

import (
	"encoding/base64"
	"fmt"
	"io"

	yaml "gopkg.in/yaml.v3"
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

func (e *base64Encoder) PrintDocumentSeparator(writer io.Writer) error {
	return nil
}

func (e *base64Encoder) PrintLeadingContent(writer io.Writer, content string) error {
	return nil
}

func (e *base64Encoder) Encode(writer io.Writer, originalNode *yaml.Node) error {
	node := unwrapDoc(originalNode)
	if guessTagFromCustomType(node) != "!!str" {
		return fmt.Errorf("cannot encode %v as base64, can only operate on strings. Please first pipe through another encoding operator to convert the value to a string", node.Tag)
	}
	_, err := writer.Write([]byte(e.encoding.EncodeToString([]byte(originalNode.Value))))
	return err
}
