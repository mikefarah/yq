# How it works

In `yq` expressions are made up of operators and pipes. A context of nodes is passed through the expression and each operation takes the context as input and returns a new context as output. That output is piped in as input for the next operation in the expression. To begin with, the context is set to the first yaml document of the first yaml file (if processing in sequence using eval).

Lets look at a couple of examples.

## Example with a simple operator

Given a document like:

```yaml
- [a]
- "cat"
```

with an expression:

```
.[] | length
```

`yq` will initially set the context as single node of the entire yaml document, an array of two elements.

```yaml
- [a]
- "cat"
```

This gets piped into the splat operator `.[]` which will split out the context into a collection of two nodes `[a]` and `"cat"`. Note that this is _not_ a yaml array.

The `length` operator take no arguments, and will simply return the length of _each_ matching node in the context. So for the context of `[a]` and `"cat"`, it will return a new context of `1` and `3`.

This being the last operation in the expression, the results will be printed out:

```
1
3
```

## Example with an operator that takes arguments.

Given a document like:

```yaml
a: cat
b: dog
```

with an expression:

```
.a = .b
```

The `=` operator takes two arguments, a `lhs` expression, which in this case is `.a` and `rhs` expression which is `.b`. 

It pipes the current, lets call it 'root' context through the `lhs` expression of `.a` to return the node 

```yaml
cat
```

Note that this node holds not only its value 'cat', but comments and metadata too, including path and parent information.

The `=` operator then pipes the 'root' context through the `rhs` expression of `.b` to return the node

```yaml
dog
```

Both sides have now been evaluated, so now the operator copies across the value from the RHS to the value on the LHS, and it returns the now updated context:

```yaml
a: dog
b: dog
```

## Relative update (e.g. `|=`)
There is another form of the `=` operator which we call the relative form. It's very similar to `=` but with one key difference when evaluating the RHS expression.

In the plain form, we pass in the 'root' level context to the RHS expression. In relative form, we pass in _each result of the LHS_ to the RHS expression. Let's go through an example.

Given a document like:

```yaml
a: 1
b: thing
```

with an expression:

```
.a |= . + 1
```

Similar to the `=` operator, `|=` takes two operands, the LHS and RHS.

It pipes the current context (the whole document) through the LHS expression of `.a` to get the node value:

```
1
```

Now it pipes _that LHS context_ into the RHS expression `. + 1` (whereas in the `=` plain form it piped the original document context into the RHS) to yield:


```
2
```

The assignment operator then copies across the value from the RHS to the value on the LHS, and it returns the now updated 'root' context:

```yaml
a: 2
b: thing
```