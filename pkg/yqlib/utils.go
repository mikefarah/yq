package yqlib

import (
	"bufio"
	"container/list"
	"io"
	"os"
	"path/filepath"


	yaml "gopkg.in/yaml.v3"
)

func readStream(filename string) (io.Reader, error) {
	if filename == "-" {
		return bufio.NewReader(os.Stdin), nil
	} else {
		return os.Open(filepath.Clean(filename)) // nolint gosec
	}
}

func readDocuments(reader io.Reader, filename string, fileIndex int) (*list.List, error) {
	decoder := yaml.NewDecoder(reader)
	inputList := list.New()
	var currentIndex uint = 0

	for {
		var dataBucket yaml.Node
		errorReading := decoder.Decode(&dataBucket)

		if errorReading == io.EOF {
			switch reader := reader.(type) {
			case *os.File:
				safelyCloseFile(reader)
			}
			return inputList, nil
		} else if errorReading != nil {
			return nil, errorReading
		}
		candidateNode := &CandidateNode{
			Document:         currentIndex,
			Filename:         filename,
			Node:             &dataBucket,
			FileIndex:        fileIndex,
			EvaluateTogether: true,
		}

		inputList.PushBack(candidateNode)

		currentIndex = currentIndex + 1
	}
}
