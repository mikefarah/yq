package cmd

import (
	"github.com/mikefarah/yq/v3/pkg/yqlib"
	errors "github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func createDeleteCmd() *cobra.Command {
	var cmdDelete = &cobra.Command{
		Use:     "delete [yaml_file] [path_expression]",
		Aliases: []string{"d"},
		Short:   "yq d [--inplace/-i] [--doc/-d index] sample.yaml 'b.e(name==fred)'",
		Example: `
yq delete things.yaml 'a.b.c'
yq delete things.yaml 'a.*.c'
yq delete things.yaml 'a.(child.subchild==co*).c'
yq delete things.yaml 'a.**'
yq delete --inplace things.yaml 'a.b.c'
yq delete --inplace -- things.yaml '--key-starting-with-dash' # need to use '--' to stop processing arguments as flags
yq d -i things.yaml 'a.b.c'
	`,
		Long: `Deletes the nodes matching the given path expression from the YAML file.
Outputs to STDOUT unless the inplace flag is used, in which case the file is updated instead.
`,
		RunE: deleteProperty,
	}
	cmdDelete.PersistentFlags().BoolVarP(&writeInplace, "inplace", "i", false, "update the yaml file inplace")
	cmdDelete.PersistentFlags().StringVarP(&docIndex, "doc", "d", "0", "process document index number (0 based, * for all documents)")
	return cmdDelete
}

func deleteProperty(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("Must provide <filename> <path_to_delete>")
	}
	var updateCommands []yqlib.UpdateCommand = make([]yqlib.UpdateCommand, 1)
	updateCommands[0] = yqlib.UpdateCommand{Command: "delete", Path: args[1]}

	return updateDoc(args[0], updateCommands, cmd.OutOrStdout())
}
