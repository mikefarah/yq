package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/spf13/cobra"
)

func createEvaluateReduceCommand() *cobra.Command {
	var cmdEvalReduce = &cobra.Command{
		Use:     "eval-reduce [reduce expression] [yaml_file1]...",
		Aliases: []string{"er"},
		Short:   "Runs a reduce expression sequentially against each document of each file given. More memory efficient than using eval-all if you can get away with it.",
		Example: `
# Merge f2.yml into f1.yml (inplace)
yq eval-reduce --inplace '{} ; . * $doc' f1.yml f2.yml
`,
		Long: `yq is a portable command-line YAML processor (https://github.com/mikefarah/yq/) 
See https://mikefarah.gitbook.io/yq/ for detailed documentation and examples.

## Evaluate Reduce ##
This command runs the reduce expression against each document of each file given, accumulating the results.
It is most useful when merging multiple files (but isn't as flexible as using eval-all with ireduce, as you
can only merge the top level nodes).
`,
		RunE: evaluateReduce,
	}
	return cmdEvalReduce
}
func evaluateReduce(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	// 2+ args, [0] = expression, file the rest

	var err error

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

	if frontMatter != "" {
		frontMatterHandler := yqlib.NewFrontMatterHandler(args[firstFileIndex])
		err = frontMatterHandler.Split()
		if err != nil {
			return err
		}
		args[firstFileIndex] = frontMatterHandler.GetYamlFrontMatterFilename()

		if frontMatter == "process" {
			reader := frontMatterHandler.GetContentReader()
			printer.SetAppendix(reader)
			defer yqlib.SafelyCloseReader(reader)
		}
		defer frontMatterHandler.CleanUp()
	}

	reduceEvaluator := yqlib.NewReduceEvaluator()
	switch len(args) {
	case 0:
		cmd.Println(cmd.UsageString())
		return nil
	case 1:
		cmd.Println(cmd.UsageString())
		return nil
	default:
		err = reduceEvaluator.EvaluateFiles(processExpression(args[0]), args[1:], printer, leadingContentPreProcessing)
	}

	completedSuccessfully = err == nil

	if err == nil && exitStatus && !printer.PrintedAnything() {
		return errors.New("no matches found")
	}

	return err
}
