//go:build yq_noshell

package yqlib

func NewShellVariablesEncoder() Encoder {
	return nil
}
