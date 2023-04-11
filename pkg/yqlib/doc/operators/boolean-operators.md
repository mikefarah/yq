# Boolean Operators

The `or` and `and` operators take two parameters and return a boolean result. 

`not` flips a boolean from true to false, or vice versa. 

`any` will return `true` if there are any `true` values in a array sequence, and `all` will return true if _all_ elements in an array are true.

`any_c(condition)` and `all_c(condition)` are like `any` and `all` but they take a condition expression that is used against each element to determine if it's `true`. Note: in `jq` you can simply pass a condition to `any` or `all` and it simply works - `yq` isn't that clever..yet

These are most commonly used with the `select` operator to filter particular nodes.

## Related Operators

- equals / not equals (`==`, `!=`) operators [here](https://mikefarah.gitbook.io/yq/operators/equals)
- comparison (`>=`, `<` etc) operators [here](https://mikefarah.gitbook.io/yq/operators/compare)
- select operator [here](https://mikefarah.gitbook.io/yq/operators/select)

