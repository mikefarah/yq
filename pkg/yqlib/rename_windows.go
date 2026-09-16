//go:build windows

package yqlib

import (
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modkernel32      = windows.NewLazySystemDLL("kernel32.dll")
	procReplaceFileW = modkernel32.NewProc("ReplaceFileW")
)

// renameFile moves 'from' onto 'to'. When 'to' already exists, the Win32 ReplaceFile
// API is used instead of a plain os.Rename (which uses MoveFileEx under the hood):
// ReplaceFile preserves the replaced file's ACLs, alternate data streams and other
// NTFS-specific attributes, which is what an atomic in-place edit is expected to do.
// Without this, yq -i silently drops any ACEs on the original file that grant access
// to users other than the one running yq, since the temp file created for the edit
// only inherits the ACL of its own (unrelated) parent directory.
//
// golang.org/x/sys/windows does not wrap ReplaceFile, so it is called directly.
func renameFile(from string, to string) error {
	if _, err := os.Lstat(to); err != nil {
		// Target doesn't exist yet (or isn't statable) - nothing to preserve, a plain
		// rename is sufficient and ReplaceFile requires the target to already exist.
		return os.Rename(from, to)
	}

	replacedName, err := windows.UTF16PtrFromString(to)
	if err != nil {
		return err
	}
	replacementName, err := windows.UTF16PtrFromString(from)
	if err != nil {
		return err
	}

	ret, _, callErr := procReplaceFileW.Call(
		uintptr(unsafe.Pointer(replacedName)),
		uintptr(unsafe.Pointer(replacementName)),
		0, // lpBackupFileName - no backup kept
		0, // dwReplaceFlags
		0, // lpExclude - reserved, must be NULL
		0, // lpReserved - reserved, must be NULL
	)
	if ret == 0 {
		return callErr
	}
	return nil
}
