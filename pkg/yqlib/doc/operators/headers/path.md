# Path

The `path` operator can be used to get the traversal paths of matching nodes in an expression. The path is returned as an array, which if traversed in order will lead to the matching node.

You can get the key/index of matching nodes by using the `path` operator to return the path array then piping that through `.[-1]` to get the last element of that array, the key.

Use `setpath` to set a value to the path array returned by `path`, and similarly `delpaths` for an array of path arrays.

