package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var shellVariant = "bash"

func createBashCompletionCmd(rootCmd *cobra.Command) *cobra.Command {
	var completionCmd = &cobra.Command{
		Use:   "shell-completion",
		Short: "Generates shell completion scripts",
		Long: `To load completion for:
bash:
	Run	
	. <(yq shell-completion)
	
	To configure your bash shell to load completions for each session add to
	your bashrc
	
	# ~/.bashrc or ~/.profile
	. <(yq shell-completion)

zsh:
	The generated completion script should be put somewhere in your $fpath named _yq

powershell:
	Users need PowerShell version 5.0 or above, which comes with Windows 10 and 
	can be downloaded separately for Windows 7 or 8.1. They can then write the 
	completions to a file and source this file from their PowerShell profile, 
	which is referenced by the $Profile environment variable.

fish:
	Save the output to a fish file and add it to your completions directory.

	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch shellVariant {
			case "bash", "":
				return rootCmd.GenBashCompletion(os.Stdout)
			case "zsh":
				return rootCmd.GenZshCompletion(os.Stdout)
			case "fish":
				return rootCmd.GenFishCompletion(os.Stdout, true)
			case "powershell":
				return rootCmd.GenPowerShellCompletion(os.Stdout)
			default:
				return fmt.Errorf("Unknown variant %v", shellVariant)
			}
		},
	}
	completionCmd.PersistentFlags().StringVarP(&shellVariant, "variation", "V", "", "shell variation: bash (default), zsh, fish, powershell")
	return completionCmd
}
