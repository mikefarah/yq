package yqlib

import (
	"bufio"
	"container/list"
	"io"
	"os"

	yaml "gopkg.in/yaml.v3"
)

//TODO: convert to interface + struct

var treeNavigator = NewDataTreeNavigator()
var treeCreator = NewPathTreeCreator()

func readStream(filename string) (io.Reader, error) {
	if filename == "-" {
		return bufio.NewReader(os.Stdin), nil
	} else {
		return os.Open(filename) // nolint gosec
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
			Document:  currentIndex,
			Filename:  filename,
			Node:      &dataBucket,
			FileIndex: fileIndex,
		}

		inputList.PushBack(candidateNode)

		currentIndex = currentIndex + 1
	}
}

// func safelyRenameFile(from string, to string) {
// 	if renameError := os.Rename(from, to); renameError != nil {
// 		log.Debugf("Error renaming from %v to %v, attempting to copy contents", from, to)
// 		log.Debug(renameError.Error())
// 		// can't do this rename when running in docker to a file targeted in a mounted volume,
// 		// so gracefully degrade to copying the entire contents.
// 		if copyError := copyFileContents(from, to); copyError != nil {
// 			log.Errorf("Failed copying from %v to %v", from, to)
// 			log.Error(copyError.Error())
// 		} else {
// 			removeErr := os.Remove(from)
// 			if removeErr != nil {
// 				log.Errorf("failed removing original file: %s", from)
// 			}
// 		}
// 	}
// }

// thanks https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
// func copyFileContents(src, dst string) (err error) {
// 	in, err := os.Open(src) // nolint gosec
// 	if err != nil {
// 		return err
// 	}
// 	defer safelyCloseFile(in)
// 	out, err := os.Create(dst)
// 	if err != nil {
// 		return err
// 	}
// 	defer safelyCloseFile(out)
// 	if _, err = io.Copy(out, in); err != nil {
// 		return err
// 	}
// 	return out.Sync()
// }

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
