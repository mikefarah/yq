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
		log.Debug("writing doc sep")
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
		if strings.Contains(readline, "$yqDocSeparator$") {
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

func (ye *goccyYamlEncoder) Encode(writer io.Writer, node *CandidateNode) error {
	log.Debug("goccyYamlEncoder - going to print %v", NodeToString(node))
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

	trailingContent := node.FootComment
	if target != nil {
		// Store and clear foot comment to handle it separately
		if astNode, ok := target.(interface{ GetComment() interface{} }); ok {
			_ = astNode // placeholder for potential future foot comment handling
		}
	}

	if err := encoder.Encode(target); err != nil {
		return err
	}

	if err := ye.PrintLeadingContent(destination, trailingContent); err != nil {
		return err
	}

	if ye.prefs.ColorsEnabled {
		return colorizeAndPrint(tempBuffer.Bytes(), writer)
	}
	return nil
}
