package yqlib

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

type yamlEncoder struct {
	prefs YamlPreferences
}

func NewYamlEncoder(prefs YamlPreferences) Encoder {
	return &yamlEncoder{prefs}
}

func (ye *yamlEncoder) CanHandleAliases() bool {
	return true
}

func (ye *yamlEncoder) PrintDocumentSeparator(writer io.Writer) error {
	if ye.prefs.PrintDocSeparators {
		log.Debug("writing doc sep")
		if err := writeString(writer, "---\n"); err != nil {
			return err
		}
	}
	return nil
}

func (ye *yamlEncoder) PrintLeadingContent(writer io.Writer, content string) error {
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

func (ye *yamlEncoder) Encode(writer io.Writer, node *CandidateNode) error {
	log.Debug("encoderYaml - going to print %v", NodeToString(node))
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

	var encoder = yaml.NewEncoder(destination)

	encoder.SetIndent(ye.prefs.Indent)

	target, err := node.MarshalYAML()

	if err != nil {
		return err
	}

	trailingContent := target.FootComment
	target.FootComment = ""

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
