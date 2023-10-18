package yqlib

import "testing"

const validationExpression = `
	with(env(numberOfCats); select(tag == "!!int") or error("numberOfCats is not a number :(")) | 
	.numPets = env(numberOfCats)
`

var errorOperatorScenarios = []expressionScenario{
	{
		description:   "Validate a particular value",
		document:      `a: hello`,
		expression:    `select(.a == "howdy") or error(".a [" + .a + "] is not howdy!")`,
		expectedError: ".a [hello] is not howdy!",
	},
	{
		description:          "Validate the environment variable is a number - invalid",
		environmentVariables: map[string]string{"numberOfCats": "please"},
		expression:           `env(numberOfCats) | select(tag == "!!int") or error("numberOfCats is not a number :(")`,
		expectedError:        "numberOfCats is not a number :(",
	},
	{
		description:          "Validate the environment variable is a number - valid",
		subdescription:       "`with` can be a convenient way of encapsulating validation.",
		environmentVariables: map[string]string{"numberOfCats": "3"},
		document:             "name: Bob\nfavouriteAnimal: cat\n",
		expression:           validationExpression,
		expected: []string{
			"D0, P[], (!!map)::name: Bob\nfavouriteAnimal: cat\nnumPets: 3\n",
		},
	},
}

func TestErrorOperatorScenarios(t *testing.T) {
	for _, tt := range errorOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "error", errorOperatorScenarios)
}
