# Date Time

Various operators for parsing and manipulating dates. 

## Date time formattings
This uses the golangs built in time library for parsing and formatting date times.

When not specified, the RFC3339 standard is assumed `2006-01-02T15:04:05Z07:00`.

To use a custom format, use the `with_dtformat` operator to set the formatting context. Expressions in the second parameter then assume that date format.

```bash
yq 'with_dtformat("myformat"; .a + "3h" | tz("Australia/Melbourne"))'
```

See https://pkg.go.dev/time#pkg-constants for more examples.



## Timezones
This uses golangs built in LoadLocation function to parse timezones strings. See https://pkg.go.dev/time#LoadLocation for more details.


## Durations
Durations are parsed using golangs built in [ParseDuration](https://pkg.go.dev/time#ParseDuration) function.

You can durations to time using the `+` operator.
