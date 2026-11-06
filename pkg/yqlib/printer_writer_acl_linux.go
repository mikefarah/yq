//go:build linux

package yqlib

import (
	"errors"
	"os"

	"golang.org/x/sys/unix"
)

func openSplitSource(root *os.Root, name string) (*os.File, error) {
	return root.OpenFile(name, os.O_RDONLY|unix.O_NOFOLLOW|unix.O_NONBLOCK, 0)
}

func copySplitAccessControl(source, destination *os.File) error {
	const accessACL = "system.posix_acl_access"

	size, err := unix.Fgetxattr(int(source.Fd()), accessACL, nil)
	if errors.Is(err, unix.ENODATA) {
		err = unix.Fremovexattr(int(destination.Fd()), accessACL)
		if errors.Is(err, unix.ENODATA) || errors.Is(err, unix.ENOTSUP) || errors.Is(err, unix.EOPNOTSUPP) {
			return nil
		}
		return err
	}
	if errors.Is(err, unix.ENOTSUP) || errors.Is(err, unix.EOPNOTSUPP) {
		return nil
	}
	if err != nil {
		return err
	}
	acl := make([]byte, size)
	if _, err := unix.Fgetxattr(int(source.Fd()), accessACL, acl); err != nil {
		return err
	}
	return unix.Fsetxattr(int(destination.Fd()), accessACL, acl, 0)
}
