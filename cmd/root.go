package cmd

import (
	"os"

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
			cmd.Println(cmd.UsageString())

			return nil
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
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
	rootCmd.PersistentFlags().BoolVarP(&outputToJSON, "tojson", "j", false, "output as json")
	rootCmd.PersistentFlags().BoolVarP(&prettyPrint, "prettyPrint", "P", false, "pretty print")
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "Print version information and quit")

	rootCmd.AddCommand(
		createReadCmd(),
		createCompareCmd(),
		createValidateCmd(),
		createWriteCmd(),
		createPrefixCmd(),
		createDeleteCmd(),
		createNewCmd(),
		createMergeCmd(),
	)
	rootCmd.SetOutput(os.Stdout)

	return rootCmd
}
