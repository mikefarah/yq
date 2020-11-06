package yqlib

import (
	"bufio"
	"container/list"
	"errors"
	"io"
	"os"

	yaml "gopkg.in/yaml.v3"
)

var treeNavigator = NewDataTreeNavigator(NavigationPrefs{})

func readStream(filename string) (io.Reader, error) {
	if filename == "" {
		return nil, errors.New("Must provide filename")
	}

	var stream io.Reader
	if filename == "-" {
		stream = bufio.NewReader(os.Stdin)
	} else {
		file, err := os.Open(filename) // nolint gosec
		if err != nil {
			return nil, err
		}
		defer safelyCloseFile(file)
		stream = file
	}
	return stream, nil
}

func EvaluateStream(filename string, reader io.Reader, node *PathTreeNode) (*list.List, error) {
	var matchingNodes = list.New()

	var currentIndex uint = 0

	decoder := yaml.NewDecoder(reader)
	for {
		var dataBucket yaml.Node
		errorReading := decoder.Decode(&dataBucket)

		if errorReading == io.EOF {
			return matchingNodes, nil
		} else if errorReading != nil {
			return nil, errorReading
		}
		candidateNode := &CandidateNode{
			Document: currentIndex,
			Filename: filename,
			Node:     &dataBucket,
		}
		inputList := list.New()
		inputList.PushBack(candidateNode)

		newMatches, errorParsing := treeNavigator.GetMatchingNodes(inputList, node)
		if errorParsing != nil {
			return nil, errorParsing
		}
		matchingNodes.PushBackList(newMatches)
		currentIndex = currentIndex + 1
	}
}

func Evaluate(filename string, node *PathTreeNode) (*list.List, error) {

	var reader, err = readStream(filename)
	if err != nil {
		return nil, err
	}
	return EvaluateStream(filename, reader, node)

}

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
	in, err := os.Open(src) // nolint gosec
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

func safelyFlush(writer *bufio.Writer) {
	if err := writer.Flush(); err != nil {
		log.Error("Error flushing writer!")
		log.Error(err.Error())
	}

}
func safelyCloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Error("Error closing file!")
		log.Error(err.Error())
	}
}
