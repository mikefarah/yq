# How it works

In `yq` expressions are made up of operators and pipes. A context of nodes is passed through the expression and each operation takes the context as input and returns a new context as output. That output is piped in as input for the next operation in the expression. To begin with, the context is set to the first yaml document of the first yaml file (if processing in sequence using eval).

Lets look at a couple of examples.

## Simple assignment example

Given a document like:

```yaml
a: cat
b: dog
```

with an expression:

```
.a = .b
```

Like math expressions - operator precedence is important. 

The `=` operator takes two arguments, a `lhs` expression, which in this case is `.a` and `rhs` expression which is `.b`. 

It pipes the current, lets call it 'root' context through the `lhs` expression of `.a` to return the node 

```yaml
cat
```

Sidenote: this node holds not only its value 'cat', but comments and metadata too, including path and parent information.

The `=` operator then pipes the 'root' context through the `rhs` expression of `.b` to return the node

```yaml
dog
```

Both sides have now been evaluated, so now the operator copies across the value from the RHS (`.b`) to the LHS (`.a`), and it returns the now updated context:

```yaml
a: dog
b: dog
```


## Complex assignment, operator precedence rules

Just like math expressions - `yq` expressions have an order of precedence. The pipe `|` operator has a low order of precedence, so operators with higher precedence will get evaluated first. 

Most of the time, this is intuitively what you'd want, for instance `.a = "cat" | .b = "dog"` is effectively: `(.a = "cat") | (.b = "dog")`.

However, this is not always the case, particularly if you have a complex LHS or RHS expression, for instance if you want to select particular nodes to update. 

Lets say you had:

```yaml
- name: bob
  fruit: apple
- name: sally
  fruit: orange

```

Lets say you wanted to update the `sally` entry to have fruit: 'mango'. The _incorrect_ way to do that is:
`.[] | select(.name == "sally") | .fruit = "mango"`.

Because `|` has a low operator precedence, this will be evaluated (_incorrectly_) as : `(.[]) | (select(.name == "sally")) | (.fruit = "mango")`. What you'll see is only the updated segment returned:

```yaml
name: sally
fruit: mango
```

To properly update this yaml, you will need to use brackets (think BODMAS from maths) and wrap the entire LHS:
`(.[] | select(.name == "sally") | .fruit) = "mango"`


Now that entire LHS expression is passed to the 'assign' (`=`) operator, and the yaml is correctly updated and returned:


```yaml
- name: bob
  fruit: apple
- name: sally
  fruit: mango

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
