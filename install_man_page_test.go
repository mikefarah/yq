package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// TestInstallManPageResolvesBinSymlink ensures that when yq is found through
// a /bin style directory that is a symlink to /usr/bin, the man page is
// installed under the physical /usr prefix rather than the symlink's parent.
func TestInstallManPageResolvesBinSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("requires POSIX sh and symbolic links")
	}

	tmp := t.TempDir()
	usrBin := filepath.Join(tmp, "usr", "bin")
	if err := os.MkdirAll(usrBin, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(usrBin, "yq"), []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join("usr", "bin"), filepath.Join(tmp, "bin")); err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile(filepath.Join("scripts", "install-man-page.sh"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "install-man-page.sh"), script, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "yq.1"), []byte("man page"), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("sh", "install-man-page.sh")
	cmd.Dir = tmp
	cmd.Env = append(os.Environ(), "PATH="+filepath.Join(tmp, "bin")+":/usr/bin:/bin")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("install-man-page.sh failed: %v\n%s", err, output)
	}

	installed := filepath.Join(tmp, "usr", "share", "man", "man1", "yq.1")
	if _, err := os.Stat(installed); err != nil {
		t.Fatalf("expected man page at %s: %v", installed, err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "share")); !os.IsNotExist(err) {
		t.Fatalf("man page was installed under the symlink parent instead of usr: %v", err)
	}
}
