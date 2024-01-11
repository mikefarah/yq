package yqlib

import (
	"fmt"
	"io"
)

type tomlEncoder struct {
}

func NewTomlEncoder() Encoder {
	return &tomlEncoder{}
}

func (te *tomlEncoder) Encode(writer io.Writer, node *CandidateNode) error {
	if node.Kind == ScalarNode {
		return writeString(writer, node.Value+"\n")
	}
	return fmt.Errorf("only scalars (e.g. strings, numbers, booleans) are supported for TOML output at the moment. Please use yaml output format (-oy) until the encoder has been fully implemented")
}

func (te *tomlEncoder) PrintDocumentSeparator(_ io.Writer) error {
	return nil
}

func (te *tomlEncoder) PrintLeadingContent(_ io.Writer, _ string) error {
	return nil
}

func (te *tomlEncoder) CanHandleAliases() bool {
	return false
}
