package yqlib

import (
	"os"
)

type writeInPlaceHandler interface {
	CreateTempFile() (*os.File, error)
	FinishWriteInPlace(evaluatedSuccessfully bool)
}

type writeInPlaceHandlerImpl struct {
	inputFilename string
	tempFile      *os.File
}

func NewWriteInPlaceHandler(inputFile string) writeInPlaceHandler {
	return &writeInPlaceHandlerImpl{inputFile, nil}
}

func (w *writeInPlaceHandlerImpl) CreateTempFile() (*os.File, error) {
	file, err := createTempFile()
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
	log.Debug("WriteInPlaceHandler: writing to tempfile: %v", file.Name())
	w.tempFile = file
	return file, err
}

func (w *writeInPlaceHandlerImpl) FinishWriteInPlace(evaluatedSuccessfully bool) {
	log.Debug("Going to write-inplace, evaluatedSuccessfully=%v, target=%v", evaluatedSuccessfully, w.inputFilename)
	safelyCloseFile(w.tempFile)
	if evaluatedSuccessfully {
		log.Debug("Moving temp file to target")
		safelyRenameFile(w.tempFile.Name(), w.inputFilename)
	} else {
		tryRemoveFile(w.tempFile.Name())
	}
}
