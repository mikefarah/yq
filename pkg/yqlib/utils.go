package yqlib

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

func readStream(filename string, leadingContentPreProcessing bool) (io.Reader, string, error) {
	var reader *bufio.Reader
	if filename == "-" {
		reader = bufio.NewReader(os.Stdin)
	} else {
		// ignore CWE-22 gosec issue - that's more targeted for http based apps that run in a public directory,
		// and ensuring that it's not possible to give a path to a file outside thar directory.
		file, err := os.Open(filename) // #nosec
		if err != nil {
			return nil, "", err
		}
		reader = bufio.NewReader(file)
	}

	if !leadingContentPreProcessing {
		return reader, "", nil
	}
	return processReadStream(reader)
}

func writeString(writer io.Writer, txt string) error {
	_, errorWriting := writer.Write([]byte(txt))
	return errorWriting
}

func processLeadingContent(mappedDoc *CandidateNode, writer io.Writer, printDocSeparators bool, outputFormat PrinterOutputFormat) error {
	log.Debug("headcommentwas %v", mappedDoc.LeadingContent)
	log.Debug("finished headcomment")
	reader := bufio.NewReader(strings.NewReader(mappedDoc.LeadingContent))

	for {

		readline, errReading := reader.ReadString('\n')
		if errReading != nil && !errors.Is(errReading, io.EOF) {
			return errReading
		}
		if strings.Contains(readline, "$yqDocSeperator$") {
			if printDocSeparators {
				if err := writeString(writer, "---\n"); err != nil {
					return err
				}
			}
		} else if outputFormat == YamlOutputFormat {
			if err := writeString(writer, readline); err != nil {
				return err
			}
		}

		if errors.Is(errReading, io.EOF) {
			if readline != "" {
				// the last comment we read didn't have a new line, put one in
				if err := writeString(writer, "\n"); err != nil {
					return err
				}
			}
			break
		}
	}

	return nil
}

func processReadStream(reader *bufio.Reader) (io.Reader, string, error) {
	var commentLineRegEx = regexp.MustCompile(`^\s*#`)
	var sb strings.Builder
	for {
		peekBytes, err := reader.Peek(3)
		if errors.Is(err, io.EOF) {
			// EOF are handled else where..
			return reader, sb.String(), nil
		} else if err != nil {
			return reader, sb.String(), err
		} else if string(peekBytes) == "---" {
			_, err := reader.ReadString('\n')
			sb.WriteString("$yqDocSeperator$\n")
			if errors.Is(err, io.EOF) {
				return reader, sb.String(), nil
			} else if err != nil {
				return reader, sb.String(), err
			}
		} else if commentLineRegEx.MatchString(string(peekBytes)) {
			line, err := reader.ReadString('\n')
			sb.WriteString(line)
			if errors.Is(err, io.EOF) {
				return reader, sb.String(), nil
			} else if err != nil {
				return reader, sb.String(), err
			}
		} else {
			return reader, sb.String(), nil
		}
	}
}

func readDocuments(reader io.Reader, filename string, fileIndex int) (*list.List, error) {
	decoder := yaml.NewDecoder(reader)
	inputList := list.New()
	var currentIndex uint

	for {
		var dataBucket yaml.Node
		errorReading := decoder.Decode(&dataBucket)

		if errors.Is(errorReading, io.EOF) {
			switch reader := reader.(type) {
			case *os.File:
				safelyCloseFile(reader)
			}
			return inputList, nil
		} else if errorReading != nil {
			return nil, fmt.Errorf("bad file '%v': %w", filename, errorReading)
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
