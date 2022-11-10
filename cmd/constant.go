package cmd

import (
	"strconv"

	"github.com/spf13/pflag"
)

type boolFlag interface {
	pflag.Value
	IsExplicitySet() bool
	IsSet() bool
}

type unwrapScalarFlagStrc struct {
	explicitySet bool
	value        bool
}

func newFlag() boolFlag {
	return &unwrapScalarFlagStrc{value: true}
}

func (f *unwrapScalarFlagStrc) IsExplicitySet() bool {
	return f.explicitySet
}

func (f *unwrapScalarFlagStrc) IsSet() bool {
	return f.value
}

func (f *unwrapScalarFlagStrc) String() string {
	return strconv.FormatBool(f.value)
}

func (f *unwrapScalarFlagStrc) Set(value string) error {

	v, err := strconv.ParseBool(value)
	f.value = v
	f.explicitySet = true
	return err
}

func (*unwrapScalarFlagStrc) Type() string {
	return "bool"
}

var unwrapScalarFlag = newFlag()

var unwrapScalar = false

var writeInplace = false
var outputToJSON = false
var outputFormat = "yaml"
var inputFormat = "yaml"

var exitStatus = false
var forceColor = false
var forceNoColor = false
var colorsEnabled = false
var indent = 2
var noDocSeparators = false
var nullInput = false
var verbose = false
var version = false
var prettyPrint = false

// can be either "" (off), "extract" or "process"
var frontMatter = ""

var splitFileExp = ""
var splitFileExpFile = ""

var completedSuccessfully = false

var forceExpression = ""

var expressionFile = ""
