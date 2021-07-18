package yqlib

import (
	"bufio"
	"io"
	"os"
)

type frontMatterHandler interface {
	Split() error
	GetYamlFrontMatterFilename() string
	GetContentFilename() string
	CleanUp()
}

type frontMatterHandlerImpl struct {
	originalFilename        string
	yamlFrontMatterFilename string
	contentFilename         string
}

func NewFrontMatterHandler(originalFilename string) frontMatterHandler {
	return &frontMatterHandlerImpl{originalFilename, "", ""}
}

func (f *frontMatterHandlerImpl) GetYamlFrontMatterFilename() string {
	return f.yamlFrontMatterFilename
}

func (f *frontMatterHandlerImpl) GetContentFilename() string {
	return f.contentFilename
}

func (f *frontMatterHandlerImpl) CleanUp() {
	tryRemoveFile(f.yamlFrontMatterFilename)
	tryRemoveFile(f.contentFilename)
}

// Splits the given file by yaml front matter
// yaml content will be saved to first temporary file
// remaining content will be saved to second temporary file
func (f *frontMatterHandlerImpl) Split() error {
	var reader io.Reader
	var err error
	if f.originalFilename == "-" {
		reader = bufio.NewReader(os.Stdin)
	} else {
		reader, err = os.Open(f.originalFilename) // #nosec
		if err != nil {
			return err
		}
	}

	yamlTempFile, err := createTempFile()
	if err != nil {
		return err
	}
	f.yamlFrontMatterFilename = yamlTempFile.Name()
	log.Debug("yamlTempFile: %v", yamlTempFile.Name())

	contentTempFile, err := createTempFile()
	if err != nil {
		return err
	}
	f.contentFilename = contentTempFile.Name()
	log.Debug("contentTempFile: %v", contentTempFile.Name())

	scanner := bufio.NewScanner(reader)

	lineCount := 0
	yamlContentBlock := true

	for scanner.Scan() {
		line := scanner.Text()

		if lineCount > 0 && line == "---" {
			//we've finished reading the yaml content
			yamlContentBlock = false
		}
		if yamlContentBlock {
			_, err = yamlTempFile.Write([]byte(line + "\n"))
		} else {
			_, err = contentTempFile.Write([]byte(line + "\n"))
		}
		if err != nil {
			return err
		}
		lineCount = lineCount + 1
	}

	safelyCloseFile(yamlTempFile)
	safelyCloseFile(contentTempFile)

	return scanner.Err()

}
