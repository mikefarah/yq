package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/spf13/cobra"
	"gopkg.in/op/go-logging.v1"
)

func initCommand(cmd *cobra.Command, args []string) (firstFileIndex int, err error) {
	cmd.SilenceUsage = true

	fileInfo, _ := os.Stdout.Stat()

	if forceColor || (!forceNoColor && (fileInfo.Mode()&os.ModeCharDevice) != 0) {
		colorsEnabled = true
	}

	firstFileIndex = -1
	if !nullInput && len(args) == 1 {
		firstFileIndex = 0
	} else if len(args) > 1 {
		firstFileIndex = 1
	}

	// backwards compatibility
	if outputToJSON {
		outputFormat = "json"
	}

	if writeInplace && (firstFileIndex == -1) {
		return 0, fmt.Errorf("write inplace flag only applicable when giving an expression and at least one file")
	}

	if writeInplace && splitFileExp != "" {
		return 0, fmt.Errorf("write inplace cannot be used with split file")
	}

	if nullInput && len(args) > 1 {
		return 0, fmt.Errorf("cannot pass files in when using null-input flag")
	}

	return firstFileIndex, nil
}

func configureDecoder() (yqlib.Decoder, error) {
	yqlibInputFormat, err := yqlib.InputFormatFromString(inputFormat)
	if err != nil {
		return nil, err
	}
	switch yqlibInputFormat {
	case yqlib.XmlInputFormat:
		return yqlib.NewXmlDecoder(xmlAttributePrefix, xmlContentName), nil
	}
	return yqlib.NewYamlDecoder(), nil
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
	case yqlib.JsonOutputFormat:
		return yqlib.NewJsonEncoder(indent)
	case yqlib.PropsOutputFormat:
		return yqlib.NewPropertiesEncoder()
	case yqlib.CsvOutputFormat:
		return yqlib.NewCsvEncoder(',')
	case yqlib.TsvOutputFormat:
		return yqlib.NewCsvEncoder('\t')
	case yqlib.YamlOutputFormat:
		return yqlib.NewYamlEncoder(indent, colorsEnabled, !noDocSeparators, unwrapScalar)
	case yqlib.XmlOutputFormat:
		return yqlib.NewXmlEncoder(indent, xmlAttributePrefix, xmlContentName)
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

func processArgs(pipingStdin bool, args []string) []string {
	if !pipingStdin {
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
