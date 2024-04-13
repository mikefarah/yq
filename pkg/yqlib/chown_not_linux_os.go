//go:build !linux

package yqlib

import (
	"io/fs"
	"os"
)

func changeOwner(_ fs.FileInfo, _ *os.File) error {
	return nil
}
