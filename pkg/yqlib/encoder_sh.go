package yqlib

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

var pattern = regexp.MustCompile(`[^\w@%+=:,./-]`)

type shEncoder struct {
}

func NewShEncoder() Encoder {
	return &shEncoder{}
}

func (e *shEncoder) CanHandleAliases() bool {
	return false
}

func (e *shEncoder) PrintDocumentSeparator(writer io.Writer) error {
	return nil
}

func (e *shEncoder) PrintLeadingContent(writer io.Writer, content string) error {
	return nil
}

func (e *shEncoder) Encode(writer io.Writer, originalNode *yaml.Node) error {
	node := unwrapDoc(originalNode)
	if guessTagFromCustomType(node) != "!!str" {
		return fmt.Errorf("cannot encode %v as URI, can only operate on strings. Please first pipe through another encoding operator to convert the value to a string", node.Tag)
	}

	value := originalNode.Value
	if pattern.MatchString(value) {
		value = "'" + strings.ReplaceAll(value, "'", "\\'") + "'"
	}
	return writeString(writer, value)
}
