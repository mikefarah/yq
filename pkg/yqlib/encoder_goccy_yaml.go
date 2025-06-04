package yqlib

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"regexp"
	"strings"

	"github.com/fatih/color"
	yaml "github.com/goccy/go-yaml"
)

// Note: The constant 'docSeparatorPlaceholder' is defined in decoder_goccy_yaml.go
// and is used here as they are in the same package.

type goccyYamlEncoder struct {
	prefs YamlPreferences
}

func NewGoccyYamlEncoder(prefs YamlPreferences) Encoder {
	return &goccyYamlEncoder{prefs}
}

func (ye *goccyYamlEncoder) CanHandleAliases() bool {
	return true
}

func (ye *goccyYamlEncoder) PrintDocumentSeparator(writer io.Writer) error {
	if ye.prefs.PrintDocSeparators {
		// log.Debug("writing doc sep") // Commented out
		if err := writeString(writer, "---\n"); err != nil {
			return err
		}
	}
	return nil
}

func (ye *goccyYamlEncoder) PrintLeadingContent(writer io.Writer, content string) error {
	reader := bufio.NewReader(strings.NewReader(content))

	var commentLineRegEx = regexp.MustCompile(`^\s*#`)

	for {
		readline, errReading := reader.ReadString('\n')
		if errReading != nil && !errors.Is(errReading, io.EOF) {
			return errReading
		}
		if strings.Contains(readline, docSeparatorPlaceholder) { // Use the constant from decoder_goccy_yaml.go
			if err := ye.PrintDocumentSeparator(writer); err != nil {
				return err
			}
		} else {
			if len(readline) > 0 && readline != "\n" && readline[0] != '%' && !commentLineRegEx.MatchString(readline) {
				readline = "# " + readline
			}
			if ye.prefs.ColorsEnabled && strings.TrimSpace(readline) != "" {
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

// PrintTrailingComment writes the given comment string to the writer.
// If the comment is not empty and does not end with a newline, one is added.
// It does not add any comment prefixes (#).
func (ye *goccyYamlEncoder) PrintTrailingComment(writer io.Writer, comment string) error {
	if comment == "" {
		return nil
	}
	if err := writeString(writer, comment); err != nil {
		return err
	}
	if !strings.HasSuffix(comment, "\n") {
		if err := writeString(writer, "\n"); err != nil {
			return err
		}
	}
	return nil
}

func (ye *goccyYamlEncoder) Encode(writer io.Writer, node *CandidateNode) error {
	// log.Debug("goccyYamlEncoder - going to print %v", NodeToString(node)) // Commented out
	if node.Kind == ScalarNode && ye.prefs.UnwrapScalar {
		valueToPrint := node.Value
		if node.LeadingContent == "" || valueToPrint != "" {
			valueToPrint = valueToPrint + "\n"
		}
		return writeString(writer, valueToPrint)
	}

	destination := writer
	tempBuffer := bytes.NewBuffer(nil)
	if ye.prefs.ColorsEnabled {
		destination = tempBuffer
	}

	// Create encoder with indent option
	var encoder *yaml.Encoder
	if ye.prefs.Indent > 0 {
		encoder = yaml.NewEncoder(destination, yaml.Indent(ye.prefs.Indent))
	} else {
		encoder = yaml.NewEncoder(destination)
	}

	target, err := node.MarshalGoccyYAML()
	if err != nil {
		return err
	}

	if err := encoder.Encode(target); err != nil {
		return err
	}

	// Use PrintTrailingComment for foot comments.
	if node.FootComment != "" { // Check if there's a foot comment to print
		if err := ye.PrintTrailingComment(destination, node.FootComment); err != nil {
			return err
		}
	}

	if ye.prefs.ColorsEnabled {
		return colorizeAndPrint(tempBuffer.Bytes(), writer)
	}
	return nil
}
