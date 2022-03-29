package cmd

import (
	"os"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/spf13/cobra"
	logging "gopkg.in/op/go-logging.v1"
)

func New() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "yq",
		Short: "yq is a lightweight and portable command-line YAML processor.",
		Long: `yq is a portable command-line YAML processor (https://github.com/mikefarah/yq/) 
See https://mikefarah.gitbook.io/yq/ for detailed documentation and examples.`,
		Example: `
# yq defaults to 'eval' command if no command is specified. See "yq eval --help" for more examples.

# read the "stuff" node from "myfile.yml"
yq '.stuff' < myfile.yml

# update myfile.yml in place
yq -i '.stuff = "foo"' myfile.yml

# print contents of sample.json as idiomatic YAML
yq -P sample.json
`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if version {
				cmd.Print(GetVersionDisplay())
				return nil
			}
			return evaluateSequence(cmd, args)

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
			yqlib.InitExpressionParser()
			yqlib.XMLPreferences.AttributePrefix = xmlAttributePrefix
			yqlib.XMLPreferences.ContentName = xmlContentName
			yqlib.XMLPreferences.StrictMode = xmlStrictMode
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")

	rootCmd.PersistentFlags().BoolVarP(&outputToJSON, "tojson", "j", false, "(deprecated) output as json. Set indent to 0 to print json in one line.")
	err := rootCmd.PersistentFlags().MarkDeprecated("tojson", "please use -o=json instead")
	if err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output-format", "o", "yaml", "[yaml|y|json|j|props|p|xml|x] output format type.")
	rootCmd.PersistentFlags().StringVarP(&inputFormat, "input-format", "p", "yaml", "[yaml|y|props|p|xml|x] parse format for input. Note that json is a subset of yaml.")

	rootCmd.PersistentFlags().StringVar(&xmlAttributePrefix, "xml-attribute-prefix", "+", "prefix for xml attributes")
	rootCmd.PersistentFlags().StringVar(&xmlContentName, "xml-content-name", "+content", "name for xml content (if no attribute name is present).")
	rootCmd.PersistentFlags().BoolVar(&xmlStrictMode, "xml-strict-mode", false, "enables strict parsing of XML. See https://pkg.go.dev/encoding/xml for more details.")

	rootCmd.PersistentFlags().BoolVarP(&nullInput, "null-input", "n", false, "Don't read input, simply evaluate the expression given. Useful for creating docs from scratch.")
	rootCmd.PersistentFlags().BoolVarP(&noDocSeparators, "no-doc", "N", false, "Don't print document separators (---)")

	rootCmd.PersistentFlags().IntVarP(&indent, "indent", "I", 2, "sets indent level for output")
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "Print version information and quit")
	rootCmd.PersistentFlags().BoolVarP(&writeInplace, "inplace", "i", false, "update the file inplace of first file given.")
	rootCmd.PersistentFlags().BoolVarP(&unwrapScalar, "unwrapScalar", "", true, "unwrap scalar, print the value with no quotes, colors or comments")
	rootCmd.PersistentFlags().BoolVarP(&prettyPrint, "prettyPrint", "P", false, "pretty print, shorthand for '... style = \"\"'")
	rootCmd.PersistentFlags().BoolVarP(&exitStatus, "exit-status", "e", false, "set exit status if there are no matches or null or false is returned")

	rootCmd.PersistentFlags().BoolVarP(&forceColor, "colors", "C", false, "force print with colors")
	rootCmd.PersistentFlags().BoolVarP(&forceNoColor, "no-colors", "M", false, "force print with no colors")
	rootCmd.PersistentFlags().StringVarP(&frontMatter, "front-matter", "f", "", "(extract|process) first input as yaml front-matter. Extract will pull out the yaml content, process will run the expression against the yaml content, leaving the remaining data intact")
	rootCmd.PersistentFlags().StringVarP(&forceExpression, "expression", "", "", "forcibly set the expression argument. Useful when yq argument detection thinks your expression is a file.")
	rootCmd.PersistentFlags().BoolVarP(&leadingContentPreProcessing, "header-preprocess", "", true, "Slurp any header comments and separators before processing expression.")

	rootCmd.PersistentFlags().StringVarP(&splitFileExp, "split-exp", "s", "", "print each result (or doc) into a file named (exp). [exp] argument must return a string. You can use $index in the expression as the result counter.")

	rootCmd.PersistentFlags().StringVarP(&expressionFile, "from-file", "", "", "Load expression from specified file.")

	rootCmd.AddCommand(
		createEvaluateSequenceCommand(),
		createEvaluateAllCommand(),
		completionCmd,
	)
	return rootCmd
}
