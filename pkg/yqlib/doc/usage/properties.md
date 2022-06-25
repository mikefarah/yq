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
    name: Mike Wazowski # comments on values appear
    pets: 
    - cat # comments on array values appear
    food: [pizza] # comments on arrays do not
emptyArray: []
emptyMap: []

```
then
```bash
yq -o=props sample.yml
```
will output
```properties
# comments on values appear
person.name = Mike Wazowski

# comments on array values appear
person.pets.0 = cat
person.food.0 = pizza
```

## Encode properties: scalar encapsulation
Note that string values with blank characters in them are encapsulated with double quotes

Given a sample.yml file of:
```yaml
# block comments don't come through
person: # neither do comments on maps
    name: Mike Wazowski # comments on values appear
    pets: 
    - cat # comments on array values appear
    food: [pizza] # comments on arrays do not
emptyArray: []
emptyMap: []

```
then
```bash
yq -o=props --unwrapScalar=false sample.yml
```
will output
```properties
# comments on values appear
person.name = "Mike Wazowski"

# comments on array values appear
person.pets.0 = cat
person.food.0 = pizza
```

## Encode properties: no comments
Given a sample.yml file of:
```yaml
# block comments don't come through
person: # neither do comments on maps
    name: Mike Wazowski # comments on values appear
    pets: 
    - cat # comments on array values appear
    food: [pizza] # comments on arrays do not
emptyArray: []
emptyMap: []

```
then
```bash
yq -o=props '... comments = ""' sample.yml
```
will output
```properties
person.name = Mike Wazowski
person.pets.0 = cat
person.food.0 = pizza
```

## Encode properties: include empty maps and arrays
Use a yq expression to set the empty maps and sequences to your desired value.

Given a sample.yml file of:
```yaml
# block comments don't come through
person: # neither do comments on maps
    name: Mike Wazowski # comments on values appear
    pets: 
    - cat # comments on array values appear
    food: [pizza] # comments on arrays do not
emptyArray: []
emptyMap: []

```
then
```bash
yq -o=props '(.. | select( (tag == "!!map" or tag =="!!seq") and length == 0)) = ""' sample.yml
```
will output
```properties
# comments on values appear
person.name = Mike Wazowski

# comments on array values appear
person.pets.0 = cat
person.food.0 = pizza
emptyArray = 
emptyMap = 
```

## Decode properties
Given a sample.properties file of:
```properties
# comments on values appear
person.name = Mike Wazowski

# comments on array values appear
person.pets.0 = cat
person.food.0 = pizza

```
then
```bash
yq -p=props sample.properties
```
will output
```yaml
person:
    name: Mike Wazowski # comments on values appear
    pets:
        - cat # comments on array values appear
    food:
        - pizza
```

## Roundtrip
Given a sample.properties file of:
```properties
# comments on values appear
person.name = Mike Wazowski

# comments on array values appear
person.pets.0 = cat
person.food.0 = pizza

```
then
```bash
yq -p=props -o=props '.person.pets.0 = "dog"' sample.properties
```
will output
```properties
# comments on values appear
person.name = Mike Wazowski

# comments on array values appear
person.pets.0 = dog
person.food.0 = pizza
```

