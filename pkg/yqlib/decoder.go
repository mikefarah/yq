package yqlib

import (
	"fmt"
	"io"
	"strings"
)

type InputFormat uint

const (
	YamlInputFormat = 1 << iota
	XMLInputFormat
	PropertiesInputFormat
	Base64InputFormat
	JsonInputFormat
	CSVObjectInputFormat
	TSVObjectInputFormat
	UriInputFormat
)

type Decoder interface {
	Init(reader io.Reader) error
	Decode() (*CandidateNode, error)
}

func InputFormatFromString(format string) (InputFormat, error) {
	switch format {
	case "yaml", "yml", "y":
		return YamlInputFormat, nil
	case "xml", "x":
		return XMLInputFormat, nil
	case "properties", "props", "p":
		return PropertiesInputFormat, nil
	case "json", "ndjson", "j":
		return JsonInputFormat, nil
	case "csv", "c":
		return CSVObjectInputFormat, nil
	case "tsv", "t":
		return TSVObjectInputFormat, nil
	default:
		return 0, fmt.Errorf("unknown format '%v' please use [yaml|xml|props]", format)
	}
}

func InputFormatFromFilename(filename string, defaultFormat string) string {
	if filename != "" {
		GetLogger().Debugf("checking filename '%s' for inputFormat", filename)
		nPos := strings.LastIndex(filename, ".")
		if nPos > -1 {
			inputFormat := filename[nPos+1:]
			GetLogger().Debugf("detected inputFormat '%s'", inputFormat)
			return inputFormat
		}
	}

	GetLogger().Debugf("using default inputFormat '%s'", defaultFormat)
	return defaultFormat
}
