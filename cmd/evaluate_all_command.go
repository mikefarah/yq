package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/spf13/cobra"
)

func createEvaluateAllCommand() *cobra.Command {
	var cmdEvalAll = &cobra.Command{
		Use:     "eval-all [expression] [yaml_file1]...",
		Aliases: []string{"ea"},
		Short:   "Loads _all_ yaml documents of _all_ yaml files and runs expression once",
		Example: `
# merges f2.yml into f1.yml (inplace)
yq eval-all --inplace 'select(fileIndex == 0) * select(fileIndex == 1)' f1.yml f2.yml		
`,
		Long: "Evaluate All:\nUseful when you need to run an expression across several yaml documents or files. Consumes more memory than eval",
		RunE: evaluateAll,
	}
	return cmdEvalAll
}
func evaluateAll(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	// 0 args, read std in
	// 1 arg, null input, process expression
	// 1 arg, read file in sequence
	// 2+ args, [0] = expression, file the rest

	var err error
	stat, _ := os.Stdin.Stat()
	pipingStdIn := (stat.Mode() & os.ModeCharDevice) == 0

	out := cmd.OutOrStdout()

	fileInfo, _ := os.Stdout.Stat()

	if forceColor || (!forceNoColor && (fileInfo.Mode()&os.ModeCharDevice) != 0) {
		colorsEnabled = true
	}

	firstFileIndex := -1
	if !nullInput && len(args) == 1 {
		firstFileIndex = 0
	} else if len(args) > 1 {
		firstFileIndex = 1
	}

	if writeInplace && (firstFileIndex == -1) {
		return fmt.Errorf("Write inplace flag only applicable when giving an expression and at least one file")
	}

	if writeInplace {
		// only use colors if its forced
		colorsEnabled = forceColor
		writeInPlaceHandler := yqlib.NewWriteInPlaceHandler(args[firstFileIndex])
		out, err = writeInPlaceHandler.CreateTempFile()
		if err != nil {
			return err
		}
		// need to indirectly call the function so  that completedSuccessfully is
		// passed when we finish execution as opposed to now
		defer func() { writeInPlaceHandler.FinishWriteInPlace(completedSuccessfully) }()
	}

	if nullInput && len(args) > 1 {
		return errors.New("Cannot pass files in when using null-input flag")
	}

	printer := yqlib.NewPrinter(out, outputToJSON, unwrapScalar, colorsEnabled, indent, !noDocSeparators)

	allAtOnceEvaluator := yqlib.NewAllAtOnceEvaluator()
	switch len(args) {
	case 0:
		if pipingStdIn {
			err = allAtOnceEvaluator.EvaluateFiles(processExpression(""), []string{"-"}, printer)
		} else {
			cmd.Println(cmd.UsageString())
			return nil
		}
	case 1:
		if nullInput {
			err = yqlib.NewStreamEvaluator().EvaluateNew(processExpression(args[0]), printer)
		} else {
			err = allAtOnceEvaluator.EvaluateFiles(processExpression(""), []string{args[0]}, printer)
		}
	default:
		err = allAtOnceEvaluator.EvaluateFiles(processExpression(args[0]), args[1:], printer)
	}

	completedSuccessfully = err == nil

	if err == nil && exitStatus && !printer.PrintedAnything() {
		return errors.New("no matches found")
	}

	return err
}
