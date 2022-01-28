# Working with Properties

## Yaml to Properties

To convert to property file format, use `--outputformat=props` or `-o=p`. Porting the comments from the `yaml` file to the property file is still in progress, currently line comments on the node values will be copied across.&#x20;

Given a sample file of:

```yaml
# block comments don't come through
person: # neither do comments on maps
    name: Mike # comments on values appear
    pets: 
    - cat # comments on array values appear
    food: [pizza] # comments on arrays do not
```

then

```bash
yq -o=p sample.yaml
```

will output:

```
# comments on values appear
person.name = Mike

# comments on array values appear
person.pets.0 = cat
person.food.0 = pizza
```

