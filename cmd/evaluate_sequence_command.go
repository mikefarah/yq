package cmd

import (
	"errors"
	"fmt"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/spf13/cobra"
)

func createEvaluateSequenceCommand() *cobra.Command {
	var cmdEvalSequence = &cobra.Command{
		Use:     "eval [expression] [yaml_file1]...",
		Aliases: []string{"e"},
		Short:   "(default) Apply the expression to each document in each yaml file in sequence",
		ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveDefault
		},
		Example: `
# Reads field under the given path for each file
yq e '.a.b' f1.yml f2.yml 

# Prints out the file
yq e sample.yaml 

# Pipe from STDIN
## use '-' as a filename to pipe from STDIN
cat file2.yml | yq e '.a.b' file1.yml - file3.yml

# Creates a new yaml document
## Note that editing an empty file does not work.
yq e -n '.a.b.c = "cat"' 

# Update a file in place
yq e '.a.b = "cool"' -i file.yaml 
`,
		Long: `yq is a portable command-line data file processor (https://github.com/mikefarah/yq/) 
See https://mikefarah.gitbook.io/yq/ for detailed documentation and examples.

## Evaluate Sequence ##
This command iterates over each yaml document from each given file, applies the 
expression and prints the result in sequence.`,
		RunE: evaluateSequence,
	}
	return cmdEvalSequence
}

func processExpression(expression string) string {

	if prettyPrint && expression == "" {
		return yqlib.PrettyPrintExp
	} else if prettyPrint {
		return fmt.Sprintf("%v | %v", expression, yqlib.PrettyPrintExp)
	}
	return expression
}

func evaluateSequence(cmd *cobra.Command, args []string) (cmdError error) {
	// 0 args, read std in
	// 1 arg, null input, process expression
	// 1 arg, read file in sequence
	// 2+ args, [0] = expression, file the rest

	out := cmd.OutOrStdout()

	var err error

	expression, args, err := initCommand(cmd, args)
	if err != nil {
		return err
	}

	if writeInplace {
		// only use colours if its forced
		colorsEnabled = forceColor
		writeInPlaceHandler := yqlib.NewWriteInPlaceHandler(args[0])
		out, err = writeInPlaceHandler.CreateTempFile()
		if err != nil {
			return err
		}
		// need to indirectly call the function so  that completedSuccessfully is
		// passed when we finish execution as opposed to now
		defer func() {
			if cmdError == nil {
				cmdError = writeInPlaceHandler.FinishWriteInPlace(completedSuccessfully)
			}
		}()
	}

	format, err := yqlib.FormatFromString(outputFormat)
	if err != nil {
		return err
	}

	printerWriter, err := configurePrinterWriter(format, out)
	if err != nil {
		return err
	}
	encoder, err := configureEncoder()
	if err != nil {
		return err
	}

	printer := yqlib.NewPrinter(encoder, printerWriter)

	if printNodeInfo {
		printer = yqlib.NewNodeInfoPrinter(printerWriter)
	}

	if nulSepOutput {
		printer.SetNulSepOutput(true)
	}

	decoder, err := configureDecoder(false)
	if err != nil {
		return err
	}
	streamEvaluator := yqlib.NewStreamEvaluator()

	if frontMatter != "" {
		yqlib.GetLogger().Debug("using front matter handler")
		frontMatterHandler := yqlib.NewFrontMatterHandler(args[0])
		err = frontMatterHandler.Split()
		if err != nil {
			return err
		}
		args[0] = frontMatterHandler.GetYamlFrontMatterFilename()

		if frontMatter == "process" {
			reader := frontMatterHandler.GetContentReader()
			printer.SetAppendix(reader)
			defer yqlib.SafelyCloseReader(reader)
		}
		defer frontMatterHandler.CleanUp()
	}

	switch len(args) {
	case 0:
		if nullInput {
			err = streamEvaluator.EvaluateNew(processExpression(expression), printer)
		} else {
			cmd.Println(cmd.UsageString())
			return nil
		}
	default:
		err = streamEvaluator.EvaluateFiles(processExpression(expression), args, printer, decoder)
	}
	completedSuccessfully = err == nil

	if err == nil && exitStatus && !printer.PrintedAnything() {
		return errors.New("no matches found")
	}

	return err
}
