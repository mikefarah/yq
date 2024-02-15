# CSV
Encode/Decode/Roundtrip CSV and TSV files.

## Encode 
Currently supports arrays of homogeneous flat objects, that is: no nesting and it assumes the _first_ object has all the keys required:

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

## Encode array of objects to csv - missing fields behaviour
First entry is used to determine the headers, and it is missing 'likesApples', so it is not included in the csv. Second entry does not have 'numberOfCats' so that is blank

Given a sample.yml file of:
```yaml
- name: Gary
  numberOfCats: 1
  height: 168.8
- name: Samantha's Rabbit
  height: -188.8
  likesApples: false

```
then
```bash
yq -o=csv sample.yml
```
will output
```csv
name,numberOfCats,height
Gary,1,168.8
Samantha's Rabbit,,-188.8
```

## Parse CSV into an array of objects
First row is assumed to be the header row. By default, entries with YAML/JSON formatting will be parsed!

Given a sample.csv file of:
```csv
name,numberOfCats,likesApples,height,facts
Gary,1,true,168.8,cool: true
Samantha's Rabbit,2,false,-188.8,tall: indeed

```
then
```bash
yq -p=csv sample.csv
```
will output
```yaml
- name: Gary
  numberOfCats: 1
  likesApples: true
  height: 168.8
  facts:
    cool: true
- name: Samantha's Rabbit
  numberOfCats: 2
  likesApples: false
  height: -188.8
  facts:
    tall: indeed
```

## Parse CSV into an array of objects, no auto-parsing
First row is assumed to be the header row. Entries with YAML/JSON will be left as strings.

Given a sample.csv file of:
```csv
name,numberOfCats,likesApples,height,facts
Gary,1,true,168.8,cool: true
Samantha's Rabbit,2,false,-188.8,tall: indeed

```
then
```bash
yq -p=csv --csv-auto-parse=f sample.csv
```
will output
```yaml
- name: Gary
  numberOfCats: 1
  likesApples: true
  height: 168.8
  facts: 'cool: true'
- name: Samantha's Rabbit
  numberOfCats: 2
  likesApples: false
  height: -188.8
  facts: 'tall: indeed'
```

## Parse TSV into an array of objects
First row is assumed to be the header row.

Given a sample.tsv file of:
```tsv
name	numberOfCats	likesApples	height
Gary	1	true	168.8
Samantha's Rabbit	2	false	-188.8

```
then
```bash
yq -p=tsv sample.tsv
```
will output
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

## Round trip
Given a sample.csv file of:
```csv
name,numberOfCats,likesApples,height
Gary,1,true,168.8
Samantha's Rabbit,2,false,-188.8

```
then
```bash
yq -p=csv -o=csv '(.[] | select(.name == "Gary") | .numberOfCats) = 3' sample.csv
```
will output
```csv
name,numberOfCats,likesApples,height
Gary,3,true,168.8
Samantha's Rabbit,2,false,-188.8
```

