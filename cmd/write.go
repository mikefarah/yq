package cmd

import (
	"github.com/spf13/cobra"
)

func createWriteCmd() *cobra.Command {
	var cmdWrite = &cobra.Command{
		Use:     "write [yaml_file] [path_expression] [value]",
		Aliases: []string{"w"},
		Short:   "yq w [--inplace/-i] [--script/-s script_file] [--doc/-d index] sample.yaml 'b.e(name==fr*).value' newValue",
		Example: `
yq write things.yaml 'a.b.c' true
yq write things.yaml 'a.*.c' true
yq write things.yaml 'a.**' true
yq write things.yaml 'a.(child.subchild==co*).c' true
yq write things.yaml 'a.b.c' --tag '!!str' true # force 'true' to be interpreted as a string instead of bool
yq write things.yaml 'a.b.c' --tag '!!float' 3
yq write --inplace -- things.yaml 'a.b.c' '--cat' # need to use '--' to stop processing arguments as flags
yq w -i things.yaml 'a.b.c' cat
yq w -i -s update_script.yaml things.yaml
yq w things.yaml 'a.b.d[+]' foo # appends a new node to the 'd' array
yq w --doc 2 things.yaml 'a.b.d[+]' foo # updates the 3rd document of the yaml file
      `,
		Long: `Updates the yaml file w.r.t the given path and value.
Outputs to STDOUT unless the inplace flag is used, in which case the file is updated instead.

Append value to array adds the value to the end of array.

Update Scripts:
Note that you can give an update script to perform more sophisticated update. Update script
format is list of update commands (update or delete) like so:
---
- command: update
  path: b.c
  value:
    #great
    things: frog # wow!
- command: delete
  path: b.d
`,
		RunE: writeProperty,
	}
	cmdWrite.PersistentFlags().BoolVarP(&writeInplace, "inplace", "i", false, "update the yaml file inplace")
	cmdWrite.PersistentFlags().StringVarP(&writeScript, "script", "s", "", "yaml script for updating yaml")
	cmdWrite.PersistentFlags().StringVarP(&sourceYamlFile, "from", "f", "", "yaml file for updating yaml (as-is)")
	cmdWrite.PersistentFlags().StringVarP(&customTag, "tag", "t", "", "set yaml tag (e.g. !!int)")
	cmdWrite.PersistentFlags().StringVarP(&docIndex, "doc", "d", "0", "process document index number (0 based, * for all documents)")
	cmdWrite.PersistentFlags().StringVarP(&customStyle, "style", "", "", "formatting style of the value: single, double, folded, flow, literal, tagged")
	cmdWrite.PersistentFlags().StringVarP(&anchorName, "anchorName", "", "", "anchor name")
	cmdWrite.PersistentFlags().BoolVarP(&makeAlias, "makeAlias", "", false, "create an alias using the value as the anchor name")
	return cmdWrite
}

func writeProperty(cmd *cobra.Command, args []string) error {
	var updateCommands, updateCommandsError = readUpdateCommands(args, 3, "Must provide <filename> <path_to_update> <value>", true)
	if updateCommandsError != nil {
		return updateCommandsError
	}
	return updateDoc(args[0], updateCommands, cmd.OutOrStdout())
}
