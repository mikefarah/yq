# System Operators

The `system` operator allows you to run an external command and use its output as a value in your expression.

**Security warning**: The system operator is disabled by default. You must explicitly pass `--security-enable-system-operator` to use it.

**Note:** When enabled, the system operator can replicate the functionality of `env` and `load`
operators via external commands. Enabling it effectively overrides `--security-disable-env-ops`
and `--security-disable-file-ops`.

## Usage

```bash
yq --security-enable-system-operator --null-input '.field = system("command"; "arg1")'
```

The operator takes:
- A command string (required)
- An argument (or an array of arguments), separated from the command by `;` (optional)

The current matched node's value is serialised and piped to the command via stdin. The command's stdout (with trailing newline stripped) is returned as a string.

## Disabling the system operator

The system operator is disabled by default. When disabled, an error is returned instead of running the command, consistent with `--security-disable-env-ops` and `--security-disable-file-ops`.

Use `--security-enable-system-operator` flag to enable it.
