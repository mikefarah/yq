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
	TomlInputFormat
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
	case "toml":
		return TomlInputFormat, nil
	default:
		return 0, fmt.Errorf("unknown format '%v' please use [yaml|json|props|csv|tsv|xml|toml]", format)
	}
}

func FormatFromFilename(filename string) string {

	if filename != "" {
		GetLogger().Debugf("checking file extension '%s' for auto format detection", filename)
		nPos := strings.LastIndex(filename, ".")
		if nPos > -1 {
			format := filename[nPos+1:]
			GetLogger().Debugf("detected format '%s'", format)
			return format
		}
	}

	GetLogger().Debugf("using default inputFormat 'yaml'")
	return "yaml"
}
