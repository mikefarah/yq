# Error

Use this operation to short-circuit expressions. Useful for validation.

## Validate a particular value
Given a sample.yml file of:
```yaml
a: hello
```
then
```bash
yq 'select(.a == "howdy") or error(".a [" + .a + "] is not howdy!")' sample.yml
```
will output
```bash
Error: .a [hello] is not howdy!
```

## Validate the environment variable is a number - invalid
Running
```bash
numberOfCats="please" yq --null-input 'env(numberOfCats) | select(tag == "!!int") or error("numberOfCats is not a number :(")'
```
will output
```bash
Error: numberOfCats is not a number :(
```

## Validate the environment variable is a number - valid
`with` can be a convenient way of encapsulating validation.

Given a sample.yml file of:
```yaml
name: Bob
favouriteAnimal: cat
```
then
```bash
numberOfCats="3" yq '
	with(env(numberOfCats); select(tag == "!!int") or error("numberOfCats is not a number :(")) | 
	.numPets = env(numberOfCats)
' sample.yml
```
will output
```yaml
name: Bob
favouriteAnimal: cat
numPets: 3
```

