//go:build linux

package yqlib

import (
	"errors"
	"fmt"
	"math"
	"os"

	"golang.org/x/sys/unix"
)

func openSplitSource(root *os.Root, name string) (*os.File, error) {
	return root.OpenFile(name, os.O_WRONLY|unix.O_NOFOLLOW|unix.O_NONBLOCK, 0)
}

func copySplitAccessControl(source, destination *os.File) error {
	const accessACL = "system.posix_acl_access"

	sourceFD, err := splitFileDescriptor(source)
	if err != nil {
		return err
	}
	destinationFD, err := splitFileDescriptor(destination)
	if err != nil {
		return err
	}

	size, err := unix.Fgetxattr(sourceFD, accessACL, nil)
	if errors.Is(err, unix.ENODATA) {
		err = unix.Fremovexattr(destinationFD, accessACL)
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
	if _, err := unix.Fgetxattr(sourceFD, accessACL, acl); err != nil {
		return err
	}
	return unix.Fsetxattr(destinationFD, accessACL, acl, 0)
}

func splitFileDescriptor(file *os.File) (int, error) {
	fd := file.Fd()
	if fd > math.MaxInt32 {
		return 0, fmt.Errorf("invalid file descriptor for %q: %w", file.Name(), os.ErrInvalid)
	}
	return int(fd), nil
}
