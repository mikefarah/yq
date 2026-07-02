package yqlib

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
)

var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

// stripUTF8BOM returns a reader that skips a leading UTF-8 BOM, if present.
func stripUTF8BOM(r io.Reader) io.Reader {
	br := bufio.NewReader(r)

	peek, err := br.Peek(3)
	if err == nil && bytes.Equal(peek, utf8BOM) {
		_, _ = br.Discard(3)
	}

	return br
}

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
		reader = bufio.NewReader(stripUTF8BOM(os.Stdin))
	} else {
		file, err := os.Open(f.originalFilename) // #nosec
		if err != nil {
			return err
		}
		reader = bufio.NewReader(stripUTF8BOM(file))
	}
	f.contentReader = reader

	yamlTempFile, err := createTempFile()
	if err != nil {
		return err
	}
	f.yamlFrontMatterFilename = yamlTempFile.Name()
	log.Debugf("yamlTempFile: %v", yamlTempFile.Name())

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
