# System Operators

The `system` operator allows you to run an external command and use its output as a value in your expression.

**Security warning**: The system operator is disabled by default. You must explicitly pass `--security-enable-system-operator` to use it.

## Usage

```bash
yq --security-enable-system-operator --null-input '.field = system("command"; "arg1")'
```

The operator takes:
- A command string (required)
- An argument or array of arguments separated by `;` (optional)

The current matched node's value is serialised and piped to the command via stdin. The command's stdout (with trailing newline stripped) is returned as a string.

## Disabling the system operator

The system operator is disabled by default. When disabled, a warning is logged and `null` is returned instead of running the command.

Use `--security-enable-system-operator` flag to enable it.
