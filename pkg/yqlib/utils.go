package yqlib

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"io"
	"os"
)

func readStream(filename string) (io.Reader, error) {
	var reader *bufio.Reader
	if filename == "-" {
		reader = bufio.NewReader(os.Stdin)
	} else {
		// ignore CWE-22 gosec issue - that's more targeted for http based apps that run in a public directory,
		// and ensuring that it's not possible to give a path to a file outside that directory.
		file, err := os.Open(filename) // #nosec
		if err != nil {
			return nil, err
		}
		reader = bufio.NewReader(file)
	}
	return reader, nil

}

func writeString(writer io.Writer, txt string) error {
	_, errorWriting := writer.Write([]byte(txt))
	return errorWriting
}

func ReadDocuments(reader io.Reader, decoder Decoder) (*list.List, error) {
	return readDocuments(reader, "", 0, decoder)
}

func readDocuments(reader io.Reader, filename string, fileIndex int, decoder Decoder) (*list.List, error) {
	err := decoder.Init(reader)
	if err != nil {
		return nil, err
	}
	inputList := list.New()
	var currentIndex uint

	for {
		candidateNode, errorReading := decoder.Decode()

		if errors.Is(errorReading, io.EOF) {
			switch reader := reader.(type) {
			case *os.File:
				safelyCloseFile(reader)
			}
			return inputList, nil
		} else if errorReading != nil {
			return nil, fmt.Errorf("bad file '%v': %w", filename, errorReading)
		}
		candidateNode.document = currentIndex
		candidateNode.filename = filename
		candidateNode.fileIndex = fileIndex
		candidateNode.EvaluateTogether = true

		inputList.PushBack(candidateNode)

		currentIndex = currentIndex + 1
	}
}
