package yqlib

import (
	"fmt"
	"os"
	"path/filepath"
)

type writeInPlaceHandler interface {
	CreateTempFile() (*os.File, error)
	FinishWriteInPlace(evaluatedSuccessfully bool) error
}

type writeInPlaceHandlerImpl struct {
	inputFilename string
	tempFile      *os.File
}

func NewWriteInPlaceHandler(inputFile string) writeInPlaceHandler {

	return &writeInPlaceHandlerImpl{inputFile, nil}
}

func (w *writeInPlaceHandlerImpl) CreateTempFile() (*os.File, error) {
	entryInfo, err := os.Lstat(w.inputFilename)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(w.inputFilename)
	if err != nil {
		return nil, err
	}
	if entryInfo.Mode()&os.ModeSymlink == 0 && info.Mode().Perm()&0o222 == 0 {
		return nil, fmt.Errorf("cannot write in place to read-only file %q", w.inputFilename)
	}

	file, err := os.CreateTemp(filepath.Dir(w.inputFilename), ".yq-inplace-")
	if err != nil {
		return nil, err
	}
	err = os.Chmod(file.Name(), info.Mode())

	if err != nil {
		_ = file.Close()
		_ = os.Remove(file.Name())
		return nil, err
	}

	if err = changeOwner(info, file); err != nil {
		_ = file.Close()
		_ = os.Remove(file.Name())
		return nil, err
	}
	log.Debugf("WriteInPlaceHandler: writing to tempfile: %v", file.Name())
	w.tempFile = file
	return file, err
}

func (w *writeInPlaceHandlerImpl) FinishWriteInPlace(evaluatedSuccessfully bool) error {
	log.Debugf("Going to write in place, evaluatedSuccessfully=%v, target=%v", evaluatedSuccessfully, w.inputFilename)
	safelyCloseFile(w.tempFile)
	if evaluatedSuccessfully {
		log.Debug("Moving temp file to target")
		if err := tryRenameFile(w.tempFile.Name(), w.inputFilename); err != nil {
			tryRemoveTempFile(w.tempFile.Name())
			return err
		}
		return nil
	}
	tryRemoveTempFile(w.tempFile.Name())

	return nil
}
