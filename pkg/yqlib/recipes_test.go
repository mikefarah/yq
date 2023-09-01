package yqlib

import (
	"testing"
)

var recipes = []expressionScenario{
	{
		description: "Find items in an array",
		explanation: []string{
			"`.[]` splats the array, and puts all the items in the context.",
			"These items are then piped (`|`) into `select(.name == \"Foo\")` which will select all the nodes that have a name property set to 'Foo'.",
			"See the [select](https://mikefarah.gitbook.io/yq/operators/select) operator for more information.",
		},
		document:   `[{name: Foo, numBuckets: 0}, {name: Bar, numBuckets: 0}]`,
		expression: `.[] | select(.name == "Foo")`,
		expected: []string{
			"D0, P[0], (!!map)::{name: Foo, numBuckets: 0}\n",
		},
	},
	{
		description: "Find and update items in an array",
		document:    `[{name: Foo, numBuckets: 0}, {name: Bar, numBuckets: 0}]`,
		expression:  `(.[] | select(.name == "Foo") | .numBuckets) |= . + 1`,
		explanation: []string{
			"Following from the example above`.[]` splats the array, selects filters the items.",
			"We then pipe (`|`) that into `.numBuckets`, which will select that field from all the matching items",
			"Splat, select and the field are all in brackets, that whole expression is passed to the `|=` operator as the left hand side expression, with `. + 1` as the right hand side expression.",
			"`|=` is the operator that updates fields relative to their own value, which is referenced as dot (`.`).",
			"The expression `. + 1` increments the numBuckets counter.",
			"See the [assign](https://mikefarah.gitbook.io/yq/operators/assign-update) and [add](https://mikefarah.gitbook.io/yq/operators/add) operators for more information.",
		},
		expected: []string{
			"D0, P[], (doc)::[{name: Foo, numBuckets: 1}, {name: Bar, numBuckets: 0}]\n",
		},
	},
}

func TestRecipes(t *testing.T) {
	for _, tt := range recipes {
		testScenario(t, &tt)
	}
	genericScenarios := make([]interface{}, len(recipes))
	for i, s := range recipes {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "recipes", genericScenarios, documentOperatorScenario)
}
