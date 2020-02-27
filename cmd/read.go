package cmd

import (
	"bytes"
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/printer"
	errors "github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const escape = "\x1b"

func format(attr color.Attribute) string {
	return fmt.Sprintf("%s[%dm", escape, attr)
}

func colorizeAndPrint(bytes []byte, writer io.Writer) error {
	tokens := lexer.Tokenize(string(bytes))
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
			Prefix: format(color.FgHiCyan),
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
			Prefix: format(color.FgHiGreen),
			Suffix: format(color.Reset),
		}
	}
	_, err := writer.Write([]byte(p.PrintTokens(tokens) + "\n"))
	return err
}

func createReadCmd() *cobra.Command {
	var cmdRead = &cobra.Command{
		Use:     "read [yaml_file] [path_expression]",
		Aliases: []string{"r"},
		Short:   "yq r [--printMode/-p pv] sample.yaml 'b.e(name==fr*).value'",
		Example: `
yq read things.yaml 'a.b.c'
yq r - 'a.b.c' # reads from stdin
yq r things.yaml 'a.*.c'
yq r things.yaml 'a.**.c' # deep splat
yq r things.yaml 'a.(child.subchild==co*).c'
yq r -d1 things.yaml 'a.array[0].blah'
yq r things.yaml 'a.array[*].blah'
yq r -- things.yaml '--key-starting-with-dashes.blah'
      `,
		Long: "Outputs the value of the given path in the yaml file to STDOUT",
		RunE: readProperty,
	}
	cmdRead.PersistentFlags().StringVarP(&docIndex, "doc", "d", "0", "process document index number (0 based, * for all documents)")
	cmdRead.PersistentFlags().StringVarP(&printMode, "printMode", "p", "v", "print mode (v (values, default), p (paths), pv (path and value pairs)")
	cmdRead.PersistentFlags().StringVarP(&defaultValue, "defaultValue", "D", "", "default value printed when there are no results")
	cmdRead.PersistentFlags().BoolVarP(&explodeAnchors, "explodeAnchors", "X", false, "explode anchors")
	cmdRead.PersistentFlags().BoolVarP(&colorsEnabled, "colorsEnabled", "C", false, "enable colors")
	return cmdRead
}

func readProperty(cmd *cobra.Command, args []string) error {
	var path = ""

	if len(args) < 1 {
		return errors.New("Must provide filename")
	} else if len(args) > 1 {
		path = args[1]
	}

	var updateAll, docIndexInt, errorParsingDocIndex = parseDocumentIndex()
	if errorParsingDocIndex != nil {
		return errorParsingDocIndex
	}

	matchingNodes, errorReadingStream := readYamlFile(args[0], path, updateAll, docIndexInt)

	if errorReadingStream != nil {
		return errorReadingStream
	}
	out := cmd.OutOrStdout()

	if colorsEnabled && !outputToJSON {
		buffer := bytes.NewBuffer(nil)
		if err := printResults(matchingNodes, buffer); err != nil {
			return err
		}
		return colorizeAndPrint(buffer.Bytes(), out)
	}
	return printResults(matchingNodes, out)
}
