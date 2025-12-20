package cmd

import (
	"errors"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/spf13/cobra"
)

func createEvaluateAllCommand() *cobra.Command {
	var cmdEvalAll = &cobra.Command{
		Use:     "eval-all [expression] [yaml_file1]...",
		Aliases: []string{"ea"},
		Short:   "Loads _all_ yaml documents of _all_ yaml files and runs expression once",
		ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveDefault
		},
		Example: `
# Merge f2.yml into f1.yml (in place)
yq eval-all --inplace 'select(fileIndex == 0) * select(fileIndex == 1)' f1.yml f2.yml
## the same command and expression using shortened names:
yq ea -i 'select(fi == 0) * select(fi == 1)' f1.yml f2.yml


# Merge all given files
yq ea '. as $item ireduce ({}; . * $item )' file1.yml file2.yml ...

# Pipe from STDIN
## use '-' as a filename to pipe from STDIN
cat file2.yml | yq ea '.a.b' file1.yml - file3.yml
`,
		Long: `yq is a portable command-line data file processor (https://github.com/mikefarah/yq/) 
See https://mikefarah.gitbook.io/yq/ for detailed documentation and examples.

## Evaluate All ##
This command loads _all_ yaml documents of _all_ yaml files and runs expression once
Useful when you need to run an expression across several yaml documents or files (like merge).
Note that it consumes more memory than eval.
`,
		RunE: evaluateAll,
	}
	return cmdEvalAll
}
func evaluateAll(cmd *cobra.Command, args []string) (cmdError error) {
	// 0 args, read std in
	// 1 arg, null input, process expression
	// 1 arg, read file in sequence
	// 2+ args, [0] = expression, file the rest

	var err error

	expression, args, err := initCommand(cmd, args)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()

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

	decoder, err := configureDecoder(true)
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
	if nulSepOutput {
		printer.SetNulSepOutput(true)
	}

	if frontMatter != "" {
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

	allAtOnceEvaluator := yqlib.NewAllAtOnceEvaluator()

	switch len(args) {
	case 0:
		if nullInput {
			err = yqlib.NewStreamEvaluator().EvaluateNew(processExpression(expression), printer)
		} else {
			cmd.Println(cmd.UsageString())
			return nil
		}
	default:
		err = allAtOnceEvaluator.EvaluateFiles(processExpression(expression), args, printer, decoder)
	}

	completedSuccessfully = err == nil

	if err == nil && exitStatus && !printer.PrintedAnything() {
		return errors.New("no matches found")
	}

	return err
}
