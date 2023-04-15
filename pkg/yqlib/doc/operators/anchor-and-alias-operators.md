# Anchor and Alias Operators

Use the `alias` and `anchor` operators to read and write yaml aliases and anchors. The `explode` operator normalises a yaml file (dereference (or expands) aliases and remove anchor names).

`yq` supports merge aliases (like `<<: *blah`) however this is no longer in the standard yaml spec (1.2) and so `yq` will automatically add the `!!merge` tag to these nodes as it is effectively a custom tag.


## Dereference and update a field
Use explode with multiply to dereference an object

Given a sample.yml file of:
```yaml
item_value: &item_value
  value: true
thingOne:
  name: item_1
  !!merge <<: *item_value
thingTwo:
  name: item_2
  !!merge <<: *item_value
```
then
```bash
yq '.thingOne |= explode(.) * {"value": false}' sample.yml
```
will output
```yaml
item_value: &item_value
  value: true
thingOne: false
thingTwo:
  name: item_2
  !!merge <<: *item_value
```

