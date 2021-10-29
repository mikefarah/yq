package cmd

import (
	"io"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
)

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
