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

		err := os.Chown(file.Name(), uid, gid)
		if err != nil {
			// this happens with snap confinement
			// not really a big issue as users can chown
			// the file themselves if required.
			log.Info("Skipping chown: %v", err)
		}
	}
	return nil
}
