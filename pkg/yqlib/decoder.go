package yqlib

import (
	"fmt"
	"io"
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
)

type Decoder interface {
	Init(reader io.Reader) error
	Decode() (*CandidateNode, error)
}

func InputFormatFromString(format string) (InputFormat, error) {
	switch format {
	case "yaml", "y":
		return YamlInputFormat, nil
	case "xml", "x":
		return XMLInputFormat, nil
	case "props", "p":
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
