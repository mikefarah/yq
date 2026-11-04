//go:build !windows

package yqlib

import "os"

func renameFile(from string, to string) error {
	return os.Rename(from, to)
}
