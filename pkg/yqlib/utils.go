package yqlib

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// filenameAliases maps real file paths to display names.
// Used by front matter handling to preserve original filenames
// when the actual content is read from temporary files.
var filenameAliases = map[string]string{}

// SetFilenameAlias registers a display name for a file path so that
// the filename operator returns the original name instead of a temp path.
func SetFilenameAlias(realPath string, displayName string) {
	filenameAliases[realPath] = displayName
}

// ClearFilenameAliases removes all filename aliases.
func ClearFilenameAliases() {
	filenameAliases = map[string]string{}
}

func resolveFilename(filename string) string {
	if alias, ok := filenameAliases[filename]; ok {
		return alias
	}
	return filename
}

// readStream returns a reader for the given file, along with a cleanup function
// that must be called once the reader is no longer needed. The cleanup is a no-op
// for stdin.
func readStream(filename string) (io.Reader, func(), error) {
	if filename == "-" {
		return bufio.NewReader(os.Stdin), func() {}, nil
	}
	// ignore CWE-22 gosec issue - that's more targeted for http based apps that run in a public directory,
	// and ensuring that it's not possible to give a path to a file outside that directory.
	file, err := os.Open(filename) // #nosec
	if err != nil {
		return nil, nil, err
	}
	return bufio.NewReader(file), func() { safelyCloseFile(file) }, nil
}

func writeString(writer io.Writer, txt string) error {
	_, errorWriting := io.WriteString(writer, txt)
	return errorWriting
}

func ReadDocuments(reader io.Reader, decoder Decoder) (*list.List, error) {
	return readDocuments(reader, "", 0, decoder)
}

func readDocuments(reader io.Reader, filename string, fileIndex int, decoder Decoder) (*list.List, error) {
	filename = resolveFilename(filename)
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

// sanitizeControlChars removes non-printable ASCII control characters
// (code points < 32 and 127) except tab, newline, and carriage return
// to prevent terminal escape sequence injection in direct output encoders.
func sanitizeControlChars(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '\t' || r == '\n' || r == '\r' {
			return r
		}
		if (r >= 0 && r < 32) || r == 127 {
			return -1
		}
		return r
	}, s)
}

