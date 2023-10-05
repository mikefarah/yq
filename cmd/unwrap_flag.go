package cmd

import (
	"strconv"

	"github.com/spf13/pflag"
)

type boolFlag interface {
	pflag.Value
	IsExplicitlySet() bool
	IsSet() bool
}

type unwrapScalarFlagStrc struct {
	explicitlySet bool
	value         bool
}

func newUnwrapFlag() boolFlag {
	return &unwrapScalarFlagStrc{value: true}
}

func (f *unwrapScalarFlagStrc) IsExplicitlySet() bool {
	return f.explicitlySet
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
	f.explicitlySet = true
	return err
}

func (*unwrapScalarFlagStrc) Type() string {
	return "bool"
}
