package yqlib

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type yamlEncoder struct {
	indent   int
	colorise bool
	prefs    YamlPreferences
}

func NewYamlEncoder(indent int, colorise bool, prefs YamlPreferences) Encoder {
	if indent < 0 {
		indent = 0
	}
	return &yamlEncoder{indent, colorise, prefs}
}

func (ye *yamlEncoder) CanHandleAliases() bool {
	return true
}

func (ye *yamlEncoder) PrintDocumentSeparator(writer io.Writer) error {
	if ye.prefs.PrintDocSeparators {
		log.Debug("-- writing doc sep")
		if err := writeString(writer, "---\n"); err != nil {
			return err
		}
	}
	return nil
}

func (ye *yamlEncoder) PrintLeadingContent(writer io.Writer, content string) error {
	// log.Debug("headcommentwas [%v]", content)
	reader := bufio.NewReader(strings.NewReader(content))

	for {

		readline, errReading := reader.ReadString('\n')
		if errReading != nil && !errors.Is(errReading, io.EOF) {
			return errReading
		}
		if strings.Contains(readline, "$yqDocSeperator$") {

			if err := ye.PrintDocumentSeparator(writer); err != nil {
				return err
			}

		} else {
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

func (ye *yamlEncoder) Encode(writer io.Writer, node *yaml.Node) error {

	if node.Kind == yaml.ScalarNode && ye.prefs.UnwrapScalar {
		return writeString(writer, node.Value+"\n")
	}

	destination := writer
	tempBuffer := bytes.NewBuffer(nil)
	if ye.colorise {
		destination = tempBuffer
	}

	var encoder = yaml.NewEncoder(destination)

	encoder.SetIndent(ye.indent)

	if err := encoder.Encode(node); err != nil {
		return err
	}

	if ye.colorise {
		return colorizeAndPrint(tempBuffer.Bytes(), writer)
	}
	return nil
}
