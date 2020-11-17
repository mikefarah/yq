# Mulitply Operator
## Examples
### Merge objects together
sample.yml:
```yaml
{a: {also: me}, b: {also: {g: wizz}}}
```
Expression
```bash
yq '. * {"a":.b}' < sample.yml
```
Result
```yaml
{a: {also: {g: wizz}}, b: {also: {g: wizz}}}
```
### Merge keeps style of LHS
sample.yml:
```yaml
a: {things: great}
b:
  also: "me"

```
Expression
```bash
yq '. * {"a":.b}' < sample.yml
```
Result
```yaml
a: {things: great, also: "me"}
b:
  also: "me"
```
### Merge arrays
sample.yml:
```yaml
{a: [1,2,3], b: [3,4,5]}
```
Expression
```bash
yq '. * {"a":.b}' < sample.yml
```
Result
```yaml
{a: [3, 4, 5], b: [3, 4, 5]}
```
### Merge to prefix an element
sample.yml:
```yaml
{a: cat, b: dog}
```
Expression
```bash
yq '. * {"a": {"c": .a}}' < sample.yml
```
Result
```yaml
{a: {c: cat}, b: dog}
```
### Merge with simple aliases
sample.yml:
```yaml
{a: &cat {c: frog}, b: {f: *cat}, c: {g: thongs}}
```
Expression
```bash
yq '.c * .b' < sample.yml
```
Result
```yaml
{g: thongs, f: *cat}
```
### Merge does not copy anchor names
sample.yml:
```yaml
{a: {c: &cat frog}, b: {f: *cat}, c: {g: thongs}}
```
Expression
```bash
yq '.c * .a' < sample.yml
```
Result
```yaml
{g: thongs, c: frog}
```
### Merge with merge anchors
sample.yml:
```yaml

foo: &foo
  a: foo_a
  thing: foo_thing
  c: foo_c

bar: &bar
  b: bar_b
  thing: bar_thing
  c: bar_c

foobarList:
  b: foobarList_b
  <<: [*foo,*bar]
  c: foobarList_c

foobar:
  c: foobar_c
  <<: *foo
  thing: foobar_thing

```
Expression
```bash
yq '.foobar * .foobarList' < sample.yml
```
Result
```yaml
c: foobarList_c
<<: [*foo, *bar]
thing: foobar_thing
b: foobarList_b
```
