# If Then Else

`if A then B else C end` evaluates `A` against each input, and returns `B` if the result is truthy, otherwise `C`. Like `select`, `null` and `false` are falsy, and everything else is truthy.

`elif` can be used to chain conditions, and `else` is optional - if it is left out and the condition is falsy the input is returned unchanged (like jq 1.7).

## Related Operators

- select operator [here](https://mikefarah.gitbook.io/yq/operators/select)
- alternative (`//`) operator [here](https://mikefarah.gitbook.io/yq/operators/alternative-default-value)
- boolean operators (`and`, `or`, `any` etc) [here](https://mikefarah.gitbook.io/yq/operators/boolean-operators)
