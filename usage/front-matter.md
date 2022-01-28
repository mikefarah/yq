# Front Matter

`yq` can process files with `yaml` front matter (e.g. jekyll, assemble and others) - this is done via the `--front-matter/-f` flag.

## Process front matter

Use `--front-matter=process` to process the front matter, that is run the expression against the `yaml` content, and output back the entire file, included the non-yaml content block. For example:

File:

```
---
a: apple
b: bannana
---
<h1>I like {{a}} and {{b}} </h1>
```

The running

```
yq --front-matter=process '.a="chocolate"' file.jekyll
```

Will yield:

```
---
a: chocolate
b: bannana
---
<h1>I like {{a}} and {{b}} </h1>
```

## Extract front matter

Running with `--front-matter=extract` will only output the yaml contents and ignore the rest. From the previous example, if you were to instead run:

```
yq --front-matter=extract '.a="chocolate"' file.jekyll
```

Then this would yield:

```
---
a: chocolate
b: bannana
```

