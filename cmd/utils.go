package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/spf13/cobra"
	"gopkg.in/op/go-logging.v1"
)

func isAutomaticOutputFormat() bool {
	return outputFormat == "" || outputFormat == "auto" || outputFormat == "a"
}

func initCommand(cmd *cobra.Command, args []string) (string, []string, error) {
	cmd.SilenceUsage = true

	fileInfo, _ := os.Stdout.Stat()

	if forceColor || (!forceNoColor && (fileInfo.Mode()&os.ModeCharDevice) != 0) {
		colorsEnabled = true
	}

	expression, args, err := processArgs(args)
	if err != nil {
		return "", nil, err
	}

	if splitFileExpFile != "" {
		splitExpressionBytes, err := os.ReadFile(splitFileExpFile)
		if err != nil {
			return "", nil, err
		}
		splitFileExp = string(splitExpressionBytes)
	}

	// backwards compatibility
	if outputToJSON {
		outputFormat = "json"
	}

	if writeInplace && (len(args) == 0 || args[0] == "-") {
		return "", nil, fmt.Errorf("write inplace flag only applicable when giving an expression and at least one file")
	}

	if frontMatter != "" && len(args) == 0 {
		return "", nil, fmt.Errorf("front matter flag only applicable when giving an expression and at least one file")
	}

	if writeInplace && splitFileExp != "" {
		return "", nil, fmt.Errorf("write inplace cannot be used with split file")
	}

	if nullInput && len(args) > 0 {
		return "", nil, fmt.Errorf("cannot pass files in when using null-input flag")
	}

	inputFilename := ""
	if len(args) > 0 {
		inputFilename = args[0]
	}
	if inputFormat == "" || inputFormat == "auto" || inputFormat == "a" {

		inputFormat = yqlib.FormatFromFilename(inputFilename)

		_, err := yqlib.InputFormatFromString(inputFormat)
		if err != nil {
			// unknown file type, default to yaml
			yqlib.GetLogger().Debug("Unknown file format extension '%v', defaulting to yaml", inputFormat)
			inputFormat = "yaml"
			if isAutomaticOutputFormat() {
				outputFormat = "yaml"
			}
		} else if isAutomaticOutputFormat() {
			// automatic input worked, we can do it for output too unless specified
			if inputFormat == "json" {
				yqlib.GetLogger().Warning("JSON file output is now JSON by default (instead of yaml). Use '-oy' or '--output-format=yaml' for yaml output")
			}
			outputFormat = inputFormat
		}
	} else if isAutomaticOutputFormat() {
		// backwards compatibility -
		// before this was introduced, `yq -pcsv things.csv`
		// would produce *yaml* output.
		//
		outputFormat = yqlib.FormatFromFilename(inputFilename)
		if inputFilename != "-" {
			yqlib.GetLogger().Warning("yq default output is now 'auto' (based on the filename extension). Normally yq would output '%v', but for backwards compatibility 'yaml' has been set. Please use -oy to specify yaml, or drop the -p flag.", outputFormat)
		}
		outputFormat = "yaml"
	}

	outputFormatType, err := yqlib.OutputFormatFromString(outputFormat)

	if err != nil {
		return "", nil, err
	}
	yqlib.GetLogger().Debug("Using input format %v", inputFormat)
	yqlib.GetLogger().Debug("Using output format %v", outputFormat)

	if outputFormatType == yqlib.YamlOutputFormat ||
		outputFormatType == yqlib.PropsOutputFormat {
		unwrapScalar = true
	}
	if unwrapScalarFlag.IsExplicitySet() {
		unwrapScalar = unwrapScalarFlag.IsSet()
	}

	//copy preference form global setting
	yqlib.ConfiguredYamlPreferences.UnwrapScalar = unwrapScalar

	yqlib.ConfiguredYamlPreferences.PrintDocSeparators = !noDocSeparators

	return expression, args, nil
}

func configureDecoder(evaluateTogether bool) (yqlib.Decoder, error) {
	yqlibInputFormat, err := yqlib.InputFormatFromString(inputFormat)
	if err != nil {
		return nil, err
	}
	yqlibDecoder, err := createDecoder(yqlibInputFormat, evaluateTogether)
	if yqlibDecoder == nil {
		return nil, fmt.Errorf("no support for %s input format", inputFormat)
	}
	return yqlibDecoder, err
}

func createDecoder(format yqlib.InputFormat, evaluateTogether bool) (yqlib.Decoder, error) {
	switch format {
	case yqlib.XMLInputFormat:
		return yqlib.NewXMLDecoder(yqlib.ConfiguredXMLPreferences), nil
	case yqlib.PropertiesInputFormat:
		return yqlib.NewPropertiesDecoder(), nil
	case yqlib.JsonInputFormat:
		return yqlib.NewJSONDecoder(), nil
	case yqlib.CSVObjectInputFormat:
		return yqlib.NewCSVObjectDecoder(','), nil
	case yqlib.TSVObjectInputFormat:
		return yqlib.NewCSVObjectDecoder('\t'), nil
	case yqlib.TomlInputFormat:
		return yqlib.NewTomlDecoder(), nil
	case yqlib.YamlInputFormat:
		prefs := yqlib.ConfiguredYamlPreferences
		prefs.EvaluateTogether = evaluateTogether
		return yqlib.NewYamlDecoder(prefs), nil
	}
	return nil, fmt.Errorf("invalid decoder: %v", format)
}

func configurePrinterWriter(format yqlib.PrinterOutputFormat, out io.Writer) (yqlib.PrinterWriter, error) {

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
	yqlibOutputFormat, err := yqlib.OutputFormatFromString(outputFormat)
	if err != nil {
		return nil, err
	}
	yqlibEncoder, err := createEncoder(yqlibOutputFormat)
	if yqlibEncoder == nil {
		return nil, fmt.Errorf("no support for %s output format", outputFormat)
	}
	return yqlibEncoder, err
}

func createEncoder(format yqlib.PrinterOutputFormat) (yqlib.Encoder, error) {
	switch format {
	case yqlib.JSONOutputFormat:
		return yqlib.NewJSONEncoder(indent, colorsEnabled, unwrapScalar), nil
	case yqlib.PropsOutputFormat:
		return yqlib.NewPropertiesEncoder(unwrapScalar), nil
	case yqlib.CSVOutputFormat:
		return yqlib.NewCsvEncoder(','), nil
	case yqlib.TSVOutputFormat:
		return yqlib.NewCsvEncoder('\t'), nil
	case yqlib.YamlOutputFormat:
		return yqlib.NewYamlEncoder(indent, colorsEnabled, yqlib.ConfiguredYamlPreferences), nil
	case yqlib.XMLOutputFormat:
		return yqlib.NewXMLEncoder(indent, yqlib.ConfiguredXMLPreferences), nil
	case yqlib.TomlOutputFormat:
		return yqlib.NewTomlEncoder(), nil
	case yqlib.ShellVariablesOutputFormat:
		return yqlib.NewShellVariablesEncoder(), nil
	}
	return nil, fmt.Errorf("invalid encoder: %v", format)
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
	stat, _ := os.Stdin.Stat()
	pipingStdin := (stat.Mode() & os.ModeCharDevice) == 0

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
	if expressionFile != "" {
		expressionBytes, err := os.ReadFile(expressionFile)
		if err != nil {
			return "", nil, err
		}
		expression = string(expressionBytes)
	}

	args := processStdInArgs(originalArgs)
	yqlib.GetLogger().Debugf("processed args: %v", args)
	if expression == "" && len(args) > 0 && args[0] != "-" && !maybeFile(args[0]) {
		yqlib.GetLogger().Debug("assuming expression is '%v'", args[0])
		expression = args[0]
		args = args[1:]
	}
	return expression, args, nil
}
