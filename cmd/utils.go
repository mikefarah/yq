package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/spf13/cobra"
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

func configurePrinterWriter(format yqlib.PrinterOutputFormat, out io.Writer) yqlib.PrinterWriter {

	var printerWriter yqlib.PrinterWriter

	if splitFileExp != "" {
		colorsEnabled = forceColor
		splitExp, err := yqlib.NewExpressionParser().ParseExpression(splitFileExp)
		if err != nil {
			return nil
		}
		printerWriter = yqlib.NewMultiPrinterWriter(splitExp, format)
	} else {
		printerWriter = yqlib.NewSinglePrinterWriter(out)
	}
	return printerWriter
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
