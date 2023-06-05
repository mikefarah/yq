package yqlib

import (
	"testing"
)

var dateTimeOperatorScenarios = []expressionScenario{
	{
		description:    "Format: from standard RFC3339 format",
		subdescription: "Providing a single parameter assumes a standard RFC3339 datetime format. If the target format is not a valid yaml datetime format, the result will be a string tagged node.",
		document:       `a: 2001-12-15T02:59:43.1Z`,
		expression:     `.a |= format_datetime("Monday, 02-Jan-06 at 3:04PM")`,
		expected: []string{
			"D0, P[], (!!map)::a: Saturday, 15-Dec-01 at 2:59AM\n",
		},
	},
	{
		description:    "Format: from custom date time",
		subdescription: "Use with_dtf to set a custom datetime format for parsing.",
		document:       `a: Saturday, 15-Dec-01 at 2:59AM`,
		expression:     `.a |= with_dtf("Monday, 02-Jan-06 at 3:04PM"; format_datetime("2006-01-02"))`,
		expected: []string{
			"D0, P[], (!!map)::a: 2001-12-15\n",
		},
	},
	{
		description: "Format: get the day of the week",
		document:    `a: 2001-12-15`,
		expression:  `.a | format_datetime("Monday")`,
		expected: []string{
			"D0, P[a], (!!str)::Saturday\n",
		},
	},
	{
		description: "Now",
		document:    "a: cool",
		expression:  `.updated = now`,
		expected: []string{
			"D0, P[], (!!map)::a: cool\nupdated: 2021-05-19T01:02:03Z\n",
		},
	},
	{
		description:    "From Unix",
		subdescription: "Converts from unix time. Note, you don't have to pipe through the tz operator :)",
		expression:     `1675301929 | from_unix | tz("UTC")`,
		expected: []string{
			"D0, P[], (!!timestamp)::2023-02-02T01:38:49Z\n",
		},
	},
	{
		description:    "To Unix",
		subdescription: "Converts to unix time",
		expression:     `now | to_unix`,
		expected: []string{
			"D0, P[], (!!int)::1621386123\n",
		},
	},
	{
		description:    "Timezone: from standard RFC3339 format",
		subdescription: "Returns a new datetime in the specified timezone. Specify standard IANA Time Zone format or 'utc', 'local'. When given a single parameter, this assumes the datetime is in RFC3339 format.",

		document:   "a: cool",
		expression: `.updated = (now | tz("Australia/Sydney"))`,
		expected: []string{
			"D0, P[], (!!map)::a: cool\nupdated: 2021-05-19T11:02:03+10:00\n",
		},
	},
	{
		description:    "Timezone: with custom format",
		subdescription: "Specify standard IANA Time Zone format or 'utc', 'local'",
		document:       "a: Saturday, 15-Dec-01 at 2:59AM GMT",
		expression:     `.a |= with_dtf("Monday, 02-Jan-06 at 3:04PM MST"; tz("Australia/Sydney"))`,
		expected: []string{
			"D0, P[], (!!map)::a: Saturday, 15-Dec-01 at 1:59PM AEDT\n",
		},
	},
	{
		description:    "Add and tz custom format",
		subdescription: "Specify standard IANA Time Zone format or 'utc', 'local'",
		document:       "a: Saturday, 15-Dec-01 at 2:59AM GMT",
		expression:     `.a |= with_dtf("Monday, 02-Jan-06 at 3:04PM MST"; tz("Australia/Sydney"))`,
		expected: []string{
			"D0, P[], (!!map)::a: Saturday, 15-Dec-01 at 1:59PM AEDT\n",
		},
	},
	{
		description: "Date addition",
		document:    `a: 2021-01-01T00:00:00Z`,
		expression:  `.a += "3h10m"`,
		expected: []string{
			"D0, P[], (!!map)::a: 2021-01-01T03:10:00Z\n",
		},
	},
	{
		description:    "Date subtraction",
		subdescription: "You can subtract durations from dates. Assumes RFC3339 date time format, see [date-time operators](https://mikefarah.gitbook.io/yq/operators/datetime#date-time-formattings) for more information.",
		document:       `a: 2021-01-01T03:10:00Z`,
		expression:     `.a -= "3h10m"`,
		expected: []string{
			"D0, P[], (!!map)::a: 2021-01-01T00:00:00Z\n",
		},
	},
	{
		description: "Date addition - custom format",
		document:    `a: Saturday, 15-Dec-01 at 2:59AM GMT`,
		expression:  `with_dtf("Monday, 02-Jan-06 at 3:04PM MST"; .a += "3h1m")`,
		expected: []string{
			"D0, P[], (!!map)::a: Saturday, 15-Dec-01 at 6:00AM GMT\n",
		},
	},
	{
		description:    "Date script with custom format",
		subdescription: "You can embed full expressions in with_dtf if needed.",
		document:       `a: Saturday, 15-Dec-01 at 2:59AM GMT`,
		expression:     `with_dtf("Monday, 02-Jan-06 at 3:04PM MST"; .a = (.a + "3h1m" | tz("Australia/Perth")))`,
		expected: []string{
			"D0, P[], (!!map)::a: Saturday, 15-Dec-01 at 2:00PM AWST\n",
		},
	},
	{
		description: "allow comma",
		skipDoc:     true,
		document:    "a: Saturday, 15-Dec-01 at 2:59AM GMT",
		expression:  `.a |= with_dtf("Monday, 02-Jan-06 at 3:04PM MST", tz("Australia/Sydney"))`,
		expected: []string{
			"D0, P[], (!!map)::a: Saturday, 15-Dec-01 at 1:59PM AEDT\n",
		},
	},
}

func TestDatetimeOperatorScenarios(t *testing.T) {
	for _, tt := range dateTimeOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "datetime", dateTimeOperatorScenarios)
}
