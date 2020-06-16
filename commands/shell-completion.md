---
description: >-
  Generate a shell completion file for supported shells
  (bash/fish/zsh/powershell)
---

# Shell Completion

```bash
yq shell-completion --variation=zsh
```

Prints to StdOut a shell completion script for zsh shell.

### Bash \(default\)

```bash
yq shell-completion
```

To configure your bash shell to load completions for each session add to your bashrc

```text
# ~/.bashrc or ~/.profile
. <(yq shell-completion)
```

### zsh

```bash
yq shell-completion --variation=zsh
```

The generated completion script should be put somewhere in your $fpath named \_yq

### fish

```bash
yq shell-completion --variation=fish
```

Save the output to a '.fish' file and add it to your completions directory.

### PowerShell

```bash
yq shell-completion --variation=powershell
```

Users need PowerShell version 5.0 or above, which comes with Windows 10 and can be downloaded separately for Windows 7 or 8.1. They can then write the completions to a file and source this file from their PowerShell profile, which is referenced by the $Profile environment variable.

