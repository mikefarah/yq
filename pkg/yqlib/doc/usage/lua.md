
## Basic example
Given a sample.yml file of:
```yaml
---
country: Australia # this place
cities:
- Sydney
- Melbourne
- Brisbane
- Perth
```
then
```bash
yq -o=lua '.' sample.yml
```
will output
```lua
return {
	["country"] = "Australia"; -- this place
	["cities"] = {
		"Sydney",
		"Melbourne",
		"Brisbane",
		"Perth",
	};
};
```

## Unquoted keys
Uses the `--lua-unquoted` option to produce a nicer-looking output.

Given a sample.yml file of:
```yaml
---
country: Australia # this place
cities:
- Sydney
- Melbourne
- Brisbane
- Perth
```
then
```bash
yq -o=lua '.' sample.yml
```
will output
```lua
return {
	country = "Australia"; -- this place
	cities = {
		"Sydney",
		"Melbourne",
		"Brisbane",
		"Perth",
	};
};
```

## Globals
Uses the `--lua-globals` option to export the values into the global scope.

Given a sample.yml file of:
```yaml
---
country: Australia # this place
cities:
- Sydney
- Melbourne
- Brisbane
- Perth
```
then
```bash
yq -o=lua '.' sample.yml
```
will output
```lua
country = "Australia"; -- this place
cities = {
	"Sydney",
	"Melbourne",
	"Brisbane",
	"Perth",
};
```

## Elaborate example
Given a sample.yml file of:
```yaml
---
hello: world
tables:
  like: this
  keys: values
  ? look: non-string keys
  : True
numbers:
  - decimal: 12345
  - hex: 0x7fabc123
  - octal: 0o30
  - float: 123.45
  - infinity: .inf
  - not: .nan

```
then
```bash
yq -o=lua '.' sample.yml
```
will output
```lua
return {
	["hello"] = "world";
	["tables"] = {
		["like"] = "this";
		["keys"] = "values";
		[{
			["look"] = "non-string keys";
		}] = true;
	};
	["numbers"] = {
		{
			["decimal"] = 12345;
		},
		{
			["hex"] = 0x7fabc123;
		},
		{
			["octal"] = 24;
		},
		{
			["float"] = 123.45;
		},
		{
			["infinity"] = (1/0);
		},
		{
			["not"] = (0/0);
		},
	};
};
```

