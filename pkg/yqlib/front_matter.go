package yqlib

import (
	"bufio"
	"errors"
	"io"
	"os"
)

type frontMatterHandler interface {
	Split() error
	GetYamlFrontMatterFilename() string
	GetContentReader() io.Reader
	CleanUp()
}

type frontMatterHandlerImpl struct {
	originalFilename        string
	yamlFrontMatterFilename string
	contentReader           io.Reader
}

func NewFrontMatterHandler(originalFilename string) frontMatterHandler {
	return &frontMatterHandlerImpl{originalFilename, "", nil}
}

func (f *frontMatterHandlerImpl) GetYamlFrontMatterFilename() string {
	return f.yamlFrontMatterFilename
}

func (f *frontMatterHandlerImpl) GetContentReader() io.Reader {
	return f.contentReader
}

func (f *frontMatterHandlerImpl) CleanUp() {
	tryRemoveTempFile(f.yamlFrontMatterFilename)
}

// Splits the given file by yaml front matter
// yaml content will be saved to first temporary file
// remaining content will be saved to second temporary file
func (f *frontMatterHandlerImpl) Split() error {
	var reader *bufio.Reader
	var err error
	if f.originalFilename == "-" {
		reader = bufio.NewReader(os.Stdin)
	} else {
		file, err := os.Open(f.originalFilename) // #nosec
		if err != nil {
			return err
		}
		reader = bufio.NewReader(file)
	}
	f.contentReader = reader

	yamlTempFile, err := createTempFile()
	if err != nil {
		return err
	}
	f.yamlFrontMatterFilename = yamlTempFile.Name()
	log.Debug("yamlTempFile: %v", yamlTempFile.Name())

	lineCount := 0

	for {
		peekBytes, err := reader.Peek(3)
		if errors.Is(err, io.EOF) {
			// we've finished reading the yaml content..I guess
			break
		} else if err != nil {
			return err
		}
		if lineCount > 0 && string(peekBytes) == "---" {
			// we've finished reading the yaml content..
			break
		}
		line, errReading := reader.ReadString('\n')
		lineCount = lineCount + 1
		if errReading != nil && !errors.Is(errReading, io.EOF) {
			return errReading
		}

		_, errWriting := yamlTempFile.WriteString(line)

		if errWriting != nil {
			return errWriting
		}
	}

	safelyCloseFile(yamlTempFile)

	return nil

}
