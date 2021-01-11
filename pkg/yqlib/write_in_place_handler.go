package yqlib

import (
	"io/ioutil"
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
	info, err := os.Stat(w.inputFilename)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(os.TempDir())
	if os.IsNotExist(err) {
		err = os.Mkdir(os.TempDir(), 0700)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	file, err := ioutil.TempFile("", "temp")
	if err != nil {
		return nil, err
	}

	err = os.Chmod(file.Name(), info.Mode())
	log.Debug("writing to tempfile: %v", file.Name())
	w.tempFile = file
	return file, err
}

func (w *writeInPlaceHandlerImpl) FinishWriteInPlace(evaluatedSuccessfully bool) {
	log.Debug("Going to write-inplace, evaluatedSuccessfully=%v, target=%v", evaluatedSuccessfully, w.inputFilename)
	safelyCloseFile(w.tempFile)
	if evaluatedSuccessfully {
		log.Debug("moved temp file to target")
		safelyRenameFile(w.tempFile.Name(), w.inputFilename)
	} else {
		log.Debug("removed temp file")
		os.Remove(w.tempFile.Name())
	}
}
