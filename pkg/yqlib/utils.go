package yqlib

import (
	"bufio"
	"container/list"
	"io"
	"os"

	yaml "gopkg.in/yaml.v3"
)

func readStream(filename string) (io.Reader, bool, error) {

	if filename == "-" {
		reader := bufio.NewReader(os.Stdin)

		seperatorBytes, err := reader.Peek(3)
		return reader, string(seperatorBytes) == "---", err
	} else {
		// ignore CWE-22 gosec issue - that's more targetted for http based apps that run in a public directory,
		// and ensuring that it's not possible to give a path to a file outside thar directory.
		reader, err := os.Open(filename)
		if err != nil {
			return nil, false, err
		}
		seperatorBytes := make([]byte, 3)
		_, err = reader.Read(seperatorBytes)
		if err != nil {
			return nil, false, err
		}
		_, err = reader.Seek(0, 0)

		return reader, string(seperatorBytes) == "---", err
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
