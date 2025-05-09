package yqlib

// func TestSingScenarioForDebugging(t *testing.T) {
// 	logging.SetLevel(logging.DEBUG, "")
// 	testScenario(t, &expressionScenario{
// 		description: "Dont explode alias and anchor - check alias parent",
// 		skipDoc:     true,
// 		document:    `{a: &a [1], b: *a}`,
// 		expression:  `.b[0]`,
// 		expected: []string{
// 			"D0, P[a 0], (!!int)::1\n",
// 		},
// 	})
// 	logging.SetLevel(logging.ERROR, "")
// }
