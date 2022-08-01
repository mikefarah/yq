# CSV
Encode/Decode/Roundtrip CSV and TSV files.

## Encode 
Currently supports arrays of homogenous flat objects, that is: no nesting and it assumes the _first_ object has all the keys required:

```yaml
- name: Bobo
  type: dog
- name: Fifi
  type: cat
```

As well as arrays of arrays of scalars (strings/numbers/booleans):

```yaml
- [Bobo, dog]
- [Fifi, cat]
```

## Decode
Decode assumes the first CSV/TSV row is the header row, and all rows beneath are the entries.
The data will be coded into an array of objects, using the header rows as keys.

```csv
name,type
Bobo,dog
Fifi,cat
```

