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
		Long: `Example usage (for bash): to load completion run
	
	. <(yq bash-completion)
	
	To configure your bash shell to load completions for each session add to your bashrc
	
	# ~/.bashrc or ~/.profile
	. <(yq bash-completion)
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
