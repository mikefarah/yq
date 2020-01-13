package cmd

import (
	"github.com/mikefarah/yq/v3/pkg/yqlib"
	errors "github.com/pkg/errors"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

func createPrefixCmd() *cobra.Command {
	var cmdPrefix = &cobra.Command{
		Use:     "prefix [yaml_file] [path]",
		Aliases: []string{"p"},
		Short:   "yq p [--inplace/-i] [--doc/-d index] sample.yaml a.b.c",
		Example: `
yq prefix things.yaml 'a.b.c'
yq prefix --inplace things.yaml 'a.b.c'
yq prefix --inplace -- things.yaml '--key-starting-with-dash' # need to use '--' to stop processing arguments as flags
yq p -i things.yaml 'a.b.c'
yq p --doc 2 things.yaml 'a.b.d'
yq p -d2 things.yaml 'a.b.d'
      `,
		Long: `Prefixes w.r.t to the yaml file at the given path.
Outputs to STDOUT unless the inplace flag is used, in which case the file is updated instead.
`,
		RunE: prefixProperty,
	}
	cmdPrefix.PersistentFlags().BoolVarP(&writeInplace, "inplace", "i", false, "update the yaml file inplace")
	cmdPrefix.PersistentFlags().StringVarP(&docIndex, "doc", "d", "0", "process document index number (0 based, * for all documents)")
	return cmdPrefix
}

func prefixProperty(cmd *cobra.Command, args []string) error {

	if len(args) < 2 {
		return errors.New("Must provide <filename> <prefixed_path>")
	}
	updateCommand := yqlib.UpdateCommand{Command: "update", Path: args[1]}
	log.Debugf("args %v", args)

	var updateAll, docIndexInt, errorParsingDocIndex = parseDocumentIndex()
	if errorParsingDocIndex != nil {
		return errorParsingDocIndex
	}

	var updateData = func(dataBucket *yaml.Node, currentIndex int) error {
		return prefixDocument(updateAll, docIndexInt, currentIndex, dataBucket, updateCommand)
	}
	return readAndUpdate(cmd.OutOrStdout(), args[0], updateData)
}
