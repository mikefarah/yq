package cmd

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const unsafeFishCompletionRequest = `    # Disable ActiveHelp which is not supported for fish shell
    set -l requestComp "YQ_ACTIVE_HELP=0 $args[1] __complete $args[2..-1] $lastArg"

    __yq_debug "Calling $requestComp"
    set -l results (eval $requestComp 2> /dev/null)`

const safeFishCompletionRequest = `    # Disable ActiveHelp which is not supported for fish shell
    set -lx YQ_ACTIVE_HELP 0
    set -l requestComp $args[1] __complete $args[2..-1] $lastArg

    __yq_debug "Calling $requestComp"
    set -l results ($requestComp 2> /dev/null)`

var completionCmd = &cobra.Command{
	Use:     "completion [bash|zsh|fish|powershell]",
	Aliases: []string{"shell-completion"},
	Short:   "Generate the autocompletion script for the specified shell",
	Long: `To load completions:

Bash:

$ source <(yq completion bash)

# To load completions for each session, execute once:
Linux:
  $ yq completion bash > /etc/bash_completion.d/yq
MacOS:
  $ yq completion bash > /usr/local/etc/bash_completion.d/yq

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ yq completion zsh > "${fpath[1]}/_yq"

# You will need to start a new shell for this setup to take effect.

Fish:

$ yq completion fish | source

# To load completions for each session, execute once:
$ yq completion fish > ~/.config/fish/completions/yq.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error = nil
		switch args[0] {
		case "bash":
			err = cmd.Root().GenBashCompletionV2(os.Stdout, true)
		case "zsh":
			err = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			err = writeFishCompletion(cmd.Root(), os.Stdout)
		case "powershell":
			err = cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
		return err

	},
}

func writeFishCompletion(root *cobra.Command, writer io.Writer) error {
	var script bytes.Buffer
	if err := root.GenFishCompletion(&script, true); err != nil {
		return err
	}

	patchedScript, err := patchFishCompletionRequest(script.String())
	if err != nil {
		return err
	}

	_, err = io.WriteString(writer, patchedScript)
	return err
}

func patchFishCompletionRequest(script string) (string, error) {
	patchedScript := strings.Replace(script, unsafeFishCompletionRequest, safeFishCompletionRequest, 1)
	if patchedScript == script {
		return "", errors.New("failed to patch fish completion request")
	}
	return patchedScript, nil
}
