Add behaves differently according to the type of the LHS:
- arrays: concatenate
- number scalars: arithmetic addition
- string scalars: concatenate

Use `+=` as append assign for things like increment. Note that `.a += .x` is equivalent to running `.a = .a + .x`.
