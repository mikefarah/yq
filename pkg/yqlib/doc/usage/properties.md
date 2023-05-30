# Properties

Encode/Decode/Roundtrip to/from a property file. Line comments on value nodes will be copied across.

By default, empty maps and arrays are not encoded - see below for an example on how to encode a value for these.

## Encode properties
Note that empty arrays and maps are not encoded by default.

Given a sample.yml file of:
```yaml
# block comments come through
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
# block comments come through
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
# block comments come through
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
# block comments come through
# comments on values appear
person.name = "Mike Wazowski"

# comments on array values appear
person.pets.0 = cat
person.food.0 = pizza
```

## Encode properties: no comments
Given a sample.yml file of:
```yaml
# block comments come through
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
# block comments come through
person.name = Mike Wazowski
person.pets.0 = cat
person.food.0 = pizza
```

## Encode properties: include empty maps and arrays
Use a yq expression to set the empty maps and sequences to your desired value.

Given a sample.yml file of:
```yaml
# block comments come through
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
# block comments come through
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
# block comments come through
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
    # block comments come through
    # comments on values appear
    name: Mike Wazowski
    pets:
        # comments on array values appear
        - cat
    food:
        - pizza
```

## Decode properties - array should be a map
If you have a numeric map key in your property files, use array_to_map to convert them to maps.

Given a sample.properties file of:
```properties
things.10 = mike
```
then
```bash
yq -p=props '.things |= array_to_map' sample.properties
```
will output
```yaml
things:
    10: mike
```

## Roundtrip
Given a sample.properties file of:
```properties
# block comments come through
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
# block comments come through
# comments on values appear
person.name = Mike Wazowski

# comments on array values appear
person.pets.0 = dog
person.food.0 = pizza
```

