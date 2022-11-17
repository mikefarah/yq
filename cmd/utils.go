package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/spf13/cobra"
	"gopkg.in/op/go-logging.v1"
)

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

	return expression, args, nil
}

func configureDecoder(evaluateTogether bool) (yqlib.Decoder, error) {
	yqlibInputFormat, err := yqlib.InputFormatFromString(inputFormat)
	if err != nil {
		return nil, err
	}
	switch yqlibInputFormat {
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
	}
	prefs := yqlib.ConfiguredYamlPreferences
	prefs.EvaluateTogether = evaluateTogether
	return yqlib.NewYamlDecoder(prefs), nil
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

func configureEncoder(format yqlib.PrinterOutputFormat) yqlib.Encoder {
	switch format {
	case yqlib.JSONOutputFormat:
		return yqlib.NewJSONEncoder(indent, colorsEnabled, unwrapScalar)
	case yqlib.PropsOutputFormat:
		return yqlib.NewPropertiesEncoder(unwrapScalar)
	case yqlib.CSVOutputFormat:
		return yqlib.NewCsvEncoder(',')
	case yqlib.TSVOutputFormat:
		return yqlib.NewCsvEncoder('\t')
	case yqlib.YamlOutputFormat:
		return yqlib.NewYamlEncoder(indent, colorsEnabled, yqlib.ConfiguredYamlPreferences)
	case yqlib.XMLOutputFormat:
		return yqlib.NewXMLEncoder(indent, yqlib.ConfiguredXMLPreferences)
	}
	panic("invalid encoder")
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
