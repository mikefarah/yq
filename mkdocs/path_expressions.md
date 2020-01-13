Path expressions are used to deeply navigate and match particular yaml nodes.

_As a general rule, you should wrap paths in quotes to prevent your CLI from processing '*, []' and other special characters._

## Simple expressions

### Maps

a.b.c

```yaml
a:
  b:
    c: thing # MATCHES
```

### Arrays

a.b[1].c

```yaml
a:
  b:
  - c: thing0 
  - c: thing1 # MATCHES
  - c: thing2
```

#### Appending to arrays
(e.g. when using the write command)

a.b[+].c


```yaml
a:
  b:
  - c: thing0 
```

Will add a new entry:

```yaml
a:
  b:
  - c: thing0 
  - c: thing1 # NEW entry from [+] on B array.
```


## Splat

### Maps
a.*.c

```yaml
a:
  b1:
    c: thing # MATCHES
    d: whatever
  b2:
    c: thing # MATCHES
    f: something irrelevant
```

#### Prefix splat

bob.item*.cats

```yaml
bob:
  item:
    cats: bananas # MATCHES
  something:
    cats: lemons
  itemThing:
    cats: more bananas # MATCHES
  item2:
    cats: apples # MATCHES
  thing:
    cats: oranges
```

### Arrays
a.b[*].c

```yaml
a:
  b:
  - c: thing0 # MATCHES
    d: what..ever
  - c: thing1 # MATCHES
    d: blarh
  - c: thing2 # MATCHES
    f: thingamabob
```

## Deep Splat

'**' will match arbitrary nodes for both maps and arrays:

a.**.c

```yaml
a:
  b1:
    c: thing1 # MATCHES
    d: cat cat
  b2:
    c: thing2 # MATCHES
    d: dog dog
  b3:
    d:
    - f:
        c: thing3 # MATCHES
        d: beep
    - f:
        g:
          c: thing4 # MATCHES
          d: boop
    - d: mooo
```


## Search by children nodes

a.(b.d==cat).b.c

```yaml
a:
  - b:
      c: thing0
      d: leopard
    ba: fast
  - b:
      c: thing1 # MATCHES
      d: cat
    ba: meowy
  - b:
      c: thing2
      d: caterpillar
    ba: icky
  - b:
      c: thing3 # MATCHES
      d: cat
    ba: also meowy
```

### With prefixes

a.(b.d==cat*).c

```yaml
a:
  - b:
      c: thing0
      d: leopard
    ba: fast
  - b:
      c: thing1 # MATCHES
      d: cat
    ba: meowy
  - b:
      c: thing2 # MATCHES
      d: caterpillar
    ba: icky
  - b:
      c: thing3 # MATCHES
      d: cat
    ba: also meowy
```


## Special Characters


### Keys with dots
When specifying a key that has a dot use key lookup indicator.

```yaml
b:
  foo.bar: 7
```

```bash
yaml r sample.yaml 'b[foo.bar]'
```

```bash
yaml w sample.yaml 'b[foo.bar]' 9
```

Any valid yaml key can be specified as part of a key lookup.

Note that the path is in quotes to avoid the square brackets being interpreted by your shell.

### Keys (and values) with leading dashes
The flag terminator needs to be used to stop the app from attempting to parse the subsequent arguments as flags, if they start if a dash.

```bash
yq n -j -- --key --value
```

Will result in

```
--key: --value
```