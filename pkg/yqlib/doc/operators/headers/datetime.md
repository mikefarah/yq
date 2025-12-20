# Date Time

Various operators for parsing and manipulating dates. 

## Date time formatting
This uses Golang's built in time library for parsing and formatting date times.

When not specified, the RFC3339 standard is assumed `2006-01-02T15:04:05Z07:00` for parsing.

To specify a custom parsing format, use the `with_dtf` operator. The first parameter sets the datetime parsing format for the expression in the second parameter. The expression can be any valid `yq` expression tree.

```bash
yq 'with_dtf("myformat"; .a + "3h" | tz("Australia/Melbourne"))'
```

See the [library docs](https://pkg.go.dev/time#pkg-constants) for examples of formatting options.


## Timezones
This uses Golang's built in LoadLocation function to parse timezones strings. See the [library docs](https://pkg.go.dev/time#LoadLocation) for more details.


## Durations
Durations are parsed using Golang's built in [ParseDuration](https://pkg.go.dev/time#ParseDuration) function.

You can add durations to time using the `+` operator.
