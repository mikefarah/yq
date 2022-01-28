# Tips, Tricks, Troubleshooting

## Validating yaml files

Yaml files can be surprisingly lenient in what can be parsed as a yaml file. A reasonable way of validation a yaml file is to ensure the top level is a map or array (although it is valid yaml to have scalars at the top level, but often this is not what you want). This can be done by:

```
yq --exit-status 'tag == "!!map" or tag== "!!seq"' file.txt > /dev/null
```

## Split expressions over multiple lines to improve readability

Feel free to use multiple lines in your expression to improve readability.

Use `with` if you need to make several updates to the same path.

```bash
yq --inplace '
  with(.a.deeply.nested; 
    . = "newValue" | . style="single") |
  with(.b.another.nested; 
    . = "cool" | . style="folded")
' my_file.yaml
```

## Create bash array

Given a yaml file like

```yaml
coolActions:
  - create
  - edit
  - delete
```

You can create a bash array named `actions` by:

```bash
> readarray actions < <(yq '.coolActions[]' sample.yaml)
> echo "${actions[1]}"
edit
```

## Set contents from another file

Use an environment variable with the `strenv` operator to inject the contents from an environment variable.

```bash
LICENSE=$(cat LICENSE) yq -n '.a = strenv(LICENSE)'
```

Note that `bash` substitution "$(..)" trims newlines, this will cause string blocks to start with `|-` instead of `|`. If you want to keep your nice trailing newline, read more [here](https://mikefarah.gitbook.io/yq/operators/string-operators#string-blocks-bash-and-newlines)


## Special characters in strings

The `strenv` operator is a great way to handle special characters in strings:

```bash
VAL='.a |!@  == "string2"' yq '.a = strenv(VAL)' example.yaml
```

## Quotes in Windows Powershell

Powershell has its [own](https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about\_quoting\_rules?view=powershell-7.1) way of handling quotes:

```bash
PS > yq -n '.test = ""something""'
test: something
PS >
```

See [https://github.com/mikefarah/yq/issues/747](https://github.com/mikefarah/yq/issues/747) for more trickery.

## Merge / combine all documents into one

To merge all given yaml files into one, use the `reduce` operator with the `*` (multiply) operator. Note the use of `ea` or `eval-all` to load all files into memory so that they can be merged.

```
yq ea '. as $item ireduce ({}; . * $item )' file1.yml file2.yml ...
```

## Merge an array of objects by key

See [here](https://mikefarah.gitbook.io/yq/operators/multiply-merge#merge-arrays-of-objects-together-matching-on-a-key) for a working example.

## Creating a new file / working with blank documents

To create a new `yaml` file simply:

```
yq -n '.someNew="content"' > newfile.yml
```

## Comparing yaml files

The best way to run a diff is to use `yq` to normalise the yaml files and then just use diff. Here is a simple example of using pretty print `-P` to normalise the styling and running diff:

```
diff <(yq -P 'sort_keys(..)' file1.yaml) <(yq -P 'sort_keys(..)' file2.yaml)
```

This way you can use the full power of `diff` and normalise the yaml files as you like.

You may also want to remove all comments using `... comments=""`

## Reading multiple streams (STDINs)

Like `diff` and other bash commands, you can use `<(exp)` to pipe in multiple streams of data into `yq`. instance:

```
yq '.apple' <(curl -s https://somewhere/data1.yaml) <(cat file.yml)
```

## Updating deeply selected paths
### or why is yq only returning the updated yaml

The most important thing to remember to do is to have brackets around the LHS expression - otherwise what `yq` will do is first filter by the selection, and then, separately, update the filtered result and return that subset.

```
yq '(.foo.bar[] | select(name == "fred) | .apple) = "cool"'
```

## Combining multiple files into one

In order to combine multiple yaml files into a single file (with `---` separators) you can just:

```
yq '.' somewhere/*.yaml
```

## Multiple updates to the same path

You can use the [with](../operators/with.md) operator to set a nested context:

```
yq 'with(.a.deeply ; .nested = "newValue" | .other= "newThing")' sample.yml
```

The first argument expression sets the root context, and the second expression runs against that root context.

## yq adds a !!merge tag automatically

The merge functionality from yaml v1.1 (e.g. `<<:`has actually been removed in the 1.2 spec. Thankfully, `yq` underlying yaml parser still supports that tag - and it's extra nice in that it explicitly puts the `!!merge` tag on key of the map entry. This tag tells other yaml parsers that this entry is a merge entry, as opposed to a regular string key that happens to have a value of `<<:`. This is backwards compatible with the 1.1 spec of yaml, it's simply an explicit way of specifying the type (for instance, you can use a `!!str` tag to enforce a particular value to be a string.

Although this does affect the readability of the yaml to humans, it still works and processes fine with various yaml processors.
