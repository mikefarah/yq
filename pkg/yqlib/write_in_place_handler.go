package yqlib

import (
	"io/ioutil"
	"os"
)

type WriteInPlaceHandler interface {
	CreateTempFile() (*os.File, error)
	FinishWriteInPlace(evaluatedSuccessfully bool)
}

type writeInPlaceHandler struct {
	inputFilename string
	tempFile      *os.File
}

func NewWriteInPlaceHandler(inputFile string) WriteInPlaceHandler {

	return &writeInPlaceHandler{inputFile, nil}
}

func (w *writeInPlaceHandler) CreateTempFile() (*os.File, error) {
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
	w.tempFile = file
	return file, err
}

func (w *writeInPlaceHandler) FinishWriteInPlace(evaluatedSuccessfully bool) {
	safelyCloseFile(w.tempFile)
	if evaluatedSuccessfully {
		safelyRenameFile(w.tempFile.Name(), w.inputFilename)
	} else {
		os.Remove(w.tempFile.Name())
	}
}
