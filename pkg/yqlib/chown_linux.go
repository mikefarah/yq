//go:build linux

package yqlib

import (
	"io/fs"
	"os"
	"syscall"
)

func changeOwner(info fs.FileInfo, file *os.File) error {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		uid := int(stat.Uid)
		gid := int(stat.Gid)
		return os.Chown(file.Name(), uid, gid)
	}
	return nil
}
