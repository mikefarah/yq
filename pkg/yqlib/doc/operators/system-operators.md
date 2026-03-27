# System Operators

The `system` operator allows you to run an external command and use its output as a value in your expression.

**Security warning**: The system operator is disabled by default. You must explicitly pass `--enable-system-operator` to use it.

## Usage

```bash
yq --enable-system-operator '.field = system("command"; "arg1")'
```

The operator takes:
- A command string (required)
- An argument or array of arguments separated by `;` (optional)

The current matched node's value is serialised and piped to the command via stdin. The command's stdout (with trailing newline stripped) is returned as a string.

## Disabling the system operator

The system operator is disabled by default. When disabled, a warning is logged and `null` is returned instead of running the command.

Use `--enable-system-operator` flag to enable it.

## system operator returns null when disabled
Use `--enable-system-operator` to enable the system operator.

Given a sample.yml file of:
```yaml
country: Australia
```
then
```bash
yq '.country = system("/usr/bin/echo"; "test")' sample.yml
```
will output
```yaml
country: null
```

## Run a command with an argument
Use `--enable-system-operator` to enable the system operator.

Given a sample.yml file of:
```yaml
country: Australia
```
then
```bash
yq '.country = system("/usr/bin/echo"; "test")' sample.yml
```
will output
```yaml
country: test
```

## Run a command without arguments
Omit the semicolon and args to run the command with no extra arguments.

Given a sample.yml file of:
```yaml
a: hello
```
then
```bash
yq '.a = system("/bin/echo")' sample.yml
```
will output
```yaml
a: ""
```

