package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/mikefarah/yq/v4/pkg/yqlib/treeops"
	"github.com/spf13/cobra"
	logging "gopkg.in/op/go-logging.v1"
)

func New() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "yq",
		Short: "yq is a lightweight and portable command-line YAML processor.",
		Long:  `yq is a lightweight and portable command-line YAML processor. It aims to be the jq or sed of yaml files.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if version {
				cmd.Print(GetVersionDisplay())
				return nil
			}
			if shellCompletion != "" {
				switch shellCompletion {
				case "bash", "":
					return cmd.GenBashCompletion(os.Stdout)
				case "zsh":
					return cmd.GenZshCompletion(os.Stdout)
				case "fish":
					return cmd.GenFishCompletion(os.Stdout, true)
				case "powershell":
					return cmd.GenPowerShellCompletion(os.Stdout)
				default:
					return fmt.Errorf("Unknown variant %v", shellCompletion)
				}
			}
			// if len(args) == 0 {
			// 	cmd.Println(cmd.UsageString())
			// 	return nil
			// }
			cmd.SilenceUsage = true

			var treeCreator = treeops.NewPathTreeCreator()

			expression := ""
			if len(args) > 0 {
				expression = args[0]
			}

			pathNode, err := treeCreator.ParsePath(expression)
			if err != nil {
				return err
			}

			if outputToJSON {
				explodeOp := treeops.Operation{OperationType: treeops.Explode}
				explodeNode := treeops.PathTreeNode{Operation: &explodeOp}
				pipeOp := treeops.Operation{OperationType: treeops.Pipe}
				pathNode = &treeops.PathTreeNode{Operation: &pipeOp, Lhs: pathNode, Rhs: &explodeNode}
			}

			matchingNodes, err := evaluate("-", pathNode)
			if err != nil {
				return err
			}

			if exitStatus && matchingNodes.Len() == 0 {
				cmd.SilenceUsage = true
				return errors.New("No matches found")
			}

			out := cmd.OutOrStdout()

			return printResults(matchingNodes, out)
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.SetOut(cmd.OutOrStdout())
			var format = logging.MustStringFormatter(
				`%{color}%{time:15:04:05} %{shortfunc} [%{level:.4s}]%{color:reset} %{message}`,
			)
			var backend = logging.AddModuleLevel(
				logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), format))

			if verbose {
				backend.SetLevel(logging.DEBUG, "")
			} else {
				backend.SetLevel(logging.ERROR, "")
			}

			logging.SetBackend(backend)
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")
	rootCmd.PersistentFlags().BoolVarP(&outputToJSON, "tojson", "j", false, "output as json. By default it prints a json document in one line, use the prettyPrint flag to print a formatted doc.")
	rootCmd.PersistentFlags().BoolVarP(&prettyPrint, "prettyPrint", "P", false, "pretty print")
	rootCmd.PersistentFlags().IntVarP(&indent, "indent", "I", 2, "sets indent level for output")
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "Print version information and quit")

	rootCmd.Flags().StringVarP(&shellCompletion, "shellCompletion", "", "", "[bash/zsh/powershell/fish] prints shell completion script")

	rootCmd.PersistentFlags().BoolVarP(&forceColor, "colors", "C", false, "force print with colors")
	rootCmd.PersistentFlags().BoolVarP(&forceNoColor, "no-colors", "M", false, "force print with no colors")

	// rootCmd.PersistentFlags().StringVarP(&docIndex, "doc", "d", "0", "process document index number (0 based, * for all documents)")
	rootCmd.PersistentFlags().StringVarP(&printMode, "printMode", "p", "v", "print mode (v (values, default), p (paths), pv (path and value pairs)")
	rootCmd.PersistentFlags().StringVarP(&defaultValue, "defaultValue", "D", "", "default value printed when there are no results")

	return rootCmd
}
