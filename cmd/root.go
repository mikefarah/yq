package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/spf13/cobra"
	logging "gopkg.in/op/go-logging.v1"
)

type runeValue rune

func newRuneVar(p *rune) *runeValue {
	return (*runeValue)(p)
}

func (r *runeValue) String() string {
	return string(*r)
}

func (r *runeValue) Set(rawVal string) error {
	val := strings.ReplaceAll(rawVal, "\\n", "\n")
	val = strings.ReplaceAll(val, "\\t", "\t")
	val = strings.ReplaceAll(val, "\\r", "\r")
	val = strings.ReplaceAll(val, "\\f", "\f")
	val = strings.ReplaceAll(val, "\\v", "\v")
	if len(val) != 1 {
		return fmt.Errorf("[%v] is not a valid character. Must be length 1 was %v", val, len(val))
	}
	*r = runeValue(rune(val[0]))
	return nil
}

func (r *runeValue) Type() string {
	return "char"
}

func New() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "yq",
		Short: "yq is a lightweight and portable command-line data file processor.",
		Long: `yq is a portable command-line data file processor (https://github.com/mikefarah/yq/) 
See https://mikefarah.gitbook.io/yq/ for detailed documentation and examples.`,
		Example: `
# yq tries to auto-detect the file format based off the extension, and defaults to YAML if it's unknown (or piping through STDIN)
# Use the '-p/--input-format' flag to specify a format type.
cat file.xml | yq -p xml

# read the "stuff" node from "myfile.yml"
yq '.stuff' < myfile.yml

# update myfile.yml in place
yq -i '.stuff = "foo"' myfile.yml

# print contents of sample.json as idiomatic YAML
yq -P -oy sample.json
`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if version {
				cmd.Print(GetVersionDisplay())
				return nil
			}
			return evaluateSequence(cmd, args)

		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.SetOut(cmd.OutOrStdout())
			level := logging.WARNING
			stringFormat := `[%{level}] %{color}%{time:15:04:05}%{color:reset} %{message}`

			if verbose {
				level = logging.DEBUG
				stringFormat = `[%{level:5.5s}] %{color}%{time:15:04:05}%{color:bold} %{shortfile:-33s} %{shortfunc:-25s}%{color:reset} %{message}`
			}

			var format = logging.MustStringFormatter(stringFormat)
			var backend = logging.AddModuleLevel(
				logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), format))

			backend.SetLevel(level, "")

			logging.SetBackend(backend)
			yqlib.InitExpressionParser()

			return nil
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")

	rootCmd.PersistentFlags().BoolVarP(&outputToJSON, "tojson", "j", false, "(deprecated) output as json. Set indent to 0 to print json in one line.")
	err := rootCmd.PersistentFlags().MarkDeprecated("tojson", "please use -o=json instead")
	if err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output-format", "o", "auto", "[auto|a|yaml|y|json|j|props|p|xml|x|tsv|t|csv|c] output format type.")
	rootCmd.PersistentFlags().StringVarP(&inputFormat, "input-format", "p", "auto", "[auto|a|yaml|y|props|p|xml|x|tsv|t|csv|c|toml] parse format for input. Note that json is a subset of yaml.")

	rootCmd.PersistentFlags().StringVar(&yqlib.ConfiguredXMLPreferences.AttributePrefix, "xml-attribute-prefix", yqlib.ConfiguredXMLPreferences.AttributePrefix, "prefix for xml attributes")
	rootCmd.PersistentFlags().StringVar(&yqlib.ConfiguredXMLPreferences.ContentName, "xml-content-name", yqlib.ConfiguredXMLPreferences.ContentName, "name for xml content (if no attribute name is present).")
	rootCmd.PersistentFlags().BoolVar(&yqlib.ConfiguredXMLPreferences.StrictMode, "xml-strict-mode", yqlib.ConfiguredXMLPreferences.StrictMode, "enables strict parsing of XML. See https://pkg.go.dev/encoding/xml for more details.")
	rootCmd.PersistentFlags().BoolVar(&yqlib.ConfiguredXMLPreferences.KeepNamespace, "xml-keep-namespace", yqlib.ConfiguredXMLPreferences.KeepNamespace, "enables keeping namespace after parsing attributes")
	rootCmd.PersistentFlags().BoolVar(&yqlib.ConfiguredXMLPreferences.UseRawToken, "xml-raw-token", yqlib.ConfiguredXMLPreferences.UseRawToken, "enables using RawToken method instead Token. Commonly disables namespace translations. See https://pkg.go.dev/encoding/xml#Decoder.RawToken for details.")
	rootCmd.PersistentFlags().StringVar(&yqlib.ConfiguredXMLPreferences.ProcInstPrefix, "xml-proc-inst-prefix", yqlib.ConfiguredXMLPreferences.ProcInstPrefix, "prefix for xml processing instructions (e.g. <?xml version=\"1\"?>)")
	rootCmd.PersistentFlags().StringVar(&yqlib.ConfiguredXMLPreferences.DirectiveName, "xml-directive-name", yqlib.ConfiguredXMLPreferences.DirectiveName, "name for xml directives (e.g. <!DOCTYPE thing cat>)")
	rootCmd.PersistentFlags().BoolVar(&yqlib.ConfiguredXMLPreferences.SkipProcInst, "xml-skip-proc-inst", yqlib.ConfiguredXMLPreferences.SkipProcInst, "skip over process instructions (e.g. <?xml version=\"1\"?>)")
	rootCmd.PersistentFlags().BoolVar(&yqlib.ConfiguredXMLPreferences.SkipDirectives, "xml-skip-directives", yqlib.ConfiguredXMLPreferences.SkipDirectives, "skip over directives (e.g. <!DOCTYPE thing cat>)")

	rootCmd.PersistentFlags().BoolVar(&yqlib.ConfiguredCsvPreferences.AutoParse, "csv-auto-parse", yqlib.ConfiguredCsvPreferences.AutoParse, "parse CSV YAML/JSON values")
	rootCmd.PersistentFlags().Var(newRuneVar(&yqlib.ConfiguredCsvPreferences.Separator), "csv-separator", "CSV Separator character")

	rootCmd.PersistentFlags().BoolVar(&yqlib.ConfiguredTsvPreferences.AutoParse, "tsv-auto-parse", yqlib.ConfiguredTsvPreferences.AutoParse, "parse TSV YAML/JSON values")

	rootCmd.PersistentFlags().StringVar(&yqlib.ConfiguredLuaPreferences.DocPrefix, "lua-prefix", yqlib.ConfiguredLuaPreferences.DocPrefix, "prefix")
	rootCmd.PersistentFlags().StringVar(&yqlib.ConfiguredLuaPreferences.DocSuffix, "lua-suffix", yqlib.ConfiguredLuaPreferences.DocSuffix, "suffix")
	rootCmd.PersistentFlags().BoolVar(&yqlib.ConfiguredLuaPreferences.UnquotedKeys, "lua-unquoted", yqlib.ConfiguredLuaPreferences.UnquotedKeys, "output unquoted string keys (e.g. {foo=\"bar\"})")
	rootCmd.PersistentFlags().BoolVar(&yqlib.ConfiguredLuaPreferences.Globals, "lua-globals", yqlib.ConfiguredLuaPreferences.Globals, "output keys as top-level global variables")

	rootCmd.PersistentFlags().BoolVarP(&nullInput, "null-input", "n", false, "Don't read input, simply evaluate the expression given. Useful for creating docs from scratch.")
	rootCmd.PersistentFlags().BoolVarP(&noDocSeparators, "no-doc", "N", false, "Don't print document separators (---)")

	rootCmd.PersistentFlags().IntVarP(&indent, "indent", "I", 2, "sets indent level for output")
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "Print version information and quit")
	rootCmd.PersistentFlags().BoolVarP(&writeInplace, "inplace", "i", false, "update the file in place of first file given.")
	rootCmd.PersistentFlags().VarP(unwrapScalarFlag, "unwrapScalar", "r", "unwrap scalar, print the value with no quotes, colors or comments. Defaults to true for yaml")
	rootCmd.PersistentFlags().Lookup("unwrapScalar").NoOptDefVal = "true"
	rootCmd.PersistentFlags().BoolVarP(&nulSepOutput, "nul-output", "0", false, "Use NUL char to separate values. If unwrap scalar is also set, fail if unwrapped scalar contains NUL char.")

	rootCmd.PersistentFlags().BoolVarP(&prettyPrint, "prettyPrint", "P", false, "pretty print, shorthand for '... style = \"\"'")
	rootCmd.PersistentFlags().BoolVarP(&exitStatus, "exit-status", "e", false, "set exit status if there are no matches or null or false is returned")

	rootCmd.PersistentFlags().BoolVarP(&forceColor, "colors", "C", false, "force print with colors")
	rootCmd.PersistentFlags().BoolVarP(&forceNoColor, "no-colors", "M", false, "force print with no colors")
	rootCmd.PersistentFlags().StringVarP(&frontMatter, "front-matter", "f", "", "(extract|process) first input as yaml front-matter. Extract will pull out the yaml content, process will run the expression against the yaml content, leaving the remaining data intact")
	rootCmd.PersistentFlags().StringVarP(&forceExpression, "expression", "", "", "forcibly set the expression argument. Useful when yq argument detection thinks your expression is a file.")
	rootCmd.PersistentFlags().BoolVarP(&yqlib.ConfiguredYamlPreferences.LeadingContentPreProcessing, "header-preprocess", "", true, "Slurp any header comments and separators before processing expression.")

	rootCmd.PersistentFlags().StringVarP(&splitFileExp, "split-exp", "s", "", "print each result (or doc) into a file named (exp). [exp] argument must return a string. You can use $index in the expression as the result counter.")
	rootCmd.PersistentFlags().StringVarP(&splitFileExpFile, "split-exp-file", "", "", "Use a file to specify the split-exp expression.")

	rootCmd.PersistentFlags().StringVarP(&expressionFile, "from-file", "", "", "Load expression from specified file.")

	rootCmd.AddCommand(
		createEvaluateSequenceCommand(),
		createEvaluateAllCommand(),
		completionCmd,
	)
	return rootCmd
}
