package yqlib

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
)

type EncoderFactoryFunction func() Encoder
type DecoderFactoryFunction func() Decoder

type Format struct {
	FormalName     string
	Names          []string
	EncoderFactory EncoderFactoryFunction
	DecoderFactory DecoderFactoryFunction
}

var YamlFormat = &Format{"yaml", []string{"y", "yml"},
	func() Encoder { return NewYamlEncoder(ConfiguredYamlPreferences) },
	func() Decoder { return NewYamlDecoder(ConfiguredYamlPreferences) },
}

var KYamlFormat = &Format{"kyaml", []string{"ky"},
	func() Encoder { return NewKYamlEncoder(ConfiguredKYamlPreferences) },
	// KYaml is stricter YAML
	func() Decoder { return NewYamlDecoder(ConfiguredYamlPreferences) },
}

var JSONFormat = &Format{"json", []string{"j"},
	func() Encoder { return NewJSONEncoder(ConfiguredJSONPreferences) },
	func() Decoder { return NewJSONDecoder() },
}

var PropertiesFormat = &Format{"props", []string{"p", "properties"},
	func() Encoder { return NewPropertiesEncoder(ConfiguredPropertiesPreferences) },
	func() Decoder { return NewPropertiesDecoder() },
}

var CSVFormat = &Format{"csv", []string{"c"},
	func() Encoder { return NewCsvEncoder(ConfiguredCsvPreferences) },
	func() Decoder { return NewCSVObjectDecoder(ConfiguredCsvPreferences) },
}

var TSVFormat = &Format{"tsv", []string{"t"},
	func() Encoder { return NewCsvEncoder(ConfiguredTsvPreferences) },
	func() Decoder { return NewCSVObjectDecoder(ConfiguredTsvPreferences) },
}

var XMLFormat = &Format{"xml", []string{"x"},
	func() Encoder { return NewXMLEncoder(ConfiguredXMLPreferences) },
	func() Decoder { return NewXMLDecoder(ConfiguredXMLPreferences) },
}

var Base64Format = &Format{"base64", []string{},
	func() Encoder { return NewBase64Encoder() },
	func() Decoder { return NewBase64Decoder() },
}

var UriFormat = &Format{"uri", []string{},
	func() Encoder { return NewUriEncoder() },
	func() Decoder { return NewUriDecoder() },
}

var ShFormat = &Format{"", nil,
	func() Encoder { return NewShEncoder() },
	nil,
}

var TomlFormat = &Format{"toml", []string{},
	func() Encoder { return NewTomlEncoderWithPrefs(ConfiguredTomlPreferences) },
	func() Decoder { return NewTomlDecoder() },
}

var HclFormat = &Format{"hcl", []string{"h", "tf"},
	func() Encoder { return NewHclEncoder(ConfiguredHclPreferences) },
	func() Decoder { return NewHclDecoder() },
}

var ShellVariablesFormat = &Format{"shell", []string{"s", "sh"},
	func() Encoder { return NewShellVariablesEncoder() },
	nil,
}

var LuaFormat = &Format{"lua", []string{"l"},
	func() Encoder { return NewLuaEncoder(ConfiguredLuaPreferences) },
	func() Decoder { return NewLuaDecoder(ConfiguredLuaPreferences) },
}

var INIFormat = &Format{"ini", []string{"i"},
	func() Encoder { return NewINIEncoder() },
	func() Decoder { return NewINIDecoder() },
}

var Formats = []*Format{
	YamlFormat,
	KYamlFormat,
	JSONFormat,
	PropertiesFormat,
	CSVFormat,
	TSVFormat,
	XMLFormat,
	Base64Format,
	UriFormat,
	ShFormat,
	TomlFormat,
	HclFormat,
	ShellVariablesFormat,
	LuaFormat,
	INIFormat,
}

func (f *Format) MatchesName(name string) bool {
	if f.FormalName == name {
		return true
	}
	return slices.Contains(f.Names, name)
}

func (f *Format) GetConfiguredEncoder() Encoder {
	return f.EncoderFactory()
}

func FormatStringFromFilename(filename string) string {
	if filename != "" {
		GetLogger().Debugf("checking filename '%s' for auto format detection", filename)
		ext := filepath.Ext(filename)
		if len(ext) >= 2 && ext[0] == '.' {
			format := strings.ToLower(ext[1:])
			GetLogger().Debugf("detected format '%s'", format)
			return format
		}
	}

	GetLogger().Debugf("using default inputFormat 'yaml'")
	return "yaml"
}

func FormatFromString(format string) (*Format, error) {
	if format != "" {
		for _, printerFormat := range Formats {
			if printerFormat.MatchesName(format) {
				return printerFormat, nil
			}
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
