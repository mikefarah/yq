package yqlib

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/fatih/color"
)

type Encoder interface {
	Encode(writer io.Writer, node *CandidateNode) error
	PrintDocumentSeparator(writer io.Writer) error
	PrintLeadingContent(writer io.Writer, content string) error
	CanHandleAliases() bool
}

func mapKeysToStrings(node *CandidateNode) {

	if node.Kind == MappingNode {
		for index, child := range node.Content {
			if index%2 == 0 { // its a map key
				child.Tag = "!!str"
			}
		}
	}

	for _, child := range node.Content {
		mapKeysToStrings(child)
	}
}

// Some funcs are shared between encoder_yaml and encoder_kyaml
func PrintYAMLDocumentSeparator(writer io.Writer, PrintDocSeparators bool) error {
	if PrintDocSeparators {
		log.Debug("writing doc sep")
		if err := writeString(writer, "---\n"); err != nil {
			return err
		}
	}
	return nil
}
func PrintYAMLLeadingContent(writer io.Writer, content string, PrintDocSeparators bool, ColorsEnabled bool) error {
	reader := bufio.NewReader(strings.NewReader(content))

	// reuse precompiled package-level regex
	// (declared in decoder_yaml.go)

	for {

		readline, errReading := reader.ReadString('\n')
		if errReading != nil && !errors.Is(errReading, io.EOF) {
			return errReading
		}
		if strings.Contains(readline, "$yqDocSeparator$") {
			// Preserve the original line ending (CRLF or LF)
			lineEnding := "\n"
			if strings.HasSuffix(readline, "\r\n") {
				lineEnding = "\r\n"
			}
			if PrintDocSeparators {
				if err := writeString(writer, "---"+lineEnding); err != nil {
					return err
				}
			}

		} else {
			if len(readline) > 0 && readline != "\n" && readline[0] != '%' && !commentLineRe.MatchString(readline) {
				readline = "# " + readline
			}
			if ColorsEnabled && strings.TrimSpace(readline) != "" {
				readline = format(color.FgHiBlack) + readline + format(color.Reset)
			}
			if err := writeString(writer, readline); err != nil {
				return err
			}
		}

		if errors.Is(errReading, io.EOF) {
			if readline != "" {
				// the last comment we read didn't have a newline, put one in
				if err := writeString(writer, "\n"); err != nil {
					return err
				}
			}
			break
		}
	}

	return nil
}
