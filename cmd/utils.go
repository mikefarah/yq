package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/spf13/cobra"
	"gopkg.in/op/go-logging.v1"
)

func isAutomaticOutputFormat() bool {
	return outputFormat == "" || outputFormat == "auto" || outputFormat == "a"
}

func initCommand(cmd *cobra.Command, args []string) (string, []string, error) {
	cmd.SilenceUsage = true

	setupColors()

	expression, args, err := processArgs(args)
	if err != nil {
		return "", nil, err
	}

	if err := loadSplitFileExpression(); err != nil {
		return "", nil, err
	}

	handleBackwardsCompatibility()

	if err := validateCommandFlags(args); err != nil {
		return "", nil, err
	}

	if err := configureFormats(args); err != nil {
		return "", nil, err
	}

	configureUnwrapScalar()

	return expression, args, nil
}

func setupColors() {
	fileInfo, _ := os.Stdout.Stat()

	if forceColor || (!forceNoColor && (fileInfo.Mode()&os.ModeCharDevice) != 0) {
		colorsEnabled = true
	}
}

func loadSplitFileExpression() error {
	if splitFileExpFile != "" {
		splitExpressionBytes, err := os.ReadFile(splitFileExpFile)
		if err != nil {
			return err
		}
		splitFileExp = string(splitExpressionBytes)
	}
	return nil
}

func handleBackwardsCompatibility() {
	// backwards compatibility
	if outputToJSON {
		outputFormat = "json"
	}
}

func validateCommandFlags(args []string) error {
	if writeInplace && (len(args) == 0 || args[0] == "-") {
		return fmt.Errorf("write in place flag only applicable when giving an expression and at least one file")
	}

	if frontMatter != "" && len(args) == 0 {
		return fmt.Errorf("front matter flag only applicable when giving an expression and at least one file")
	}

	if writeInplace && splitFileExp != "" {
		return fmt.Errorf("write in place cannot be used with split file")
	}

	if nullInput && len(args) > 0 {
		return fmt.Errorf("cannot pass files in when using null-input flag")
	}

	return nil
}

func configureFormats(args []string) error {
	inputFilename := ""
	if len(args) > 0 {
		inputFilename = args[0]
	}

	if err := configureInputFormat(inputFilename); err != nil {
		return err
	}

	if err := configureOutputFormat(); err != nil {
		return err
	}

	yqlib.GetLogger().Debug("Using input format %v", inputFormat)
	yqlib.GetLogger().Debug("Using output format %v", outputFormat)

	return nil
}

func configureInputFormat(inputFilename string) error {
	if inputFormat == "" || inputFormat == "auto" || inputFormat == "a" {
		inputFormat = yqlib.FormatStringFromFilename(inputFilename)

		_, err := yqlib.FormatFromString(inputFormat)
		if err != nil {
			// unknown file type, default to yaml
			yqlib.GetLogger().Debug("Unknown file format extension '%v', defaulting to yaml", inputFormat)
			inputFormat = "yaml"
			if isAutomaticOutputFormat() {
				outputFormat = "yaml"
			}
		} else if isAutomaticOutputFormat() {
			outputFormat = inputFormat
		}
	} else if isAutomaticOutputFormat() {
		// backwards compatibility -
		// before this was introduced, `yq -pcsv things.csv`
		// would produce *yaml* output.
		//
		outputFormat = yqlib.FormatStringFromFilename(inputFilename)
		if inputFilename != "-" {
			yqlib.GetLogger().Warning("yq default output is now 'auto' (based on the filename extension). Normally yq would output '%v', but for backwards compatibility 'yaml' has been set. Please use -oy to specify yaml, or drop the -p flag.", outputFormat)
		}
		outputFormat = "yaml"
	}
	return nil
}

func configureOutputFormat() error {
	outputFormatType, err := yqlib.FormatFromString(outputFormat)
	if err != nil {
		return err
	}

	if outputFormatType == yqlib.YamlFormat ||
		outputFormatType == yqlib.PropertiesFormat {
		unwrapScalar = true
	}

	return nil
}

func configureUnwrapScalar() {
	if unwrapScalarFlag.IsExplicitlySet() {
		unwrapScalar = unwrapScalarFlag.IsSet()
	}
}

func configureDecoder(evaluateTogether bool) (yqlib.Decoder, error) {
	format, err := yqlib.FormatFromString(inputFormat)
	if err != nil {
		return nil, err
	}
	yqlib.ConfiguredYamlPreferences.EvaluateTogether = evaluateTogether

	if format.DecoderFactory == nil {
		return nil, fmt.Errorf("no support for %s input format", inputFormat)
	}
	yqlibDecoder := format.DecoderFactory()
	if yqlibDecoder == nil {
		return nil, fmt.Errorf("no support for %s input format", inputFormat)
	}
	return yqlibDecoder, nil
}

func configurePrinterWriter(format *yqlib.Format, out io.Writer) (yqlib.PrinterWriter, error) {

	var printerWriter yqlib.PrinterWriter

	if splitFileExp != "" {
		colorsEnabled = forceColor
		splitExp, err := yqlib.ExpressionParser.ParseExpression(splitFileExp)
		if err != nil {
			return nil, fmt.Errorf("bad split document expression: %w", err)
		}
		printerWriter = yqlib.NewMultiPrinterWriter(splitExp, format)
	} else {
		printerWriter = yqlib.NewSinglePrinterWriter(out)
	}
	return printerWriter, nil
}

func configureEncoder() (yqlib.Encoder, error) {
	yqlibOutputFormat, err := yqlib.FormatFromString(outputFormat)
	if err != nil {
		return nil, err
	}
	yqlib.ConfiguredXMLPreferences.Indent = indent
	yqlib.ConfiguredYamlPreferences.Indent = indent
	yqlib.ConfiguredKYamlPreferences.Indent = indent
	yqlib.ConfiguredJSONPreferences.Indent = indent

	yqlib.ConfiguredYamlPreferences.UnwrapScalar = unwrapScalar
	yqlib.ConfiguredKYamlPreferences.UnwrapScalar = unwrapScalar
	yqlib.ConfiguredPropertiesPreferences.UnwrapScalar = unwrapScalar
	yqlib.ConfiguredJSONPreferences.UnwrapScalar = unwrapScalar
	yqlib.ConfiguredShellVariablesPreferences.UnwrapScalar = unwrapScalar

	yqlib.ConfiguredYamlPreferences.ColorsEnabled = colorsEnabled
	yqlib.ConfiguredKYamlPreferences.ColorsEnabled = colorsEnabled
	yqlib.ConfiguredJSONPreferences.ColorsEnabled = colorsEnabled
	yqlib.ConfiguredHclPreferences.ColorsEnabled = colorsEnabled
	yqlib.ConfiguredTomlPreferences.ColorsEnabled = colorsEnabled

	yqlib.ConfiguredYamlPreferences.PrintDocSeparators = !noDocSeparators
	yqlib.ConfiguredKYamlPreferences.PrintDocSeparators = !noDocSeparators

	encoder := yqlibOutputFormat.EncoderFactory()

	if encoder == nil {
		return nil, fmt.Errorf("no support for %s output format", outputFormat)
	}
	return encoder, err
}

// this is a hack to enable backwards compatibility with githubactions (which pipe /dev/null into everything)
// and being able to call yq with the filename as a single parameter
//
// without this - yq detects there is stdin (thanks githubactions),
// then tries to parse the filename as an expression
func maybeFile(str string) bool {
	yqlib.GetLogger().Debugf("checking '%v' is a file", str)
	stat, err := os.Stat(str) // #nosec
	result := err == nil && !stat.IsDir()
	if yqlib.GetLogger().IsEnabledFor(logging.DEBUG) {
		if err != nil {
			yqlib.GetLogger().Debugf("error: %v", err)
		} else {
			yqlib.GetLogger().Debugf("error: %v, dir: %v", err, stat.IsDir())
		}
		yqlib.GetLogger().Debugf("result: %v", result)
	}
	return result
}

func processStdInArgs(args []string) []string {
	stat, err := os.Stdin.Stat()
	if err != nil {
		yqlib.GetLogger().Debugf("error getting stdin: %v", err)
	}
	pipingStdin := stat != nil && (stat.Mode()&os.ModeCharDevice) == 0

	// if we've been given a file, don't automatically
	// read from stdin.
	// this happens if there is more than one argument
	// or only one argument and its a file
	if nullInput || !pipingStdin || len(args) > 1 || (len(args) > 0 && maybeFile(args[0])) {
		return args
	}

	for _, arg := range args {
		if arg == "-" {
			return args
		}
	}
	yqlib.GetLogger().Debugf("missing '-', adding it to the end")

	// we're piping from stdin, but there's no '-' arg
	// lets add one to the end
	return append(args, "-")
}

func processArgs(originalArgs []string) (string, []string, error) {
	expression := forceExpression
	args := processStdInArgs(originalArgs)
	maybeFirstArgIsAFile := len(args) > 0 && maybeFile(args[0])

	if expressionFile == "" && maybeFirstArgIsAFile && strings.HasSuffix(args[0], ".yq") {
		// lets check if an expression file was given
		yqlib.GetLogger().Debug("Assuming arg %v is an expression file", args[0])
		expressionFile = args[0]
		args = args[1:]
	}

	if expressionFile != "" {
		expressionBytes, err := os.ReadFile(expressionFile)
		if err != nil {
			return "", nil, err
		}
		//replace \r\n (windows) with good ol' unix file endings.
		expression = strings.ReplaceAll(string(expressionBytes), "\r\n", "\n")
	}

	yqlib.GetLogger().Debugf("processed args: %v", args)
	if expression == "" && len(args) > 0 && args[0] != "-" && !maybeFile(args[0]) {
		yqlib.GetLogger().Debug("assuming expression is '%v'", args[0])
		expression = args[0]
		args = args[1:]
	}
	return expression, args, nil
}
