package yqlib

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/printer"
)

// Thanks @risentveber!

const escape = "\x1b"

func format(attr color.Attribute) string {
	return fmt.Sprintf("%s[%dm", escape, attr)
}

func colorizeAndPrint(yamlBytes []byte, writer io.Writer) error {
	tokens := lexer.Tokenize(string(yamlBytes))
	var p printer.Printer
	p.Bool = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiMagenta),
			Suffix: format(color.Reset),
		}
	}
	p.Number = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiMagenta),
			Suffix: format(color.Reset),
		}
	}
	p.MapKey = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgCyan),
			Suffix: format(color.Reset),
		}
	}
	p.Anchor = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiYellow),
			Suffix: format(color.Reset),
		}
	}
	p.Alias = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiYellow),
			Suffix: format(color.Reset),
		}
	}
	p.String = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgGreen),
			Suffix: format(color.Reset),
		}
	}
	p.Comment = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiBlack),
			Suffix: format(color.Reset),
		}
	}
	_, err := writer.Write([]byte(p.PrintTokens(tokens) + "\n"))
	return err
}
