package cmd

import (
	errors "github.com/pkg/errors"
	"github.com/spf13/cobra"
)

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
	cmdRead.PersistentFlags().BoolVarP(&printLength, "length", "l", false, "print length of results")
	cmdRead.PersistentFlags().BoolVarP(&collectIntoArray, "collect", "c", false, "collect results into array")
	cmdRead.PersistentFlags().BoolVarP(&unwrapScalar, "unwrapScalar", "", true, "unwrap scalar, print the value with no quotes, colors or comments")
	cmdRead.PersistentFlags().BoolVarP(&stripComments, "stripComments", "", false, "print yaml without any comments")
	cmdRead.PersistentFlags().BoolVarP(&explodeAnchors, "explodeAnchors", "X", false, "explode anchors")
	cmdRead.PersistentFlags().BoolVarP(&exitStatus, "exitStatus", "e", false, "set exit status if no matches are found")
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

	if exitStatus && len(matchingNodes) == 0 {
		cmd.SilenceUsage = true
		return errors.New("No matches found")
	}

	if errorReadingStream != nil {
		cmd.SilenceUsage = true
		return errorReadingStream
	}
	out := cmd.OutOrStdout()

	return printResults(matchingNodes, out)
}
