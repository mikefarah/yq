Yaml files can be created using the 'new' command. This works in the same way as the write command, but you don't pass in an existing Yaml file.

```
yaml n <path> <new value>
```

### Creating a simple yaml file
```bash
yaml n b.c cat
```
will output:
```yaml
b:
  c: cat
```

### Creating using a create script
Create scripts follow the same format as the update scripts.

Given a script create_instructions.yaml of:
```yaml
b.c: 3
b.e[0].name: Howdy Partner
```
then

```bash
yaml n -s create_instructions.yaml
```
will output:
```yaml
b:
  c: 3
  e:
    - name: Howdy Partner
```

You can also pipe the instructions in:

```bash
cat create_instructions.yaml | yaml n -s -
```