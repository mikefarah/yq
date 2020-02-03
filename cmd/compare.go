package cmd

import (
	"bufio"
	"bytes"

	"github.com/kylelemons/godebug/diff"
	errors "github.com/pkg/errors"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

func createCompareCmd() *cobra.Command {
	var cmdCompare = &cobra.Command{
		Use:     "compare [yaml_file_a] [yaml_file_b]",
		Aliases: []string{"x"},
		Short:   "yq x data1.yml data2.yml",
		Example: `
yq x - data2.yml # reads from stdin
`,
		Long: "Compares two yaml files, prints the difference",
		RunE: compareDocuments,
	}
	cmdCompare.PersistentFlags().StringVarP(&docIndex, "doc", "d", "0", "process document index number (0 based)")
	cmdCompare.PersistentFlags().BoolVarP(&prettyPrint, "prettyPrint", "P", false, "pretty print (does not have an affect with json output)")
	return cmdCompare
}

func compareDocuments(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("Must provide at 2 yaml files")
	}
	if docIndex == "*" {
		return errors.New("Document splat for compare not yet supported")
	}

	var _, docIndexInt, errorParsingDocIndex = parseDocumentIndex()
	if errorParsingDocIndex != nil {
		return errorParsingDocIndex
	}

	var dataBucketA yaml.Node
	var dataBucketB yaml.Node
	var errorReadingStream error
	errorReadingStream = readData(args[0], docIndexInt, &dataBucketA)
	if errorReadingStream != nil {
		return errorReadingStream
	}

	errorReadingStream = readData(args[1], docIndexInt, &dataBucketB)
	if errorReadingStream != nil {
		return errorReadingStream
	}

	if prettyPrint {
		updateStyleOfNode(&dataBucketA, 0)
		updateStyleOfNode(&dataBucketB, 0)
	}

	if errorReadingStream != nil {
		return errorReadingStream
	}

	var dataBufferA bytes.Buffer
	printNode(&dataBucketA, bufio.NewWriter(&dataBufferA))

	var dataBufferB bytes.Buffer
	printNode(&dataBucketB, bufio.NewWriter(&dataBufferB))

	cmd.Print(diff.Diff(dataBufferA.String(), dataBufferB.String()))
	return nil
}
