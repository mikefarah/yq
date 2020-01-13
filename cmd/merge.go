package cmd

import (
	"github.com/mikefarah/yq/v3/pkg/yqlib"
	errors "github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

func createMergeCmd() *cobra.Command {
	var cmdMerge = &cobra.Command{
		Use:     "merge [initial_yaml_file] [additional_yaml_file]...",
		Aliases: []string{"m"},
		Short:   "yq m [--inplace/-i] [--doc/-d index] [--overwrite/-x] [--append/-a] sample.yaml sample2.yaml",
		Example: `
yq merge things.yaml other.yaml
yq merge --inplace things.yaml other.yaml
yq m -i things.yaml other.yaml
yq m --overwrite things.yaml other.yaml
yq m -i -x things.yaml other.yaml
yq m -i -a things.yaml other.yaml
yq m -i --autocreate=false things.yaml other.yaml
      `,
		Long: `Updates the yaml file by adding/updating the path(s) and value(s) from additional yaml file(s).
Outputs to STDOUT unless the inplace flag is used, in which case the file is updated instead.

If overwrite flag is set then existing values will be overwritten using the values from each additional yaml file.
If append flag is set then existing arrays will be merged with the arrays from each additional yaml file.
`,
		RunE: mergeProperties,
	}
	cmdMerge.PersistentFlags().BoolVarP(&writeInplace, "inplace", "i", false, "update the yaml file inplace")
	cmdMerge.PersistentFlags().BoolVarP(&overwriteFlag, "overwrite", "x", false, "update the yaml file by overwriting existing values")
	cmdMerge.PersistentFlags().BoolVarP(&autoCreateFlag, "autocreate", "c", true, "automatically create any missing entries")
	cmdMerge.PersistentFlags().BoolVarP(&appendFlag, "append", "a", false, "update the yaml file by appending array values")
	cmdMerge.PersistentFlags().BoolVarP(&allowEmptyFlag, "allow-empty", "e", false, "allow empty yaml files")
	cmdMerge.PersistentFlags().StringVarP(&docIndex, "doc", "d", "0", "process document index number (0 based, * for all documents)")
	return cmdMerge
}

func mergeProperties(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("Must provide at least 2 yaml files")
	}
	// first generate update commands from the file
	var filesToMerge = args[1:]
	var updateCommands []yqlib.UpdateCommand = make([]yqlib.UpdateCommand, 0)

	for _, fileToMerge := range filesToMerge {
		matchingNodes, errorProcessingFile := readYamlFile(fileToMerge, "**", false, 0)
		if errorProcessingFile != nil && (!allowEmptyFlag || !strings.HasPrefix(errorProcessingFile.Error(), "Could not process document index")) {
			return errorProcessingFile
		}
		for _, matchingNode := range matchingNodes {
			mergePath := lib.MergePathStackToString(matchingNode.PathStack, appendFlag)
			updateCommands = append(updateCommands, yqlib.UpdateCommand{Command: "update", Path: mergePath, Value: matchingNode.Node, Overwrite: overwriteFlag})
		}
	}

	return updateDoc(args[0], updateCommands, cmd.OutOrStdout())
}
