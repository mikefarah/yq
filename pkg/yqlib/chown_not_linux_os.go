//go:build !linux

package yqlib

import (
	"io/fs"
	"os"
)

func changeOwner(info fs.FileInfo, file *os.File) error {
	return nil
}
