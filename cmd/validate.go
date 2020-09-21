package cmd

import (
	errors "github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func createValidateCmd() *cobra.Command {
	var cmdRead = &cobra.Command{
		Use:     "validate [yaml_file]",
		Aliases: []string{"v"},
		Short:   "yq v sample.yaml",
		Example: `
yq v - # reads from stdin
`,
		RunE:          validateProperty,
		SilenceUsage:  true,
		SilenceErrors: false,
	}
	cmdRead.PersistentFlags().StringVarP(&docIndex, "doc", "d", "0", "process document index number (0 based, * for all documents)")
	return cmdRead
}

func validateProperty(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("Must provide filename")
	}

	var updateAll, docIndexInt, errorParsingDocIndex = parseDocumentIndex()
	if errorParsingDocIndex != nil {
		return errorParsingDocIndex
	}

	_, errorReadingStream := readYamlFile(args[0], "", updateAll, docIndexInt)

	return errorReadingStream
}
