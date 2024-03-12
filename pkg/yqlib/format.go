package yqlib

import (
	"fmt"
	"strings"
)

type EncoderFactoryFunction func() Encoder
type InDocumentEncoderFactoryFunction func(indent int) Encoder
type DecoderFactoryFunction func() Decoder

type Format struct {
	FormalName       string
	Names            []string
	DefaultExtension string
	EncoderFactory   EncoderFactoryFunction
	DecoderFactory   DecoderFactoryFunction

	/**
	* Like the Encoder Factory, but for encoding content within the document itself.
	* Should turn off colors and other settings to ensure the content comes out right.
	* If this function is not configured, it will default to the EncoderFactory.
	**/
	InDocumentEncoderFactory InDocumentEncoderFactoryFunction
}

var YamlFormat = &Format{"yaml", []string{"y", "yml"}, "yml",
	func() Encoder { return NewYamlEncoder(ConfiguredYamlPreferences) },
	func() Decoder { return NewYamlDecoder(ConfiguredYamlPreferences) },
	func(indent int) Encoder {
		prefs := ConfiguredYamlPreferences.Copy()
		prefs.Indent = indent
		prefs.ColorsEnabled = false
		return NewYamlEncoder(prefs)
	},
}

var JSONFormat = &Format{"json", []string{"j"}, "json",
	func() Encoder { return NewJSONEncoder(ConfiguredJSONPreferences) },
	func() Decoder { return NewJSONDecoder() },
	func(indent int) Encoder {
		prefs := ConfiguredJSONPreferences.Copy()
		prefs.Indent = indent
		prefs.ColorsEnabled = false
		prefs.UnwrapScalar = false
		return NewJSONEncoder(prefs)
	},
}

var CSVFormat = &Format{"csv", []string{"c"}, "csv",
	func() Encoder { return NewCsvEncoder(ConfiguredCsvPreferences) },
	func() Decoder { return NewCSVObjectDecoder(ConfiguredCsvPreferences) },
	nil,
}

var TSVFormat = &Format{"tsv", []string{"t"}, "tsv",
	func() Encoder { return NewCsvEncoder(ConfiguredTsvPreferences) },
	func() Decoder { return NewCSVObjectDecoder(ConfiguredTsvPreferences) },
	nil,
}

var Base64Format = &Format{"base64", []string{}, "txt",
	func() Encoder { return NewBase64Encoder() },
	func() Decoder { return NewBase64Decoder() },
	nil,
}

var UriFormat = &Format{"uri", []string{}, "txt",
	func() Encoder { return NewUriEncoder() },
	func() Decoder { return NewUriDecoder() },
	nil,
}

var ShFormat = &Format{"", nil, "sh",
	func() Encoder { return NewShEncoder() },
	nil,
	nil,
}

var TomlFormat = &Format{"toml", []string{}, "toml",
	func() Encoder { return NewTomlEncoder() },
	func() Decoder { return NewTomlDecoder() },
	nil,
}

var ShellVariablesFormat = &Format{"shell", []string{"s", "sh"}, "sh",
	func() Encoder { return NewShellVariablesEncoder() },
	nil,
	nil,
}

var LuaFormat = &Format{"lua", []string{"l"}, "lua",
	func() Encoder { return NewLuaEncoder(ConfiguredLuaPreferences) },
	func() Decoder { return NewLuaDecoder(ConfiguredLuaPreferences) },
	nil,
}

var Formats = []*Format{
	YamlFormat,
	JSONFormat,
	CSVFormat,
	TSVFormat,
	Base64Format,
	UriFormat,
	ShFormat,
	TomlFormat,
	ShellVariablesFormat,
	LuaFormat,
}

func RegisterFormat(f *Format) {
	Formats = append(Formats, f)
}

func (f *Format) MatchesName(name string) bool {
	if f.FormalName == name {
		return true
	}
	for _, n := range f.Names {
		if n == name {
			return true
		}
	}
	return false
}

func (f *Format) GetInDocumentEncoder(indent int) Encoder {
	if f.InDocumentEncoderFactory != nil {
		return f.InDocumentEncoderFactory(indent)
	}
	return f.EncoderFactory()
}

func FormatStringFromFilename(filename string) string {

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

func FormatFromString(format string) (*Format, error) {
	for _, printerFormat := range Formats {
		if printerFormat.MatchesName(format) {
			return printerFormat, nil
		}
	}
	return nil, fmt.Errorf("unknown format '%v' please use [%v]", format, GetAvailableOutputFormatString())
}

func GetAvailableOutputFormats() []*Format {
	var formats = []*Format{}
	for _, printerFormat := range Formats {
		if printerFormat.EncoderFactory != nil {
			formats = append(formats, printerFormat)
		}
	}
	return formats
}

func GetAvailableOutputFormatString() string {
	var formats = []string{}
	for _, printerFormat := range GetAvailableOutputFormats() {

		if printerFormat.FormalName != "" {
			formats = append(formats, printerFormat.FormalName)
		}
		if len(printerFormat.Names) >= 1 {
			formats = append(formats, printerFormat.Names[0])
		}
	}
	return strings.Join(formats, "|")
}

func GetAvailableInputFormats() []*Format {
	var formats = []*Format{}
	for _, printerFormat := range Formats {
		if printerFormat.DecoderFactory != nil {
			formats = append(formats, printerFormat)
		}
	}
	return formats
}

func GetAvailableInputFormatString() string {
	var formats = []string{}
	for _, printerFormat := range GetAvailableInputFormats() {

		if printerFormat.FormalName != "" {
			formats = append(formats, printerFormat.FormalName)
		}
		if len(printerFormat.Names) >= 1 {
			formats = append(formats, printerFormat.Names[0])
		}
	}
	return strings.Join(formats, "|")
}
