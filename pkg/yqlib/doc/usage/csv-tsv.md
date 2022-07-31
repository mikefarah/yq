
{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

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
Add the header row manually, then the we convert each object into an array of values - resulting in an array of arrays. Nice thing about this method is you can pick the columns and call the header whatever you like.

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

## Encode array of objects to csv - generic
This is a little trickier than the previous example - we dynamically work out the $header, and use that to automatically create the value arrays.

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
yq -o=csv '(.[0] | keys | .[] ) as $header |  [[$header]] +  [.[] | [ .[$header] ]]' sample.yml
```
will output
```csv
name,numberOfCats,likesApples,height
Gary,1,true,168.8
Samantha's Rabbit,2,false,-188.8
```

## Parse CSV into an array of objects
First row is assumed to define the fields

Given a sample.csv file of:
```csv
name,numberOfCats,likesApples,height
Gary,1,true,168.8
Samantha's Rabbit,2,false,-188.8

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
- name: Samantha's Rabbit
  numberOfCats: 2
  likesApples: false
  height: -188.8
```

## Parse TSV into an array of objects
First row is assumed to define the fields

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

