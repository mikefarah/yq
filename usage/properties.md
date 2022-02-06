# Properties

Encode to a property file (decode not yet supported). Line comments on value nodes will be copied across.

By default, empty maps and arrays are not encoded - see below for an example on how to encode a value for these.

{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Encode properties
Note that empty arrays and maps are not encoded by default.

Given a sample.yml file of:
```yaml
# block comments don't come through
person: # neither do comments on maps
    name: Mike # comments on values appear
    pets: 
    - cat # comments on array values appear
    food: [pizza] # comments on arrays do not
emptyArray: []
emptyMap: []

```
then
```bash
yq -o=props -I=0 '.' sample.yml
```
will output
```properties
# comments on values appear
person.name = Mike

# comments on array values appear
person.pets.0 = cat
person.food.0 = pizza
```

## Encode properties: no comments
Given a sample.yml file of:
```yaml
# block comments don't come through
person: # neither do comments on maps
    name: Mike # comments on values appear
    pets: 
    - cat # comments on array values appear
    food: [pizza] # comments on arrays do not
emptyArray: []
emptyMap: []

```
then
```bash
yq -o=props -I=0 '... comments = ""' sample.yml
```
will output
```properties
person.name = Mike
person.pets.0 = cat
person.food.0 = pizza
```

## Encode properties: include empty maps and arrays
Use a yq expression to set the empty maps and sequences to your desired value.

Given a sample.yml file of:
```yaml
# block comments don't come through
person: # neither do comments on maps
    name: Mike # comments on values appear
    pets: 
    - cat # comments on array values appear
    food: [pizza] # comments on arrays do not
emptyArray: []
emptyMap: []

```
then
```bash
yq -o=props -I=0 '(.. | select( (tag == "!!map" or tag =="!!seq") and length == 0)) = ""' sample.yml
```
will output
```properties
# comments on values appear
person.name = Mike

# comments on array values appear
person.pets.0 = cat
person.food.0 = pizza
emptyArray = 
emptyMap = 
```

