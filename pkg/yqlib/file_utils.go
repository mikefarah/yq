package yqlib

import (
	"io"
	"os"
	"path/filepath"
)

func safelyRenameFile(from string, to string) {
	if renameError := os.Rename(from, to); renameError != nil {
		log.Debugf("Error renaming from %v to %v, attempting to copy contents", from, to)
		log.Debug(renameError.Error())
		// can't do this rename when running in docker to a file targeted in a mounted volume,
		// so gracefully degrade to copying the entire contents.
		if copyError := copyFileContents(from, to); copyError != nil {
			log.Errorf("Failed copying from %v to %v", from, to)
			log.Error(copyError.Error())
		} else {
			removeErr := os.Remove(from)
			if removeErr != nil {
				log.Errorf("failed removing original file: %s", from)
			}
		}
	}
}

// thanks https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(filepath.Clean(src))
	if err != nil {
		return err
	}
	defer safelyCloseFile(in)
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer safelyCloseFile(out)
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func safelyCloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Error("Error closing file!")
		log.Error(err.Error())
	}
}
