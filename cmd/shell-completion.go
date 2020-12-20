package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "shell-completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

$ source <(yq shell-completion bash)

# To load completions for each session, execute once:
Linux:
  $ yq shell-completion bash > /etc/bash_completion.d/yq
MacOS:
  $ yq shell-completion bash > /usr/local/etc/bash_completion.d/yq

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ yq shell-completion zsh > "${fpath[1]}/_yq"

# You will need to start a new shell for this setup to take effect.

Fish:

$ yq shell-completion fish | source

# To load completions for each session, execute once:
$ yq shell-completion fish > ~/.config/fish/completions/yq.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error = nil
		switch args[0] {
		case "bash":
			err = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			err = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			err = cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			err = cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
		return err

	},
}
