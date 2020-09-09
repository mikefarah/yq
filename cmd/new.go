package cmd

import (
	"github.com/mikefarah/yq/v3/pkg/yqlib"
	"github.com/spf13/cobra"
)

func createNewCmd() *cobra.Command {
	var cmdNew = &cobra.Command{
		Use:     "new [path] [value]",
		Aliases: []string{"n"},
		Short:   "yq n [--script/-s script_file] a.b.c newValue",
		Example: `
yq new 'a.b.c' cat
yq n 'a.b.c' --tag '!!str' true # force 'true' to be interpreted as a string instead of bool
yq n 'a.b[+]' cat
yq n -- '--key-starting-with-dash' cat # need to use '--' to stop processing arguments as flags
yq n --script create_script.yaml
      `,
		Long: `Creates a new yaml w.r.t the given path and value.
Outputs to STDOUT

Create Scripts:
Note that you can give a create script to perform more sophisticated yaml. This follows the same format as the update script.
`,
		RunE: newProperty,
	}
	cmdNew.PersistentFlags().StringVarP(&writeScript, "script", "s", "", "yaml script for creating yaml")
	cmdNew.PersistentFlags().StringVarP(&customTag, "tag", "t", "", "set yaml tag (e.g. !!int)")
	cmdNew.PersistentFlags().StringVarP(&customStyle, "style", "", "", "formatting style of the value: single, double, folded, flow, literal, tagged")
	cmdNew.PersistentFlags().StringVarP(&anchorName, "anchorName", "", "", "anchor name")
	cmdNew.PersistentFlags().BoolVarP(&makeAlias, "makeAlias", "", false, "create an alias using the value as the anchor name")
	return cmdNew
}

func newProperty(cmd *cobra.Command, args []string) error {
	var badArgsMessage = "Must provide <path_to_update> <value>"
	var updateCommands, updateCommandsError = readUpdateCommands(args, 2, badArgsMessage, false)
	if updateCommandsError != nil {
		return updateCommandsError
	}
	newNode := lib.New(updateCommands[0].Path)

	for _, updateCommand := range updateCommands {

		errorUpdating := lib.Update(&newNode, updateCommand, true)

		if errorUpdating != nil {
			return errorUpdating
		}
	}

	var encoder = yqlib.NewYamlEncoder(cmd.OutOrStdout(), indent, colorsEnabled)
	return encoder.Encode(&newNode)
}
