package yqlib

import (
	yaml "github.com/mikefarah/yaml/v2"
	errors "github.com/pkg/errors"
	"strings"
)

func YamlToString(context interface{}, trimOutput bool) (string, error) {
	switch context := context.(type) {
	case string:
		return context, nil
	default:
		return marshalContext(context, trimOutput)
	}
}

func marshalContext(context interface{}, trimOutput bool) (string, error) {
	out, err := yaml.Marshal(context)

	if err != nil {
		return "", errors.Wrap(err, "error printing yaml")
	}

	outStr := string(out)
	// trim the trailing new line as it's easier for a script to add
	// it in if required than to remove it
	if trimOutput {
		return strings.Trim(outStr, "\n "), nil
	}
	return outStr, nil
}