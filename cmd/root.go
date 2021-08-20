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
		Long: `yq is a portable command-line YAML processor (https://github.com/mikefarah/yq/) 
See https://mikefarah.gitbook.io/yq/ for detailed documentation and examples.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if version {
				cmd.Print(GetVersionDisplay())
				return nil
			}
			cmd.Println(cmd.UsageString())
			return nil

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

	rootCmd.PersistentFlags().BoolVarP(&outputToJSON, "tojson", "j", false, "(deprecated) output as json. Set indent to 0 to print json in one line.")
	err := rootCmd.PersistentFlags().MarkDeprecated("tojson", "please use -o=json instead")
	if err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output-format", "o", "yaml", "[yaml|json|props] output format type.")
	rootCmd.PersistentFlags().BoolVarP(&nullInput, "null-input", "n", false, "Don't read input, simply evaluate the expression given. Useful for creating yaml docs from scratch.")
	rootCmd.PersistentFlags().BoolVarP(&noDocSeparators, "no-doc", "N", false, "Don't print document separators (---)")

	rootCmd.PersistentFlags().IntVarP(&indent, "indent", "I", 2, "sets indent level for output")
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "Print version information and quit")
	rootCmd.PersistentFlags().BoolVarP(&writeInplace, "inplace", "i", false, "update the yaml file inplace of first yaml file given.")
	rootCmd.PersistentFlags().BoolVarP(&unwrapScalar, "unwrapScalar", "", true, "unwrap scalar, print the value with no quotes, colors or comments")
	rootCmd.PersistentFlags().BoolVarP(&prettyPrint, "prettyPrint", "P", false, "pretty print, shorthand for '... style = \"\"'")
	rootCmd.PersistentFlags().BoolVarP(&exitStatus, "exit-status", "e", false, "set exit status if there are no matches or null or false is returned")

	rootCmd.PersistentFlags().BoolVarP(&forceColor, "colors", "C", false, "force print with colors")
	rootCmd.PersistentFlags().BoolVarP(&forceNoColor, "no-colors", "M", false, "force print with no colors")
	rootCmd.PersistentFlags().StringVarP(&frontMatter, "front-matter", "f", "", "(extract|process) first input as yaml front-matter. Extract will pull out the yaml content, process will run the expression against the yaml content, leaving the remaining data intact")
	rootCmd.PersistentFlags().BoolVarP(&leadingContentPreProcessing, "header-preprocess", "", true, "Slurp any header comments and seperators before processing expression. This is a workaround for go-yaml to persist header content. This flag will be removed once this feature has been out in the wild for a while.")
	rootCmd.AddCommand(
		createEvaluateSequenceCommand(),
		createEvaluateAllCommand(),
		completionCmd,
	)
	return rootCmd
}
