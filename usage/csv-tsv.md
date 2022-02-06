# Working with CSV, TSV

{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Yaml to CSV/TSV

You can convert compatible yaml structures to CSV or TSV by using:

* `--outputformat=csv` or `-o=c` for csv (comma separated values)
* `--outputformat=tsv` or `-o=t` for tsv (tab separated values)

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
