package yqlib

import (
	"fmt"
	"io"
	"os"
)

func tryRenameFile(from string, to string) error {
	if info, err := os.Lstat(to); err == nil && info.Mode()&os.ModeSymlink != 0 {
		log.Debug("Target file is symlink, skipping rename and attempting to copy contents")

		if copyError := copyFileContents(from, to); copyError != nil {
			return fmt.Errorf("failed copying from %v to %v: %w", from, to, copyError)
		}
		tryRemoveTempFile(from)
		return nil
	} else if renameError := os.Rename(from, to); renameError != nil {
		log.Debugf("Error renaming from %v to %v, attempting to copy contents", from, to)
		log.Debug(renameError.Error())
		log.Debug("going to try copying instead")
		// can't do this rename when running in docker to a file targeted in a mounted volume,
		// so gracefully degrade to copying the entire contents.
		if copyError := copyFileContents(from, to); copyError != nil {
			return fmt.Errorf("failed copying from %v to %v: %w", from, to, copyError)
		}
		tryRemoveTempFile(from)
	}
	return nil
}

func tryRemoveTempFile(filename string) {
	log.Debug("Removing temp file: %v", filename)
	removeErr := os.Remove(filename)
	if removeErr != nil {
		log.Errorf("Failed to remove temp file: %v", filename)
	}
}

// thanks https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
func copyFileContents(src, dst string) (err error) {
	// ignore CWE-22 gosec issue - that's more targeted for http based apps that run in a public directory,
	// and ensuring that it's not possible to give a path to a file outside thar directory.

	in, err := os.Open(src) // #nosec
	if err != nil {
		return err
	}
	defer safelyCloseFile(in)
	out, err := os.Create(dst) // #nosec
	if err != nil {
		return err
	}
	defer safelyCloseFile(out)
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func SafelyCloseReader(reader io.Reader) {
	switch reader := reader.(type) {
	case *os.File:
		safelyCloseFile(reader)
	}
}

func safelyCloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Error("Error closing file!")
		log.Error(err.Error())
	}
}

func createTempFile() (*os.File, error) {
	_, err := os.Stat(os.TempDir())
	if os.IsNotExist(err) {
		err = os.Mkdir(os.TempDir(), 0700)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	file, err := os.CreateTemp("", "temp")
	if err != nil {
		return nil, err
	}

	return file, err
}
