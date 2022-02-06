# Union

This operator is used to combine different results together.

{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Combine scalars
Running
```bash
yq --null-input '1, true, "cat"'
```
will output
```yaml
1
true
cat
```

## Combine selected paths
Given a sample.yml file of:
```yaml
a: fieldA
b: fieldB
c: fieldC
```
then
```bash
yq '.a, .c' sample.yml
```
will output
```yaml
fieldA
fieldC
```

