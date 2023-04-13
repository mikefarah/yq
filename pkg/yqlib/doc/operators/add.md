# Add

Add behaves differently according to the type of the LHS:
* arrays: concatenate
* number scalars: arithmetic addition
* string scalars: concatenate
* maps: shallow merge (use the multiply operator (`*`) to deeply merge)

Use `+=` as a relative append assign for things like increment. Note that `.a += .x` is equivalent to running `.a = .a + .x`.


## 
Given a sample.yml file of:
```yaml
a: hello
```
then
```bash
yq sample.yml
```
will output
```yaml
a: hello
```

