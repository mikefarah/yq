package main

import (
	"bytes"
	"fmt"
	"io"

	yaml "gopkg.in/yaml.v2"
)

func format(context interface{}) (string, error) {
	out, err := yaml.Marshal(context)
	if err != nil {
		return "", fmt.Errorf("error printing yaml: %v", err)
	}
	var dst bytes.Buffer
	// add error checking if indent expanded to ever return an error
	_ = indent(&dst, out)

	return dst.String(), nil
}

func indent(dst io.ByteWriter, src []byte) error {
	for _, c := range src {
		switch c {
		case '-':
			_ = dst.WriteByte(' ')
			_ = dst.WriteByte(' ')
			_ = dst.WriteByte(c)

		default:
			_ = dst.WriteByte(c)
		}
	}
	return nil
}
