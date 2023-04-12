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


## Encode CSV simple
Given a sample.yml file of:
```yaml
- [i, like, csv]
- [because, excel, is, cool]
```
then
```bash
yq -o=csv sample.yml
```
will output
```csv
i,like,csv
because,excel,is,cool
```

## Encode TSV simple
Given a sample.yml file of:
```yaml
- [i, like, csv]
- [because, excel, is, cool]
```
then
```bash
yq -o=tsv sample.yml
```
will output
```tsv
i	like	csv
because	excel	is	cool
```

## Encode array of objects to csv
Given a sample.yml file of:
```yaml
- name: Gary
  numberOfCats: 1
  likesApples: true
  height: 168.8
- name: Samantha's Rabbit
  numberOfCats: 2
  likesApples: false
  height: -188.8

```
then
```bash
yq -o=csv sample.yml
```
will output
```csv
name,numberOfCats,likesApples,height
Gary,1,true,168.8
Samantha's Rabbit,2,false,-188.8
```

## Encode array of objects to custom csv format
Add the header row manually, then the we convert each object into an array of values - resulting in an array of arrays. Pick the columns and call the header whatever you like.

Given a sample.yml file of:
```yaml
- name: Gary
  numberOfCats: 1
  likesApples: true
  height: 168.8
- name: Samantha's Rabbit
  numberOfCats: 2
  likesApples: false
  height: -188.8

```
then
```bash
yq -o=csv '[["Name", "Number of Cats"]] +  [.[] | [.name, .numberOfCats ]]' sample.yml
```
will output
```csv
Name,Number of Cats
Gary,1
Samantha's Rabbit,2
```

