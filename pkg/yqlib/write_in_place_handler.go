package yqlib

import (
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
	// Create the temp file alongside the target so it inherits the target
	// directory's permissions/ACLs (see createTempFile); this keeps the
	// original file's permissions intact after the rename, notably on Windows.
	file, err := createTempFile(filepath.Dir(w.inputFilename))

	if err != nil {
		return nil, err
	}
	info, err := os.Stat(w.inputFilename)
	if err != nil {
		return nil, err
	}
	err = os.Chmod(file.Name(), info.Mode())

	if err != nil {
		return nil, err
	}

	if err = changeOwner(info, file); err != nil {
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
		return tryRenameFile(w.tempFile.Name(), w.inputFilename)
	}
	tryRemoveTempFile(w.tempFile.Name())

	return nil
}
