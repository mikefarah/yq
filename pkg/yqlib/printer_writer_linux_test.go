//go:build linux

package yqlib

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
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
	reader, err := os.OpenFile("pipe.txt", os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		t.Fatal(err)
	}
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
	if err := os.WriteFile("output.txt", []byte("before"), 0600); err != nil {
		t.Fatal(err)
	}
	file, err := os.OpenFile("output.txt", os.O_RDONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	acl := restrictiveAccessACL()
	fd, err := splitFileDescriptor(file)
	if err != nil {
		t.Fatal(err)
	}
	err = unix.Fsetxattr(fd, "system.posix_acl_access", acl, 0)
	if closeErr := file.Close(); closeErr != nil {
		t.Fatal(closeErr)
	}
	if err != nil {
		if errors.Is(err, unix.ENOTSUP) || errors.Is(err, unix.EOPNOTSUPP) || errors.Is(err, unix.EPERM) {
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
	fd, err = splitFileDescriptor(output)
	if err != nil {
		t.Fatal(err)
	}
	n, err := unix.Fgetxattr(fd, "system.posix_acl_access", actualACL)
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

func TestSplitFileDescriptorRejectsClosedFile(t *testing.T) {
	file, err := os.CreateTemp(t.TempDir(), "closed")
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := splitFileDescriptor(file); !errors.Is(err, os.ErrInvalid) {
		t.Fatalf("expected invalid file descriptor error, got %v", err)
	}
}

func TestMultiPrinterWriterPermissionHelper(t *testing.T) {
	if os.Getenv("YQ_SPLIT_PERMISSION_HELPER") != "1" {
		return
	}
	writer, err := getTestSplitWriter(t, "output.txt")
	expectWrite := os.Getenv("YQ_SPLIT_EXPECT_WRITE") == "true"
	if expectWrite && err != nil {
		t.Fatalf("expected write-only destination to be accepted: %v", err)
	}
	if !expectWrite && err == nil {
		t.Fatal("expected read-only destination to be rejected")
	}
	if err == nil {
		if _, err := writer.WriteString("updated"); err != nil {
			t.Fatal(err)
		}
		if err := writer.Flush(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestMultiPrinterWriterRespectsNonRootWritePermissions(t *testing.T) {
	tempDir := t.TempDir()
	workDir := filepath.Join(tempDir, "work")
	if err := os.Mkdir(workDir, 0700); err != nil {
		t.Fatal(err)
	}
	testBinary, err := os.ReadFile("/proc/self/exe")
	if err != nil {
		t.Fatal(err)
	}
	testBinaryPath := filepath.Join(workDir, "yq-test")
	if err := os.WriteFile(testBinaryPath, testBinary, 0600); err != nil { // #nosec G703 -- path is inside the private test directory
		t.Fatal(err)
	}
	if err := os.Chmod(testBinaryPath, 0700); err != nil {
		t.Fatal(err)
	}

	var credential *syscall.Credential
	if os.Geteuid() == 0 {
		nobody, err := user.LookupId("65534")
		if err != nil {
			t.Fatalf("could not find unprivileged user: %v", err)
		}
		uid, err := strconv.ParseUint(nobody.Uid, 10, 32)
		if err != nil {
			t.Fatal(err)
		}
		gid, err := strconv.ParseUint(nobody.Gid, 10, 32)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.Chown(tempDir, int(uid), int(gid)); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(filepath.Dir(tempDir), 0711); err != nil {
			t.Fatal(err)
		}
		if err := os.Chown(workDir, int(uid), int(gid)); err != nil {
			t.Fatal(err)
		}
		if err := os.Chown(testBinaryPath, int(uid), int(gid)); err != nil {
			t.Fatal(err)
		}
		credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
	}

	for _, testCase := range []struct {
		name        string
		mode        os.FileMode
		expectWrite bool
	}{
		{name: "write-only", mode: 0200, expectWrite: true},
		{name: "read-only", mode: 0400, expectWrite: false},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			outputPath := filepath.Join(workDir, "output.txt")
			if err := os.WriteFile(outputPath, []byte("original"), 0600); err != nil {
				t.Fatal(err)
			}
			if credential != nil {
				if err := os.Chown(outputPath, int(credential.Uid), int(credential.Gid)); err != nil {
					t.Fatal(err)
				}
			}
			if err := os.Chmod(outputPath, testCase.mode); err != nil {
				t.Fatal(err)
			}

			command := exec.Command(testBinaryPath, "-test.run=^TestMultiPrinterWriterPermissionHelper$")
			command.Dir = workDir
			command.Env = append(os.Environ(),
				"YQ_SPLIT_PERMISSION_HELPER=1",
				"YQ_SPLIT_EXPECT_WRITE="+strconv.FormatBool(testCase.expectWrite),
			)
			if credential != nil {
				command.SysProcAttr = &syscall.SysProcAttr{Credential: credential}
			}
			output, err := command.CombinedOutput()
			if err != nil {
				t.Fatalf("permission subprocess failed: %v\n%s", err, output)
			}

			info, err := os.Stat(outputPath)
			if err != nil {
				t.Fatal(err)
			}
			if info.Mode().Perm() != testCase.mode {
				t.Errorf("file mode changed: got %o, want %o", info.Mode().Perm(), testCase.mode)
			}
			if err := os.Chmod(outputPath, testCase.mode|0400); err != nil {
				t.Fatal(err)
			}
			content, err := os.ReadFile(outputPath)
			if err != nil {
				t.Fatal(err)
			}
			if testCase.expectWrite && string(content) != "updated" {
				t.Errorf("writable destination was not updated: %q", content)
			}
			if !testCase.expectWrite && string(content) != "original" {
				t.Errorf("read-only destination was changed: %q", content)
			}
		})
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
