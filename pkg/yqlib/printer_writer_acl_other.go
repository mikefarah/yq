//go:build !linux

package yqlib

import "os"

func openSplitSource(root *os.Root, name string) (*os.File, error) {
	return root.OpenFile(name, os.O_WRONLY, 0)
}

func copySplitAccessControl(_, _ *os.File) error {
	return nil
}
