# Equals / Not Equals

This is a boolean operator that will return `true` if the LHS is equal to the RHS and `false` otherwise.

```
.a == .b
```

It is most often used with the select operator to find particular nodes:

```
select(.a == .b)
```

The not equals `!=` operator returns `false` if the LHS is equal to the RHS.

## Related Operators

- comparison (`>=`, `<` etc) operators [here](https://mikefarah.gitbook.io/yq/operators/compare)
- boolean operators (`and`, `or`, `any` etc) [here](https://mikefarah.gitbook.io/yq/operators/boolean-operators)
- select operator [here](https://mikefarah.gitbook.io/yq/operators/select)

