# Expression Syntax: A Visual Guide
In `yq`, expressions are made up of operators and pipes. A context of nodes is passed through the expression, and each operation takes the context as input and returns a new context as output. That output is piped in as input for the next operation in the expression.

Let's break down the process step by step using a diagram. We'll start with a single YAML document, apply an expression, and observe how the context changes at each step.

Given a document like:

```yaml
root:
  items:
    - name: apple
      type: fruit
    - name: carrot
      type: vegetable
    - name: banana
      type: fruit
```

You can use dot notation to access nested structures. For example, to access the `name` of the first item, you would use the expression `.root.items[0].name`, which would return `apple`.

But lets see how we could find all the fruit under `items`

## Step 1: Initial Context
The context starts at the root of the YAML document. In this case, the entire document is the initial context. 

```
root
└── items
    ├── name: apple
    │   type: fruit
    ├── name: carrot
    │   type: vegetable
    └── name: banana
        type: fruit
```

## Step 2: Splatting the Array
Using the expression `.root.items[]`, we "splat" the items array. This means each element of the array becomes its own node in the context:

```
Node 1: { name: apple, type: fruit }
Node 2: { name: carrot, type: vegetable }
Node 3: { name: banana, type: fruit }
```

## Step 3: Filtering the Nodes
Next, we apply a filter to select only the nodes where type is fruit. The expression `.root.items[] | select(.type == "fruit")` filters the nodes:

```
Filtered Node 1: { name: apple, type: fruit }
Filtered Node 2: { name: banana, type: fruit }
```

## Step 4: Extracting a Field
Finally, we extract the name field from the filtered nodes using `.root.items[] | select(.type == "fruit") | .name` This results in:

```
apple
banana
```

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

Side note: this node holds not only its value 'cat', but comments and metadata too, including path and parent information.

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

**Important**: To properly update this YAML, you must wrap the entire LHS in parentheses. Think of it like using brackets in math to ensure the correct order of operations.
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