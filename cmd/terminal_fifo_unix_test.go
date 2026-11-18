//go:build unix

package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/sys/unix"
)

func TestSetupColorsForNamedPipe(t *testing.T) {
	testSetupColorsForFile(t, terminalCase{isTerm: false, open: func(t *testing.T) *os.File {
		path := filepath.Join(t.TempDir(), "fifo")
		if err := unix.Mkfifo(path, 0600); err != nil {
			t.Fatalf("failed to create named pipe: %v", err)
		}
		// O_RDWR avoids blocking on open: a FIFO opened only for reading
		// (or only writing) blocks until a peer opens the other end.
		f, err := os.OpenFile(path, os.O_RDWR, os.ModeNamedPipe)
		if err != nil {
			t.Fatalf("failed to open named pipe: %v", err)
		}
		return f
	}})
}
