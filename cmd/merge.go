package cmd

import (
	"github.com/mikefarah/yq/v3/pkg/yqlib"
	errors "github.com/pkg/errors"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

func createMergeCmd() *cobra.Command {
	var cmdMerge = &cobra.Command{
		Use:     "merge [initial_yaml_file] [additional_yaml_file]...",
		Aliases: []string{"m"},
		Short:   "yq m [--inplace/-i] [--doc/-d index] [--overwrite/-x] [--arrayMerge/-a strategy] sample.yaml sample2.yaml",
		Example: `
yq merge things.yaml other.yaml
yq merge --inplace things.yaml other.yaml
yq m -i things.yaml other.yaml
yq m --overwrite things.yaml other.yaml
yq m -i -x things.yaml other.yaml
yq m -i -a=append things.yaml other.yaml
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
	cmdMerge.PersistentFlags().StringVarP(&arrayMergeStrategyFlag, "arrays", "a", "update", `array merge strategy (update/append/overwrite)
update: recursively update arrays by their index
append: concatenate arrays together
overwrite: replace arrays
`)
	cmdMerge.PersistentFlags().StringVarP(&commentsMergeStrategyFlag, "comments", "", "setWhenBlank", `comments merge strategy (setWhenBlank/ignore/append/overwrite)
setWhenBlank: set comment if the original document has no comment at that node
ignore: leave comments as-is in the original
append: append comments together
overwrite: overwrite comments completely
`)
	cmdMerge.PersistentFlags().StringVarP(&docIndex, "doc", "d", "0", "process document index number (0 based, * for all documents)")
	return cmdMerge
}

/*
* We don't deeply traverse arrays when appending a merge, instead we want to
* append the entire array element.
 */
func createReadFunctionForMerge(arrayMergeStrategy yqlib.ArrayMergeStrategy) func(*yaml.Node) ([]*yqlib.NodeContext, error) {
	return func(dataBucket *yaml.Node) ([]*yqlib.NodeContext, error) {
		return lib.GetForMerge(dataBucket, "**", arrayMergeStrategy)
	}
}

func mergeProperties(cmd *cobra.Command, args []string) error {
	var updateCommands []yqlib.UpdateCommand = make([]yqlib.UpdateCommand, 0)

	if len(args) < 1 {
		return errors.New("Must provide at least 1 yaml file")
	}
	var arrayMergeStrategy yqlib.ArrayMergeStrategy

	switch arrayMergeStrategyFlag {
	case "update":
		arrayMergeStrategy = yqlib.UpdateArrayMergeStrategy
	case "append":
		arrayMergeStrategy = yqlib.AppendArrayMergeStrategy
	case "overwrite":
		arrayMergeStrategy = yqlib.OverwriteArrayMergeStrategy
	default:
		return errors.New("Array merge strategy must be one of: update/append/overwrite")
	}

	var commentsMergeStrategy yqlib.CommentsMergeStrategy

	switch commentsMergeStrategyFlag {
	case "setWhenBlank":
		commentsMergeStrategy = yqlib.SetWhenBlankCommentsMergeStrategy
	case "ignore":
		commentsMergeStrategy = yqlib.IgnoreCommentsMergeStrategy
	case "append":
		commentsMergeStrategy = yqlib.AppendCommentsMergeStrategy
	case "overwrite":
		commentsMergeStrategy = yqlib.OverwriteCommentsMergeStrategy
	default:
		return errors.New("Comments merge strategy must be one of: setWhenBlank/ignore/append/overwrite")
	}

	if len(args) > 1 {
		// first generate update commands from the file
		var filesToMerge = args[1:]

		for _, fileToMerge := range filesToMerge {
			matchingNodes, errorProcessingFile := doReadYamlFile(fileToMerge, createReadFunctionForMerge(arrayMergeStrategy), false, 0)
			if errorProcessingFile != nil {
				return errorProcessingFile
			}
			log.Debugf("finished reading for merge!")
			for _, matchingNode := range matchingNodes {
				log.Debugf("matched node %v", lib.PathStackToString(matchingNode.PathStack))
				yqlib.DebugNode(matchingNode.Node)
			}
			for _, matchingNode := range matchingNodes {
				mergePath := lib.MergePathStackToString(matchingNode.PathStack, arrayMergeStrategy)
				updateCommands = append(updateCommands, yqlib.UpdateCommand{
					Command:               "merge",
					Path:                  mergePath,
					Value:                 matchingNode.Node,
					Overwrite:             overwriteFlag,
					CommentsMergeStrategy: commentsMergeStrategy,
					// dont update the content for nodes midway, only leaf nodes
					DontUpdateNodeContent: matchingNode.IsMiddleNode && (arrayMergeStrategy != yqlib.OverwriteArrayMergeStrategy || matchingNode.Node.Kind != yaml.SequenceNode),
				})
			}
		}
	}

	return updateDoc(args[0], updateCommands, cmd.OutOrStdout())
}
