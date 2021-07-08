package yqlib

import (
	"bufio"
	"container/list"
	"io"
	"os"

	yaml "gopkg.in/yaml.v3"
)

func readStream(filename string) (io.Reader, error) {
	if filename == "-" {
		return bufio.NewReader(os.Stdin), nil
	} else {
		// ignore CWE-22 gosec issue - that's more targetted for http based apps that run in a public directory,
		// and ensuring that it's not possible to give a path to a file outside thar directory.
		return os.Open(filename) // #nosec
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
