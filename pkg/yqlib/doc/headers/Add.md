Add behaves differently according to the type of the LHS:
- arrays: concatenate
- number scalars: arithmetic addition (soon)
- string scalars: concatenate (soon)

Use `+=` as append assign for things like increment. `.a += .x` is equivalent to running `.a |= . + .x`.
