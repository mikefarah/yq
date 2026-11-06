//go:build linux

package yqlib

import (
	"bytes"
	"encoding/binary"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"golang.org/x/sys/unix"
)

func TestReplaceSplitFileUsesOpenHandleForPermissions(t *testing.T) {
	tempDir := setupSplitWriterTest(t)
	outside := filepath.Join(tempDir, "outside.txt")
	if err := os.WriteFile(outside, []byte("outside"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(outside, 0600); err != nil {
		t.Fatal(err)
	}

	root, err := os.OpenRoot(".")
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	file, err := root.OpenFile("temporary", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	if err := root.Remove("temporary"); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(outside, "temporary"); err != nil {
		t.Fatal(err)
	}
	if err := replaceSplitFile(root, file, "temporary", "output.txt", 0640, true); err != nil {
		t.Fatal(err)
	}

	outsideInfo, err := os.Stat(outside)
	if err != nil {
		t.Fatal(err)
	}
	if outsideInfo.Mode().Perm() != 0600 {
		t.Errorf("outside file permissions changed: %v", outsideInfo.Mode().Perm())
	}
	fileInfo, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if fileInfo.Mode().Perm() != 0640 {
		t.Errorf("temporary file permissions not preserved: %v", fileInfo.Mode().Perm())
	}
}

func TestMultiPrinterWriterRejectsSpecialFiles(t *testing.T) {
	setupSplitWriterTest(t)
	if err := syscall.Mkfifo("pipe.txt", 0600); err != nil {
		t.Fatal(err)
	}
	listener, err := net.Listen("unix", "socket.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	for _, path := range []string{"pipe.txt", "socket.txt"} {
		t.Run(path, func(t *testing.T) {
			before, err := os.Lstat(path)
			if err != nil {
				t.Fatal(err)
			}
			_, err = getTestSplitWriter(t, path)
			if err == nil || !strings.Contains(err.Error(), "--security-allow-write-any-path") {
				t.Fatalf("expected %q to be rejected with override guidance, got %v", path, err)
			}
			after, err := os.Lstat(path)
			if err != nil {
				t.Fatal(err)
			}
			if !os.SameFile(before, after) {
				t.Errorf("rejected destination %q was replaced", path)
			}
		})
	}
}

func TestMultiPrinterWriterOverrideStreamsToFIFO(t *testing.T) {
	setupSplitWriterTest(t)
	if err := syscall.Mkfifo("pipe.txt", 0600); err != nil {
		t.Fatal(err)
	}
	before, err := os.Lstat("pipe.txt")
	if err != nil {
		t.Fatal(err)
	}
	fd, err := syscall.Open("pipe.txt", syscall.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		t.Fatal(err)
	}
	reader := os.NewFile(uintptr(fd), "pipe.txt")
	defer reader.Close()

	ConfiguredSecurityPreferences.AllowWriteAnyPath = true
	writer, err := getTestSplitWriter(t, "pipe.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.WriteString("output"); err != nil {
		t.Fatal(err)
	}
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}
	buffer := make([]byte, len("output"))
	n, err := reader.Read(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if string(buffer[:n]) != "output" {
		t.Errorf("unexpected FIFO contents: %q", buffer[:n])
	}
	after, err := os.Lstat("pipe.txt")
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(before, after) {
		t.Error("FIFO was replaced")
	}
}

func TestMultiPrinterWriterOverrideDoesNotReplaceSocket(t *testing.T) {
	setupSplitWriterTest(t)
	listener, err := net.Listen("unix", "socket.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	before, err := os.Lstat("socket.txt")
	if err != nil {
		t.Fatal(err)
	}

	ConfiguredSecurityPreferences.AllowWriteAnyPath = true
	if _, err := getTestSplitWriter(t, "socket.txt"); err == nil {
		t.Fatal("expected legacy socket open to fail")
	}
	after, err := os.Lstat("socket.txt")
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(before, after) {
		t.Error("socket was replaced")
	}
}

func TestMultiPrinterWriterPreservesRestrictiveAccessACL(t *testing.T) {
	setupSplitWriterTest(t)
	if err := os.WriteFile("output.txt", []byte("before"), 0640); err != nil {
		t.Fatal(err)
	}
	file, err := os.OpenFile("output.txt", os.O_RDONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	acl := restrictiveAccessACL()
	err = unix.Fsetxattr(int(file.Fd()), "system.posix_acl_access", acl, 0)
	if closeErr := file.Close(); closeErr != nil {
		t.Fatal(closeErr)
	}
	if err != nil {
		if err == unix.ENOTSUP || err == unix.EOPNOTSUPP || err == unix.EPERM {
			t.Skipf("access ACLs are unavailable: %v", err)
		}
		t.Fatal(err)
	}

	writer, err := getTestSplitWriter(t, "output.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.WriteString("after"); err != nil {
		t.Fatal(err)
	}
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}

	output, err := os.Open("output.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer output.Close()
	actualACL := make([]byte, len(acl))
	n, err := unix.Fgetxattr(int(output.Fd()), "system.posix_acl_access", actualACL)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(acl, actualACL[:n]) {
		t.Errorf("access ACL changed: got %v, want %v", actualACL[:n], acl)
	}
	info, err := output.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0640 {
		t.Errorf("file mode changed: got %o, want 640", info.Mode().Perm())
	}
}

func restrictiveAccessACL() []byte {
	acl := make([]byte, 4+5*8)
	binary.LittleEndian.PutUint32(acl, 2)
	// The group entry denies access while the mask keeps the mode at 0640.
	entries := []struct {
		tag  uint16
		perm uint16
		id   uint32
	}{
		{tag: 0x01, perm: 0x06, id: ^uint32(0)},
		{tag: 0x02, perm: 0x00, id: 424242},
		{tag: 0x04, perm: 0x00, id: ^uint32(0)},
		{tag: 0x10, perm: 0x04, id: ^uint32(0)},
		{tag: 0x20, perm: 0x00, id: ^uint32(0)},
	}
	for i, entry := range entries {
		offset := 4 + i*8
		binary.LittleEndian.PutUint16(acl[offset:], entry.tag)
		binary.LittleEndian.PutUint16(acl[offset+2:], entry.perm)
		binary.LittleEndian.PutUint32(acl[offset+4:], entry.id)
	}
	return acl
}
