---
description: >-
  Generate a shell completion file for supported shells
  (bash/fish/zsh/powershell)
---

# Shell Completion

```bash
yq shell-completion zsh
```

Prints to StdOut a shell completion script for zsh shell.

### Bash (default)

```bash
source <(yq shell-completion bash)
```

#### To load completions for each session, execute once:

Linux:

```bash
yq shell-completion bash > /etc/bash_completion.d/yq 
```

MacOS:

```bash
yq shell-completion bash > /usr/local/etc/bash_completion.d/yq
```

### zsh

If shell completion is not already enabled in your environment you will need to enable it.  You can execute the following once:

```bash
echo "autoload -U compinit; compinit" >> ~/.zshrc
```

#### To load completions for each session, execute once:

```bash
yq shell-completion zsh > "${fpath[1]}/_yq"
```

You will need to start a new shell for this setup to take effect.

### fish

```bash
yq shell-completion fish | source
```

#### To load completions for each session, execute once:

```
yq shell-completion fish > ~/.config/fish/completions/yq.fish
```

### PowerShell

```bash
yq shell-completion powershell
```

Users need PowerShell version 5.0 or above, which comes with Windows 10 and can be downloaded separately for Windows 7 or 8.1. They can then write the completions to a file and source this file from their PowerShell profile, which is referenced by the $Profile environment variable.
