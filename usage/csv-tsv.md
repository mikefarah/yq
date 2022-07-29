# Working with CSV, TSV

{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Yaml to CSV/TSV

You can convert compatible yaml structures to CSV or TSV by using:

* `--output-format=csv` or `-o=c` for csv (comma separated values)
* `--output-format=tsv` or `-o=t` for tsv (tab separated values)

Compatible structures is either an array of scalars (strings/numbers/booleans), which is a single row; or an array of arrays of scalars (multiple rows).

```yaml
- [i, like, csv]
- [because, excel, is, cool]
```

then

```bash
yq '.' -o=csv sample.yaml
```

will output:

```csv
i,like,csv
because,excel,is,cool
```

Similarly, for tsv:

```bash
yq '.' -o=tsv sample.yaml
```

will output:

```
i	like	csv
because	excel	is	cool
```

### Converting an array of objects to CSV

If you have a yaml document like:

```yaml
foo: bar
items:
  - name: Tom
    species: cat
    color: blue
  - name: Jim
    species: dog
    color: brown
```

To convert to CSV, you need to transform this into a array of CSV rows. Assuming you also want a header, then you can do:

```bash
yq '[["name", "species", "color"]] + [.items[] | [.name, .species, .color]]' data.yaml -o=csv
```

to yield:

```csv
name,species,color
Tom,cat,blue
Jim,dog,brown
```
