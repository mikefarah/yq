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

## Format: from standard RFC3339 format
Providing a single parameter assumes a standard RFC3339 datetime format. If the target format is not a valid yaml datetime format, the result will be a string tagged node.

Given a sample.yml file of:
```yaml
a: 2001-12-15T02:59:43.1Z
```
then
```bash
yq '.a |= format_datetime("Monday, 02-Jan-06 at 3:04PM")' sample.yml
```
will output
```yaml
a: Saturday, 15-Dec-01 at 2:59AM
```

## Format: from custom date time
Use with_dtf to set a custom datetime format for parsing.

Given a sample.yml file of:
```yaml
a: Saturday, 15-Dec-01 at 2:59AM
```
then
```bash
yq '.a |= with_dtf("Monday, 02-Jan-06 at 3:04PM"; format_datetime("2006-01-02"))' sample.yml
```
will output
```yaml
a: 2001-12-15
```

## Format: get the day of the week
Given a sample.yml file of:
```yaml
a: 2001-12-15
```
then
```bash
yq '.a | format_datetime("Monday")' sample.yml
```
will output
```yaml
Saturday
```

## Now
Given a sample.yml file of:
```yaml
a: cool
```
then
```bash
yq '.updated = now' sample.yml
```
will output
```yaml
a: cool
updated: 2021-05-19T01:02:03Z
```

## From Unix
Converts from unix time. Note, you don't have to pipe through the tz operator :)

Running
```bash
yq --null-input '1675301929 | from_unix | tz("UTC")'
```
will output
```yaml
2023-02-02T01:38:49Z
```

## To Unix
Converts to unix time

Running
```bash
yq --null-input 'now | to_unix'
```
will output
```yaml
1621386123
```

## Timezone: from standard RFC3339 format
Returns a new datetime in the specified timezone. Specify standard IANA Time Zone format or 'utc', 'local'. When given a single parameter, this assumes the datetime is in RFC3339 format.

Given a sample.yml file of:
```yaml
a: cool
```
then
```bash
yq '.updated = (now | tz("Australia/Sydney"))' sample.yml
```
will output
```yaml
a: cool
updated: 2021-05-19T11:02:03+10:00
```

## Timezone: with custom format
Specify standard IANA Time Zone format or 'utc', 'local'

Given a sample.yml file of:
```yaml
a: Saturday, 15-Dec-01 at 2:59AM GMT
```
then
```bash
yq '.a |= with_dtf("Monday, 02-Jan-06 at 3:04PM MST"; tz("Australia/Sydney"))' sample.yml
```
will output
```yaml
a: Saturday, 15-Dec-01 at 1:59PM AEDT
```

## Add and tz custom format
Specify standard IANA Time Zone format or 'utc', 'local'

Given a sample.yml file of:
```yaml
a: Saturday, 15-Dec-01 at 2:59AM GMT
```
then
```bash
yq '.a |= with_dtf("Monday, 02-Jan-06 at 3:04PM MST"; tz("Australia/Sydney"))' sample.yml
```
will output
```yaml
a: Saturday, 15-Dec-01 at 1:59PM AEDT
```

## Date addition
Given a sample.yml file of:
```yaml
a: 2021-01-01T00:00:00Z
```
then
```bash
yq '.a += "3h10m"' sample.yml
```
will output
```yaml
a: 2021-01-01T03:10:00Z
```

## Date subtraction
You can subtract durations from dates. Assumes RFC3339 date time format, see [date-time operators](https://mikefarah.gitbook.io/yq/operators/datetime#date-time-formattings) for more information.

Given a sample.yml file of:
```yaml
a: 2021-01-01T03:10:00Z
```
then
```bash
yq '.a -= "3h10m"' sample.yml
```
will output
```yaml
a: 2021-01-01T00:00:00Z
```

## Date addition - custom format
Given a sample.yml file of:
```yaml
a: Saturday, 15-Dec-01 at 2:59AM GMT
```
then
```bash
yq 'with_dtf("Monday, 02-Jan-06 at 3:04PM MST"; .a += "3h1m")' sample.yml
```
will output
```yaml
a: Saturday, 15-Dec-01 at 6:00AM GMT
```

## Date script with custom format
You can embed full expressions in with_dtf if needed.

Given a sample.yml file of:
```yaml
a: Saturday, 15-Dec-01 at 2:59AM GMT
```
then
```bash
yq 'with_dtf("Monday, 02-Jan-06 at 3:04PM MST"; .a = (.a + "3h1m" | tz("Australia/Perth")))' sample.yml
```
will output
```yaml
a: Saturday, 15-Dec-01 at 2:00PM AWST
```

