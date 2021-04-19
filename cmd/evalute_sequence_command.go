package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/spf13/cobra"
)

func createEvaluateSequenceCommand() *cobra.Command {
	var cmdEvalSequence = &cobra.Command{
		Use:     "eval [expression] [yaml_file1]...",
		Aliases: []string{"e"},
		Short:   "Apply expression to each document in each yaml file given in sequence",
		Example: `
# runs the expression against each file, in series
yq e '.a.b | length' f1.yml f2.yml 

# prints out the file
yq e sample.yaml 

# use '-' as a filename to read from STDIN
cat file2.yml | yq e '.a.b' file1.yml - file3.yml

# prints a new yaml document
yq e -n '.a.b.c = "cat"' 


# updates file.yaml directly
yq e '.a.b = "cool"' -i file.yaml 
`,
		Long: "Evaluate Sequence:\nIterate over each yaml document, apply the expression and print the results, in sequence.",
		RunE: evaluateSequence,
	}
	return cmdEvalSequence
}

func processExpression(expression string) string {
	if prettyPrint && expression == "" {
		return `... style=""`
	} else if prettyPrint {
		return fmt.Sprintf("%v | ... style= \"\"", expression)
	}
	return expression
}

func evaluateSequence(cmd *cobra.Command, args []string) error {
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

	printer := yqlib.NewPrinter(out, outputToJSON, unwrapScalar, colorsEnabled, indent, !noDocSeparators)

	streamEvaluator := yqlib.NewStreamEvaluator()

	if nullInput && len(args) > 1 {
		return errors.New("Cannot pass files in when using null-input flag")
	}

	switch len(args) {
	case 0:
		if pipingStdIn {
			err = streamEvaluator.EvaluateFiles(processExpression(""), []string{"-"}, printer)
		} else {
			cmd.Println(cmd.UsageString())
			return nil
		}
	case 1:
		if nullInput {
			err = streamEvaluator.EvaluateNew(processExpression(args[0]), printer)
		} else {
			err = streamEvaluator.EvaluateFiles(processExpression(""), []string{args[0]}, printer)
		}
	default:
		err = streamEvaluator.EvaluateFiles(processExpression(args[0]), args[1:], printer)
	}
	completedSuccessfully = err == nil

	if err == nil && exitStatus && !printer.PrintedAnything() {
		return errors.New("no matches found")
	}

	return err
}
