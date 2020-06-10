package cmd

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	"github.com/kylelemons/godebug/diff"
	"github.com/mikefarah/yq/v3/pkg/yqlib"
	errors "github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// turn off for unit tests :(
var forceOsExit = true

func createCompareCmd() *cobra.Command {
	var cmdCompare = &cobra.Command{
		Use:     "compare [yaml_file_a] [yaml_file_b]",
		Aliases: []string{"x"},
		Short:   "yq x [--prettyPrint/-P] dataA.yaml dataB.yaml 'b.e(name==fr*).value'",
		Example: `
yq x - data2.yml # reads from stdin
yq x -pp dataA.yaml dataB.yaml '**' # compare paths
yq x -d1 dataA.yaml dataB.yaml 'a.b.c'
`,
		Long: "Deeply compares two yaml files, prints the difference. Use with prettyPrint flag to ignore formatting differences.",
		RunE: compareDocuments,
	}
	cmdCompare.PersistentFlags().StringVarP(&docIndex, "doc", "d", "0", "process document index number (0 based, * for all documents)")
	cmdCompare.PersistentFlags().StringVarP(&printMode, "printMode", "p", "v", "print mode (v (values, default), p (paths), pv (path and value pairs)")
	cmdCompare.PersistentFlags().StringVarP(&defaultValue, "defaultValue", "D", "", "default value printed when there are no results")
	cmdCompare.PersistentFlags().BoolVarP(&stripComments, "stripComments", "", false, "strip comments out before comparing")
	cmdCompare.PersistentFlags().BoolVarP(&explodeAnchors, "explodeAnchors", "X", false, "explode anchors")
	return cmdCompare
}

func compareDocuments(cmd *cobra.Command, args []string) error {
	var path = ""

	if len(args) < 2 {
		return errors.New("Must provide at 2 yaml files")
	} else if len(args) > 2 {
		path = args[2]
	}

	var updateAll, docIndexInt, errorParsingDocIndex = parseDocumentIndex()
	if errorParsingDocIndex != nil {
		return errorParsingDocIndex
	}

	var matchingNodesA []*yqlib.NodeContext
	var matchingNodesB []*yqlib.NodeContext
	var errorDoingThings error

	matchingNodesA, errorDoingThings = readYamlFile(args[0], path, updateAll, docIndexInt)

	if errorDoingThings != nil {
		return errorDoingThings
	}

	matchingNodesB, errorDoingThings = readYamlFile(args[1], path, updateAll, docIndexInt)
	if errorDoingThings != nil {
		return errorDoingThings
	}

	var dataBufferA bytes.Buffer
	var dataBufferB bytes.Buffer
	errorDoingThings = printResults(matchingNodesA, bufio.NewWriter(&dataBufferA))
	if errorDoingThings != nil {
		return errorDoingThings
	}
	errorDoingThings = printResults(matchingNodesB, bufio.NewWriter(&dataBufferB))
	if errorDoingThings != nil {
		return errorDoingThings
	}

	diffString := diff.Diff(strings.TrimSuffix(dataBufferA.String(), "\n"), strings.TrimSuffix(dataBufferB.String(), "\n"))

	if len(diffString) > 1 {
		cmd.Print(diffString)
		cmd.Print("\n")
		if forceOsExit {
			os.Exit(1)
		}
	}
	return nil
}
