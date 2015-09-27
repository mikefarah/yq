# yaml
yaml command line tool written in go

Allows you to read (and soon update) yaml files given a yaml path.

Usage:
```
yaml <yaml file> <path>
```

E.g.:
```
yaml sample.yaml b.c
```
will output the value of '2'.

Arrays:
Just use the index to access a specific element:
e.g.: given
```
b:
  e:
    - name: fred
      value: 3
    - name: sam
      value: 4
```
then
```
yaml sample.yaml b.e.1.name
```
will output 'sam'
